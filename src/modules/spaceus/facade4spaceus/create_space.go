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

// CreateSpaceResult is a result of CreateSpace
type CreateSpaceResult struct {
	Space          dbo4spaceus.SpaceEntry
	ContactusSpace dal4contactus.ContactusSpaceEntry
	Member         dal4contactus.ContactEntry
}

func (result CreateSpaceResult) Records() []dal.Record {
	return []dal.Record{result.Space.Record, result.ContactusSpace.Record, result.Member.Record}
}

// CreateSpace creates SpaceIDs record
func CreateSpace(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.CreateSpaceRequest,
) (
	result CreateSpaceResult,
	err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
		if result, err = CreateSpaceTxWorker(ctx, tx, time.Now(), request, params); err != nil {
			return
		}
		recordsToInsert := result.Records()
		if err = tx.InsertMulti(ctx, recordsToInsert); err != nil {
			return
		}
		return
	})
	return
}

func CreateSpaceTxWorker(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	createdAt time.Time,
	request dto4spaceus.CreateSpaceRequest,
	params *dal4userus.UserWorkerParams,
) (
	result CreateSpaceResult,
	err error,
) {
	if request.Title == "" {
		result.Space.ID, _ = params.User.Data.GetFirstSpaceBriefBySpaceType(request.Type)
		if result.Space.ID != "" {
			result.Space = dbo4spaceus.NewSpaceEntry(result.Space.ID)
			result.ContactusSpace = dal4contactus.NewContactusSpaceEntry(result.Space.ID)
			err = tx.GetMulti(ctx, []dal.Record{result.Space.Record, result.ContactusSpace.Record})
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

	//if request.Type == "family" && request.Title == "" {
	//	request.Title = "Family"
	//}
	result.Space.Data = &dbo4spaceus.SpaceDbo{
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
	result.Space.Data.IncreaseVersion(createdAt, params.User.ID)
	result.Space.Data.CountryID = params.User.Data.CountryID

	if request.Type == "work" {
		zero := 0
		hundred := 100
		result.Space.Data.Metrics = []*dbo4spaceus.SpaceMetric{
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

	if err = result.Space.Data.Validate(); err != nil {
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

	result.Space = dbo4spaceus.NewSpaceEntryWithDbo(spaceID, result.Space.Data)
	//if err = tx.Insert(ctx, space.Record); err != nil {
	//	err = fmt.Errorf("failed to insert a new spaceDbo record: %w", err)
	//	return
	//}

	result.ContactusSpace = dal4contactus.NewContactusSpaceEntry(spaceID)

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
	result.ContactusSpace.Data.AddContact(userSpaceContactID, &spaceMember)

	//if err = tx.Insert(ctx, contactusSpace.Record); err != nil {
	//	err = fmt.Errorf("failed to insert a new spaceDbo contactus record: %w", err)
	//	return
	//}

	userSpaceBrief := dbo4userus.UserSpaceBrief{
		SpaceBrief:    result.Space.Data.SpaceBrief,
		UserContactID: userSpaceContactID,
		Roles:         roles,
	}

	params.UserUpdates = append(params.UserUpdates, params.User.Data.SetSpaceBrief(spaceID, &userSpaceBrief)...)

	params.UserUpdates = append(params.UserUpdates, dbo4linkage.UpdateRelatedIDs(&params.User.Data.WithRelated, &params.User.Data.WithRelatedIDs)...)

	if err = params.User.Data.Validate(); err != nil {
		err = fmt.Errorf("user record is not valid after adding new space info: %v", err)
		return
	}
	params.User.Record.MarkAsChanged()
	//if params.User.Record.Exists() {
	//	// Will be updated by RunUserWorker
	//} else {
	//	if err = tx.Insert(ctx, params.User.Record); err != nil {
	//		err = fmt.Errorf("failed to insert new user record: %w", err)
	//		return
	//	}
	//}

	spaceMember.Roles = roles
	if result.Member, err = CreateMemberEntryFromBrief(spaceID, userSpaceContactID, spaceMember, createdAt, params.User.ID); err != nil {
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
