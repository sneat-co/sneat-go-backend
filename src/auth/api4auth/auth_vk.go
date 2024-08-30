package api4auth

import (
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"net/http"

	"context"
	//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/dtdal/gaedal"
	//"google.golang.org/appengine/v2/datastore"
	//"errors"
	//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/auth"
	//"strconv"
	//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/dtdal"
	//"github.com/strongo/nds"
	//"github.com/strongo/vk"
	//"strings"
	//"fmt"
	//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
)

//const VK_USER_ALEXT = 7631716

func HandleSignedWithVK(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	panic("disabled")
	//err := r.ParseForm()
	//if err != nil {
	//	BadRequestError(ctx, hashedWriter, err)
	//	return
	//}
	//
	//var (
	//	vkUserID int64
	//	vkUserEntity models.UserVkEntity
	//	//userAuth
	//)
	//if vkUserID, err = strconv.ParseInt(r.PostFormValue("vkUserID"), 10, 64); err != nil {
	//	BadRequestError(ctx, hashedWriter, errors.Wrap(err, "Missing or bad vkUserID"))
	//	return
	//}
	//
	//logTodoMergeUsers := func(userID int64) {
	//	m := fmt.Sprintf("TODO: Merge users: userID=%v, authInfo.AppUserIntID=%v", userID, authInfo.UserID)
	//	//logus.Errorf(ctx, m)
	//	hashedWriter.WriteHeader(http.StatusInternalServerError)
	//	hashedWriter.Write([]byte(m))
	//}
	//
	//
	//// Try to get UserVk entity by key and if it has AppUserIntID we can create and return token right away
	//vkUser, err := dtdal.UserVk.GetUserVkByID(ctx, vkUserID)
	//if err != nil {
	//	if dal.IsNotFound(err) {  // It's OK if UserVk entity not found.
	//		logus.Debugf(ctx, "UserVk entity not found by vkUserID=%d", vkUserID)
	//	} else {  // For other errors fail gracefully.
	//		InternalError(ctx, hashedWriter, err)
	//		return
	//	}
	//}
	//if vkUserEntity.UserID == 0 {
	//	// For some reason we have UserVk entity without associated AppUser
	//	logus.Warningf(ctx, "vkUserEntity.AppUserIntID == 0 - TOOD: Create user?")// TODO: Create user?
	//} else {
	//	if authInfo.UserID != 0 && vkUserEntity.UserID != authInfo.UserID {
	//		logTodoMergeUsers(vkUserEntity.UserID)
	//		return
	//	}
	//	logus.Debugf(ctx, "UserVk entity found by key and has AppUserIntID=%v", vkUserEntity.UserID)
	//	ReturnToken(ctx, hashedWriter, vkUserEntity.UserID, vkUserID == VK_USER_ALEXT)
	//	return
	//}
	//
	//accessToken := r.PostFormValue("vkAccessToken")
	//vkLanguage := r.PostFormValue("vkLanguage")
	//
	//if accessToken == "" {
	//	BadRequestError(ctx, hashedWriter, errors.New("Missing accessToken"))
	//	return
	//}
	//if vkLanguage == "" {
	//	BadRequestError(ctx, hashedWriter, errors.New("Missing vkLanguage"))
	//	return
	//}
	//
	//vkApi := vk.NewApiWithAccessToken(dtdal.HttpClient(ctx), accessToken)
	//
	//vkUserInfo, err := vkApi.GetUserByIntID(ctx, vkUserID, "nom", vk.FieldFirstName, vk.FieldLastName, vk.FieldNickname, vk.FieldScreenName)
	//if err != nil {
	//	if vkErr, ok := err.(vk.VkError); ok && vkErr.VkErrorCode() == 5 && strings.Contains(err.Error(), "access_token was given to another ip address") {
	//		// Good access token
	//	} else {
	//		err = errors.Wrap(err, "Failed to get verify VK access token")
	//		logus.Warningf(ctx, err.Error())
	//		hashedWriter.WriteHeader(http.StatusInternalServerError)
	//		hashedWriter.Write([]byte(err.Error()))
	//		return
	//	}
	//}
	//vkUser.FirstName = r.PostFormValue("firstName")
	//vkUser.LastName = r.PostFormValue("lastName")
	//
	//userID, user, err := dtdal.User.GetUserByVkUserID(ctx, vkUserID)
	//if err == nil && user.VkUserID == vkUserID {
	//	// For some reason we have a user with VkUserID but without UserVk entity
	//	err = dtdal.DB.RunInTransaction(ctx, func(tctx context.Context) error {
	//		if err = nds.Get(tc, vkUserKey, &vkUserEntity); err != nil {
	//			if err == datastore.ErrNoSuchEntity {
	//				vkUserEntity = models.UserVkEntity{
	//					UserID: userID,
	//					FirstName: vkUser.FirstName,
	//					LastName: vkUser.LastName,
	//					ScreenName: vkUser.ScreenName,
	//					Nickname: vkUser.Nickname,
	//				}
	//				if _, err = nds.Put(tctx, vkUserKey, &vkUserEntity); err != nil {
	//					err = errors.Wrap(err, "Failed to create a UserVk entity")
	//					return err
	//				}
	//				return nil
	//			} else {
	//				err = errors.Wrapf(err, "Failed to get UserVk entity by key=%v", vkUserKey)
	//				return err
	//			}
	//		}
	//		return nil
	//	}, nil)
	//} else if err == nil && err != datastore.ErrNoSuchEntity {
	//	err = errors.Wrap(err, "Failed to get user by VkUserID")
	//	InternalError(ctx, hashedWriter, err)
	//	return
	//}
	//
	//updateUser := func() (err error) {
	//	if (user.FirstName == "" && vkUser.FirstName != "") || (user.LastName == "" && vkUser.LastName != "") || (user.ScreenName == "" && vkUser.ScreenName != "") || (user.Nickname == "" && vkUser.Nickname != "") {
	//		var changed bool
	//		err = dtdal.DB.gaedb.RunInTransaction(ctx, func(ctx context.Context) error {
	//			user, err := userDal.GetUserByIdOBSOLETE(ctx, userID)
	//			if err != nil {
	//				return err
	//			}
	//			if user.VkUserID == 0 {
	//				user.VkUserID = vkUserID
	//				changed = true
	//			}
	//			if user.FirstName == "" && vkUser.FirstName != "" {
	//				user.FirstName = vkUser.FirstName
	//				changed = true
	//			}
	//			if user.LastName == "" && vkUser.LastName != "" {
	//				user.LastName = vkUser.LastName
	//				changed = true
	//			}
	//			if user.ScreenName == "" && vkUser.ScreenName != "" {
	//				user.ScreenName = vkUser.ScreenName
	//				changed = true
	//			}
	//			if user.Nickname == "" && vkUser.Nickname != "" {
	//				user.Nickname = vkUser.Nickname
	//				changed = true
	//			}
	//			if changed {
	//				_, err = nds.Put(ctx, gaedal.NewAppUserKey(ctx, userID), user)
	//			}
	//			return err
	//		}, nil)
	//		if err != nil {
	//			err = errors.Wrap(err, "Failed to update user with VkUserID")
	//			InternalError(ctx, hashedWriter, err)
	//			return
	//		}
	//		if changed {
	//			logus.Infof(ctx, "User update with VK info")
	//		}
	//	}
	//	return nil
	//}
	//
	//if userID != 0 {
	//	if authInfo.UserID != 0 && userID != authInfo.UserID {
	//		logTodoMergeUsers(userID)
	//		return
	//	}
	//	if err = updateUser(); err != nil {
	//		return
	//	}
	//	ReturnToken(ctx, hashedWriter, userID, vkUserID == VK_USER_ALEXT)
	//	return
	//}
	//
	//if authInfo.UserID != 0 {
	//	if user, err = dal4userus.GetUserByID(ctx, authInfo.UserID); err != nil {
	//		if err == datastore.ErrNoSuchEntity {
	//			logus.Warningf(ctx, "User not found ContactID=%v", authInfo.UserID)
	//		} else {
	//			logus.Errorf(ctx, err.Error())
	//		}
	//		InternalError(ctx, hashedWriter, err)
	//		return
	//	}
	//	if user.VkUserID != 0 && user.VkUserID != vkUserID {
	//		logTodoMergeUsers(userID)
	//		return
	//	} else if err = updateUser(); err != nil {
	//		return
	//	}
	//	ReturnToken(ctx, hashedWriter, userID, vkUserID == VK_USER_ALEXT)
	//	return
	//}
	//
	//createUserData := dtdal.CreateUserData{
	//	VkUserID: vkUserID,
	//	FirstName: vkUser.FirstName,
	//	LastName: vkUser.LastName,
	//	ScreenName: vkUser.ScreenName,
	//	Nickname: vkUser.Nickname,
	//}
	//
	//userID, _, err = facade4debtus.User.GetOrCreateUserByEmail(ctx, "", false, &createUserData)
	//if err != nil {
	//	InternalError(ctx, hashedWriter, err)
	//	return
	//}
	//ReturnToken(ctx, hashedWriter, userID, vkUserID == VK_USER_ALEXT)
}
