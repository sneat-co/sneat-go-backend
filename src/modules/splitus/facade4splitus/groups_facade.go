package facade4splitus

import (
	"context"
	"errors"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
)

func CreateGroup(c context.Context,
	groupEntity *models4splitus.GroupDbo,
	tgBotCode string,
	beforeGroupInsert func(tc context.Context, groupEntity *models4splitus.GroupDbo) (group models4splitus.GroupEntry, err error),
	afterGroupInsert func(c context.Context, group models4splitus.GroupEntry, user dbo4userus.UserEntry) (err error),
) (group models4splitus.GroupEntry, groupMember models4splitus.GroupMember, err error) {
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return errors.New("CreateGroup is not implemented")
		//user, err := dal4userus.GetUserByID(c, tx, groupEntity.CreatorUserID)
		//if err != nil {
		//	return err
		//}
		//existingGroups := user.Data.ActiveGroups()
		//
		//if beforeGroupInsert != nil {
		//	if group, err = beforeGroupInsert(c, groupEntity); err != nil {
		//		return err
		//	}
		//}
		//
		//var groupMembersChanged bool
		//groupMembersChanged, _, memberIndex, member, members := groupEntity.AddOrGetMember(groupEntity.CreatorUserID, "", user.Data.FullName())
		//member.Shares = 1
		//members[memberIndex] = member
		//groupEntity.SetGroupMembers(members)
		//
		//if group.ContactID == "" {
		//	for _, existingGroup := range existingGroups {
		//		if existingGroup.Name == groupEntity.Name {
		//			return errors.New("Duplicate group name")
		//		}
		//	}
		//	if group, err = dtdal.Group.InsertGroup(c, tx, groupEntity); err != nil {
		//		return err
		//	}
		//} else if groupMembersChanged {
		//	if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
		//		return err
		//	}
		//}
		//
		//groupJson := models4debtus.UserGroupJson{
		//	ContactID:           group.ContactID,
		//	Name:         group.Data.Name,
		//	Note:         group.Data.Note,
		//	MembersCount: group.Data.MembersCount,
		//}
		//
		//if tgBotCode != "" {
		//	for _, tgGroupBot := range groupJson.TgBots {
		//		if tgGroupBot == tgBotCode {
		//			goto botFound
		//		}
		//	}
		//	groupJson.TgBots = append(groupJson.TgBots, tgBotCode)
		//botFound:
		//}
		//
		//user.Data.SetActiveGroups(append(existingGroups, groupJson))
		//
		//if afterGroupInsert != nil {
		//	if err = afterGroupInsert(c, group, user); err != nil {
		//		return err
		//	}
		//}
		//
		//if err = facade4debtus.UserEntry.SaveUserOBSOLETE(c, tx, user); err != nil {
		//	return err
		//}
		//if err = groupFacade.DelayUpdateGroupUsers(c, group.ContactID); err != nil {
		//	return err
		//}
		//return err
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}
	logus.Infof(c, "GroupEntry created, ContactID=%v", group.ID)
	return
}

type NewUser struct {
	Name string
	botsfwmodels.PlatformUserData
	ChatMember botsfw.WebhookActor
}

func AddUsersToTheGroupAndOutstandingBills(c context.Context, spaceID string, newUsers []NewUser) (splitusSpace models4splitus.SplitusSpaceEntry, newUsers2 []NewUser, err error) {
	logus.Debugf(c, "groupFacade.AddUsersToTheGroupAndOutstandingBills(spaceID=%v, newUsers=%v)", spaceID, newUsers)
	splitusSpace = models4splitus.NewSplitusSpaceEntry(spaceID)
	if len(newUsers) == 0 {
		panic("len(newUsers) == 0")
	}
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		changed := false
		if err = tx.Get(c, splitusSpace.Record); err != nil {
			return
		}
		logus.Debugf(c, "splitusSpace: %+v", splitusSpace.Data)
		j := 0
		for _, newUser := range newUsers {
			_, isChanged, _, _, groupMembers := splitusSpace.Data.AddOrGetMember(newUser.GetAppUserID(), "", newUser.Name)
			changed = changed || isChanged
			if isChanged {
				splitusSpace.Data.SetGroupMembers(groupMembers)
				newUsers[j] = newUser
				j += 1
			}
		}
		newUsers = newUsers[:j]
		if changed {
			logus.Debugf(c, "splitusSpace: %+v", splitusSpace.Data)
			if err = tx.Set(c, splitusSpace.Record); err != nil {
				return
			}
			if err = DelayUpdateGroupUsers(c, splitusSpace.ID); err != nil {
				return err
			}
		}
		return
	}); err != nil {
		return splitusSpace, newUsers, err
	}
	return splitusSpace, newUsers, err
}

