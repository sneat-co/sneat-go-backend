package unsorted

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	common4all2 "github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
	"net/http"
	"strings"
)

func HandlerCreateGroup(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo, user dbo4userus.UserEntry) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	name := strings.TrimSpace(r.PostForm.Get("name"))
	note := strings.TrimSpace(r.PostForm.Get("note"))

	groupEntity := models4splitus.GroupDbo{
		CreatorUserID: authInfo.UserID,
		Name:          name,
	}
	if len(note) > 0 {
		groupEntity.Note = note
	}

	group, _, err := facade4splitus.CreateGroup(ctx, &groupEntity, "", nil, nil)
	if err != nil {
		common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}
	logus.Infof(ctx, "GroupEntry created, ContactID: %v", group.ID)
	if err = groupToResponse(ctx, w, group, user); err != nil {
		common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}
}

func HandlerGetGroup(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo, user dbo4userus.UserEntry) {
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		common4all2.BadRequestError(ctx, w, errors.New("missing id parameter: id"))
		return
	}
	common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.New("not implemented yet"))

	//db, err := facade.GetSneatDB(ctx)
	//if err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
	//group, err := dtdal.Group.GetGroupByID(ctx, db, groupID)
	//if err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
	//if err = groupToResponse(ctx, w, group, user); err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
}

func groupToResponse(_ context.Context, w http.ResponseWriter, group models4splitus.GroupEntry, user dbo4userus.UserEntry) error {
	if jsons, err := groupsToJson([]models4splitus.GroupEntry{group}, user); err != nil {
		return err
	} else {
		common4all2.MarkResponseAsJson(w.Header())
		_, _ = w.Write(jsons[0])
		return nil
	}
}

func groupsToJson(_ []models4splitus.GroupEntry, _ dbo4userus.UserEntry) (result [][]byte, err error) {
	err = errors.New("groupsToJson not implemented yet")
	return
	//result = make([][]byte, len(groups))
	//
	//groupStatuses := make(map[string]string, len(groups))
	//
	//for _, group := range user.Data.ActiveGroups() {
	//	groupStatuses[group.ContactID] = const4debtus.StatusActive
	//}
	//
	//for i, group := range groups {
	//	groupDto := dto4debtus.GroupDto{
	//		ContactID:           group.ContactID,
	//		Name:         group.Data.Name,
	//		Note:         group.Data.Note,
	//		MembersCount: group.Data.MembersCount,
	//	}
	//	if status, ok := groupStatuses[group.ContactID]; ok {
	//		groupDto.Status = status
	//	} else {
	//		groupDto.Status = const4debtus.StatusArchived
	//	}
	//	contactsByID := user.Data.ContactsByID()
	//	if group.Data.MembersJson != "" {
	//		for _, member := range group.Data.GetGroupMembers() {
	//			memberDto := dto4debtus.GroupMemberDto{
	//				ContactID:   member.ContactID,
	//				Name: member.Name,
	//			}
	//			if member.UserID == user.ContactID {
	//				memberDto.Name = ""
	//				memberDto.UserID = member.UserID
	//			} else if member.Name == "" {
	//				err = fmt.Errorf("group(%v) has member(id=%v) without UserID and without Name", group.ContactID, member.ContactID)
	//				return
	//			}
	//			for _, contactID := range member.ContactIDs {
	//				if _, ok := contactsByID[contactID]; ok {
	//					memberDto.ContactID = contactID
	//					break
	//				}
	//			}
	//			groupDto.Members = append(groupDto.Members, memberDto)
	//		}
	//	}
	//	if result[i], err = ffjson.MarshalFast(&groupDto); err != nil {
	//		return
	//	}
	//}
	//return
}

