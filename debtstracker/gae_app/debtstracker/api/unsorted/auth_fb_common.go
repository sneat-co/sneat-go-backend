package unsorted

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	fb "github.com/strongo/facebook"
	"github.com/strongo/logus"
	"net/http"
)

var ErrUnauthorized = errors.New("Unauthorized")
var ErrBadRequest = errors.New("Bad request")

func signInFbUser(c context.Context, fbAppID, fbUserID string, r *http.Request, authInfo auth.AuthInfo) (
	user models.AppUser, isNewUser bool, userFacebook models.UserFacebook, fbApp *fb.App, fbSession *fb.Session, err error,
) {
	logus.Debugf(c, "api.signInFbUser()")

	if fbAppID == "" {
		panic("fbAppID is empty string")
	}
	if fbUserID == "" {
		panic("fbUserID is empty string")
	}

	signedRequest := r.PostFormValue("signed_request")
	accessToken := r.PostFormValue("access_token")

	//var isFbm bool

	// Create FB API Session
	{
		if signedRequest != "" && accessToken != "" {
			err = fmt.Errorf("%w: Parameters access_token & signed_request should not be passed together", ErrBadRequest)
			return
		} else if accessToken != "" {
			panic("not imlemented")
			//_, fbSession, err = fbmbots.FbAppAndSessionFromAccessToken(c, r, accessToken)
		} else if signedRequest != "" {
			panic("not imlemented")
			//var (
			//	signedData fb.Result
			//)
			//if fbApp, _, err = fbmbots.GetFbAppAndHost(r); err != nil {
			//	return
			//}
			//if signedData, err = fbApp.ParseSignedRequest(signedRequest); err != nil {
			//	return
			//}
			//psID := signedData.Get("psid").(string)
			//if psID != "" {
			//	if fbUserID == "" {
			//		fbUserID = psID
			//	} else if fbUserID != psID {
			//		err = fmt.Errorf("%w: fbUserID != psID", ErrBadRequest)
			//		return
			//	}
			//	var (
			//		pageID string
			//		ok     bool
			//	)
			//	if pageID, ok = signedData.Get("page_id").(string); !ok {
			//		pageID = strconv.FormatFloat(signedData.Get("page_id").(float64), 'f', 0, 64)
			//	}
			//
			//	logus.Debugf(c, "pageID: %v, signedData: %v", pageID, signedData)
			//	if fbmBot, ok := fbmbots.Bots(c).ByID[pageID]; !ok {
			//		err = errors.New("ReferredTo settings not found by page ID=" + pageID)
			//	} else {
			//		isFbm = true
			//		_, fbSession, err = fbmbots.FbAppAndSessionFromAccessToken(c, r, fbmBot.Token)
			//	}
			//} else {
			//	err = fmt.Errorf("Not implemented for signed request: %v", signedData)
			//	return
			//}
		} else {
			err = fmt.Errorf("%w: Either access_token or signed_request should be passed", ErrBadRequest)
			return
		}
		//if err != nil {
		//	err = fmt.Errorf("%w: %v", ErrUnauthorized, err.Error())
		//	return
		//}
	}

	//if userFacebook, err = dtdal.UserFacebook.GetFbUserByFbID(c, fbAppID, fbUserID); err != nil && !dal.IsNotFound(err) {
	//	err = fmt.Errorf("%w: Failed to get UserFacebook record by ID", err)
	//	return
	//} else if !dal.IsNotFound(err) && fbUserID != "" && fbUserID != userFacebook.FbUserOrPageScopeID {
	//	err = fmt.Errorf("%w: fbUserID:%v != userFacebook.ID:%v", ErrUnauthorized, fbUserID, userFacebook.FbUserOrPageScopeID)
	//	return
	//}
	//
	//if accessToken != "" || userFacebook.Data == nil || userFacebook.Data.DtUpdated.Before(time.Now().Add(-1*time.Hour)) {
	//	if user, userFacebook, isNewUser, err = createOrUpdateFbUserDbRecord(c, isFbm, fbAppID, fbUserID, fbSession, authInfo, models.NewClientInfoFromRequest(r)); err != nil {
	//		return
	//	}
	//} else {
	//	logus.Debugf(c, "Not updating FB user db record as last updated less then an hour ago")
	//}
	//
	//if err != nil {
	//	return
	//} else if user.ID == 0 {
	//	panic("userID == 0")
	//} else if user.Data == nil {
	//	panic("user.DebutsAppUserDataOBSOLETE == nil")
	//}
	//
	//return
}

//func getFbUserInfo(c context.Context, fbSession *fb.Session, isFbm bool, fbUserID string,
//) (
//	emailConfirmed bool, email, firstName, lastName string, err error,
//) {
//	var (
//		endPoint string
//		fields   string
//	)
//	if isFbm {
//		endPoint = "/" + fbUserID
//		fields = "first_name,last_name,profile_pic,locale,timezone,gender,is_payment_enabled,last_ad_referral"
//	} else {
//		endPoint = "/me"
//		fields = "email,first_name,last_name" //TODO: Try to add fields matching FBM case above. profile_pic is not OK :(
//	}
//	fbResp, err := fbSession.Get(endPoint, fb.Params{
//		"fields": fields,
//	})
//
//	if err != nil {
//		err = fmt.Errorf("%w: Failed to call Facebook API", err)
//		return
//	}
//
//	logus.Debugf(c, "Facebook API response: %v", fbResp)
//
//	var ok bool
//	if email, ok = fbResp["email"].(string); ok && email != "" {
//		emailConfirmed = true
//	} else {
//		email = fmt.Sprintf("%v@fb", fbUserID)
//	}
//
//	firstName, _ = fbResp["first_name"].(string)
//	lastName, _ = fbResp["last_name"].(string)
//	return
//}
//
//func createOrUpdateFbUserDbRecord(c context.Context, isFbm bool, fbAppID, fbUserID string, fbSession *fb.Session, authInfo auth.AuthInfo, clientInfo models.ClientInfo) (user models.AppUser, userFacebook models.UserFacebook, isNewUser bool, err error) {
//	var (
//		emailConfirmed             bool
//		email, firstName, lastName string
//	)
//	emailConfirmed, email, firstName, lastName, err = getFbUserInfo(c, fbSession, isFbm, fbUserID)
//
//	userFacebook, user, err = facade.User.GetOrCreateUserFacebookOnSignIn(c, authInfo.UserID, fbAppID, fbUserID, firstName, lastName, email, emailConfirmed, clientInfo)
//	if err != nil {
//		return
//	}
//	return
//}

func authWriteResponseForAuthFailed(c context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ErrUnauthorized) {
		w.WriteHeader(http.StatusUnauthorized)
		logus.Debugf(c, err.Error())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		logus.Errorf(c, "Auth failed: %v", err.Error())
	}
	_, _ = w.Write([]byte(err.Error()))
}

func authWriteResponseForUser(c context.Context, w http.ResponseWriter, user models.AppUser, isNewUser bool) {
	ReturnToken(c, w, user.ID, isNewUser, user.Data.EmailConfirmed && IsAdmin(user.Data.EmailAddress))
}
