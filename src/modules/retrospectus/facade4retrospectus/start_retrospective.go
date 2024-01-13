package facade4retrospectus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/models4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dal4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/models4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// StartRetrospective starts retro
func StartRetrospective(ctx context.Context, userContext facade.User, request StartRetrospectiveRequest) (response *RetrospectiveResponse, isNewRetrospective bool, err error) {
	uid := userContext.GetID()

	teamKey := newTeamKey(request.TeamID)

	retrospective := new(models4retrospectus.Retrospective)

	err = dal4contactus.RunContactusTeamWorker(ctx, userContext, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams) (err error) {
			team := params.Team

			if err := tx.GetMulti(ctx, []dal.Record{params.Team.Record, params.TeamModuleEntry.Record}); err != nil {
				return err
			}

			if !params.Team.Data.HasUserID(uid) {
				return errors.New("current user does not belong to team (uid is missing in team.userIDs)")
			}
			if !params.TeamModuleEntry.Data.HasUserID(uid) {
				return errors.New("current user does not belong to team (uid is missing in contactusTeam.userIDs)")
			}
			if _, b := params.TeamModuleEntry.Data.GetContactBriefByUserID(uid); b == nil {
				return errors.New("current user does not have contact brief in contactusTeam.Data.Contacts")
			}

			var retroTeam dal4retrospectus.RetroTeam
			retroTeam, err = dal4retrospectus.GetRetroTeam(ctx, tx, request.TeamID)
			if err != nil && !dal.IsNotFound(err) {
				return err
			}

			teamChanged := false // All reads should be before any write in transaction
			if request.MeetingID == UpcomingRetrospectiveID {
				var activeRetroID string
				if retroTeam.Data.Active != nil {
					activeRetroID = retroTeam.Data.Active.ID
				}
				if activeRetroID == "" {
					request.MeetingID = params.Started.Format("2006-01-02")

					retroTeam.Data.Active = &models4teamus.TeamMeetingInfo{
						ID:      request.MeetingID,
						Started: &params.Started,
					}
					retroTeam.Data.UpcomingRetro = nil
					teamChanged = true
				} else {
					request.MeetingID = activeRetroID
				}
			} else if activeRetrospective := retroTeam.Data.ActiveRetro(); activeRetrospective.ID == request.MeetingID {
				return nil
			} else if activeRetrospective.ID == "" {
				retroTeam.Data.Active = &models4teamus.TeamMeetingInfo{
					ID:      request.MeetingID,
					Started: &params.Started,
				}
				teamChanged = true
			} else {
				return fmt.Errorf("an attempt to start a new retrospective while another one in progress (new: %v, active: %v)", request.MeetingID, activeRetrospective.ID)
			}

			byUser := dbmodels.ByUser{UID: uid}
			timer := models4meetingus.Timer{
				By:     byUser,
				At:     params.Started,
				Status: models4meetingus.TimerStatusActive,
			}

			retrospectiveKey := models4retrospectus.NewRetrospectiveKey(request.MeetingID, teamKey)
			retrospectiveRecord := dal.NewRecordWithData(retrospectiveKey, retrospective)
			if err = tx.Get(ctx, retrospectiveRecord); err != nil {
				if dal.IsNotFound(err) {
					isNewRetrospective = true
				} else {
					return fmt.Errorf("failed to check retrospetive record for existence: %w", err)
				}
			} else if err = retrospectiveRecord.Error(); err != nil {
				if dal.IsNotFound(err) {
					isNewRetrospective = true
				} else {
					return fmt.Errorf("retrospectiveRecord.Error(): %w", err)
				}
			} else if !retrospectiveRecord.Exists() {
				isNewRetrospective = true
			}

			//if err = txGetRetrospective(ctx, tx, retrospectiveRecord); err != nil {
			//	return err
			//}

			isNewRetrospective = !retrospectiveRecord.Exists()

			var usersWithRetroItems map[string]userRetroItems

			if isNewRetrospective && retroTeam.Data.UpcomingRetro != nil {
				if usersWithRetroItems, err = getUsersWithRetroItems(ctx, tx, team, retroTeam); err != nil {
					return err
				}
			}

			if teamChanged { // All reads should be before any write in transaction
				if err = txUpdateTeam(ctx, tx, params.Started, team, []dal.Update{
					{Field: "activeMeetings.retrospective", Value: request.MeetingID},
				}); err != nil {
					return err
				}
			}

			var retrospectiveUpdates []dal.Update

			if isNewRetrospective {
				retrospective = &models4retrospectus.Retrospective{
					Stage:          models4retrospectus.StageFeedback,
					StartedBy:      &byUser,
					TimeStarted:    &params.Started,
					TimeLastAction: &params.Started,
					Meeting: models4meetingus.Meeting{
						Version: 1,
						WithUserIDs: dbmodels.WithUserIDs{
							UserIDs: team.Data.UserIDs,
						},
						Timer: &timer,
					},
					Settings: models4retrospectus.RetrospectiveSettings{
						MaxVotesPerUser: models4retrospectus.DefaultMaxVotesPerUser,
					},
				}
				for contactID, contact := range params.TeamModuleEntry.Data.GetContactBriefsByRoles(const4contactus.TeamMemberRoleMember) {
					retrospective.AddContact(team.ID, contactID, &models4meetingus.MeetingMemberBrief{
						ContactBrief: *contact,
					})
				}
			} else {
				if retrospective != nil && retrospective.TimeStarted != nil { // Already started
					return nil
				}
				if retrospective == nil {
					retrospective = &models4retrospectus.Retrospective{}
				}
				retrospective.Timer = &timer
				retrospective.StartedBy = &byUser
				retrospective.Stage = models4retrospectus.StageFeedback

				retrospectiveUpdates = []dal.Update{
					{Field: "stage", Value: retrospective.Stage},
					{Field: "timer", Value: retrospective.Timer},
					{Field: "startedBy", Value: retrospective.StartedBy},
				}
			}

			if len(usersWithRetroItems) > 0 {
				for userID, userRetroItems := range usersWithRetroItems { // TODO: use go routine to run in parallel?
					if len(userRetroItems.byType) == 0 {
						continue
					}
					retroUserCounts := retrospective.CountsByMemberAndType[userID]
					if retroUserCounts == nil {
						retroUserCounts = make(map[string]int, len(userRetroItems.byType))
					}
					memberID, _ := retrospective.GetContactBriefByUserID(userID)
					for itemType, items := range userRetroItems.byType {
						itemsCount := len(items)
						retroUserCounts[itemType] = itemsCount
						retrospectiveUpdates = append(retrospectiveUpdates, dal.Update{
							Field: fmt.Sprintf("countsByMemberAndType.%v.%v", memberID, itemType),
							Value: itemsCount,
						})
					}
				}
			}

			if isNewRetrospective {
				if err = txCreateRetrospective(ctx, tx, retrospectiveKey, retrospective); err != nil {
					return err
				}
			} else {
				if err = txUpdateRetrospective(ctx, tx, retrospectiveKey, retrospective, retrospectiveUpdates); err != nil {
					return err
				}
			}
			return nil
		})
	if err != nil {
		return
	}
	return &RetrospectiveResponse{ID: request.MeetingID, Data: retrospective}, isNewRetrospective, err
}

