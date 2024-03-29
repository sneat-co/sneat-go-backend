package apigaedepended

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
	"github.com/strongo/strongoapp/appuser"
	"google.golang.org/appengine/v2/user"
	"net/http"
	"net/url"
	"strings"
)

var handleFunc = http.HandleFunc

func InitApiGaeDepended() {
	handleFunc("/auth/google/signin", dtdal.HttpAppHost.HandleWithContext(handleSigninWithGoogle))
	handleFunc("/auth/google/signed", dtdal.HttpAppHost.HandleWithContext(handleSignedWithGoogle))
}

const REDIRECT_DESTINATION_PARAM_NAME = "redirect-to"

func handleSigninWithGoogle(c context.Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	destinationUrl := query.Get("to")
	if destinationUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Missing 'to' parameter"))
		return
	}

	callbackUrl := fmt.Sprintf("/auth/google/signed?%v=%v", REDIRECT_DESTINATION_PARAM_NAME, url.QueryEscape(destinationUrl))
	if secret := query.Get("secret"); secret != "" {
		callbackUrl += "&secret=" + secret
	}

	loginUrl, err := user.LoginURL(c, callbackUrl)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, loginUrl, http.StatusFound)
}

func handleSignedWithGoogle(c context.Context, w http.ResponseWriter, r *http.Request) {
	var userID string
	if authInfo, _, err := auth.Authenticate(w, r, false); err != nil {
		if err != auth.ErrNoToken {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	} else {
		userID = authInfo.UserID
	}

	clientInfo := models.NewClientInfoFromRequest(r)

	googleUser := user.Current(c)
	if googleUser == nil {
		err := errors.New("handleSignedWithGoogle(): googleUser == nil")
		log.Errorf(c, err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	userGoogle, _, err := facade.User.GetOrCreateUserGoogleOnSignIn(c, googleUser, userID, clientInfo)
	if err != nil {
		log.Errorf(c, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}

	if !userGoogle.Record.Exists() {
		log.Errorf(c, "userGoogle.UserGoogleData == nil")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("userGoogle.UserGoogleData == nil"))
	}

	accountData := userGoogle.Data().(*appuser.AccountDataBase)
	appUserID := userGoogle.Data().GetAppUserID()
	log.Debugf(c, "userGoogle.AppUserIntID: %s", appUserID)
	token := auth.IssueToken(appUserID, "web", accountData.EmailLowerCase == "alexander.trakhimenok@gmail.com")
	destinationUrl := r.URL.Query().Get(REDIRECT_DESTINATION_PARAM_NAME)

	var delimiter string
	if strings.Contains(destinationUrl, "#") {
		delimiter = "&"
	} else {
		delimiter = "#"
	}
	destinationUrl += delimiter + "signed-in-with=google"
	destinationUrl += "&secret=" + token
	http.Redirect(w, r, destinationUrl, http.StatusFound)
}
