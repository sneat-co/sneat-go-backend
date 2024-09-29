package unsorted4auth

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/anybot"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
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

// CreateUserEntity
// Deprecated
func CreateUserEntity(createUserData CreateUserData) (user dbo4userus.UserEntry) {
	return
	//return &models4debtus.DebutsAppUserDataOBSOLETE{
	//	//FbUserID: createUserData.FbUserID,
	//	//VkUserID: createUserData.VkUserID,
	//	//GoogleUniqueUserID: createUserData.GoogleUserID,
	//	//ContactDetails: dto4contactus.ContactDetails{
	//	//	NameFields: person.NameFields{
	//	//		FirstName:  createUserData.FirstName,
	//	//		LastName:   createUserData.LastName,
	//	//		ScreenName: createUserData.ScreenName,
	//	//		NickName:   createUserData.Nickname,
	//	//	},
	//	//},
	//}
}

type UserDal interface {
	GetUserByStrID(ctx context.Context, userID string) (dbo4userus.UserEntry, error)
	GetUserByVkUserID(ctx context.Context, vkUserID int64) (dbo4userus.UserEntry, error)
	CreateAnonymousUser(ctx context.Context) (dbo4userus.UserEntry, error)
	CreateUser(ctx context.Context, userEntity *dbo4userus.UserDbo) (dbo4userus.UserEntry, error)
	DelaySetUserPreferredLocale(ctx context.Context, delay time.Duration, userID string, localeCode5 string) error
}

type PasswordResetDal interface {
	GetPasswordResetByID(ctx context.Context, tx dal.ReadSession, id int) (models4auth.PasswordReset, error)
	CreatePasswordResetByID(ctx context.Context, tx dal.ReadwriteTransaction, entity *models4auth.PasswordResetData) (models4auth.PasswordReset, error)
	SavePasswordResetByID(ctx context.Context, tx dal.ReadwriteTransaction, record models4auth.PasswordReset) (err error)
}

type UserGoogleDal interface {
	GetUserGoogleByID(ctx context.Context, googleUserID string) (userGoogle models4auth.UserAccountEntry, err error)
	DeleteUserGoogle(ctx context.Context, googleUserID string) (err error)
}

type UserVkDal interface {
	GetUserVkByID(ctx context.Context, vkUserID int64) (userGoogle models4auth.UserVk, err error)
	SaveUserVk(ctx context.Context, userVk models4auth.UserVk) (err error)
}

type UserEmailDal interface {
	GetUserEmailByID(ctx context.Context, tx dal.ReadSession, email string) (userEmail models4auth.UserEmailEntry, err error)
	SaveUserEmail(ctx context.Context, tx dal.ReadwriteTransaction, userEmail models4auth.UserEmailEntry) (err error)
}

type UserGooglePlusDal interface {
	GetUserGooglePlusByID(ctx context.Context, id string) (userGooglePlus models4auth.UserGooglePlus, err error)
	//SaveUserGooglePlusByID(ctx context.Context, userGooglePlus models4auth.UserGooglePlus) (err error)
}

type UserFacebookDal interface {
	GetFbUserByFbID(ctx context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (fbUser models4auth.UserFacebook, err error)
	SaveFbUser(ctx context.Context, tx dal.ReadwriteTransaction, fbUser models4auth.UserFacebook) (err error)
	DeleteFbUser(ctx context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (err error)
	//CreateFbUserRecord(ctx context.Context, fbUserID string, appUserID int64) (fbUser models.UserFacebook, err error)
}

type LoginPinDal interface {
	GetLoginPinByID(ctx context.Context, tx dal.ReadSession, loginID int) (loginPin models4auth.LoginPin, err error)
	SaveLoginPin(ctx context.Context, tx dal.ReadwriteTransaction, loginPin models4auth.LoginPin) (err error)
	CreateLoginPin(ctx context.Context, tx dal.ReadwriteTransaction, channel, gaClientID string, createdUserID string) (loginPin models4auth.LoginPin, err error)
}

type LoginCodeDal interface {
	NewLoginCode(ctx context.Context, userID string) (code int, err error)
	ClaimLoginCode(ctx context.Context, code int) (userID string, err error)
}

type TgChatDal interface {
	GetTgChatByID(ctx context.Context, tgBotID string, tgChatID int64) (tgChat anybot.SneatAppTgChatEntry, err error)
	DoSomething( // TODO: WTF name?
		ctx context.Context,
		userTask *sync.WaitGroup,
		currency string,
		tgChatID int64,
		authInfo token4auth.AuthInfo,
		user dbo4userus.UserEntry,
		sendToTelegram func(tgChat botsfwtgmodels.TgChatData) error,
	) (err error)
}

type TgUserDal interface {
	FindByUserName(ctx context.Context, tx dal.ReadSession, userName string) (tgUsers []botsfwtgmodels.TgPlatformUser, err error)
}

var User UserDal

var UserFacebook UserFacebookDal

var UserGooglePlus UserGooglePlusDal

var PasswordReset PasswordResetDal

var UserEmail UserEmailDal

var TgChat TgChatDal

var TgUser TgUserDal
