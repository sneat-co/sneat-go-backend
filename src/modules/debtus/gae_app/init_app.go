package gaeapp

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/anybot/facade4anybot"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/emailing"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_transfer"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus/apigaedepended"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus/unsorted"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/apimapping"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/apps/vkapp"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal/gaedal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/maintainance"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/reminders"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/support"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/webhooks"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/website"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/botcmds4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/strongo/delaying"
	"net/http"
)

// Init initializes debts tracker server
func Init(botHost botsfw.BotHost) {
	if botHost == nil {
		panic("botHost parameter is required")
	}

	initDelaying()

	gaedal.RegisterDal()
	apigaedepended.InitApiGaeDepended()

	httpRouter := httprouter.New()
	http.Handle("/", httpRouter)

	apimapping.InitApi(httpRouter)
	website.InitWebsite(httpRouter)
	webhooks.InitWebhooks(httpRouter)
	vkapp.InitVkIFrameApp(httpRouter)
	support.InitSupportHandlers(httpRouter)

	InitCronHandlers(httpRouter)
	InitTaskQueueHandlers(httpRouter)

	//InitBots(httpRouter, botHost, nil)

	//httpRouter.GET("/test-pointer", testModelPointer)
	httpRouter.GET("/Users/astec/", NotFoundSilent)

	maintainance.RegisterMappers()
}

func initDelaying() {
	delaying.Init(delaying.VoidWithLog)
	gaedal.InitDelaying(delaying.MustRegisterFunc)
	facade4anybot.InitDelaying(delaying.MustRegisterFunc)
	facade4splitus.InitDelaying(delaying.MustRegisterFunc)
	emailing.InitDelaying(delaying.MustRegisterFunc)
	dtb_transfer.InitDelaying(delaying.MustRegisterFunc)
	botcmds4splitus.InitDelaying(delaying.MustRegisterFunc)
	reminders.InitDelaying(delaying.MustRegisterFunc)
	unsorted.InitDelaying(delaying.MustRegisterFunc)
}

func NotFoundSilent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNotFound)
}

func InitCronHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/cron/send-reminders", dtdal.HttpAppHost.HandleWithContext(reminders.CronSendReminders))
}

func InitTaskQueueHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/taskqueu/send-reminder", dtdal.HttpAppHost.HandleWithContext(reminders.SendReminderHandler)) // TODO: Remove obsolete!
	router.HandlerFunc(http.MethodPost, "/task-queue/send-reminder", dtdal.HttpAppHost.HandleWithContext(reminders.SendReminderHandler))
}

type TestTransferCounterparty struct {
	UserID   int64  `firestore:"userID"`
	UserName string `firestore:"userName"`
	Comment  string `firestore:"comment"`
}

type TestTransfer struct {
	From TestTransferCounterparty
	To   TestTransferCounterparty
}

//func testModelPointer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	c := appengine.NewContext(r)
//	testTransfer := TestTransfer{
//		From: TestTransferCounterparty{UserID: 1, UserName: "First"},
//		To:   TestTransferCounterparty{UserID: 2, UserName: "Second"},
//	}
//	key := datastore.NewKey(c, "TestTransfer", "", 1, nil)
//	if _, err := datastore.Put(c, key, &testTransfer); err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte(err.Error()))
//		return
//	}
//	var testTransfer2 TestTransfer
//	datastore.Get(c, key, &testTransfer2)
//	logus.Debugf(c, "testTransfer2: %v", testTransfer2)
//	logus.Debugf(c, "testTransfer2.From: %v", testTransfer2.From)
//	logus.Debugf(c, "testTransfer2.To: %v", testTransfer2.To)
//	testTransfer2.From.Comment = "Comment #1"
//	if _, err := datastore.Put(c, key, &testTransfer); err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte(err.Error()))
//		return
//	}
//	var testTransfer3 TestTransfer
//	datastore.Get(c, key, &testTransfer3)
//	logus.Debugf(c, "testTransfer2: %v", testTransfer3)
//	logus.Debugf(c, "testTransfer2.From: %v", testTransfer3.From)
//	logus.Debugf(c, "testTransfer2.To: %v", testTransfer3.To)
//}
