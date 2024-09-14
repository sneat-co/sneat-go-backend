package cmds4anybot

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/facade4anybot"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/platforms/debtustgbots/tgsharedcommands"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"net/url"
	"strings"
)

func StartBotLink(botID, command string, params ...string) string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "https://t.me/%v?start=%v", botID, command)
	for _, p := range params {
		buf.WriteString("__")
		buf.WriteString(p)
	}
	return buf.String()
}

const StartCommandCode = "start"

func createStartCommand(
	startInBotAction StartInBotActionFunc,
	startInGroupAction botsfw.CommandAction,
	getWelcomeMessageText WelcomeMessageProvider,
	setMainMenu SetMainMenuFunc,

) botsfw.Command {
	return botsfw.Command{
		Code:     StartCommandCode,
		Commands: []string{"/start"},
		InputTypes: []botinput.WebhookInputType{
			botinput.WebhookInputText,
			botinput.WebhookInputReferral,            // FBM
			botinput.WebhookInputConversationStarted, // Viber
		},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return sharedStartCommandAction(whc, startInBotAction, startInGroupAction, getWelcomeMessageText)
		},
		CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			return sharedStartCommandCallbackAction(whc, callbackUrl, getWelcomeMessageText, setMainMenu, startInBotAction)
		},
	}
}

func sharedStartCommandCallbackAction(
	whc botsfw.WebhookContext,
	callbackUrl *url.URL,
	getWelcomeMessageText WelcomeMessageProvider,
	setMainMenu SetMainMenuFunc,
	startInBotAction StartInBotActionFunc,
) (m botsfw.MessageFromBot, err error) {
	q := callbackUrl.Query()
	if localeCode := q.Get("locale"); localeCode != "" {
		// TODO: Need refactoring - duplicate welcome message rendering in setPreferredLocaleAction() & runBotSpecificStartCommand()
		if m, err = runBotSpecificStartCommand(whc, startInBotAction, []string{}, getWelcomeMessageText); err != nil {
			return
		}
		if m, err = setPreferredLocaleAction(whc, localeCode, setPreferredLocaleModeStart, setMainMenu, getWelcomeMessageText); err != nil {
			return m, fmt.Errorf("failed to setPreferredLocaleAction(): %w", err)
		}
		return
	}
	m.Text = fmt.Sprintf("Unknown callback parameters: %s", callbackUrl)
	m.IsEdit = false
	return
}

func sharedStartCommandAction(
	whc botsfw.WebhookContext,
	startInBotAction StartInBotActionFunc,
	startInGroupAction botsfw.CommandAction,
	getWelcomeMessage WelcomeMessageProvider,
) (
	m botsfw.MessageFromBot, err error,
) {
	whc.Input().LogRequest()
	ctx := whc.Context()
	text := whc.Input().(botinput.WebhookTextMessage).Text()
	logus.Debugf(ctx, "createStartCommand.Action() => text: "+text)

	startParam, startParams := tgsharedcommands.ParseStartCommand(whc)

	var isInGroup bool
	if isInGroup, err = whc.IsInGroup(); err != nil {
		return
	} else if isInGroup {
		return startInGroupAction(whc)
	}
	chatEntity := whc.ChatData()
	chatEntity.SetAwaitingReplyTo("")

	switch {
	case startParam == "help_inline":
		return startInlineHelp(whc)
	case strings.HasPrefix(startParam, "login-"):
		loginID, err := common4debtus.DecodeIntID(startParam[len("login-"):])
		if err != nil {
			return m, err
		}
		return startLoginGac(whc, loginID)
		//case strings.HasPrefix(textToMatchNoStart, JOIN_BILL_COMMAND):
		//	return JoinBillCommand.Action(whc)
	case strings.HasPrefix(startParam, "refbytguser-") && startParam != "refbytguser-YOUR_CHANNEL":
		facade4anybot.Referer.AddTelegramReferrer(ctx, whc.AppUserID(), strings.TrimPrefix(startParam, "refbytguser-"), whc.GetBotCode())
	}
	//if m.Text, err = getWelcomeMessage(whc); err != nil {
	//	return
	//} else if m.Text != "" {
	//	responder := whc.Responder()
	//	if _, err = responder.SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
	//		return
	//	}
	//}
	if m.Text, err = getWelcomeMessage(whc); err != nil {
		return
	}
	var user dbo4userus.UserEntry
	if user, err = GetUser(whc); err != nil {
		return
	}
	if user.Data.PreferredLocale == "" {
		var localesMsg botsfw.MessageFromBot
		if localesMsg, err = onStartAskLocaleAction(whc, nil, getWelcomeMessage); err != nil {
			return
		}
		if localesMsg.Text = strings.TrimSpace(localesMsg.Text); localesMsg.Text != "" {
			m.Text += "\n" + localesMsg.Text
			m.Keyboard = localesMsg.Keyboard
			m.Format = botsfw.MessageFormatHTML
		}
		return
	}
	if m, err = runBotSpecificStartCommand(whc, startInBotAction, startParams, getWelcomeMessage); err != nil {
		return
	}
	return
}

