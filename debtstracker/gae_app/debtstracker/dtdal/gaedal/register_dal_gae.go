package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-core/facade"
	apphostgae "github.com/strongo/app-host-gae"
	"google.golang.org/appengine/v2"
	"net/http"
)

func RegisterDal() {

	//dtdal.DB = gaedb.NewDatabase()
	//telegramBot.Init(facade2debtus.GetDatabase)
	//
	dtdal.Contact = NewContactDalGae()
	dtdal.Transfer = NewTransferDalGae()
	//dtdal.Reward = NewRewardDalGae()
	dtdal.User = NewUserDalGae()
	dtdal.Bill = newBillDalGae()
	//dtdal.Split = splitDalGae{}
	dtdal.TgGroup = newTgGroupDalGae()
	//dtdal.BillSchedule = NewBillScheduleDalGae()
	dtdal.Receipt = NewReceiptDalGae()
	dtdal.Reminder = NewReminderDalGae()
	dtdal.UserBrowser = NewUserBrowserDalGae()
	dtdal.UserGoogle = NewUserGoogleDalGae()
	dtdal.PasswordReset = NewPasswordResetDalGae()
	dtdal.Email = NewEmailDalGae()
	dtdal.UserGooglePlus = NewUserGooglePlusDalGae()
	dtdal.UserEmail = NewUserEmailGaeDal()
	dtdal.UserFacebook = NewUserFacebookDalGae()
	dtdal.LoginPin = NewLoginPinDalGae()
	dtdal.LoginCode = NewLoginCodeDalGae()
	dtdal.Twilio = NewTwilioDalGae()
	dtdal.Invite = NewInviteDalGae()
	dtdal.Admin = NewAdminDalGae()
	dtdal.TgChat = NewTgChatDalGae()
	dtdal.TgUser = NewTgUserDalGae()
	dtdal.Group = NewGroupDalGae()
	dtdal.UserOneSignal = NewUserOneSignalDalGae()
	dtdal.UserGaClient = NewUserGaClientDalGae()
	dtdal.Feedback = NewFeedbackDalGae()
	//dtdal.UserVk = NewUserVkDalGae()
	//dtdal.GroupMember = NewGroupMemberDalGae()
	dtdal.HttpClient = func(c context.Context) *http.Client {
		return http.DefaultClient
		//return urlfetch.Client(c)
	}
	dtdal.HttpAppHost = apphostgae.NewHttpAppHostGAE()

	//dtdal.HandleWithContext = func(handler strongoapp.HttpHandlerWithContext) func(w http.ResponseWriter, r *http.Request) {
	//	return func(w http.ResponseWriter, r *http.Request) {
	//		handler(appengine.NewContext(r), w, r)
	//	}
	//}
	//dtdal.TaskQueue = TaskQueueDalGae{}
	dtdal.BotHost = ApiBotHost{}
}

type ApiBotHost struct {
}

func (h ApiBotHost) Context(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func (h ApiBotHost) GetHTTPClient(c context.Context) *http.Client {
	return dtdal.HttpClient(c)
}

//func (h ApiBotHost) GetBotCoreStores(platform string, appContext botsfw.BotAppContext, r *http.Request) botsfwdal.DataAccess {
//	panic("Not implemented")
//}

func (h ApiBotHost) DB(c context.Context) (dal.DB, error) {
	return facade.GetDatabase(c)
}