func DelayUpdateGroupUsers(c context.Context, groupID string) error { // TODO: Move to DAL?
	if groupID == "" {
		panic("groupID is empty string")
	}
	return delayerUpdateGroupUsers.EnqueueWork(c, delaying.With(const4userus.QueueUsers, "update-group-users", 0), groupID)
}

func delayedUpdateGroupUsers(c context.Context, spaceID string) (err error) {
	if spaceID == "" {
		logus.Criticalf(c, "spaceID is empty string")
		return nil
	}

	logus.Debugf(c, "delayedUpdateGroupUsers(spaceID=%v)", spaceID)
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		splitusSpace := models4splitus.NewSplitusSpaceEntry(spaceID)
		if err = tx.Get(c, splitusSpace.Record); err != nil {
			return err
		}
		for _, member := range splitusSpace.Data.GetGroupMembers() {
			if member.UserID != "" {
				if err = delayUpdateUserWithGroups(c, member.UserID, []string{spaceID}, []string{}); err != nil {
					return err
				}
			}
		}
		return err
	}); err != nil {
		return err
	}
	return err
}

func delayedUpdateUserWithGroups(c context.Context, userID string, groupIDs2add, groupIDs2remove []string) (err error) {
	logus.Debugf(c, "delayedUpdateUserWithGroups(userID=%s, groupIDs2add=%+v, groupIDs2remove=%+v)", userID, groupIDs2add, groupIDs2remove)
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var splitusSpaceRecords []dal.Record
		groups2add := make([]models4splitus.SplitusSpaceEntry, len(groupIDs2add))
		for i, spaceID := range groupIDs2add {
			groups2add[i] = models4splitus.NewSplitusSpaceEntry(spaceID)
		}
		if err = tx.GetMulti(c, splitusSpaceRecords); err != nil {
			return err
		}
		for _, group := range groups2add {
			if err = group.Record.Error(); err != nil {
				return err
			}
		}
		return errors.New("not implemented")
		//var user models4debtus.AppUserOBSOLETE
		//if user, err = facade4auth.UserEntry.GetUserByStrID(c, userID); err != nil {
		//	return
		//}
		//return UserEntry.UpdateUserWithGroups(c, tx, user, groups2add, groupIDs2remove)
	}); err != nil {
		return err
	}
	return err
}

func UpdateUserWithGroups(c context.Context, _ dal.ReadwriteTransaction, user dbo4userus.UserEntry, groups2add []models4splitus.GroupEntry, groups2remove []string) (err error) {
	logus.Debugf(c, "updateUserWithGroup(user.ContactID=%s, len(groups2add)=%d, groups2remove=%+v)", user.ID, len(groups2add), groups2remove)
	return errors.New("not implemented")
	//groups := user.Data.ActiveGroups()
	//updated := false
	//for _, group2add := range groups2add {
	//	updated = user.Data.AddGroup(group2add, "") || updated
	//}
	//for _, group2remove := range groups2remove {
	//	for i, group := range groups {
	//		if group.ContactID == group2remove {
	//			groups = append(groups[:i], groups[i+1:]...)
	//			updated = true
	//			continue
	//		}
	//	}
	//}
	//if !updated {
	//	logus.Debugf(c, "UserEntry is not update with groups")
	//	return
	//}
	//user.Data.SetActiveGroups(groups)
	//if err = UserEntry.SaveUserOBSOLETE(c, tx, user); err != nil {
	//	return
	//}
	//return
}

