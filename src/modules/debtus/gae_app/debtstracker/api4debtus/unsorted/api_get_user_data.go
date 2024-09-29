package unsorted

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"net/http"
)

//func panicUnknownStatus(status string) {
//	panic("Unknown status: " + status)
//}

func HandleGetUserData(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	panic("not implemented")
	//logus.Debugf(ctx, "HandleGetUserData(authInfo.UserID: %s)", authInfo.UserID)
	//user, err := getApiUser(ctx, w, r, authInfo)
	//if err != nil {
	//	return
	//}
	//api.MarkResponseAsJson(w.Header())
	//
	//rPath := r.URL.Path
	//
	////getQueryValue := r.URL.Query().Get
	//getQueryValue := func(name string) string {
	//	prefix := "/" + name + "-"
	//	start := strings.Index(rPath, prefix) + len(prefix)
	//	if start < 0 {
	//		return ""
	//	}
	//	end := strings.Index(rPath[start:], "/")
	//	if end < 0 {
	//		end = len(rPath)
	//	} else {
	//		end += start
	//	}
	//	return rPath[start:end]
	//}
	//
	//status := getQueryValue("status")
	//
	//if status != "" && status != models4debtus.STATUS_ACTIVE && status != models4debtus.STATUS_ARCHIVED {
	//	api.BadRequestMessage(ctx, w, "Unknown status: "+status)
	//	return
	//}
	//
	//dataCodes := strings.Split(getQueryValue("load"), ",")
	//if len(dataCodes) == 0 {
	//	api.BadRequestMessage(ctx, w, "Missing `load` parameter value")
	//	return
	//}
	//
	////logus.Debugf(ctx, "load: %v", dataCodes)
	//
	//dataResults := make([]*bytes.Buffer, len(dataCodes))
	//
	//hasContent := false
	//for i, dataCode := range dataCodes {
	//	//logus.Debugf(ctx, "i=%d, dataCode=%v", i, dataCode)
	//	dataResults[i] = &bytes.Buffer{}
	//	switch dataCode {
	//	case "Contacts":
	//		if status == models4debtus.STATUS_ACTIVE || status == models4debtus.STATUS_ARCHIVED {
	//			hasContent = writeUserContactsToJson(ctx, dataResults[i], status, user) || hasContent
	//		} else {
	//			panicUnknownStatus(status)
	//		}
	//	case "Groups":
	//		if status == models4debtus.STATUS_ACTIVE || status == models4debtus.STATUS_ARCHIVED {
	//			hasContent = writeUserGroupsToJson(ctx, dataResults[i], status, user) || hasContent
	//		} else {
	//			panicUnknownStatus(status)
	//		}
	//	case "Bills":
	//		switch status {
	//		case models4debtus.STATUS_ACTIVE:
	//			hasContent = writeUserActiveBillsToJson(ctx, dataResults[i], user) || hasContent
	//		default:
	//			panicUnknownStatus(status)
	//		}
	//	case "BillSchedules":
	//		switch status {
	//		case models4debtus.STATUS_ACTIVE:
	//			hasContent = writeUserActiveBillSchedulesToJson(ctx, dataResults[i], user) || hasContent
	//		default:
	//			panicUnknownStatus(status)
	//		}
	//	default:
	//		api.BadRequestMessage(ctx, w, "Unknown data code: "+dataCode)
	//		return
	//	}
	//}
	//
	//if !hasContent {
	//	w.WriteHeader(http.StatusNoContent)
	//	return
	//}
	//_, _ = w.Write(([]byte)("{"))
	//needComma := false
	//for _, dataResult := range dataResults {
	//	if dataResult.Len() > 0 {
	//		if needComma {
	//			_, _ = w.Write([]byte(","))
	//		} else {
	//			needComma = true
	//		}
	//		_, _ = w.Write([]byte("\n"))
	//		_, _ = w.Write(dataResult.Bytes())
	//	}
	//}
	//_, _ = w.Write(([]byte)("\n}"))
}

//func writeUserGroupsToJson(_ context.Context, w io.Writer, status string, user models4debtus.AppUser) bool {
//	//logus.Debugf(c, "writeUserGroupsToJson(status=%v)", status)
//	var jsonVal string
//	switch status {
//	case models4debtus.STATUS_ACTIVE:
//		jsonVal = user.Data.GroupsJsonActive
//	case models4debtus.STATUS_ARCHIVED:
//		jsonVal = user.Data.GroupsJsonArchived
//	default:
//		panicUnknownStatus(status)
//	}
//	if jsonVal != "" {
//		_, _ = w.Write(([]byte)(`"Groups":`))
//		_, _ = w.Write([]byte(jsonVal))
//		return true
//	}
//	return false
//}
//
//func writeUserContactsToJson(ctx context.Context, w io.Writer, status string, user models4debtus.AppUser) bool {
//	//logus.Debugf(c, "writeUserContactsToJson(status=%v)", status)
//	var jsonVal string
//	switch status {
//	case models4debtus.STATUS_ACTIVE:
//		jsonVal = user.Data.ContactsJsonActive
//	case models4debtus.STATUS_ARCHIVED:
//		jsonVal = user.Data.ContactsJsonArchived
//	default:
//		panicUnknownStatus(status)
//	}
//
//	if jsonVal != "" {
//		_, _ = w.Write(([]byte)(`"Contacts":`))
//		_, _ = w.Write([]byte(jsonVal))
//		return true
//	}
//	return false
//}
//
//func writeUserActiveBillsToJson(ctx context.Context, w io.Writer, user models4debtus.AppUser) bool {
//	if user.Data.BillsJsonActive != "" {
//		logus.Debugf(c, "User has BillsJsonActive")
//		if user.Data.BillsCountActive == 0 {
//			logus.Warningf(c, "User(id=%s).BillsJsonActive is not empty && BillsCountActive == 0", user.ContactID)
//		}
//		_, _ = w.Write(([]byte)(`"Bills":`))
//		_, _ = w.Write([]byte(user.Data.BillsJsonActive))
//		return true
//	}
//	return false
//}
//
//func writeUserActiveBillSchedulesToJson(ctx context.Context, w io.Writer, user models4debtus.AppUser) bool {
//	if user.Data.BillSchedulesJsonActive != "" {
//		logus.Debugf(c, "User has BillSchedulesJsonActive")
//		if user.Data.BillSchedulesCountActive == 0 {
//			logus.Warningf(c, "User(id=%s).BillSchedulesJsonActive is not empty && BillSchedulesCountActive == 0", user.ContactID)
//		}
//		_, _ = w.Write(([]byte)(`"BillSchedules":`))
//		_, _ = w.Write([]byte(user.Data.BillSchedulesJsonActive))
//		return true
//	}
//	return false
//}
