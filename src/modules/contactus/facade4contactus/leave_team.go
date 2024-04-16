package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// LeaveTeam leaves team
func LeaveTeam(ctx context.Context, userContext facade.User, request dto4contactus.ContactRequestWithOptionalMessage) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return dal4contactus.RunContactusTeamWorker(ctx, userContext, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams) (err error) {

			if err := tx.GetMulti(ctx, []dal.Record{params.TeamModuleEntry.Record}); err != nil {
				return fmt.Errorf("failed to get team record: %w", err)
			}

			uid := userContext.GetID()
			user := new(models4userus.UserDto)
			userKey := dal.NewKeyWithID(models4userus.UsersCollection, uid)
			userRecord := dal.NewRecordWithData(userKey, nil)
			if err = tx.Get(ctx, userRecord); err != nil {
				return
			}

			// Update team record
			{
				team := params.Team
				var updates []dal.Update
				var memberUserID string
				memberUserID, updates, err = removeTeamMember(team, params.TeamModuleEntry,
					func(_ string, m *briefs4contactus.ContactBrief) bool {
						return m.UserID == uid
					})
				if err != nil || len(updates) == 0 {
					return
				}
				if memberUserID != uid {
					err = fmt.Errorf("user ID does not match members record: memberUserID[%v] != ctx.UserID[%v]: %w",
						memberUserID, uid, facade.ErrBadRequest)
					return
				}
				if err = team.Data.Validate(); err != nil {
					return fmt.Errorf("team reacord is not valid: %v", err)
				}
				if len(params.TeamModuleEntry.Data.Contacts) == 0 || len(team.Data.UserIDs) == 0 {
					if err = tx.Delete(ctx, team.Key); err != nil {
						return err
					}
				} else {
					if err = txUpdateMemberGroup(ctx, tx, params.Started, uid, params.Team.Data, params.Team.Key, updates); err != nil {
						return
					}
				}
			}

			// Update user
			{
				if userRecord.Exists() {
					update := updateUserRecordOnTeamMemberRemoved(user, request.TeamID)
					if update != nil {
						if err = txUpdate(ctx, tx, userKey, []dal.Update{*update}); err != nil {
							return err
						}
					}
				}
			}
			return err
		})
}
