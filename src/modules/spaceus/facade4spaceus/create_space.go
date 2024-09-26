package facade4spaceus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
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

// CreateSpaceParams is a result of CreateSpace
type CreateSpaceParams struct {
	UserUpdates    []dal.Update
	User           dbo4userus.UserEntry
	Space          dbo4spaceus.SpaceEntry
	ContactusSpace dal4contactus.ContactusSpaceEntry
	Member         dal4contactus.ContactEntry
	*record.WithRecordChanges
}

// CreateSpace creates SpaceIDs record
func CreateSpace(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.CreateSpaceRequest,
) (
	createSpaceParams CreateSpaceParams,
	err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
		createSpaceParams = CreateSpaceParams{
			User:              params.User,
			WithRecordChanges: new(record.WithRecordChanges),
		}
		now := time.Now()
		if err = CreateSpaceTxWorker(ctx, tx, now, request, &createSpaceParams); err != nil {
			return
		}
		if err = createSpaceParams.ApplyChanges(ctx, tx); err != nil {
			err = fmt.Errorf("failed to apply changes returned by CreateSpaceTxWorker(): %w", err)
		}
		params.UserUpdates = append(params.UserUpdates, createSpaceParams.UserUpdates...)
		return
	})
	return
}

func CreateSpaceTxWorker(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	createdAt time.Time,
	request dto4spaceus.CreateSpaceRequest,
	params *CreateSpaceParams,
) (
	err error,
) {
	if request.Title == "" {
		params.Space.ID, _ = params.User.Data.GetFirstSpaceBriefBySpaceType(request.Type)
		if params.Space.ID != "" {
			params.Space = dbo4spaceus.NewSpaceEntry(params.Space.ID)
			params.ContactusSpace = dal4contactus.NewContactusSpaceEntry(params.Space.ID)
			err = tx.GetMulti(ctx, []dal.Record{params.Space.Record, params.ContactusSpace.Record})
			return
		}
	}

	var userSpaceContactID string

	if userSpaceContactID, err = person.GenerateIDFromNameOrRandom(params.User.Data.Names, nil); err != nil {
		err = fmt.Errorf("failed to generate  member ContactID: %w", err)
		return
	}

	//if request.Type == "family" && request.Title == "" {
	//	request.Title = "Family"
	//}
	params.Space.Data = &dbo4spaceus.SpaceDbo{
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
				CreatedAt: createdAt,
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
	params.Space.Data.IncreaseVersion(createdAt, params.User.ID)
	params.Space.Data.CountryID = params.User.Data.CountryID

	if request.Type == "work" {
		zero := 0
		hundred := 100
		params.Space.Data.Metrics = []*dbo4spaceus.SpaceMetric{
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

	if err = params.Space.Data.Validate(); err != nil {
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
		err = fmt.Errorf("failed to get an unique spaceID for a new spaceDbo: %w", err)
		return
	}

	params.Space = dbo4spaceus.NewSpaceEntryWithDbo(spaceID, params.Space.Data)

	params.ContactusSpace = dal4contactus.NewContactusSpaceEntry(spaceID)

	spaceContactBrief := params.User.Data.ContactBrief // This should copy data from user's contact brief as it's not a pointer

	spaceContactBrief.UserID = params.User.ID

	roles := []string{
		const4contactus.SpaceMemberRoleMember,
		const4contactus.SpaceMemberRoleCreator,
		const4contactus.SpaceMemberRoleOwner,
		const4contactus.SpaceMemberRoleContributor,
	}
	if request.Type == core4spaceus.SpaceTypeFamily {
		roles = append(roles, const4contactus.SpaceMemberRoleAdult)
	}
	spaceContactBrief.Roles = roles

	if spaceContactBrief.Gender == "" {
		spaceContactBrief.Gender = "unknown"
	}
	if params.User.Data.Defaults != nil && len(params.User.Data.Defaults.ShortNames) > 0 {
		spaceContactBrief.ShortTitle = params.User.Data.Defaults.ShortNames[0].Name
	}
	params.ContactusSpace.Data.AddContact(userSpaceContactID, &spaceContactBrief)
	if err = params.ContactusSpace.Data.Validate(); err != nil {
		err = fmt.Errorf("newly created contactus space record is not valid: %w", err)
		return
	}
	params.ContactusSpace.Record.MarkAsChanged()
	params.ContactusSpace.Record.SetError(nil)
	params.QueueForInsert(params.ContactusSpace.Record)

	params.UserUpdates = append(params.UserUpdates, updateUserWithSpaceBrief(params.User, spaceID, userSpaceContactID, params.Space.Data.SpaceBrief, roles)...)

	if err = params.User.Data.Validate(); err != nil {
		err = fmt.Errorf("user record is not valid after adding new space info: %v", err)
		return
	}
	params.User.Record.MarkAsChanged()

	if params.Member, err = NewMemberContactEntryFromContactBrief(spaceID, userSpaceContactID, spaceContactBrief, createdAt, params.User.ID); err != nil {
		err = fmt.Errorf("failed to create member's record: %w", err)
		return
	}
	params.QueueForInsert(params.Member.Record)

	if err = params.Space.Data.Validate(); err != nil {
		params.Space.Record.SetError(err)
		return fmt.Errorf("newly created space data is not valid: %w", err)
	}
	params.Space.Record.SetError(nil)
	params.QueueForInsert(params.Space.Record)

	return
}

func updateUserWithSpaceBrief(user dbo4userus.UserEntry, spaceID, userSpaceContactID string, spaceBrief dbo4spaceus.SpaceBrief, spaceUserRoles []string) (updates []dal.Update) {
	userSpaceBrief := dbo4userus.UserSpaceBrief{
		SpaceBrief:    spaceBrief,
		UserContactID: userSpaceContactID,
		Roles:         spaceUserRoles,
	}
	updates = append(updates, user.Data.SetSpaceBrief(spaceID, &userSpaceBrief)...)
	updates = append(updates, dbo4linkage.UpdateRelatedIDs(&user.Data.WithRelated, &user.Data.WithRelatedIDs)...)
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