func DelayUpdateContactWithGroups(c context.Context, contactID string, addGroupIDs, removeGroupIDs []string) error {
	return delayerUpdateContactWithGroups.EnqueueWork(c, delaying.With(const4userus.QueueUsers, "update-contact-groups", 0), contactID, addGroupIDs, removeGroupIDs)
}

func delayedUpdateContactWithGroup(c context.Context, contactID string, addGroupIDs, removeGroupIDs []string) (err error) {
	logus.Debugf(c, "delayedUpdateContactWithGroup(contactID=%s, addGroupIDs=%v, removeGroupIDs=%v)", contactID, addGroupIDs, removeGroupIDs)
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		//if _, err = facade4debtus.GetContactByID(c, tx, contactID); err != nil {
		//	return err
		//}
		return UpdateContactWithGroups(c, contactID, addGroupIDs, removeGroupIDs)
	}); err != nil {
		return
	}
	return
}

func UpdateContactWithGroups(c context.Context, contactID string, addGroupIDs, removeGroupIDs []string) error {
	logus.Debugf(c, "UpdateContactWithGroups(contactID=%s, addGroupIDs=%+v, removeGroupIDs=%+v)", contactID, addGroupIDs, removeGroupIDs)
	return errors.New("UpdateContactWithGroups not implemented")
	//if contact, err := facade4debtus.GetContactByID(c, nil, contactID); err != nil {
	//	return err
	//} else {
	//	var isAdded bool
	//	contact.Data.SpaceIDs, isAdded = slices.MergeStrings(contact.Data.SpaceIDs, addGroupIDs)
	//	var removedCount int
	//	contact.Data.SpaceIDs, removedCount = slices.RemoveStrings(contact.Data.SpaceIDs, removeGroupIDs)
	//	if isAdded || removedCount > 0 {
	//		return facade4debtus.SaveContact(c, contact)
	//	}
	//	return nil
	//}
}

//var ErrAttemptToLeaveUnsettledGroup = errors.New("an attempt to leave unsettled group")

//func LeaveGroup(c context.Context, groupID string, userID string) (splitusSpace models4splitus.SplitusSpaceEntry, user dbo4userus.UserEntry, err error) {
//	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
//		splitusSpace.ID = groupID
//		user.ID = userID
//		if err = tx.GetMulti(c, []dal.Record{splitusSpace.Record, user.Record}); err != nil {
//			return
//		}
//		//if splitusSpace, err = dtdal.GroupEntry.GetGroupByID(c, groupID); err != nil {
//		//	return
//		//}
//		//if user, err = dtdal.UserEntry.GetUserByStrID(c, userID); err != nil {
//		//	return
//		//}
//
//		{ // Update splitusSpace
//			members := splitusSpace.Data.GetGroupMembers()
//			for i, m := range members {
//				if m.UserID == userID {
//					if len(m.Balance) != 0 {
//						err = ErrAttemptToLeaveUnsettledGroup
//						return
//					}
//					members = append(members[:i], members[i+1:]...)
//					splitusSpace.Data.SetGroupMembers(members)
//					if err = tx.Set(c, splitusSpace.Record); err != nil {
//						return
//					}
//					break
//				}
//			}
//		}
//		groups := user.Data.ActiveGroups()
//		userChanged := false
//		for i, g := range groups {
//			if g.ID == groupID {
//				groups = append(groups[:i], groups[i+1:]...)
//				userChanged = true
//				break
//			}
//		}
//		if userChanged {
//			user.Data.SetActiveGroups(groups)
//			//if err = facade4debtus.UserEntry.SaveUserOBSOLETE(c, tx, user); err != nil {
//			//	return
//			//}
//			return errors.New("not implemented")
//		}
//		if err = DelayUpdateGroupUsers(c, groupID); err != nil {
//			return
//		}
//		return
//	}, dal.TxWithCrossGroup()); err != nil {
//		return
//	}
//	return
//}
