package unsorted

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade/dto"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"context"
	"errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func HandlerCreateGroup(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo, user models.AppUser) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	name := strings.TrimSpace(r.PostForm.Get("name"))
	note := strings.TrimSpace(r.PostForm.Get("note"))

	groupEntity := models.GroupDbo{
		CreatorUserID: authInfo.UserID,
		Name:          name,
	}
	if len(note) > 0 {
		groupEntity.Note = note
	}

	group, _, err := facade.Group.CreateGroup(c, &groupEntity, "", nil, nil)
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	logus.Infof(c, "GroupEntry created, ID: %v", group.ID)
	if err = groupToResponse(c, w, group, user); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func HandlerGetGroup(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo, user models.AppUser) {
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		api.BadRequestError(c, w, errors.New("missing id parameter: id"))
		return
	}
	db, err := facade.GetDatabase(c)
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	group, err := dtdal.Group.GetGroupByID(c, db, groupID)
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	if err = groupToResponse(c, w, group, user); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func groupToResponse(c context.Context, w http.ResponseWriter, group models.GroupEntry, user models.AppUser) error {
	if jsons, err := groupsToJson([]models.GroupEntry{group}, user); err != nil {
		return err
	} else {
		api.MarkResponseAsJson(w.Header())
		_, _ = w.Write(jsons[0])
		return nil
	}
}

func groupsToJson(groups []models.GroupEntry, user models.AppUser) (result [][]byte, err error) {
	result = make([][]byte, len(groups))

	groupStatuses := make(map[string]string, len(groups))

	for _, group := range user.Data.ActiveGroups() {
		groupStatuses[group.ID] = models.STATUS_ACTIVE
	}

	for i, group := range groups {
		groupDto := dto.GroupDto{
			ID:           group.ID,
			Name:         group.Data.Name,
			Note:         group.Data.Note,
			MembersCount: group.Data.MembersCount,
		}
		if status, ok := groupStatuses[group.ID]; ok {
			groupDto.Status = status
		} else {
			groupDto.Status = models.STATUS_ARCHIVED
		}
		contactsByID := user.Data.ContactsByID()
		if group.Data.MembersJson != "" {
			for _, member := range group.Data.GetGroupMembers() {
				memberDto := dto.GroupMemberDto{
					ID:   member.ID,
					Name: member.Name,
				}
				if member.UserID == user.ID {
					memberDto.Name = ""
					memberDto.UserID = member.UserID
				} else if member.Name == "" {
					err = fmt.Errorf("group(%v) has member(id=%v) without UserID and without Name", group.ID, member.ID)
					return
				}
				for _, contactID := range member.ContactIDs {
					if _, ok := contactsByID[contactID]; ok {
						memberDto.ContactID = contactID
						break
					}
				}
				groupDto.Members = append(groupDto.Members, memberDto)
			}
		}
		if result[i], err = ffjson.MarshalFast(&groupDto); err != nil {
			return
		}
	}
	return
}

func HandleJoinGroups(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	defer r.Body.Close()

	var groupIDs []string
	if body, err := io.ReadAll(r.Body); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
	} else if groupIDs = strings.Split(string(body), ","); len(groupIDs) == 0 {
		api.BadRequestError(c, w, errors.New("Missing body"))
		return
	}

	groups := make([]models.GroupEntry, len(groupIDs))
	var user models.AppUser

	db, err := facade.GetDatabase(c)
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if user, err = facade.User.GetUserByID(c, tx, authInfo.UserID); err != nil {
			return
		}
		var waitGroup sync.WaitGroup
		waitGroup.Add(len(groupIDs))

		errs := make([]error, len(groupIDs))
		for i, groupID := range groupIDs {
			go func(i int, groupID string) {
				var group models.GroupEntry
				if group, errs[i] = dtdal.Group.GetGroupByID(c, tx, groupID); errs[i] != nil {
					waitGroup.Done()
					return
				}
				groups[i] = group
				userName := user.Data.FullName()
				if userName == models.NoName {
					userName = ""
				}
				if _, changed, _, _, members := group.Data.AddOrGetMember(authInfo.UserID, "", userName); changed {
					group.Data.SetGroupMembers(members)
					if errs[i] = dtdal.Group.SaveGroup(c, tx, group); errs[i] != nil {
						waitGroup.Done()
						return
					}
				}
				if errs[i] = facade.Group.DelayUpdateGroupUsers(c, groupID); errs[i] != nil {
					waitGroup.Done()
					return
				}
				waitGroup.Done()
			}(i, groupID)
		}
		waitGroup.Wait()
		for _, err = range errs {
			if err != nil {
				return
			}
		}

		if err = facade.User.UpdateUserWithGroups(c, tx, user, groups, []string{}); err != nil {
			return
		}

		return
	}, dal.TxWithCrossGroup())

	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	jsons, err := groupsToJson(groups, user)
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
	}
	_, _ = w.Write(([]byte)("["))
	lastJsonIndex := len(jsons) - 1
	for i, json := range jsons {
		_, _ = w.Write(json)
		if i < lastJsonIndex {
			_, _ = w.Write([]byte(","))
		}
	}
	_, _ = w.Write(([]byte)("]"))
}

