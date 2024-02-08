package dtb_general

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
)

const INVITES_SHOT_COMMAND = emoji.PRESENT_ICON

// This commands are required for main menu because of circular references
var _lendCommand = botsfw.Command{Code: "lend", Title: trans.COMMAND_TEXT_GAVE, Icon: emoji.GIVE_ICON}
var _borrowCommand = botsfw.Command{Code: "borrow", Title: trans.COMMAND_TEXT_GOT, Icon: emoji.TAKE_ICON}
var _returnCommand = botsfw.Command{Code: "return", Title: trans.COMMAND_TEXT_RETURN, Icon: emoji.RETURN_BACK_ICON}

func MainMenuKeyboardOnReceiptAck(whc botsfw.WebhookContext) *tgbotapi.ReplyKeyboardMarkup {
	return mainMenuTelegramKeyboard(whc, getMainMenuParams(whc, true))
}

type mainMenuParams struct {
	showBalanceAndHistory bool
	showReturn            bool
}

func getMainMenuParams(whc botsfw.WebhookContext, onReceiptAck bool) (params mainMenuParams) {
	//var (
	//	user *models.DebutsAppUserDataOBSOLETE
	//	isAppUser bool
	//)
	//c := whc.Context()
	//if userEntity, err := whc.AppUserData(); err != nil {
	//	log.Errorf(c, "Failed to get user: %v", err)
	//} else if user, isAppUser = userEntity.(*models.DebutsAppUserDataOBSOLETE); !isAppUser {
	//	log.Errorf(c, "Failed to case user to *models.DebutsAppUserDataOBSOLETE: %T", userEntity)
	//} else if onReceiptAck || !user.Balance().IsZero() {
	//	params.showReturn = true
	//}
	params.showBalanceAndHistory = onReceiptAck //|| (user != nil && user.CountOfTransfers > 0)
	return
}

func mainMenuTelegramKeyboard(whc botsfw.WebhookContext, params mainMenuParams) *tgbotapi.ReplyKeyboardMarkup {
	firstRow := []string{
		_lendCommand.DefaultTitle(whc),
		_borrowCommand.DefaultTitle(whc),
	}

	if params.showReturn {
		firstRow = append(firstRow, _returnCommand.DefaultTitle(whc))
	}

	buttonRows := make([][]string, 0, 3)
	buttonRows = append(buttonRows, firstRow)

	if params.showBalanceAndHistory {
		buttonRows = append(buttonRows, []string{
			whc.CommandText(trans.COMMAND_TEXT_BALANCE, emoji.BALANCE_ICON),
			//whc.CommandText(trans.COMMAND_TEXT_CONTACTS, emoji.MAN_AND_WOMAN),
			whc.CommandText(trans.COMMAND_TEXT_HISTORY, emoji.HISTORY_ICON),
		})
	}

	buttonRows = append(buttonRows, []string{
		//whc.CommandText(trans.COMMAND_TEXT_SETTING, emoji.SETTINGS_ICON),
		//whc.CommandText(trans.COMMAND_TEXT_HIGH_FIVE, emoji.BULB_ICON),
		//whc.CommandText(trans.COMMAND_TEXT_HELP, emoji.HELP_ICON),
		emoji.SETTINGS_ICON,
		emoji.MAN_AND_WOMAN,
		emoji.PUBLIC_LOUDSPEAKER,
		emoji.STAR_ICON,
		emoji.HELP_ICON,
	})

	return tgbotapi.NewReplyKeyboardUsingStrings(buttonRows)
}
func SetMainMenuKeyboard(whc botsfw.WebhookContext, m *botsfw.MessageFromBot) {
	params := getMainMenuParams(whc, true)
	switch whc.BotPlatform().ID() {
	case telegram.PlatformID:
		m.Keyboard = mainMenuTelegramKeyboard(whc, params)
	//case viber.PlatformID:
	//	panic("not implemented")
	//m.Keyboard = mainMenuViberKeyboard(whc, params)
	//case fbm.PlatformID:
	//	panic("not implemented")
	//if m.Text != "" {
	//	panic("FBM does not support message text and attachments in the same request.")
	//}
	//m.FbmAttachment = mainMenuFbmAttachment(whc, params)
	default:
		panic("Unsupported platform id=" + whc.BotPlatform().ID())
	}
}

