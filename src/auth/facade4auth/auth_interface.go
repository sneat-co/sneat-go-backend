package facade4auth

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"sync"
	"time"
)

type CreateUserData struct {
	//FbUserID     string
	//GoogleUserID string
	//VkUserID     int64
	FirstName  string
	LastName   string
	ScreenName string
	Nickname   string
}

func CreateUserEntity(createUserData CreateUserData) (user *models4debtus.DebutsAppUserDataOBSOLETE) {
	return &models4debtus.DebutsAppUserDataOBSOLETE{
		//FbUserID: createUserData.FbUserID,
		//VkUserID: createUserData.VkUserID,
		//GoogleUniqueUserID: createUserData.GoogleUserID,
		//ContactDetails: dto4contactus.ContactDetails{
		//	NameFields: person.NameFields{
		//		FirstName:  createUserData.FirstName,
		//		LastName:   createUserData.LastName,
		//		ScreenName: createUserData.ScreenName,
		//		NickName:   createUserData.Nickname,
		//	},
		//},
	}
}

type UserDal interface {
	GetUserByStrID(c context.Context, userID string) (dbo4userus.UserEntry, error)
	GetUserByVkUserID(c context.Context, vkUserID int64) (dbo4userus.UserEntry, error)
	CreateAnonymousUser(c context.Context) (dbo4userus.UserEntry, error)
	CreateUser(c context.Context, userEntity *dbo4userus.UserDbo) (dbo4userus.UserEntry, error)
	DelaySetUserPreferredLocale(c context.Context, delay time.Duration, userID string, localeCode5 string) error
}

type PasswordResetDal interface {
	GetPasswordResetByID(c context.Context, tx dal.ReadSession, id int) (models4auth.PasswordReset, error)
	CreatePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, entity *models4auth.PasswordResetData) (models4auth.PasswordReset, error)
	SavePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, record models4auth.PasswordReset) (err error)
}

type UserGoogleDal interface {
	GetUserGoogleByID(c context.Context, googleUserID string) (userGoogle models4auth.UserAccountEntry, err error)
	DeleteUserGoogle(c context.Context, googleUserID string) (err error)
}

type UserVkDal interface {
	GetUserVkByID(c context.Context, vkUserID int64) (userGoogle models4auth.UserVk, err error)
	SaveUserVk(c context.Context, userVk models4auth.UserVk) (err error)
}

type UserEmailDal interface {
	GetUserEmailByID(c context.Context, tx dal.ReadSession, email string) (userEmail models4auth.UserEmailEntry, err error)
	SaveUserEmail(c context.Context, tx dal.ReadwriteTransaction, userEmail models4auth.UserEmailEntry) (err error)
}

type UserGooglePlusDal interface {
	GetUserGooglePlusByID(c context.Context, id string) (userGooglePlus models4auth.UserGooglePlus, err error)
	//SaveUserGooglePlusByID(c context.Context, userGooglePlus models4auth.UserGooglePlus) (err error)
}

type UserFacebookDal interface {
	GetFbUserByFbID(c context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (fbUser models4auth.UserFacebook, err error)
	SaveFbUser(c context.Context, tx dal.ReadwriteTransaction, fbUser models4auth.UserFacebook) (err error)
	DeleteFbUser(c context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (err error)
	//CreateFbUserRecord(c context.Context, fbUserID string, appUserID int64) (fbUser models.UserFacebook, err error)
}

type LoginPinDal interface {
	GetLoginPinByID(c context.Context, tx dal.ReadSession, loginID int) (loginPin models4auth.LoginPin, err error)
	SaveLoginPin(c context.Context, tx dal.ReadwriteTransaction, loginPin models4auth.LoginPin) (err error)
	CreateLoginPin(c context.Context, tx dal.ReadwriteTransaction, channel, gaClientID string, createdUserID string) (loginPin models4auth.LoginPin, err error)
}

type LoginCodeDal interface {
	NewLoginCode(c context.Context, userID string) (code int, err error)
	ClaimLoginCode(c context.Context, code int) (userID string, err error)
}

type TgChatDal interface {
	GetTgChatByID(c context.Context, tgBotID string, tgChatID int64) (tgChat models4debtus.DebtusTelegramChat, err error)
	DoSomething( // TODO: WTF name?
		c context.Context,
		userTask *sync.WaitGroup,
		currency string,
		tgChatID int64,
		authInfo token4auth.AuthInfo,
		user dbo4userus.UserEntry,
		sendToTelegram func(tgChat botsfwtgmodels.TgChatData) error,
	) (err error)
}

type TgUserDal interface {
	FindByUserName(c context.Context, tx dal.ReadSession, userName string) (tgUsers []botsfwtgmodels.TgPlatformUser, err error)
}

var User UserDal

var UserFacebook UserFacebookDal

var UserGooglePlus UserGooglePlusDal

var PasswordReset PasswordResetDal

var UserEmail UserEmailDal

var TgChat TgChatDal

var TgUser TgUserDal