func HandlerDeleteGroup(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {

}

func HandlerUpdateGroup(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	logus.Debugf(c, "HandlerUpdateGroup()")

	var (
		user  models.AppUser
		group models.GroupEntry
		err   error
	)

	if group.ID = r.URL.Query().Get("id"); group.ID == "" {
		api.BadRequestError(c, w, errors.New("Missing id parameter"))
		return
	}

	groupName := strings.TrimSpace(r.FormValue("name"))
	groupNote := strings.TrimSpace(r.FormValue("note"))

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if group, err = dtdal.Group.GetGroupByID(c, tx, group.ID); err != nil {
			return
		}

		if group.Data.CreatorUserID != authInfo.UserID {
			err = fmt.Errorf("user is not authrized to edit this group")
			return
		}

		changed := false
		if groupName != "" && group.Data.Name != groupName {
			group.Data.Name = groupName
			changed = true
		}
		if group.Data.Note != groupNote {
			group.Data.Note = groupNote
			changed = true
		}
		if changed {
			if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
				return
			}
		}
		if user, err = facade.User.GetUserByID(c, tx, authInfo.UserID); err != nil {
			return
		}

		if err = facade.User.UpdateUserWithGroups(c, tx, user, []models.GroupEntry{group}, nil); err != nil {
			return
		}

		if err = facade.Group.DelayUpdateGroupUsers(c, group.ID); err != nil {
			return
		}

		return
	})

	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	if err = groupToResponse(c, w, group, user); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func HandlerSetContactsToGroup(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo, user models.AppUser) {
	logus.Debugf(c, "HandlerSetContactsToGroup()")

	var (
		groupID string
		group   models.GroupEntry
		err     error
	)

	if groupID = r.URL.Query().Get("id"); groupID == "" {
		api.BadRequestError(c, w, errors.New("Missing id parameter"))
		return
	}

	var (
		addContactIDs   []string
		removeMemberIDs []string
	)
	//addContactIDs := strings.Split(r.FormValue("addContactIDs"), ",")
	removeMemberIDs = strings.Split(r.FormValue("removeMemberIDs"), ",")

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		var contacts2add []models.ContactEntry
		if contacts2add, err = facade.GetContactsByIDs(c, tx, addContactIDs); err != nil {
			return err
		}

		for _, contact := range contacts2add {
			if contact.Data.UserID != authInfo.UserID {
				return validation.NewBadRequestError(fmt.Errorf("contact %s does not belong to the user %s", contact.ID, authInfo.UserID))
			}
		}

		if group, err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
			return err
		}
		members := group.Data.GetGroupMembers()
		changed := false
		changedContactIDs := make([]string, 0, len(addContactIDs)+len(removeMemberIDs))

		var groupUserIDs []string

		addGroupUserID := func(member models.GroupMemberJson) {
			userID := user.ID
			if member.UserID != "" && member.UserID != userID {
				groupUserIDs = append(groupUserIDs, member.UserID)
			}
		}

		for _, contact2add := range contacts2add {
			var (
				isChanged bool
			)
			for _, member := range members {
				for _, mContactID := range member.ContactIDs {
					if mContactID == contact2add.ID {
						goto found
					}
				}
			}
			_, isChanged, _, _, members = group.Data.AddOrGetMember(contact2add.Data.CounterpartyUserID, contact2add.ID, contact2add.Data.FullName())
			if isChanged {
				changed = true
				changedContactIDs = append(changedContactIDs, contact2add.ID)
			}
		found:
		}

		for _, memberID := range removeMemberIDs {
			for i, member := range members {
				if member.ID == memberID {
					members = append(members[:i], members[i+1:]...)
					changed = true
					addGroupUserID(member)
					for _, contactID := range member.ContactIDs {
						for _, changedContactID := range changedContactIDs {
							if changedContactID == contactID {
								goto alreadyChanged
							}
						}
						changedContactIDs = append(changedContactIDs, contactID)
					alreadyChanged:
					}
				}
			}
		}
		if changed || len(changedContactIDs) > 0 { // Check for len(changedContactIDs) is excessive but just in case.
			group.Data.SetGroupMembers(members)
			if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
				return err
			}
		}

		{ // Executing this block outside of IF just in case for self-healing.
			if user, err = facade.User.GetUserByID(c, tx, user.ID); err != nil {
				return err
			}
			if err = facade.User.UpdateUserWithGroups(c, tx, user, []models.GroupEntry{group}, []string{}); err != nil {
				return err
			}

			for _, member := range members {
				addGroupUserID(member)
			}

			if len(groupUserIDs) > 0 {
				if err = facade.Group.DelayUpdateGroupUsers(c, groupID); err != nil {
					return err
				}
			}

			if len(changedContactIDs) == 1 {
				err = facade.User.UpdateContactWithGroups(c, changedContactIDs[0], []string{groupID}, []string{})
			} else {
				for _, contactID := range changedContactIDs {
					if err = facade.User.DelayUpdateContactWithGroups(c, contactID, []string{groupID}, []string{}); err != nil {
						return err
					}
				}
			}
		}
		return err
	}); err != nil {
		if validation.IsBadRecordError(err) {
			api.BadRequestError(c, w, err)
			return
		}
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	if err = groupToResponse(c, w, group, user); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func StringToInt64s(s, sep string) (result []int64, err error) {
	if s == "" {
		return
	}
	vals := strings.Split(s, sep)
	result = make([]int64, len(vals))
	for i, val := range vals {
		if result[i], err = strconv.ParseInt(val, 10, 64); err != nil {
			return
		}
	}
	return
}