func runBotSpecificStartCommand(whc botsfw.WebhookContext, startInBotAction StartInBotActionFunc, startParams []string, getWelcomeMessage WelcomeMessageProvider) (m botsfw.MessageFromBot, err error) {
	if m, err = startInBotAction(whc, startParams); err != nil {
		return
	}
	responder := whc.Responder()
	ctx := whc.Context()
	if _, err = responder.SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		return
	}

	if m.Text, err = getWelcomeMessage(whc); err != nil {
		return
	}
	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🛠 Settings",
				CallbackData: SettingsCommandCode,
			},
			{
				Text:         whc.CommandText(trans.COMMAND_TEXT_LANGUAGE, emoji.EARTH_ICON),
				CallbackData: SettingsLocaleListCallbackPath,
			},
		},
	)

	return
}

func startLoginGac(whc botsfw.WebhookContext, loginID int) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	var loginPin models4auth.LoginPin
	if loginPin, err = facade4anybot.AuthFacade.AssignPinCode(ctx, loginID, whc.AppUserID()); err != nil {
		return
	}
	return whc.NewMessageByCode(trans.MESSAGE_TEXT_LOGIN_CODE, models4auth.LoginCodeToString(loginPin.Data.Code)), nil
}

func startInlineHelp(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m = whc.NewMessage("<b>Help: How to use this bot in chats</b>\n\nExplain here how to use bot's inline mode.")
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 1", URL: "https://debtstracker.io/#btn=1"}},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 2", URL: "https://debtstracker.io/#btn=2"}},
		//[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonSwitch("Back to chat 1", "1")},
		//[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonSwitch("Back to chat 2", "2")},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 3", CallbackData: "help-3"}},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 4", CallbackData: "help-4"}},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 5", CallbackData: "help-5"}},
	)
	return m, err
}

func GetUser(whc botsfw.WebhookContext) (user dbo4userus.UserEntry, err error) { // TODO: Make library and use across app
	appUserID := whc.AppUserID()
	if appUserID == "" {
		return user, fmt.Errorf("%w: app user ID is empty", dal.ErrRecordNotFound)
	}
	user = dbo4userus.NewUserEntry(appUserID)
	ctx := whc.Context()
	tx := whc.Tx()
	return user, dal4userus.GetUser(ctx, tx, user)
}

//var LangKeyboard = tgbotapi.NewInlineKeyboardMarkup(
//	[]tgbotapi.InlineKeyboardButton{
//		{
//			Text:         i18n.LocaleEnUS.TitleWithIcon(),
//			CallbackData: onStartCallbackCommandCode + "?lang=" + i18n.LocaleCodeEnUS,
//		},
//		{
//			Text:         i18n.LocaleRuRu.TitleWithIcon(),
//			CallbackData: onStartCallbackCommandCode + "?lang=" + i18n.LocalCodeRuRu,
//		},
//	},
//)
