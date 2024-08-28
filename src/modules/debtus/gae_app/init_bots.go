package gaeapp

//func newTranslator(ctx context.Context) i18n.Translator {
//	return i18n.NewMapTranslator(ctx, trans.TRANS)
//}
//
//type botsHttpRouter struct {
//	r *httprouter.Router
//}
//
//func (v botsHttpRouter) Handle(method, path string, handle http.HandlerFunc) {
//	v.r.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
//		handle(writer, request)
//	})
//}

//func InitBots(httpRouter *httprouter.Router, botHost botsfw.BotHost, appContext botsfw.BotAppContext) {
//
//	driver := botswebhook.NewBotDriver( // Orchestrate requests to appropriate handlers
//		botswebhook.AnalyticsSettings{GaTrackingID: common4debtus.GA_TRACKING_ID}, // TODO: Refactor to list of analytics providers
//		appContext, // Holds User entity kind name, translator, etc.
//		botHost,    // Defines how to create context.Context, HttpClient, DB, etc...
//		"Please report any issues to @DebtsTrackerGroup", // Is it wrong place? Router has similar.
//	)
//
//	makeAppUserDto := func(botID string) (appUser botsfwmodels.AppUserData, err error) {
//		return nil, fmt.Errorf("%w: makeAppUserDto() is not implemented", botsfw.ErrNotImplemented)
//	}
//	var recordsMaker = botsfwmodels.NewBotRecordsMaker(
//		"*",
//		makeAppUserDto,
//		telegram.BaseTgUserDtoMaker,
//		telegram.BaseTgChatDtoMaker,
//	)
//
//	//var getDb dalgo4botsfw.DbProvider = func(ctx context.Context, botID string) (dal.DB, error) {
//	//	return nil, errors.New("not implemented")
//	//	//fsClient, err := firestore.NewClient(ctx, "demo-local-sneat-app")
//	//	//if err != nil {
//	//	//	return nil, err
//	//	//}
//	//	//return dalgo2firestore.NewDatabase("sneat", fsClient), nil
//	//}
//
//	//dataAccess := dalgo4botsfw.NewDataAccess(telegram.PlatformID, getDb, recordsMaker)
//
//	driver.RegisterWebhookHandlers(botsHttpRouter{httpRouter}, "/bot",
//		//telegram.NewTelegramWebhookHandler(
//		//	telegramBotsWithRouter, // Maps of botscore by code, language, token, etc...
//		//	newTranslator,          // Creates translator that gets a context.Context (for logging purpose)
//		//),
//		telegram.NewTelegramWebhookHandler(
//			//dataAccess,
//			telegramBotsWithRouter, // Maps of botscore by code, language, token, etc...
//			newTranslator,          // Creates translator that gets a context.Context (for logging purpose)
//			recordsMaker,
//			func(data botsfwmodels.AppUserData, sender botsfw.WebhookSender) error { // TODO: implement?
//				return nil
//			},
//		),
//		//viber.NewViberWebhookHandler(
//		//	debtusviberbots.Bots,
//		//	newTranslator,
//		//),
//		//fbm.NewFbmWebhookHandler(
//		//	debtusfbmbots.Bots,
//		//	newTranslator,
//		//),
//	)
//}

//func telegramBotsWithRouter(ctx context.Context) botsfw.SettingsBy {
//	return debtustgbots.Bots(dtdal.HttpAppHost.GetEnvironment(ctx, nil))
//}
