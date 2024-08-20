package shared_all

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/debtusbotconst"
	"net/url"
	"strings"
)

const HELP_COMMAND = "help"

func createHelpRootCommand(params BotParams) botsfw.Command {
	return botsfw.Command{
		Code:     HELP_COMMAND,
		Commands: []string{"/help", emoji.HELP_ICON},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			switch whc.GetBotSettings().Profile.ID() {
			case debtusbotconst.DebtusBotProfileID:
				return params.HelpCommandAction(whc)
			}
			return helpRootAction(whc, false)
		},
		CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			_ = whc.ChatData()
			q := callbackUrl.Query().Get("q")
			switch q {
			case "":
				m, err = helpRootAction(whc, true)
			case trans.HELP_HOW_TO_CREATE_BILL_Q:
				m, err = helpHowToCreateNewBill(whc)
			}
			m.Format = botsfw.MessageFormatHTML
			m.IsEdit = true
			return
		},
	}
}

func helpRootAction(whc botsfw.WebhookContext, isCallback bool) (m botsfw.MessageFromBot, err error) {
	m.Text = whc.Translate(trans.MESSAGE_TEXT_HELP_ROOT, strings.Replace(whc.GetBotCode(), "Bot", "Group", 1))
	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{{
			Text:         whc.Translate(trans.HELP_HOW_TO_CREATE_BILL_Q),
			CallbackData: "help?q=" + url.QueryEscape(trans.HELP_HOW_TO_CREATE_BILL_Q),
		}},
	)
	if isCallback {
		m.IsEdit = true
	}
	return
}

func helpHowToCreateNewBill(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	var buffer bytes.Buffer
	if err = common4debtus.HtmlTemplates.RenderTemplate(whc.Context(), &buffer, whc, trans.HELP_HOW_TO_CREATE_BILL_A, struct{ BotCode string }{whc.GetBotCode()}); err != nil {
		return
	}
	m.Text = fmt.Sprintf("<b>%v</b>", whc.Translate(trans.HELP_HOW_TO_CREATE_BILL_Q)) +
		"\n\n" + buffer.String()
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: emoji.CONTACTS_ICON + " Split bill in Telegram Group",
				URL:  "https://t.me/%v?startgroup=new-bill",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
				emoji.ROCKET_ICON+" Split bill with Telegram user(s)",
				"",
			),
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         emoji.CLIPBOARD_ICON + "New bill manually",
				CallbackData: "new-bill-manually",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: emoji.RETURN_BACK_ICON + " " + whc.Translate(trans.MESSAGE_TEXT_HELP_BACK_TO_ROOT),
				CallbackData: "help",
			},
		},
	)
	return
}
