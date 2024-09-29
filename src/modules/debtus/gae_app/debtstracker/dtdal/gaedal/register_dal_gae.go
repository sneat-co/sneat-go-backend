package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	facade4auth2 "github.com/sneat-co/sneat-core-modules/auth/facade4auth"
	unsorted4auth2 "github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func RegisterDal() {

	//dtdal.DB = gaedb.NewDatabase()
	//telegramBot.Init(facade4debtus.GetDatabase)
	//
	dtdal.Contact = NewContactDalGae()
	dtdal.Transfer = NewTransferDalGae()
	//dtdal.Reward = NewRewardDalGae()
	unsorted4auth2.User = facade4auth2.NewUserDalGae()
	//dtdal.Bill = facade4splitus.newBillDalGae()
	//dtdal.Split = splitDalGae{}
	dtdal.TgGroup = facade4auth2.NewTgGroupDalGae()
	//dtdal.BillSchedule = NewBillScheduleDalGae()
	dtdal.Receipt = NewReceiptDalGae()
	dtdal.Reminder = NewReminderDalGae()
	unsorted4auth2.UserGoogle = facade4auth2.NewUserGoogleDalGae()
	unsorted4auth2.PasswordReset = facade4auth2.NewPasswordResetDalGae()
	common4all.Email = NewEmailDalGae()
	unsorted4auth2.UserGooglePlus = facade4auth2.NewUserGooglePlusDalGae()
	unsorted4auth2.UserEmail = facade4auth2.NewUserEmailGaeDal()
	unsorted4auth2.UserFacebook = facade4auth2.NewUserFacebookDalGae()
	unsorted4auth2.LoginPin = facade4auth2.NewLoginPinDalGae()
	unsorted4auth2.LoginCode = facade4auth2.NewLoginCodeDalGae()
	dtdal.Twilio = NewTwilioDalGae()
	dtdal.Invite = NewInviteDalGae()
	dtdal.Admin = NewAdminDalGae()
	unsorted4auth2.TgChat = facade4auth2.NewTgChatDalGae()
	unsorted4auth2.TgUser = facade4auth2.NewTgUserDalGae()
	//dtdal.Group = facade4splitus.NewGroupDalGae()
	dtdal.Feedback = NewFeedbackDalGae()
	//dtdal.UserVk = NewUserVkDalGae()
	//dtdal.GroupMember = NewGroupMemberDalGae()
	dtdal.HttpClient = func(ctx context.Context) *http.Client {
		return http.DefaultClient
		//return urlfetch.Client(ctx)
	}
	//dtdal.HttpAppHost = apphostgae.NewHttpAppHostGAE()

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
	return r.Context()
}

func (h ApiBotHost) GetHTTPClient(ctx context.Context) *http.Client {
	return dtdal.HttpClient(ctx)
}

//func (h ApiBotHost) GetBotCoreStores(platform string, appContext botsfw.BotAppContext, r *http.Request) botsfwdal.DataAccess {
//	panic("Not implemented")
//}

func (h ApiBotHost) DB(ctx context.Context) (dal.DB, error) {
	return facade.GetSneatDB(ctx)
}
