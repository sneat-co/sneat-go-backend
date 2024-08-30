package facade4spaceus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/gosimple/slug"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"strings"
	"time"
)

type CreateSpaceResult struct {
	Space dbo4spaceus.SpaceEntry `json:"-"`
	User  dbo4userus.UserEntry   `json:"-"`
}

func CreateFamilySpace(
	ctx context.Context, userCtx facade.UserContext,
) (
	space dbo4spaceus.SpaceEntry, contactusSpace dal4contactus.ContactusSpaceEntry, err error,
) {
	request := dto4spaceus.CreateSpaceRequest{Type: core4spaceus.SpaceTypeFamily}
	return CreateSpace(ctx, userCtx, request)
}

// CreateSpace creates SpaceIDs record
func CreateSpace(ctx context.Context, userCtx facade.UserContext, request dto4spaceus.CreateSpaceRequest) (space dbo4spaceus.SpaceEntry, contactusSpace dal4contactus.ContactusSpaceEntry, err error) {
	err = dal4userus.RunUserWorker(ctx, userCtx, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) error {
		space, contactusSpace, err = CreateSpaceTxWorker(ctx, tx, request, params)
		return err
	})
	return
}

func CreateSpaceTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, request dto4spaceus.CreateSpaceRequest, params *dal4userus.UserWorkerParams) (space dbo4spaceus.SpaceEntry, contactusSpace dal4contactus.ContactusSpaceEntry, err error) {
	now := time.Now()
	if request.Title == "" {
		space.ID, _ = params.User.Data.GetSpaceBriefByType(request.Type)
		if space.ID != "" {
			space = dbo4spaceus.NewSpaceEntry(space.ID)
			contactusSpace = dal4contactus.NewContactusSpaceEntry(space.ID)
			err = tx.GetMulti(ctx, []dal.Record{space.Record, contactusSpace.Record})
			return
		}
	}

	var userSpaceContactID string

	if userSpaceContactID, err = person.GenerateIDFromNameOrRandom(params.User.Data.Names, nil); err != nil {
		err = fmt.Errorf("failed to generate  member ContactID: %w", err)
		return
	}

	roles := []string{
		const4contactus.SpaceMemberRoleMember,
		const4contactus.SpaceMemberRoleCreator,
		const4contactus.SpaceMemberRoleOwner,
		const4contactus.SpaceMemberRoleContributor,
	}
	if request.Type == "family" {
		roles = append(roles, const4contactus.SpaceMemberRoleAdult)
	}

	if request.Type == "family" && request.Title == "" {
		request.Title = "Family"
	}
	space.Data = &dbo4spaceus.SpaceDbo{
		SpaceBrief: dbo4spaceus.SpaceBrief{
			Type:   request.Type,
			Title:  request.Title,
			Status: dbmodels.StatusActive,
		},
		WithUserIDs: dbmodels.WithUserIDs{
			UserIDs: []string{params.User.ID},
		},
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: now,
			},
			CreatedByField: with.CreatedByField{
				CreatedBy: params.User.ID,
			},
		},
		//WithUpdated: dbmodels.WithUpdated{}, // This is updated by IncreaseVersion()
		//WithMembers: models4memberus.WithMembers{}, // Moved to contactus module
		//NumberOf: map[string]int{
		//	"members": 1,
		//},
	}
	space.Data.IncreaseVersion(now, params.User.ID)
	space.Data.CountryID = params.User.Data.CountryID

	if request.Type == "work" {
		zero := 0
		hundred := 100
		space.Data.Metrics = []*dbo4spaceus.SpaceMetric{
			{ID: "cc", Title: "Code coverage", Type: "int", Mode: "SpaceIDs", Min: &zero, Max: &hundred},
			{ID: "bb", Title: "Build is broken", Type: "bool", Mode: "SpaceIDs", Bool: &dbo4spaceus.BoolMetric{
				True:  &dbo4spaceus.BoolMetricVal{Label: "Yes", Color: "danger"},
				False: &dbo4spaceus.BoolMetricVal{Label: "No", Color: "success"},
			}},
			{ID: "wfh", Title: "Working From Home", Type: "bool", Mode: "personal", Bool: &dbo4spaceus.BoolMetric{
				True:  &dbo4spaceus.BoolMetricVal{Label: "Yes", Color: "tertiary"},
				False: &dbo4spaceus.BoolMetricVal{Label: "No", Color: "secondary"},
			}},
		}
	}

	if err = space.Data.Validate(); err != nil {
		err = fmt.Errorf("spaceDbo reacord is not valid: %w", err)
		return
	}

	var spaceID string
	title := request.Title
	if request.Type == "family" {
		title = ""
	}
	spaceID, err = getUniqueSpaceID(ctx, tx, title)
	if err != nil {
		err = fmt.Errorf("failed to get an unique ContactID for a new spaceDbo: %w", err)
		return
	}

	space = dbo4spaceus.NewSpaceEntryWithDbo(spaceID, space.Data)
	if err = tx.Insert(ctx, space.Record); err != nil {
		err = fmt.Errorf("failed to insert a new spaceDbo record: %w", err)
		return
	}

	contactusSpace = dal4contactus.NewContactusSpaceEntry(spaceID)

	spaceMember := params.User.Data.ContactBrief // This should copy data from user's contact brief as it's not a pointer

	spaceMember.UserID = params.User.ID
	spaceMember.Roles = roles
	if spaceMember.Gender == "" {
		spaceMember.Gender = "unknown"
	}
	if params.User.Data.Defaults != nil && len(params.User.Data.Defaults.ShortNames) > 0 {
		spaceMember.ShortTitle = params.User.Data.Defaults.ShortNames[0].Name
	}
	//if len(spaceMember.Emails) == 0 && len(user.Emails) > 0 {
	//	spaceMember.Emails = user.Emails
	//}
	//if len(spaceMember.Phones) == 0 && len(user.Data.Phones) > 0 {
	//	spaceMember.Phones = user.Data.Phones
	//}
	contactusSpace.Data.AddContact(userSpaceContactID, &spaceMember)

	if err = tx.Insert(ctx, contactusSpace.Record); err != nil {
		err = fmt.Errorf("failed to insert a new spaceDbo contactus record: %w", err)
		return
	}

	userSpaceBrief := dbo4userus.UserSpaceBrief{
		SpaceBrief:    space.Data.SpaceBrief,
		UserContactID: userSpaceContactID,
		Roles:         roles,
	}

	params.UserUpdates = append(params.UserUpdates, params.User.Data.SetSpaceBrief(spaceID, &userSpaceBrief)...)

	params.UserUpdates = append(params.UserUpdates, dbo4linkage.UpdateRelatedIDs(&params.User.Data.WithRelated, &params.User.Data.WithRelatedIDs)...)

	if err = params.User.Data.Validate(); err != nil {
		err = fmt.Errorf("user record is not valid after adding new space info: %v", err)
		return
	}
	if params.User.Record.Exists() {
		// Will be updated by RunUserWorker
	} else {
		if err = tx.Insert(ctx, params.User.Record); err != nil {
			err = fmt.Errorf("failed to insert new user record: %w", err)
			return
		}
	}

	spaceMember.Roles = roles
	if _, err = CreateMemberRecordFromBrief(ctx, tx, spaceID, userSpaceContactID, spaceMember, now, params.User.ID); err != nil {
		err = fmt.Errorf("failed to create member's record: %w", err)
		return
	}

	return
}

func getUniqueSpaceID(ctx context.Context, getter dal.ReadSession, title string) (spaceID string, err error) {
	if title == "" || const4contactus.IsKnownSpaceMemberRole(title, []string{}) {
		spaceID = random.ID(5)
	} else {
		spaceID = strings.Replace(slug.Make(title), "-", "", -1)
	}
	const maxAttemptsCount = 9
	for i := 0; i <= maxAttemptsCount; i++ {
		if i == maxAttemptsCount {
			return "", errors.New("too many attempts to get an unique space ContactID")
		}
		spaceID = strings.ToLower(spaceID)
		teamKey := dal.NewKeyWithID(dbo4spaceus.SpacesCollection, spaceID)
		teamRecord := dal.NewRecordWithData(teamKey, nil)
		if err = getter.Get(ctx, teamRecord); dal.IsNotFound(err) {
			return spaceID, nil
		} else if err != nil {
			return spaceID, err
		}
		if i == 0 && title != "" {
			spaceID += "_"
		}
		spaceID += random.ID(1)
	}
	return spaceID, nil
}
