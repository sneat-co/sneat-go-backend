package botscore

import (
	"context"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-fw/botswebhook"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/sneatbot"
	"github.com/strongo/i18n"
	"net/http"
)

type botsHttpRouter struct {
	r *httprouter.Router
}

func (v botsHttpRouter) Handle(method, path string, handle http.HandlerFunc) {
	v.r.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		handle(writer, request)
	})
}

//func sneatAppTgUserDboMaker(botID string) (botsfwmodels.PlatformUserData, error) {
//	panic("deprecated: use profile.NewChatData()")
//}
//
//func sneatAppTgChatDboMaker(botID string) (botChat botsfwmodels.BotChatData, err error) {
//	panic("deprecated: use profile.NewChatData()")
//}

func InitializeBots(botHost botsfw.BotHost, httpRouter *httprouter.Router) {
	if botHost == nil {
		panic("botHost == nil")
	}
	if httpRouter == nil {
		panic("httpRouter = nil")
	}
	//botsfw.SetLogger(null)

	driver := botswebhook.NewBotDriver( // Orchestrate requests to appropriate handlers
		botswebhook.AnalyticsSettings{GaTrackingID: ""}, // TODO: Refactor to list of analytics providers
		//sneatbot.NewSneatAppContextForBotsFW(),                // Holds User entity kind name, translator, etc.
		botHost, // Defines how to create context.Context, HttpClient, DB, etc...
		"Please report any issues to @trakhimenok", // TODO: Is it a wrong place? Router has similar.
	)

	//makeAppUserDto := func(botID string) (appUser botsfwmodels.AppUserData, err error) {
	//	return nil, fmt.Errorf("%w: makeAppUserDto() is not implemented", botsfw.ErrNotImplemented)
	//}
	//var recordsMaker = botsfwmodels.NewBotRecordsMaker(
	//	"*",
	//	makeAppUserDto,
	//	sneatAppTgUserDboMaker, //telegram.BaseTgUserDtoMaker,
	//	sneatAppTgChatDboMaker, //telegram.BaseTgChatDtoMaker
	//)

	botContextProvider := botsfw.NewBotContextProvider(botHost, sneatbot.NewSneatAppContextForBotsFW(), telegramBotsWithRouter)

	var tgWebhookHandler botsfw.WebhookHandler = telegram.NewTelegramWebhookHandler(
		botContextProvider,
		newTranslator, // Creates translator that gets a context.Context (for logging purpose)
		func(data botsfwmodels.AppUserData, sender botinput.WebhookSender) error { // TODO: implement?
			return nil
		},
	)
	//dataAccess := dalgo4botsfw.NewDataAccess(telegram.PlatformID, getDb, recordsMaker)
	driver.RegisterWebhookHandlers(botsHttpRouter{httpRouter}, "/bot",
		tgWebhookHandler,
		//viber.NewViberWebhookHandler(...),
		//fbm.NewFbmWebhookHandler(...),
	)
}

func newTranslator(ctx context.Context) i18n.Translator {
	return i18n.NewMapTranslator(ctx, trans.TRANS)
}

func telegramBotsWithRouter(context.Context) botsfw.BotSettingsBy {
	//env := app.Host.GetEnvironment(c, nil) // TODO: request is not being passed, needs to be fixed
	return telegramBots(botsfw.EnvLocal)
}