type userRetroItems struct {
	//user   *models4userus.UserDto
	byType models4retrospectus.RetroItemsByType
}

func getUsersWithRetroItems(ctx context.Context, tx dal.ReadwriteTransaction, team dal4teamus.TeamContext, retroTeam dal4retrospectus.RetroTeam) (usersWithRetroItemByUserID map[string]userRetroItems, err error) {
	teamUsersCount := len(team.Data.UserIDs)
	usersWithRetroItemByUserID = make(map[string]userRetroItems, teamUsersCount)
	userIDs := make([]string, 0, teamUsersCount)
	for userID := range retroTeam.Data.UpcomingRetro.ItemsByUserAndType {
		userIDs = append(userIDs, userID)
	}
	userKeys := models4userus.NewUserKeys(userIDs)
	var usersRecords []dal.Record // []*models4userus.UserDto
	users := make([]*models4userus.UserDto, len(userKeys))
	for i, userKey := range userKeys {
		users[i] = new(models4userus.UserDto)
		usersRecords[i] = dal.NewRecordWithData(userKey, users)
	}
	err = facade4userus.TxGetUsers(ctx, tx, usersRecords)
	if err != nil {
		return nil, err
	}
	for i, userRecord := range usersRecords {
		if !userRecord.Exists() {
			continue
		}
		teamInfo := users[i].GetUserTeamInfoByID(team.ID)
		if teamInfo == nil {
			continue
		}

		//if len(teamInfo.RetroItems) == 0 {
		//	continue
		//}
		//usersWithRetroItemByUserID[userIDs[i]] = userRetroItems{user: users[i], byType: teamInfo.RetroItems}
	}
	return usersWithRetroItemByUserID, err
}
