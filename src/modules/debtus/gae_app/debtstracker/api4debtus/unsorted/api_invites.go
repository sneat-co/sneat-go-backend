package unsorted

import (
	"context"
	"net/http"
)

func CreateInvite(_ context.Context, _ http.ResponseWriter, _ *http.Request) {
	panic("disabled: switch to Firestore authentication") // TODO: switch to Firestore authentication
	//gaeUser := gaeuser.Current(c)
	//if !gaeUser.Admin {
	//	w.WriteHeader(http.StatusForbidden)
	//}
	//
	//createUserData := &dtdal.CreateUserData{}
	//clientInfo := models.NewClientInfoFromRequest(r)
	//userEmail, _, err := facade4debtus.User.GetOrCreateEmailUser(c, gaeUser.Email, true, createUserData, clientInfo)
	//if err != nil {
	//	api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
	//	return
	//}
	//_, _ = w.Write([]byte(userEmail.Data.AppUserID))
}
