package facade4teamus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/gosimple/slug"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"strings"
	"time"
)

type CreateTeamResult struct {
	Team dal4teamus.TeamContext    `json:"-"`
	User models4userus.UserContext `json:"-"`
}

// CreateTeam creates TeamIDs record
func CreateTeam(ctx context.Context, userContext facade.User, request dto4teamus.CreateTeamRequest) (response CreateTeamResult, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	db := facade.GetDatabase(ctx)

	// We do not use facade4userus.RunUserWorker dues to cycle dependency
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		response, err = createTeamTxWorker(ctx, userContext, tx, request)
		return err
	})
	return response, err
}

func createTeamTxWorker(ctx context.Context, userContext facade.User, tx dal.ReadwriteTransaction, request dto4teamus.CreateTeamRequest) (response CreateTeamResult, err error) {
	now := time.Now()
	userID := userContext.GetID()
	if strings.TrimSpace(userID) == "" {
		return response, facade.ErrUnauthenticated
	}
	var userTeamContactID string

	user := models4userus.NewUserContext(userID)
	response.User = user

	if err = tx.Get(ctx, user.Record); err != nil {
		return
	}

	if request.Title == "" {
		teamID, _ := user.Dbo.GetTeamBriefByType(request.Type)
		if teamID != "" {
			response.Team.ID = teamID
			if team, err := GetTeamByID(ctx, tx, teamID); err != nil {
				return response, err
			} else {
				response.Team = team
				return response, nil
			}
		}
	}

	userTeamContactID, err = person.GenerateIDFromNameOrRandom(user.Dbo.Names, nil)
	if err != nil {
		return response, fmt.Errorf("failed to generate  member ID: %w", err)
	}

	roles := []string{
		const4contactus.TeamMemberRoleMember,
		const4contactus.TeamMemberRoleCreator,
		const4contactus.TeamMemberRoleOwner,
		const4contactus.TeamMemberRoleContributor,
	}
	if request.Type == "family" {
		roles = append(roles, const4contactus.TeamMemberRoleAdult)
	}

	if request.Type == "family" && request.Title == "" {
		request.Title = "Family"
	}
	teamDbo := &models4teamus.TeamDbo{
		TeamBrief: models4teamus.TeamBrief{
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
		NumberOf: map[string]int{
			"members": 1,
		},
	}
	teamDbo.IncreaseVersion(now, userID)
	teamDbo.CountryID = user.Dbo.CountryID
	if request.Type == "work" {
		zero := 0
		hundred := 100
		teamDbo.Metrics = []*models4teamus.TeamMetric{
			{ID: "cc", Title: "Code coverage", Type: "int", Mode: "TeamIDs", Min: &zero, Max: &hundred},
			{ID: "bb", Title: "Build is broken", Type: "bool", Mode: "TeamIDs", Bool: &models4teamus.BoolMetric{
				True:  &models4teamus.BoolMetricVal{Label: "Yes", Color: "danger"},
				False: &models4teamus.BoolMetricVal{Label: "No", Color: "success"},
			}},
			{ID: "wfh", Title: "Working From Home", Type: "bool", Mode: "personal", Bool: &models4teamus.BoolMetric{
				True:  &models4teamus.BoolMetricVal{Label: "Yes", Color: "tertiary"},
				False: &models4teamus.BoolMetricVal{Label: "No", Color: "secondary"},
			}},
		}
	}

	if err = teamDbo.Validate(); err != nil {
		return response, fmt.Errorf("teamDbo reacord is not valid: %w", err)
	}

	var teamID string
	title := request.Title
	if request.Type == "family" {
		title = ""
	}
	teamID, err = getUniqueTeamID(ctx, tx, title)
	if err != nil {
		return response, fmt.Errorf("failed to get an unique ID for a new teamDbo: %w", err)
	}
	teamKey := dal.NewKeyWithID(dal4teamus.TeamsCollection, teamID)

	teamRecord := dal.NewRecordWithData(teamKey, teamDbo)
	if err = tx.Insert(ctx, teamRecord); err != nil {
		return response, fmt.Errorf("failed to insert a new teamDbo record: %w", err)
	}

	teamContactus := dal4contactus.NewContactusTeamModuleEntry(teamID)

	teamMember := user.Dbo.ContactBrief // This should copy data from user's contact brief as it's not a pointer

	teamMember.UserID = userID
	teamMember.Roles = roles
	if teamMember.Gender == "" {
		teamMember.Gender = "unknown"
	}
	if user.Dbo.Defaults != nil && len(user.Dbo.Defaults.ShortNames) > 0 {
		teamMember.ShortTitle = user.Dbo.Defaults.ShortNames[0].Name
	}
	//if len(teamMember.Emails) == 0 && len(user.Emails) > 0 {
	//	teamMember.Emails = user.Emails
	//}
	//if len(teamMember.Phones) == 0 && len(user.Dbo.Phones) > 0 {
	//	teamMember.Phones = user.Dbo.Phones
	//}
	teamContactus.Data.AddContact(userTeamContactID, &teamMember)

	if err := tx.Insert(ctx, teamContactus.Record); err != nil {
		return response, fmt.Errorf("failed to insert a new teamDbo contactus record: %w", err)
	}

	userTeamBrief := models4userus.UserTeamBrief{
		TeamBrief:     teamDbo.TeamBrief,
		UserContactID: userTeamContactID,
		Roles:         roles,
	}

	if user.Dbo.Teams == nil {
		user.Dbo.Teams = make(map[string]*models4userus.UserTeamBrief, 1)
	}
	updates := user.Dbo.SetTeamBrief(teamID, &userTeamBrief)

	updates = append(updates, user.Dbo.UpdateRelatedIDs()...)

	if err = user.Dbo.Validate(); err != nil {
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
	if _, err = CreateMemberRecordFromBrief(ctx, tx, teamID, userTeamContactID, teamMember, now, userID); err != nil {
		return response, fmt.Errorf("failed to create member's record: %w", err)
	}

	response.Team.ID = teamID
	response.Team.Data = teamDbo
	return
}

func getUniqueTeamID(ctx context.Context, getter dal.ReadSession, title string) (teamID string, err error) {
	if title == "" || const4contactus.IsKnownTeamMemberRole(title, []string{}) {
		teamID = random.ID(5)
	} else {
		teamID = strings.Replace(slug.Make(title), "-", "", -1)
	}
	const maxAttemptsCount = 9
	for i := 0; i <= maxAttemptsCount; i++ {
		if i == maxAttemptsCount {
			return "", errors.New("too many attempts to get an unique team ID")
		}
		teamID = strings.ToLower(teamID)
		teamKey := dal.NewKeyWithID(dal4teamus.TeamsCollection, teamID)
		teamRecord := dal.NewRecordWithData(teamKey, nil)
		if err = getter.Get(ctx, teamRecord); dal.IsNotFound(err) {
			return teamID, nil
		} else if err != nil {
			return teamID, err
		}
		if i == 0 && title != "" {
			teamID += "_"
		}
		teamID += random.ID(1)
	}
	return teamID, nil
}