func HandleJoinGroups(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.New("not implemented yet"))
	//defer r.Body.Close()
	//
	//var groupIDs []string
	//if body, err := io.ReadAll(r.Body); err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//} else if groupIDs = strings.Split(string(body), ","); len(groupIDs) == 0 {
	//	api4debtus.BadRequestError(ctx, w, errors.New("Missing body"))
	//	return
	//}
	//
	//groups := make([]models4splitus.GroupEntry, len(groupIDs))
	//user := dbo4userus.NewUserEntry(authInfo.UserID)
	//
	//err := facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
	//
	//	if err = facade4userus.GetUserByIdOBSOLETE(ctx, tx, user.Record); err != nil {
	//		return
	//	}
	//	var waitGroup sync.WaitGroup
	//	waitGroup.Add(len(groupIDs))
	//
	//	errs := make([]error, len(groupIDs))
	//	for i, groupID := range groupIDs {
	//		go func(i int, groupID string) {
	//			var group models4splitus.GroupEntry
	//			if group, errs[i] = dtdal.Group.GetGroupByID(ctx, tx, groupID); errs[i] != nil {
	//				waitGroup.Done()
	//				return
	//			}
	//			groups[i] = group
	//			userName := user.Data.GetFullName()
	//			if userName == models4debtus.NoName {
	//				userName = ""
	//			}
	//			if _, changed, _, _, members := group.Data.AddOrGetMember(authInfo.UserID, "", userName); changed {
	//				group.Data.SetGroupMembers(members)
	//				if errs[i] = dtdal.Group.SaveGroup(ctx, tx, group); errs[i] != nil {
	//					waitGroup.Done()
	//					return
	//				}
	//			}
	//			if errs[i] = facade4splitus.Group.DelayUpdateGroupUsers(ctx, groupID); errs[i] != nil {
	//				waitGroup.Done()
	//				return
	//			}
	//			waitGroup.Done()
	//		}(i, groupID)
	//	}
	//	waitGroup.Wait()
	//	for _, err = range errs {
	//		if err != nil {
	//			return
	//		}
	//	}
	//
	//	if err = facade4splitus.UpdateUserWithGroups(ctx, tx, user, groups, []string{}); err != nil {
	//		return
	//	}
	//
	//	return
	//}, dal.TxWithCrossGroup())
	//
	//if err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
	//
	//jsons, err := groupsToJson(groups, user)
	//if err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//}
	//_, _ = w.Write(([]byte)("["))
	//lastJsonIndex := len(jsons) - 1
	//for i, json := range jsons {
	//	_, _ = w.Write(json)
	//	if i < lastJsonIndex {
	//		_, _ = w.Write([]byte(","))
	//	}
	//}
	//_, _ = w.Write(([]byte)("]"))
}

func HandlerDeleteGroup(_ context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {

}

func HandlerUpdateGroup(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	logus.Debugf(ctx, "HandlerUpdateGroup()")
	common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.New("not implemented yet"))
	//var (
	//	group models4splitus.GroupEntry
	//	err   error
	//)
	//
	//user := dbo4userus.NewUserEntry(authInfo.UserID)
	//
	//if group.ContactID = r.URL.Query().Get("id"); group.ContactID == "" {
	//	api4debtus.BadRequestError(ctx, w, errors.New("Missing id parameter"))
	//	return
	//}
	//
	//groupName := strings.TrimSpace(r.FormValue("name"))
	//groupNote := strings.TrimSpace(r.FormValue("note"))
	//
	//err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
	//	if group, err = dtdal.Group.GetGroupByID(ctx, tx, group.ContactID); err != nil {
	//		return
	//	}
	//
	//	if group.Data.CreatorUserID != authInfo.UserID {
	//		err = fmt.Errorf("user is not authrized to edit this group")
	//		return
	//	}
	//
	//	changed := false
	//	if groupName != "" && group.Data.Name != groupName {
	//		group.Data.Name = groupName
	//		changed = true
	//	}
	//	if group.Data.Note != groupNote {
	//		group.Data.Note = groupNote
	//		changed = true
	//	}
	//	if changed {
	//		if err = dtdal.Group.SaveGroup(ctx, tx, group); err != nil {
	//			return
	//		}
	//	}
	//	if err = facade4userus.GetUserByIdOBSOLETE(ctx, tx, user.Record); err != nil {
	//		return
	//	}
	//
	//	if err = facade4splitus.UpdateUserWithGroups(ctx, tx, user, []models4splitus.GroupEntry{group}, nil); err != nil {
	//		return
	//	}
	//
	//	if err = facade4splitus.Group.DelayUpdateGroupUsers(ctx, group.ContactID); err != nil {
	//		return
	//	}
	//
	//	return
	//})
	//
	//if err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
	//
	//if err = groupToResponse(ctx, w, group, user); err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
}

