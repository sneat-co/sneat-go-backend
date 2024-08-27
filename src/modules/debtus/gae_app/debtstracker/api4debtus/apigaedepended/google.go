package apigaedepended

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"net/http"
)

var handleFunc = http.HandleFunc

func InitApiGaeDepended() {
	handleFunc("/auth/google/signin", dtdal.HttpAppHost.HandleWithContext(handleSigninWithGoogle))
	handleFunc("/auth/google/signed", dtdal.HttpAppHost.HandleWithContext(handleSignedWithGoogle))
}

const REDIRECT_DESTINATION_PARAM_NAME = "redirect-to"

func handleSigninWithGoogle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	panic("disabled: switch to Firestore authentication")
	//query := r.URL.Query()
	//
	//destinationUrl := query.Get("to")
	//if destinationUrl == "" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	_, _ = w.Write([]byte("Missing 'to' parameter"))
	//	return
	//}
	//
	//callbackUrl := fmt.Sprintf("/auth/google/signed?%v=%v", REDIRECT_DESTINATION_PARAM_NAME, url.QueryEscape(destinationUrl))
	//if secret := query.Get("secret"); secret != "" {
	//	callbackUrl += "&secret=" + secret
	//}

	//loginUrl, err := gaeuser.LoginURL(ctx, callbackUrl)

	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_, _ = w.Write([]byte(err.Error()))
	//	return
	//}
	//http.Redirect(w, r, loginUrl, http.StatusFound)
}

func handleSignedWithGoogle(_ context.Context, _ http.ResponseWriter, _ *http.Request) {
	panic("disabled: switch to Firestore authentication")
	//var userID string
	//if authInfo, _, err := auth.Authenticate(w, r, false); err != nil {
	//	if err != auth.ErrNoToken {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		_, _ = w.Write([]byte(err.Error()))
	//	}
	//} else {
	//	userID = authInfo.UserID
	//}
	//
	//clientInfo := models.NewClientInfoFromRequest(r)

	//googleUser := gaeuser.Current(c)
	//if googleUser == nil {
	//	err := errors.New("handleSignedWithGoogle(): googleUser == nil")
	//	logus.Errorf(c, err.Error())
	//	w.WriteHeader(http.StatusUnauthorized)
	//	_, _ = w.Write([]byte(err.Error()))
	//	return
	//}
	//
	//userGoogle, _, err := facade4debtus.User.GetOrCreateUserGoogleOnSignIn(c, googleUser, userID, clientInfo)
	//if err != nil {
	//	logus.Errorf(c, err.Error())
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_, _ = w.Write([]byte(err.Error()))
	//}
	//
	//if !userGoogle.Record.Exists() {
	//	logus.Errorf(c, "userGoogle.UserGoogleData == nil")
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_, _ = w.Write([]byte("userGoogle.UserGoogleData == nil"))
	//}
	//
	//accountData := userGoogle.Data().(*appuser.AccountDataBase)
	//appUserID := userGoogle.Data().GetAppUserID()
	//logus.Debugf(c, "userGoogle.AppUserIntID: %s", appUserID)
	//token := auth.IssueToken(appUserID, "web", accountData.EmailLowerCase == "alexander.trakhimenok@gmail.com")
	//destinationUrl := r.URL.Query().Get(REDIRECT_DESTINATION_PARAM_NAME)
	//
	//var delimiter string
	//if strings.Contains(destinationUrl, "#") {
	//	delimiter = "&"
	//} else {
	//	delimiter = "#"
	//}
	//destinationUrl += delimiter + "signed-in-with=google"
	//destinationUrl += "&secret=" + token
	//http.Redirect(w, r, destinationUrl, http.StatusFound)
}
