package unsorted

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"net/http"
	//"strings"
	//"encoding/json"
	//"io/ioutil"
	//"time"
	//"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	//"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	//"errors"
	//"github.com/pquerna/ffjson/ffjson"
	//"github.com/dal-go/dalgo"
	//"github.com/strongo/logus"
)

type GoogleAuthData struct {
	UserId         string `json:"userId" datastore:",noindex"`
	DisplayName    string `json:"displayName" datastore:",noindex"`
	RefreshToken   string `json:"refreshToken" datastore:",noindex"`
	Email          string `json:"email" datastore:",noindex"`
	ServerAuthCode string `json:"serverAuthCode" datastore:",noindex"`
	AccessToken    string `json:"accessToken" datastore:",noindex"`
	ImageUrl       string `json:"imageUrl" datastore:",noindex"`
	IdToken        string `json:"idToken" datastore:",noindex"`
}

func HandleSignedInWithGooglePlus(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	panic("Not implemented")
	//decoder := ffjson.NewDecoder()
	//googleAuthData := GoogleAuthData{}
	//defer r.Body.Close()
	//if err := decoder.DecodeReader(r.Body, &googleAuthData); err != nil {
	//	BadRequestError(c, w, err)
	//	return
	//}
	//
	//if googleAuthData.UserId == "" {
	//	BadRequestMessage(c, w, "Missing required field: userId")
	//	return
	//}
	//
	//if googleAuthData.Email == "" {
	//	BadRequestMessage(c, w, "Missing required field: email")
	//	return
	//}
	//
	//tokenData := make(map[string]string, 16)
	//
	//// TODO: https://developers.google.com/identity/sign-in/web/backend-auth - verify "aud" and check "sub" fields
	//if resp, err := dtdal.HttpClient(c).Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + googleAuthData.IdToken); err != nil {
	//	ErrorAsJson(c, w, http.StatusBadRequest, fmt.Errorf("%w: Failed to call googleapis", err))
	//	return
	//} else if resp.StatusCode != 200 {
	//	defer resp.Body.Close()
	//	if body, err := ioutil.ReadAll(resp.Body); err != nil {
	//		BadRequestMessage(c, w, "Failed to verify idToken")
	//	} else {
	//		BadRequestMessage(c, w, "Failed to verify idToken: "+string(body))
	//	}
	//	return
	//} else {
	//	defer resp.Body.Close()
	//	if body, err := ioutil.ReadAll(resp.Body); err != nil {
	//		ErrorAsJson(c, w, http.StatusInternalServerError, errors.Wrap(err, "Failed to read response body"))
	//		return
	//	} else {
	//		logus.Infof(c, "idToken verified: %s", string(body))
	//		if err = json.Unmarshal(body, &tokenData); err != nil {
	//			ErrorAsJson(c, w, http.StatusInternalServerError, errors.Wrap(err, "Failed to unmarshal response body as JSON"))
	//			return
	//		}
	//		if aud, ok := tokenData["aud"]; !ok || !strings.HasPrefix(aud, "74823955721-") {
	//			BadRequestMessage(c, w, "idToken data has unexpected AUD: "+aud)
	//			return
	//		}
	//	}
	//}
	//
	//var (
	//	userGooglePlus models.UserGooglePlus
	//	isNewUser      bool
	//)
	//
	//err := dtdal.DB.RunInTransaction(
	//	c,
	//	func(c context.Context) (err error) {
	//		userGooglePlus, err = dtdal.UserGooglePlus.GetUserGooglePlusByID(c, googleAuthData.UserId)
	//		if err != nil {
	//			if dal.IsNotFound(err) {
	//				err = nil
	//				userGooglePlus.DtCreated = time.Now()
	//			} else {
	//				return
	//			}
	//		}
	//		var (
	//			changed bool
	//		)
	//
	//		userGooglePlus.EmailVerified = tokenData["email_verified"] == "true"
	//		if v, ok := tokenData["given_name"]; ok {
	//			changed = true
	//			userGooglePlus.NameFirst = v
	//		}
	//		if v, ok := tokenData["family_name"]; ok {
	//			changed = true
	//			userGooglePlus.NameLast = v
	//		}
	//		if v, ok := tokenData["locale"]; ok {
	//			changed = true
	//			userGooglePlus.Locale = v
	//		}
	//
	//		if userGooglePlus.AppUserIntID == 0 {
	//			//createUserData := dtdal.CreateUserData{
	//			//	//GoogleUserID: googleAuthData.UserId,
	//			//	FirstName:    userGooglePlus.NameFirst,
	//			//	LastName:     userGooglePlus.NameLast,
	//			//}
	//			var user models.AppUser
	//			//if user, isNewUser, err = facade.User.GetOrCreateEmailUser(c, googleAuthData.Email, userGooglePlus.EmailVerified, &createUserData); err != nil {
	//			//	return
	//			//}
	//			userGooglePlus.AppUserIntID = user.ID
	//			changed = true
	//		}
	//
	//		if googleAuthData.IdToken != "" && userGooglePlus.IdToken != googleAuthData.IdToken {
	//			userGooglePlus.IdToken = googleAuthData.IdToken
	//			changed = true
	//		}
	//
	//		if googleAuthData.AccessToken != "" && userGooglePlus.AccessToken != googleAuthData.AccessToken {
	//			userGooglePlus.AccessToken = googleAuthData.AccessToken
	//			changed = true
	//		}
	//		if googleAuthData.ServerAuthCode != "" && userGooglePlus.ServerAuthCode != googleAuthData.ServerAuthCode {
	//			userGooglePlus.ServerAuthCode = googleAuthData.ServerAuthCode
	//			changed = true
	//		}
	//		if googleAuthData.ImageURL != "" && userGooglePlus.ImageURL != googleAuthData.ImageURL {
	//			userGooglePlus.ImageURL = googleAuthData.ImageURL
	//			changed = true
	//		}
	//		if googleAuthData.RefreshToken != "" && userGooglePlus.RefreshToken != googleAuthData.RefreshToken {
	//			userGooglePlus.RefreshToken = googleAuthData.RefreshToken
	//			changed = true
	//		}
	//		if googleAuthData.DisplayName != "" && userGooglePlus.DisplayName != googleAuthData.DisplayName {
	//			userGooglePlus.DisplayName = googleAuthData.DisplayName
	//			changed = true
	//		}
	//		if googleAuthData.Email != "" && userGooglePlus.Email != googleAuthData.Email {
	//			userGooglePlus.Email = googleAuthData.Email
	//			changed = true
	//		}
	//
	//		if changed {
	//			userGooglePlus.DtUpdated = time.Now()
	//			if err = dtdal.UserGooglePlus.SaveUserGooglePlusByID(c, userGooglePlus); err != nil {
	//				return
	//			}
	//		}
	//		return nil
	//	},
	//	dtdal.CrossGroupTransaction,
	//)
	//
	//if err != nil {
	//	ErrorAsJson(c, w, http.StatusInternalServerError, err)
	//	return
	//}
	//
	//ReturnToken(c, w, userGooglePlus.AppUserIntID, isNewUser, googleAuthData.Email == "alexander.trakhimenok@gmail.com")
}