func HandlerSetContactsToGroup(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo, user dbo4userus.UserEntry) {
	logus.Debugf(ctx, "HandlerSetContactsToGroup()")
	common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.New("HandlerSetContactsToGroup() not implemented yet"))
	//var (
	//	groupID string
	//	group   models4splitus.GroupEntry
	//	err     error
	//)
	//
	//if groupID = r.URL.Query().Get("id"); groupID == "" {
	//	api4debtus.BadRequestError(ctx, w, errors.New("Missing id parameter"))
	//	return
	//}
	//
	//var (
	//	addContactIDs   []string
	//	removeMemberIDs []string
	//)
	////addContactIDs := strings.Split(r.FormValue("addContactIDs"), ",")
	//removeMemberIDs = strings.Split(r.FormValue("removeMemberIDs"), ",")
	//
	//if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
	//	var contacts2add []models4debtus.DebtusSpaceContactEntry
	//	if contacts2add, err = facade4debtus.GetContactsByIDs(ctx, tx, addContactIDs); err != nil {
	//		return err
	//	}
	//
	//	for _, contact := range contacts2add {
	//		if contact.Data.UserID != authInfo.UserID {
	//			return validation.NewBadRequestError(fmt.Errorf("contact %s does not belong to the user %s", contact.ContactID, authInfo.UserID))
	//		}
	//	}
	//
	//	if group, err = dtdal.Group.GetGroupByID(ctx, tx, groupID); err != nil {
	//		return err
	//	}
	//	members := group.Data.GetGroupMembers()
	//	changed := false
	//	changedContactIDs := make([]string, 0, len(addContactIDs)+len(removeMemberIDs))
	//
	//	var groupUserIDs []string
	//
	//	addGroupUserID := func(member briefs4splitus.GroupMemberJson) {
	//		userID := user.ContactID
	//		if member.UserID != "" && member.UserID != userID {
	//			groupUserIDs = append(groupUserIDs, member.UserID)
	//		}
	//	}
	//
	//	for _, contact2add := range contacts2add {
	//		var (
	//			isChanged bool
	//		)
	//		for _, member := range members {
	//			for _, mContactID := range member.ContactIDs {
	//				if mContactID == contact2add.ContactID {
	//					goto found
	//				}
	//			}
	//		}
	//		_, isChanged, _, _, members = group.Data.AddOrGetMember(contact2add.Data.CounterpartyUserID, contact2add.ContactID, contact2add.Data.FullName())
	//		if isChanged {
	//			changed = true
	//			changedContactIDs = append(changedContactIDs, contact2add.ContactID)
	//		}
	//	found:
	//	}
	//
	//	for _, memberID := range removeMemberIDs {
	//		for i, member := range members {
	//			if member.ContactID == memberID {
	//				members = append(members[:i], members[i+1:]...)
	//				changed = true
	//				addGroupUserID(member)
	//				for _, contactID := range member.ContactIDs {
	//					for _, changedContactID := range changedContactIDs {
	//						if changedContactID == contactID {
	//							goto alreadyChanged
	//						}
	//					}
	//					changedContactIDs = append(changedContactIDs, contactID)
	//				alreadyChanged:
	//				}
	//			}
	//		}
	//	}
	//	if changed || len(changedContactIDs) > 0 { // Check for len(changedContactIDs) is excessive but just in case.
	//		group.Data.SetGroupMembers(members)
	//		if err = dtdal.Group.SaveGroup(ctx, tx, group); err != nil {
	//			return err
	//		}
	//	}
	//
	//	{ // Executing this block outside of IF just in case for self-healing.
	//		if user, err = dal4userus.GetUserByID(ctx, tx, user.ContactID); err != nil {
	//			return err
	//		}
	//		if err = facade4splitus.UpdateUserWithGroups(ctx, tx, user, []models4splitus.GroupEntry{group}, []string{}); err != nil {
	//			return err
	//		}
	//
	//		for _, member := range members {
	//			addGroupUserID(member)
	//		}
	//
	//		if len(groupUserIDs) > 0 {
	//			if err = facade4splitus.Group.DelayUpdateGroupUsers(ctx, groupID); err != nil {
	//				return err
	//			}
	//		}
	//
	//		if len(changedContactIDs) == 1 {
	//			err = facade4splitus.UpdateContactWithGroups(ctx, changedContactIDs[0], []string{groupID}, []string{})
	//		} else {
	//			for _, contactID := range changedContactIDs {
	//				if err = facade4splitus.DelayUpdateContactWithGroups(ctx, contactID, []string{groupID}, []string{}); err != nil {
	//					return err
	//				}
	//			}
	//		}
	//	}
	//	return err
	//}); err != nil {
	//	if validation.IsBadRecordError(err) {
	//		api4debtus.BadRequestError(ctx, w, err)
	//		return
	//	}
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
	//if err = groupToResponse(ctx, w, group, user); err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
}

//func StringToInt64s(s, sep string) (result []int64, err error) {
//	if s == "" {
//		return
//	}
//	vals := strings.Split(s, sep)
//	result = make([]int64, len(vals))
//	for i, val := range vals {
//		if result[i], err = strconv.ParseInt(val, 10, 64); err != nil {
//			return
//		}
//	}
//	return
//}
