package gaeapp

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-fw/botswebhook"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/platforms/tgbots"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/collectus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/splitus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/strongo/i18n"
	"net/http"
)

func newTranslator(c context.Context) i18n.Translator {
	return i18n.NewMapTranslator(c, trans.TRANS)
}

type botsHttpRouter struct {
	r *httprouter.Router
}

func (v botsHttpRouter) Handle(method, path string, handle http.HandlerFunc) {
	v.r.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		handle(writer, request)
	})
}

func InitBots(httpRouter *httprouter.Router, botHost botsfw.BotHost, appContext botsfw.BotAppContext) {

	driver := botswebhook.NewBotDriver( // Orchestrate requests to appropriate handlers
		botswebhook.AnalyticsSettings{GaTrackingID: common.GA_TRACKING_ID}, // TODO: Refactor to list of analytics providers
		appContext, // Holds User entity kind name, translator, etc.
		botHost,    // Defines how to create context.Context, HttpClient, DB, etc...
		"Please report any issues to @DebtsTrackerGroup", // Is it wrong place? Router has similar.
	)

	makeAppUserDto := func(botID string) (appUser botsfwmodels.AppUserData, err error) {
		return nil, fmt.Errorf("%w: makeAppUserDto() is not implemented", botsfw.ErrNotImplemented)
	}
	var recordsMaker = botsfwmodels.NewBotRecordsMaker(
		"*",
		makeAppUserDto,
		telegram.BaseTgUserDtoMaker,
		telegram.BaseTgChatDtoMaker,
	)

	//var getDb dalgo4botsfw.DbProvider = func(c context.Context, botID string) (dal.DB, error) {
	//	return nil, errors.New("not implemented")
	//	//fsClient, err := firestore.NewClient(c, "demo-local-sneat-app")
	//	//if err != nil {
	//	//	return nil, err
	//	//}
	//	//return dalgo2firestore.NewDatabase("sneat", fsClient), nil
	//}

	//dataAccess := dalgo4botsfw.NewDataAccess(telegram.PlatformID, getDb, recordsMaker)

	driver.RegisterWebhookHandlers(botsHttpRouter{httpRouter}, "/bot",
		//telegram.NewTelegramWebhookHandler(
		//	telegramBotsWithRouter, // Maps of bots by code, language, token, etc...
		//	newTranslator,          // Creates translator that gets a context.Context (for logging purpose)
		//),
		telegram.NewTelegramWebhookHandler(
			//dataAccess,
			telegramBotsWithRouter, // Maps of bots by code, language, token, etc...
			newTranslator,          // Creates translator that gets a context.Context (for logging purpose)
			recordsMaker,
			func(data botsfwmodels.AppUserData, sender botsfw.WebhookSender) error { // TODO: implement?
				return nil
			},
		),
		//viber.NewViberWebhookHandler(
		//	viberbots.Bots,
		//	newTranslator,
		//),
		//fbm.NewFbmWebhookHandler(
		//	fbmbots.Bots,
		//	newTranslator,
		//),
	)
}

func telegramBotsWithRouter(c context.Context) botsfw.SettingsBy {
	return tgbots.Bots(dtdal.HttpAppHost.GetEnvironment(c, nil), func(profile string) botsfw.WebhooksRouter {
		switch profile {
		case bot.ProfileDebtus:
			return debtus.Router
		case bot.ProfileSplitus:
			return splitus.Router
		case bot.ProfileCollectus:
			return collectus.Router
		default:
			panic("Unknown bot profile: " + profile)
		}
	})
}
