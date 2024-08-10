package facade2debtus

import (
	"context"
	"errors"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"github.com/strongo/slices"
)

type groupFacade struct {
}

var Group = groupFacade{}

func (groupFacade groupFacade) CreateGroup(c context.Context,
	groupEntity *models.GroupDbo,
	tgBotCode string,
	beforeGroupInsert func(tc context.Context, groupEntity *models.GroupDbo) (group models.GroupEntry, err error),
	afterGroupInsert func(c context.Context, group models.GroupEntry, user models.AppUser) (err error),
) (group models.GroupEntry, groupMember models.GroupMember, err error) {
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		user, err := User.GetUserByID(c, tx, groupEntity.CreatorUserID)
		if err != nil {
			return err
		}
		existingGroups := user.Data.ActiveGroups()

		if beforeGroupInsert != nil {
			if group, err = beforeGroupInsert(c, groupEntity); err != nil {
				return err
			}
		}

		var groupMembersChanged bool
		groupMembersChanged, _, memberIndex, member, members := groupEntity.AddOrGetMember(groupEntity.CreatorUserID, "", user.Data.FullName())
		member.Shares = 1
		members[memberIndex] = member
		groupEntity.SetGroupMembers(members)

		if group.ID == "" {
			for _, existingGroup := range existingGroups {
				if existingGroup.Name == groupEntity.Name {
					return errors.New("Duplicate group name")
				}
			}
			if group, err = dtdal.Group.InsertGroup(c, tx, groupEntity); err != nil {
				return err
			}
		} else if groupMembersChanged {
			if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
				return err
			}
		}

		groupJson := models.UserGroupJson{
			ID:           group.ID,
			Name:         group.Data.Name,
			Note:         group.Data.Note,
			MembersCount: group.Data.MembersCount,
		}

		if tgBotCode != "" {
			for _, tgGroupBot := range groupJson.TgBots {
				if tgGroupBot == tgBotCode {
					goto botFound
				}
			}
			groupJson.TgBots = append(groupJson.TgBots, tgBotCode)
		botFound:
		}

		user.Data.SetActiveGroups(append(existingGroups, groupJson))

		if afterGroupInsert != nil {
			if err = afterGroupInsert(c, group, user); err != nil {
				return err
			}
		}

		if err = User.SaveUser(c, tx, user); err != nil {
			return err
		}
		if err = groupFacade.DelayUpdateGroupUsers(c, group.ID); err != nil {
			return err
		}
		return err
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}
	logus.Infof(c, "GroupEntry created, ID=%v", group.ID)
	return
}

type NewUser struct {
	Name string
	botsfwmodels.BotUserData
	ChatMember botsfw.WebhookActor
}

func (groupFacade) AddUsersToTheGroupAndOutstandingBills(c context.Context, groupID string, newUsers []NewUser) (group models.GroupEntry, newUsers2 []NewUser, err error) {
	logus.Debugf(c, "groupFacade.AddUsersToTheGroupAndOutstandingBills(groupID=%v, newUsers=%v)", groupID, newUsers)
	if groupID == "" {
		panic("groupID is empty string")
	}
	if len(newUsers) == 0 {
		panic("len(newUsers) == 0")
	}
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		changed := false
		if group, err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
			return
		}
		logus.Debugf(c, "group: %+v", group.Data)
		j := 0
		for _, newUser := range newUsers {
			_, isChanged, _, _, groupMembers := group.Data.AddOrGetMember(newUser.GetAppUserID(), "", newUser.Name)
			changed = changed || isChanged
			if isChanged {
				group.Data.SetGroupMembers(groupMembers)
				newUsers[j] = newUser
				j += 1
			}
		}
		newUsers = newUsers[:j]
		if changed {
			logus.Debugf(c, "group: %+v", group.Data)
			if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
				return
			}
			if err = Group.DelayUpdateGroupUsers(c, group.ID); err != nil {
				return err
			}
		}
		return
	}); err != nil {
		return group, newUsers, err
	}
	return group, newUsers, err
}

func (groupFacade) DelayUpdateGroupUsers(c context.Context, groupID string) error { // TODO: Move to DAL?
	if groupID == "" {
		panic("groupID is empty string")
	}
	return delayUpdateGroupUsers.EnqueueWork(c, delaying.With(common.QUEUE_USERS, "update-group-users", 0), groupID)
}

func updateGroupUsers(c context.Context, groupID string) (err error) {
	if groupID == "" {
		logus.Criticalf(c, "groupID is empty string")
		return nil
	}

	logus.Debugf(c, "updateGroupUsers(groupID=%v)", groupID)
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		group, err := dtdal.Group.GetGroupByID(c, tx, groupID)
		if err != nil {
			return err
		}
		var args [][]interface{}
		for _, member := range group.Data.GetGroupMembers() {
			if member.UserID != "" {
				args = append(args, []interface{}{member.UserID, []string{groupID}, []string{}})
			}
		}
		return delayUpdateUserWithGroups.EnqueueWorkMulti(c, delaying.With(common.QUEUE_USERS, "update-user-with-groups", 0), args...)
	}); err != nil {
		return err
	}
	return err
}