//func mainMenuFbmAttachment(whc botsfw.WebhookContext, params mainMenuParams) *fbmbotapi.RequestAttachment {
//	attachment := &fbmbotapi.RequestAttachment{
//		Type: fbmbotapi.RequestAttachmentTypeTemplate,
//		Payload: fbmbotapi.NewListTemplate(
//			fbmbotapi.TopElementStyleCompact,
//			fbmbotapi.NewRequestElementWithDefaultAction(
//				"DebtsTracker.io",
//				"Tracks personal debts (auto-reminders to your debtors)",
//				fbmbotapi.NewDefaultActionWithWebURL(fbmbotapi.RequestWebURLAction{MessengerExtensions: true, URL: "https://debtstracker-dev1.appspot.com/app/?page=debts&lang=ru"}),
//				fbmbotapi.NewRequestWebURLButtonWithRatio(emoji.CURRENCY_EXCAHNGE_ICON+" Record new debt", "https://debtstracker-dev1.appspot.com/app/?page=new-debt&lang=ru", "full"),
//			),
//			fbmbotapi.NewRequestElementWithDefaultAction(
//				"Current balance",
//				"You owe $100",
//				fbmbotapi.NewDefaultActionWithWebURL(fbmbotapi.RequestWebURLAction{MessengerExtensions: true, URL: "https://debtstracker-dev1.appspot.com/app/?page=debts&lang=ru"}),
//				fbmbotapi.NewRequestWebURLButtonWithRatio(emoji.BALANCE_ICON+" Record return", "https://debtstracker-dev1.appspot.com/app/?page=return&lang=ru", "full"),
//			),
//			fbmbotapi.NewRequestElementWithDefaultAction(
//				"History",
//				"Last transfer: $100 to Jack Smith",
//				fbmbotapi.NewDefaultActionWithWebURL(fbmbotapi.RequestWebURLAction{MessengerExtensions: true, URL: "https://debtstracker-dev1.appspot.com/app/?page=history&lang=ru"}),
//				fbmbotapi.NewRequestWebURLButtonWithRatio(emoji.HISTORY_ICON+" View full history", "https://debtstracker-dev1.appspot.com/app/?page=history&lang=ru", "full"),
//			),
//			fbmbotapi.NewRequestElementWithDefaultAction(
//				"Settings",
//				"You can change language, notification preferences, etc.",
//				fbmbotapi.NewDefaultActionWithWebURL(fbmbotapi.RequestWebURLAction{MessengerExtensions: true, URL: "https://debtstracker-dev1.appspot.com/app/?page=debts&lang=ru"}),
//				fbmbotapi.NewRequestWebURLButtonWithRatio(emoji.SETTINGS_ICON+" Edit my preferences", "https://debtstracker-dev1.appspot.com/app/?page=settings&lang=ru", "full"),
//			),
//		),
//	}
//	log.Debugf(whc.Context(), "First element: %v", attachment.Payload.RequestAttachmentListTemplate.Elements[0])
//	return attachment
//}

const (
	UTM_CAMPAIGN_BOT_MAIN_MENU = "bot-main-menu"
)

//func mainMenuViberKeyboard(whc botsfw.WebhookContext, params mainMenuParams) *viberinterface.Keyboard {
//	var buttons []viberinterface.Button
//	lendingText := _lendCommand.DefaultTitle(whc)
//	borrowText := _borrowCommand.DefaultTitle(whc)
//	const (
//		maxColumns = 6
//		in3columns = maxColumns / 3
//		in2columns = maxColumns / 2
//	)
//	if params.showReturn {
//		returnText := _returnCommand.DefaultTitle(whc)
//		buttons = []viberinterface.Button{
//			{
//				Columns:    in3columns,
//				BgColor:    viberbots.ButtonBgColor,
//				Text:       lendingText,
//				ActionType: viberinterface.ActionTypeOpenUrl,
//				ActionBody: common.GetNewDebtPageUrl(whc, models.TransferDirectionUser2Counterparty, UTM_CAMPAIGN_BOT_MAIN_MENU),
//			},
//			{
//				Columns:    in3columns,
//				BgColor:    viberbots.ButtonBgColor,
//				Text:       borrowText,
//				ActionType: viberinterface.ActionTypeOpenUrl,
//				ActionBody: common.GetNewDebtPageUrl(whc, models.TransferDirectionCounterparty2User, UTM_CAMPAIGN_BOT_MAIN_MENU),
//			},
//			{Columns: in3columns, ActionBody: returnText, Text: returnText, BgColor: viberbots.ButtonBgColor},
//		}
//	} else {
//		buttons = []viberinterface.Button{
//			{Columns: in2columns, ActionBody: lendingText, Text: lendingText, BgColor: viberbots.ButtonBgColor},
//			{Columns: in2columns, ActionBody: borrowText, Text: borrowText, BgColor: viberbots.ButtonBgColor},
//		}
//	}
//	if params.showBalanceAndHistory {
//		userID := whc.AppUserID()
//		locale := whc.Locale()
//		balanceUrl := common.GetBalanceUrlForUser(userID, locale, whc.BotPlatform().ID(), whc.GetBotCode())
//		historyUrl := common.GetHistoryUrlForUser(userID, locale, whc.BotPlatform().ID(), whc.GetBotCode())
//		buttons = append(buttons, []viberinterface.Button{
//			{Columns: in2columns, ActionType: "open-url", ActionBody: balanceUrl, Text: whc.CommandText(trans.COMMAND_TEXT_BALANCE, emoji.BALANCE_ICON), BgColor: viberbots.ButtonBgColor},
//			{Columns: in2columns, ActionType: "open-url", ActionBody: historyUrl, Text: whc.CommandText(trans.COMMAND_TEXT_HISTORY, emoji.HISTORY_ICON), BgColor: viberbots.ButtonBgColor},
//		}...)
//	}
//	{ // Last row
//		settings := whc.CommandText(trans.COMMAND_TEXT_SETTING, emoji.SETTINGS_ICON)
//		rate := whc.CommandText(trans.COMMAND_TEXT_HIGH_FIVE, emoji.STAR_ICON)
//		help := whc.CommandText(trans.COMMAND_TEXT_HELP, emoji.HELP_ICON)
//		buttons = append(buttons, []viberinterface.Button{
//			{Columns: in3columns, ActionBody: settings, Text: settings, BgColor: viberbots.ButtonBgColor},
//			{Columns: in3columns, ActionBody: rate, Text: rate, BgColor: viberbots.ButtonBgColor},
//			{Columns: in3columns, ActionBody: help, Text: help, BgColor: viberbots.ButtonBgColor},
//		}...)
//	}
//
//	return viberinterface.NewKeyboard(viberbots.KeyboardBgColor, false, buttons...)
//}
