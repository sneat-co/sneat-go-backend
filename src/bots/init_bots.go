package bots

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-fw/botswebhook"
	"github.com/julienschmidt/httprouter"
	"github.com/strongo/i18n"
	"github.com/strongo/log"
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

func InitializeBots(botHost botsfw.BotHost, httpRouter *httprouter.Router) {
	if botHost == nil {
		panic("botHost == nil")
	}
	if httpRouter == nil {
		panic("httpRouter = nil")
	}
	botsfw.SetLogger(log.StrongoLogger{})

	driver := botswebhook.NewBotDriver( // Orchestrate requests to appropriate handlers
		botswebhook.AnalyticsSettings{GaTrackingID: ""}, // TODO: Refactor to list of analytics providers
		sneatAppBotContext{},                            // Holds User entity kind name, translator, etc.
		botHost,                                         // Defines how to create context.Context, HttpClient, DB, etc...
		"Please report any issues to @trakhimenok",      // Is it wrong place? Router has similar.
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

	//dataAccess := dalgo4botsfw.NewDataAccess(telegram.PlatformID, getDb, recordsMaker)
	driver.RegisterWebhookHandlers(botsHttpRouter{httpRouter}, "/bot",
		telegram.NewTelegramWebhookHandler(
			//dataAccess,
			telegramBotsWithRouter, // Maps of bots by code, language, token, etc...
			newTranslator,          // Creates translator that gets a context.Context (for logging purpose)
			recordsMaker,
			func(data botsfwmodels.AppUserData, sender botsfw.WebhookSender) error { // TODO: implement?
				return nil
			},
		),
		//viber.NewViberWebhookHandler(...),
		//fbm.NewFbmWebhookHandler(...),
	)
}

func newTranslator(c context.Context) i18n.Translator {
	return i18n.NewMapTranslator(c, nil)
}

func telegramBotsWithRouter(context.Context) botsfw.SettingsBy {
	//env := app.Host.GetEnvironment(c, nil) // TODO: request is not being passed, needs to be fixed
	return telegramBots(botsfw.EnvLocal)
}
