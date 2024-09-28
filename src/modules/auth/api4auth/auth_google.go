package api4auth

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/token4auth"
	"net/http"
	//"strings"
	//"encoding/json"
	//"io/ioutil"
	//"time"
	//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/dtdal"
	//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
	//"errors"
	//"github.com/pquerna/ffjson/ffjson"
	//"github.com/dal-go/dalgo"
	//"github.com/strongo/logus"
)

type GoogleAuthData struct {
	UserId         string `json:"userId" firestore:"userId"`
	DisplayName    string `json:"displayName" firestore:"displayName"`
	RefreshToken   string `json:"refreshToken" firestore:"refreshToken"`
	Email          string `json:"email" firestore:"email"`
	ServerAuthCode string `json:"serverAuthCode" firestore:"serverAuthCode"`
	AccessToken    string `json:"accessToken" firestore:"accessToken"`
	ImageUrl       string `json:"imageUrl" firestore:"imageUrl"`
	IdToken        string `json:"idToken" firestore:"idToken"`
}

func HandleSignedInWithGooglePlus(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	panic("Not implemented")
	//decoder := ffjson.NewDecoder()
	//googleAuthData := GoogleAuthData{}
	//defer r.Body.Close()
	//if err := decoder.DecodeReader(r.Body, &googleAuthData); err != nil {
	//	BadRequestError(ctx, w, err)
	//	return
	//}
	//
	//if googleAuthData.UserId == "" {
	//	BadRequestMessage(ctx, w, "Missing required field: userId")
	//	return
	//}
	//
	//if googleAuthData.Email == "" {
	//	BadRequestMessage(ctx, w, "Missing required field: email")
	//	return
	//}
	//
	//tokenData := make(map[string]string, 16)
	//
	//// TODO: https://developers.google.com/identity/sign-in/web/backend-auth - verify "aud" and check "sub" fields
	//if resp, err := dtdal.HttpClient(ctx).Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + googleAuthData.IdToken); err != nil {
	//	ErrorAsJson(ctx, w, http.StatusBadRequest, fmt.Errorf("%w: Failed to call googleapis", err))
	//	return
	//} else if resp.StatusCode != 200 {
	//	defer resp.Body.Close()
	//	if body, err := ioutil.ReadAll(resp.Body); err != nil {
	//		BadRequestMessage(ctx, w, "Failed to verify idToken")
	//	} else {
	//		BadRequestMessage(ctx, w, "Failed to verify idToken: "+string(body))
	//	}
	//	return
	//} else {
	//	defer resp.Body.Close()
	//	if body, err := ioutil.ReadAll(resp.Body); err != nil {
	//		ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.Wrap(err, "Failed to read response body"))
	//		return
	//	} else {
	//		logus.Infof(ctx, "idToken verified: %s", string(body))
	//		if err = json.Unmarshal(body, &tokenData); err != nil {
	//			ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.Wrap(err, "Failed to unmarshal response body as JSON"))
	//			return
	//		}
	//		if aud, ok := tokenData["aud"]; !ok || !strings.HasPrefix(aud, "74823955721-") {
	//			BadRequestMessage(ctx, w, "idToken data has unexpected AUD: "+aud)
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
	//	ctx,
	//	func(ctx context.Context) (err error) {
	//		userGooglePlus, err = dtdal.UserGooglePlus.GetUserGooglePlusByID(ctx, googleAuthData.UserId)
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
	//			//if user, isNewUser, err = facade4debtus.User.GetOrCreateEmailUser(ctx, googleAuthData.Email, userGooglePlus.EmailVerified, &createUserData); err != nil {
	//			//	return
	//			//}
	//			userGooglePlus.AppUserIntID = user.ContactID
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
	//			if err = dtdal.UserGooglePlus.SaveUserGooglePlusByID(ctx, userGooglePlus); err != nil {
	//				return
	//			}
	//		}
	//		return nil
	//	},
	//	dtdal.CrossGroupTransaction,
	//)
	//
	//if err != nil {
	//	ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
	//
	//ReturnToken(ctx, w, userGooglePlus.AppUserIntID, isNewUser, googleAuthData.Email == "alexander.trakhimenok@gmail.com")
}
