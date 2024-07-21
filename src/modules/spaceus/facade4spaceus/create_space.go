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
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
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
	Space dal4spaceus.SpaceEntry `json:"-"`
	User  dbo4userus.UserEntry   `json:"-"`
}

// CreateSpace creates SpaceIDs record
func CreateSpace(ctx context.Context, userContext facade.User, request dto4spaceus.CreateSpaceRequest) (response CreateSpaceResult, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	db := facade.GetDatabase(ctx)

	// We do not use facade4userus.RunUserWorker dues to cycle dependency
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		response, err = createSpaceTxWorker(ctx, userContext, tx, request)
		return err
	})
	return response, err
}

func createSpaceTxWorker(ctx context.Context, userContext facade.User, tx dal.ReadwriteTransaction, request dto4spaceus.CreateSpaceRequest) (response CreateSpaceResult, err error) {
	now := time.Now()
	userID := userContext.GetID()
	if strings.TrimSpace(userID) == "" {
		return response, facade.ErrUnauthenticated
	}
	var userSpaceContactID string

	user := dbo4userus.NewUserEntry(userID)
	response.User = user

	if err = tx.Get(ctx, user.Record); err != nil {
		return
	}

	if request.Title == "" {
		spaceID, _ := user.Data.GetSpaceBriefByType(request.Type)
		if spaceID != "" {
			response.Space.ID = spaceID
			if space, err := GetSpaceByID(ctx, tx, spaceID); err != nil {
				return response, err
			} else {
				response.Space = space
				return response, nil
			}
		}
	}

	userSpaceContactID, err = person.GenerateIDFromNameOrRandom(user.Data.Names, nil)
	if err != nil {
		return response, fmt.Errorf("failed to generate  member ID: %w", err)
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
	teamDbo := &dbo4spaceus.SpaceDbo{
		SpaceBrief: dbo4spaceus.SpaceBrief{
			Type:   request.Type,
			Title:  request.Title,
			Status: dbmodels.StatusActive,
		},
		WithUserIDs: dbmodels.WithUserIDs{
			UserIDs: []string{userID},
		},
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: now,
			},
			CreatedByField: with.CreatedByField{
				CreatedBy: userID,
			},
		},
		//WithUpdated: dbmodels.WithUpdated{}, // This is updated by IncreaseVersion()
		//WithMembers: models4memberus.WithMembers{}, // Moved to contactus module
		//NumberOf: map[string]int{
		//	"members": 1,
		//},
	}
	teamDbo.IncreaseVersion(now, userID)
	teamDbo.CountryID = user.Data.CountryID
	if request.Type == "work" {
		zero := 0
		hundred := 100
		teamDbo.Metrics = []*dbo4spaceus.SpaceMetric{
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

	if err = teamDbo.Validate(); err != nil {
		return response, fmt.Errorf("spaceDbo reacord is not valid: %w", err)
	}

	var spaceID string
	title := request.Title
	if request.Type == "family" {
		title = ""
	}
	spaceID, err = getUniqueSpaceID(ctx, tx, title)
	if err != nil {
		return response, fmt.Errorf("failed to get an unique ID for a new teamDbo: %w", err)
	}
	teamKey := dal.NewKeyWithID(dal4spaceus.SpacesCollection, spaceID)

	teamRecord := dal.NewRecordWithData(teamKey, teamDbo)
	if err = tx.Insert(ctx, teamRecord); err != nil {
		return response, fmt.Errorf("failed to insert a new teamDbo record: %w", err)
	}

	teamContactus := dal4contactus.NewContactusSpaceModuleEntry(spaceID)

	teamMember := user.Data.ContactBrief // This should copy data from user's contact brief as it's not a pointer

	teamMember.UserID = userID
	teamMember.Roles = roles
	if teamMember.Gender == "" {
		teamMember.Gender = "unknown"
	}
	if user.Data.Defaults != nil && len(user.Data.Defaults.ShortNames) > 0 {
		teamMember.ShortTitle = user.Data.Defaults.ShortNames[0].Name
	}
	//if len(teamMember.Emails) == 0 && len(user.Emails) > 0 {
	//	teamMember.Emails = user.Emails
	//}
	//if len(teamMember.Phones) == 0 && len(user.Data.Phones) > 0 {
	//	teamMember.Phones = user.Data.Phones
	//}
	teamContactus.Data.AddContact(userSpaceContactID, &teamMember)

	if err := tx.Insert(ctx, teamContactus.Record); err != nil {
		return response, fmt.Errorf("failed to insert a new teamDbo contactus record: %w", err)
	}

	userSpaceBrief := dbo4userus.UserSpaceBrief{
		SpaceBrief:    teamDbo.SpaceBrief,
		UserContactID: userSpaceContactID,
		Roles:         roles,
	}

	if user.Data.Spaces == nil {
		user.Data.Spaces = make(map[string]*dbo4userus.UserSpaceBrief, 1)
	}
	updates := user.Data.SetSpaceBrief(spaceID, &userSpaceBrief)

	updates = append(updates, dbo4linkage.UpdateRelatedIDs(&user.Data.WithRelated, &user.Data.WithRelatedIDs)...)

	if err = user.Data.Validate(); err != nil {
		return response, fmt.Errorf("user record is not valid after adding new team info: %v", err)
	}
	if user.Record.Exists() {
		if err = tx.Update(ctx, user.Key, updates); err != nil {
			return response, fmt.Errorf("failed to update user record with a new teamDbo info: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, user.Record); err != nil {
			return response, fmt.Errorf("failed to insert new user record: %w", err)
		}
	}

	teamMember.Roles = roles
	if _, err = CreateMemberRecordFromBrief(ctx, tx, spaceID, userSpaceContactID, teamMember, now, userID); err != nil {
		return response, fmt.Errorf("failed to create member's record: %w", err)
	}

	response.Space.ID = spaceID
	response.Space.Data = teamDbo
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
			return "", errors.New("too many attempts to get an unique space ID")
		}
		spaceID = strings.ToLower(spaceID)
		teamKey := dal.NewKeyWithID(dal4spaceus.SpacesCollection, spaceID)
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
