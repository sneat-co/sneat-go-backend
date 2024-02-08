package api

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

func panicUnknownStatus(status string) {
	panic("Unknown status: " + status)
}
func handleGetUserData(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	log.Debugf(c, "handleGetUserData(authInfo.UserID: %d)", authInfo.UserID)
	user, err := getApiUser(c, w, r, authInfo)
	if err != nil {
		return
	}
	markResponseAsJson(w.Header())

	rPath := r.URL.Path

	//getQueryValue := r.URL.Query().Get
	getQueryValue := func(name string) string {
		prefix := "/" + name + "-"
		start := strings.Index(rPath, prefix) + len(prefix)
		if start < 0 {
			return ""
		}
		end := strings.Index(rPath[start:], "/")
		if end < 0 {
			end = len(rPath)
		} else {
			end += start
		}
		return rPath[start:end]
	}

	status := getQueryValue("status")

	if status != "" && status != models.STATUS_ACTIVE && status != models.STATUS_ARCHIVED {
		BadRequestMessage(c, w, "Unknown status: "+status)
		return
	}

	dataCodes := strings.Split(getQueryValue("load"), ",")
	if len(dataCodes) == 0 {
		BadRequestMessage(c, w, "Missing `load` parameter value")
		return
	}

	//log.Debugf(c, "load: %v", dataCodes)

	dataResults := make([]*bytes.Buffer, len(dataCodes))

	hasContent := false
	for i, dataCode := range dataCodes {
		//log.Debugf(c, "i=%d, dataCode=%v", i, dataCode)
		dataResults[i] = &bytes.Buffer{}
		switch dataCode {
		case "Contacts":
			if status == models.STATUS_ACTIVE || status == models.STATUS_ARCHIVED {
				hasContent = writeUserContactsToJson(c, dataResults[i], status, user) || hasContent
			} else {
				panicUnknownStatus(status)
			}
		case "Groups":
			if status == models.STATUS_ACTIVE || status == models.STATUS_ARCHIVED {
				hasContent = writeUserGroupsToJson(c, dataResults[i], status, user) || hasContent
			} else {
				panicUnknownStatus(status)
			}
		case "Bills":
			switch status {
			case models.STATUS_ACTIVE:
				hasContent = writeUserActiveBillsToJson(c, dataResults[i], user) || hasContent
			default:
				panicUnknownStatus(status)
			}
		case "BillSchedules":
			switch status {
			case models.STATUS_ACTIVE:
				hasContent = writeUserActiveBillSchedulesToJson(c, dataResults[i], user) || hasContent
			default:
				panicUnknownStatus(status)
			}
		default:
			BadRequestMessage(c, w, "Unknown data code: "+dataCode)
			return
		}
	}

	if !hasContent {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, _ = w.Write(([]byte)("{"))
	needComma := false
	for _, dataResult := range dataResults {
		if dataResult.Len() > 0 {
			if needComma {
				_, _ = w.Write([]byte(","))
			} else {
				needComma = true
			}
			_, _ = w.Write([]byte("\n"))
			_, _ = w.Write(dataResult.Bytes())
		}
	}
	_, _ = w.Write(([]byte)("\n}"))
}

func writeUserGroupsToJson(_ context.Context, w io.Writer, status string, user models.AppUser) bool {
	//log.Debugf(c, "writeUserGroupsToJson(status=%v)", status)
	var jsonVal string
	switch status {
	case models.STATUS_ACTIVE:
		jsonVal = user.Data.GroupsJsonActive
	case models.STATUS_ARCHIVED:
		jsonVal = user.Data.GroupsJsonArchived
	default:
		panicUnknownStatus(status)
	}
	if jsonVal != "" {
		_, _ = w.Write(([]byte)(`"Groups":`))
		_, _ = w.Write([]byte(jsonVal))
		return true
	}
	return false
}

func writeUserContactsToJson(c context.Context, w io.Writer, status string, user models.AppUser) bool {
	//log.Debugf(c, "writeUserContactsToJson(status=%v)", status)
	var jsonVal string
	switch status {
	case models.STATUS_ACTIVE:
		jsonVal = user.Data.ContactsJsonActive
	case models.STATUS_ARCHIVED:
		jsonVal = user.Data.ContactsJsonArchived
	default:
		panicUnknownStatus(status)
	}

	if jsonVal != "" {
		_, _ = w.Write(([]byte)(`"Contacts":`))
		_, _ = w.Write([]byte(jsonVal))
		return true
	}
	return false
}

func writeUserActiveBillsToJson(c context.Context, w io.Writer, user models.AppUser) bool {
	if user.Data.BillsJsonActive != "" {
		log.Debugf(c, "User has BillsJsonActive")
		if user.Data.BillsCountActive == 0 {
			log.Warningf(c, "User(id=%d).BillsJsonActive is not empty && BillsCountActive == 0", user.ID)
		}
		_, _ = w.Write(([]byte)(`"Bills":`))
		_, _ = w.Write([]byte(user.Data.BillsJsonActive))
		return true
	}
	return false
}

func writeUserActiveBillSchedulesToJson(c context.Context, w io.Writer, user models.AppUser) bool {
	if user.Data.BillSchedulesJsonActive != "" {
		log.Debugf(c, "User has BillSchedulesJsonActive")
		if user.Data.BillSchedulesCountActive == 0 {
			log.Warningf(c, "User(id=%d).BillSchedulesJsonActive is not empty && BillSchedulesCountActive == 0", user.ID)
		}
		_, _ = w.Write(([]byte)(`"BillSchedules":`))
		_, _ = w.Write([]byte(user.Data.BillSchedulesJsonActive))
		return true
	}
	return false
}
