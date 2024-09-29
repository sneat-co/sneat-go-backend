package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/common4all"
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
	unsorted4auth.User = facade4auth.NewUserDalGae()
	//dtdal.Bill = facade4splitus.newBillDalGae()
	//dtdal.Split = splitDalGae{}
	dtdal.TgGroup = facade4auth.NewTgGroupDalGae()
	//dtdal.BillSchedule = NewBillScheduleDalGae()
	dtdal.Receipt = NewReceiptDalGae()
	dtdal.Reminder = NewReminderDalGae()
	unsorted4auth.UserGoogle = facade4auth.NewUserGoogleDalGae()
	unsorted4auth.PasswordReset = facade4auth.NewPasswordResetDalGae()
	common4all.Email = NewEmailDalGae()
	unsorted4auth.UserGooglePlus = facade4auth.NewUserGooglePlusDalGae()
	unsorted4auth.UserEmail = facade4auth.NewUserEmailGaeDal()
	unsorted4auth.UserFacebook = facade4auth.NewUserFacebookDalGae()
	unsorted4auth.LoginPin = facade4auth.NewLoginPinDalGae()
	unsorted4auth.LoginCode = facade4auth.NewLoginCodeDalGae()
	dtdal.Twilio = NewTwilioDalGae()
	dtdal.Invite = NewInviteDalGae()
	dtdal.Admin = NewAdminDalGae()
	unsorted4auth.TgChat = facade4auth.NewTgChatDalGae()
	unsorted4auth.TgUser = facade4auth.NewTgUserDalGae()
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
