package dtb_fbm

//import (
//	"fmt"
//
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/bot/platforms/debtusfbmbots"
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/bot/profiles/debtusbot/cmd/dtb_transfer"
//	"github.com/sneat-co/debtusbot-translations/emoji"
//	"github.com/sneat-co/debtusbot-translations/trans"
//	"github.com/strongo/strongoapp"
//	"github.com/strongo/bots-api4debtus-fbm"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//)
//
//func aboutCard(whc botsfw.WebhookContext) fbmbotapi.RequestElement {
//	baseUrl := fbmAppBaseUrl(whc)
//	return fbmbotapi.NewRequestElementWithDefaultAction(
//		"More...",
//		"What can I do for you?",
//		newDefaultUrlAction(baseUrl, ""),
//		newUrlButton(emoji.HELP_ICON, "Help", baseUrl, "#help"),
//		newUrlButton(emoji.CONTACTS_ICON, "Contacts", baseUrl, "#contacts"),
//		newUrlButton(emoji.HISTORY_ICON, "History", baseUrl, "#history"),
//	)
//}
//
//func linkAccountsCard(whc botsfw.WebhookContext) fbmbotapi.RequestElement {
//	baseUrl := fbmAppBaseUrl(whc)
//	return fbmbotapi.NewRequestElementWithDefaultAction(
//		"link Accounts",
//		"to to...",
//		newDefaultUrlAction(baseUrl, ""),
//		newUrlButton(emoji.ROCKET_ICON, "Telegram", "t.me/", ""),
//	)
//}
//
//func mainMenuCard(whc botsfw.WebhookContext) fbmbotapi.RequestElement {
//	baseUrl := fbmAppBaseUrl(whc)
//	return fbmbotapi.NewRequestElementWithDefaultAction(
//		"Welcome",
//		"This is an app to split bills and track debt records.",
//		newDefaultUrlAction(baseUrl, ""),
//		newPostbackButton(emoji.MEMO_ICON, "Debts", FbmDebtsCommand.Code),
//		newPostbackButton(emoji.BILLS_ICON, "Bills", FbmBillsCommand.Code),
//		newPostbackButton(emoji.SETTINGS_ICON, "Settings", "fbm-settings"),
//	)
//}
//
////func mainMenuCard(whc botsfw.WebhookContext) fbm_api.RequestAttachmentPayload {
////	//baseUrl := fbmAppBaseUrl(whc)
////	return &fbm_api.NewListTemplate(
////		fbm_api.TopElementStyleCompact,
////		fbm_api.NewRequestElementWithDefaultAction(
////			emoji.MEMO_ICON + EM_SPACE + "Debts",
////			"Track your debts",
////			fbm_api.RequestDefaultAction{
////				ExtraType: fbm_api.
////				newPostbackButton(emoji.MEMO_ICON, "Debts", FbmDebtsCommand.Code)
////			},
////		),
////	)
////}
//
//func askLanguageCard(whc botsfw.WebhookContext) fbmbotapi.RequestAttachmentPayload {
//	fbmbotapi.NewButtonTemplate(
//		"",
//	)
//	requestElement := fbmbotapi.RequestElement{
//		Title:    whc.Translate(trans.MESSAGE_TEXT_HI),
//		Subtitle: "Please choose your language:",
//	}
//	for _, lang := range []i18n.Locale{i18n.LocaleEnUS, i18n.LocaleRuRu} {
//		requestElement.Buttons = append(requestElement.Buttons, newPostbackButton(lang.FlagIcon, lang.NativeTitle, "fbm-set-lang?code5="+lang.Code5))
//	}
//	requestElement.Buttons = append(requestElement.Buttons, newUrlButton("", "More...", fbmAppBaseUrl(whc), "#set-locale"))
//	return fbmbotapi.NewGenericTemplate(requestElement)
//}
//
//func welcomeCard(whc botsfw.WebhookContext) fbmbotapi.RequestElement {
//	baseUrl := fbmAppBaseUrl(whc)
//	return fbmbotapi.NewRequestElementWithDefaultAction(
//		"Welcome!",
//		"Have you ever used DebtsTracker.io app/bot outside of FB Messenger before?",
//		newDefaultUrlAction(baseUrl, ""),
//		newPostbackButton(emoji.MEMO_ICON, "Have not used", "fbm-debts"),
//		newPostbackButton(emoji.ROCKET_ICON, "Used @ https://debtstracker.io/", "fbm-bills"),
//		newPostbackButton(emoji.ROBOT_ICON, "Used @ Telegram", "fbm-settings"),
//	)
//}
//
//func debtsCard(whc botsfw.WebhookContext) fbmbotapi.RequestElement {
//	baseUrl := fbmAppBaseUrl(whc)
//	requestElement := fbmbotapi.NewRequestElementWithDefaultAction(
//		"Debts",
//		"Tracks personal debts (auto-reminders to your debtors)",
//		newDefaultUrlAction(baseUrl, "#debts"),
//		newPostbackButton(emoji.MEMO_ICON, whc.Translate("New record"), "new-debt-or-return"),
//		newPostbackButton(emoji.CLIPBOARD_ICON, whc.Translate(trans.COMMAND_TEXT_BALANCE), dtb_transfer.BALANCE_COMMAND),
//		newPostbackButton(emoji.HISTORY_ICON, whc.Translate(trans.COMMAND_TEXT_HISTORY), dtb_transfer.HISTORY_COMMAND),
//	)
//	//requestElement.ImageURL = ""
//	return requestElement
//}
//
//func billsCard(whc botsfw.WebhookContext) fbmbotapi.RequestElement {
//	baseUrl := fbmAppBaseUrl(whc)
//	return fbmbotapi.NewRequestElementWithDefaultAction(
//		"Bills",
//		"Split regular or single bills and get paid back",
//		newDefaultUrlAction(baseUrl, "#bills"),
//		newUrlButton(emoji.DIVIDE_ICON, "Split bill", baseUrl, "#split-bill"),
//		newUrlButton(emoji.BILLS_ICON, "Outstanding bills", baseUrl, "#bills"),
//		newUrlButton(emoji.CALENDAR_ICON, "Recurring bills", baseUrl, "#bills"),
//	)
//}
//
//func settingsCard(whc botsfw.WebhookContext) fbmbotapi.RequestElement {
//	baseUrl := fbmAppBaseUrl(whc)
//	return fbmbotapi.NewRequestElementWithDefaultAction(
//		"Settings",
//		"Adjust settings",
//		newDefaultUrlAction(baseUrl, "#bills"),
//		newUrlButton(emoji.BILLS_ICON, "Bills", baseUrl, "#bills"),
//		newUrlButton(emoji.MEMO_ICON, "Split bill", baseUrl, "#split-bill"),
//	)
//}
//
//func fbmAppBaseUrl(whc botsfw.WebhookContext) string {
//	fbApp, host, err := debtusfbmbots.GetFbAppAndHost(whc.Request())
//	if err != nil {
//		panic(err)
//	}
//	return fmt.Sprintf("https://%v/app/#fbm%v", host, fbApp.AppId)
//}
//
//func newDefaultUrlAction(baseUrl, hash string) fbmbotapi.RequestDefaultAction {
//	return fbmbotapi.NewDefaultActionWithWebURL(
//		fbmbotapi.RequestWebURLAction{
//			MessengerExtensions: true,
//			URL:                 baseUrl + hash,
//		},
//	)
//}
//
//func newUrlButton(icon, title, baseUrl, hash string) fbmbotapi.RequestButton {
//	if icon != "" {
//		title = icon + EM_SPACE + title
//	}
//	button := fbmbotapi.NewRequestWebURLButtonWithRatio(
//		title,
//		baseUrl+hash,
//		fbmbotapi.WebviewHeightRatioFull,
//	)
//	button.MessengerExtensions = true
//	return button
//}
//
//func newPostbackButton(icon, title, payload string) fbmbotapi.RequestButton {
//	if icon != "" {
//		title = icon + EM_SPACE + title
//	}
//	button := fbmbotapi.NewRequestPostbackButton(
//		title,
//		payload,
//	)
//	return button
//}
