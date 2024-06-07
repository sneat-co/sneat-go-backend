package shared_all

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/platforms/tgbots"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/i18n"
	"github.com/strongo/log"
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
			whc.LogRequest()
			c := whc.Context()
			text := whc.Input().(botsfw.WebhookTextMessage).Text()
			log.Debugf(c, "createStartCommand.Action() => text: "+text)

			startParam, startParams := tgbots.ParseStartCommand(whc)

			if whc.IsInGroup() {
				return botParams.StartInGroupAction(whc)
			} else {
				chatEntity := whc.ChatData()
				chatEntity.SetAwaitingReplyTo("")

				switch {
				case startParam == "help_inline":
					return startInlineHelp(whc)
				case strings.HasPrefix(startParam, "login-"):
					loginID, err := common.DecodeIntID(startParam[len("login-"):])
					if err != nil {
						return m, err
					}
					return startLoginGac(whc, loginID)
					//case strings.HasPrefix(textToMatchNoStart, JOIN_BILL_COMMAND):
					//	return JoinBillCommand.Action(whc)
				case strings.HasPrefix(startParam, "refbytguser-") && startParam != "refbytguser-YOUR_CHANNEL":
					facade.Referer.AddTelegramReferrer(c, whc.AppUserID(), strings.TrimPrefix(startParam, "refbytguser-"), whc.GetBotCode())
				}
				return startInBotAction(whc, startParams, botParams)
			}
		},
	}
}

func startLoginGac(whc botsfw.WebhookContext, loginID int) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	var loginPin models.LoginPin
	if loginPin, err = facade.AuthFacade.AssignPinCode(c, loginID, whc.AppUserID()); err != nil {
		return
	}
	return whc.NewMessageByCode(trans.MESSAGE_TEXT_LOGIN_CODE, models.LoginCodeToString(loginPin.Data.Code)), nil
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

func GetUser(whc botsfw.WebhookContext) (userEntity *models.DebutsAppUserDataOBSOLETE, err error) { // TODO: Make library and use across app
	panic("not implemented: obsolete")
	//var botAppUser botsfwmodels.AppUserData
	//if botAppUser, err = whc.AppUserData(); err != nil {
	//	return
	//}
	//userEntity = botAppUser.(*models.DebutsAppUserDataOBSOLETE)
	//return
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
			c := whc.Context()
			log.Debugf(c, "Locale: "+lang)

			whc.ChatData().SetPreferredLanguage(lang)

			var db dal.DB
			if db, err = facade.GetDatabase(c); err != nil {
				return
			}

			if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
				user, err := facade.User.GetUserByID(c, tx, whc.AppUserID())
				if err != nil {
					return err
				}
				if err = user.Data.SetPreferredLocale(lang); err != nil {
					return err
				}
				if err = facade.User.SaveUser(c, tx, user); err != nil {
					return err
				}
				return nil
			}, nil); err != nil {
				return
			}

			if err = whc.SetLocale(lang); err != nil {
				return
			}

			//if whc.IsInGroup() {
			//	var group models.GroupEntry
			//	if group, err = GetGroup(whc, callbackUrl); err != nil {
			//		return
			//	}
			//	return onStartCallbackInGroup(whc, group, params)
			//} else {
			//	return onStartCallbackInBot(whc, params)
			//}
			return
		},
	)
}
