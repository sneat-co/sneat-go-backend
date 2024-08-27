package shared_all

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/platforms/debtustgbots/tgsharedcommands"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"net/url"
	"strings"
)

func StartBotLink(botID, command string, params ...string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "https://t.me/%v?start=%v", botID, command)
	for _, p := range params {
		buf.WriteString("__")
		buf.WriteString(p)
	}
	return buf.String()
}

func createStartCommand(botParams BotParams) botsfw.Command {
	return botsfw.Command{
		Code:       "start",
		Commands:   []string{"/start"},
		InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputInlineQuery},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return startCommandAction(whc, botParams)
		},
	}
}

func startCommandAction(whc botsfw.WebhookContext, botParams BotParams) (m botsfw.MessageFromBot, err error) {
	whc.LogRequest()
	ctx := whc.Context()
	text := whc.Input().(botsfw.WebhookTextMessage).Text()
	logus.Debugf(ctx, "createStartCommand.Action() => text: "+text)

	startParam, startParams := tgsharedcommands.ParseStartCommand(whc)

	var isInGroup bool
	if isInGroup, err = whc.IsInGroup(); err != nil {
		return
	} else if isInGroup {
		return botParams.StartInGroupAction(whc)
	} else {
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
			facade4debtus.Referer.AddTelegramReferrer(ctx, whc.AppUserID(), strings.TrimPrefix(startParam, "refbytguser-"), whc.GetBotCode())
		}
		return startInBotAction(whc, startParams, botParams)
	}
}
func startLoginGac(whc botsfw.WebhookContext, loginID int) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	var loginPin models4auth.LoginPin
	if loginPin, err = facade4debtus.AuthFacade.AssignPinCode(ctx, loginID, whc.AppUserID()); err != nil {
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
	return user, dal4userus.GetUser(ctx, nil, user)
}

var LangKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		{
			Text:         i18n.LocaleEnUS.TitleWithIcon(),
			CallbackData: onStartCallbackCommandCode + "?lang=" + i18n.LocaleCodeEnUS,
		},
		{
			Text:         i18n.LocaleRuRu.TitleWithIcon(),
			CallbackData: onStartCallbackCommandCode + "?lang=" + i18n.LocalCodeRuRu,
		},
	},
)

const onStartCallbackCommandCode = "on-start-callback"

func onStartCallbackCommand(params BotParams) botsfw.Command {
	return botsfw.NewCallbackCommand(onStartCallbackCommandCode,
		func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			lang := callbackUrl.Query().Get("lang")
			mode := "onboarding" // TODO: should we set mode?
			return setPreferredLanguageAction(whc, lang, mode, params)
			//ctx := whc.Context()
			//if lang != "" {
			//	logus.Debugf(ctx, "Locale: "+lang)
			//
			//	whc.ChatData().SetPreferredLanguage(lang)
			//
			//	appUserID := whc.AppUserID()
			//	if err = whc.SetLocale(lang); err != nil {
			//		return
			//	}
			//	if appUserID != "" {
			//		userCtx := facade.NewUserContext(whc.AppUserID())
			//		if err = dal4userus.RunUserWorker(c, userCtx,
			//			func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) error {
			//				if params.UserUpdates, err = params.UserEntry.Data.SetPreferredLocale(lang); err != nil {
			//					return err
			//				}
			//				return nil
			//			}); err != nil {
			//			return
			//		}
			//	}
			//	m.Text = fmt.Sprintf("Language set to %s", lang)
			//}

			//if whc.IsInGroup() {
			//	var group models.GroupEntry
			//	if group, err = GetGroup(whc, callbackUrl); err != nil {
			//		return
			//	}
			//	return onStartCallbackInGroup(whc, group, params)
			//} else {
			//	return onStartCallbackInBot(whc, params)
			//}
			//return
		},
	)
}