func delayedUpdateUserWithGroups(c context.Context, userID string, groupIDs2add, groupIDs2remove []string) (err error) {
	logus.Debugf(c, "delayedUpdateUserWithGroups(userID=%s, groupIDs2add=%+v, groupIDs2remove=%+v)", userID, groupIDs2add, groupIDs2remove)
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		groups2add := make([]models.GroupEntry, len(groupIDs2add))
		for i, groupID := range groupIDs2add {
			if groups2add[i], err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
				return
			}
		}
		var user models.AppUser
		if user, err = dtdal.User.GetUserByStrID(c, userID); err != nil {
			return
		}
		return User.UpdateUserWithGroups(c, tx, user, groups2add, groupIDs2remove)
	}); err != nil {
		return err
	}
	return err
}

func (userFacade) UpdateUserWithGroups(c context.Context, tx dal.ReadwriteTransaction, user models.AppUser, groups2add []models.GroupEntry, groups2remove []string) (err error) {
	logus.Debugf(c, "updateUserWithGroup(user.ID=%s, len(groups2add)=%d, groups2remove=%+v)", user.ID, len(groups2add), groups2remove)
	groups := user.Data.ActiveGroups()
	updated := false
	for _, group2add := range groups2add {
		updated = user.Data.AddGroup(group2add, "") || updated
	}
	for _, group2remove := range groups2remove {
		for i, group := range groups {
			if group.ID == group2remove {
				groups = append(groups[:i], groups[i+1:]...)
				updated = true
				continue
			}
		}
	}
	if !updated {
		logus.Debugf(c, "User is not update with groups")
		return
	}
	user.Data.SetActiveGroups(groups)
	if err = User.SaveUser(c, tx, user); err != nil {
		return
	}
	return
}

func (userFacade) DelayUpdateContactWithGroups(c context.Context, contactID string, addGroupIDs, removeGroupIDs []string) error {
	return delayUpdateContactWithGroups.EnqueueWork(c, delaying.With(common.QUEUE_USERS, "update-contact-groups", 0), contactID, addGroupIDs, removeGroupIDs)
}

func delayedUpdateContactWithGroup(c context.Context, contactID string, addGroupIDs, removeGroupIDs []string) (err error) {
	logus.Debugf(c, "delayedUpdateContactWithGroup(contactID=%s, addGroupIDs=%v, removeGroupIDs=%v)", contactID, addGroupIDs, removeGroupIDs)
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if _, err = GetContactByID(c, tx, contactID); err != nil {
			return err
		}
		return User.UpdateContactWithGroups(c, contactID, addGroupIDs, removeGroupIDs)
	}); err != nil {
		return
	}
	return
}

func (userFacade) UpdateContactWithGroups(c context.Context, contactID string, addGroupIDs, removeGroupIDs []string) error {
	logus.Debugf(c, "UpdateContactWithGroups(contactID=%s, addGroupIDs=%+v, removeGroupIDs=%+v)", contactID, addGroupIDs, removeGroupIDs)
	if contact, err := GetContactByID(c, nil, contactID); err != nil {
		return err
	} else {
		var isAdded bool
		contact.Data.GroupIDs, isAdded = slices.MergeStrings(contact.Data.GroupIDs, addGroupIDs)
		var removedCount int
		contact.Data.GroupIDs, removedCount = slices.RemoveStrings(contact.Data.GroupIDs, removeGroupIDs)
		if isAdded || removedCount > 0 {
			return SaveContact(c, contact)
		}
		return nil
	}
}

var ErrAttemptToLeaveUnsettledGroup = errors.New("an attept to leave unsettled group")

func (groupFacade) LeaveGroup(c context.Context, groupID string, userID string) (group models.GroupEntry, user models.AppUser, err error) {
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		group.ID = groupID
		user.ID = userID
		if err = tx.GetMulti(c, []dal.Record{group.Record, user.Record}); err != nil {
			return
		}
		//if group, err = dtdal.GroupEntry.GetGroupByID(c, groupID); err != nil {
		//	return
		//}
		//if user, err = dtdal.User.GetUserByStrID(c, userID); err != nil {
		//	return
		//}

		{ // Update group
			members := group.Data.GetGroupMembers()
			for i, m := range members {
				if m.UserID == userID {
					if len(m.Balance) != 0 {
						err = ErrAttemptToLeaveUnsettledGroup
						return
					}
					members = append(members[:i], members[i+1:]...)
					group.Data.SetGroupMembers(members)
					if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
						return
					}
					break
				}
			}
		}
		groups := user.Data.ActiveGroups()
		userChanged := false
		for i, g := range groups {
			if g.ID == groupID {
				groups = append(groups[:i], groups[i+1:]...)
				userChanged = true
				break
			}
		}
		if userChanged {
			user.Data.SetActiveGroups(groups)
			if err = User.SaveUser(c, tx, user); err != nil {
				return
			}
		}
		if err = Group.DelayUpdateGroupUsers(c, groupID); err != nil {
			return
		}
		return
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}
	return
}
