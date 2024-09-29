package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
)

const HelpCommandCode = "help"

func createHelpRootCommand(helpCommandAction botsfw.CommandAction, helpCallbackAction botsfw.CallbackAction) botsfw.Command {
	return botsfw.Command{
		Code:           HelpCommandCode,
		Commands:       []string{"/help", emoji.HELP_ICON},
		InputTypes:     []botinput.WebhookInputType{botinput.WebhookInputText},
		Action:         helpCommandAction,
		CallbackAction: helpCallbackAction,
	}
}

//func helpRootAction(whc botsfw.WebhookContext, isCallback bool) (m botsfw.MessageFromBot, err error) {
//	m.Text = whc.Translate(trans.MESSAGE_TEXT_HELP_ROOT, strings.Replace(whc.GetBotCode(), "Bot", "Group", 1))
//	m.Format = botsfw.MessageFormatHTML
//	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
//		[]tgbotapi.InlineKeyboardButton{{
//			Text:         whc.Translate(trans.HELP_HOW_TO_CREATE_BILL_Q),
//			CallbackData: "help?q=" + url.QueryEscape(trans.HELP_HOW_TO_CREATE_BILL_Q),
//		}},
//	)
//	if isCallback {
//		m.IsEdit = true
//	}
//	return
//}
//
//func helpHowToCreateNewBill(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
//	var buffer bytes.Buffer
//	if err = common4all.HtmlTemplates.RenderTemplate(whc.Context(), &buffer, whc, trans.HELP_HOW_TO_CREATE_BILL_A, struct{ BotCode string }{whc.GetBotCode()}); err != nil {
//		return
//	}
//	m.Text = fmt.Sprintf("<b>%v</b>", whc.Translate(trans.HELP_HOW_TO_CREATE_BILL_Q)) +
//		"\n\n" + buffer.String()
//	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
//		[]tgbotapi.InlineKeyboardButton{
//			{
//				Text: emoji.CONTACTS_ICON + " Split bill in Telegram Group",
//				URL:  "https://t.me/%v?startgroup=new-bill",
//			},
//		},
//		[]tgbotapi.InlineKeyboardButton{
//			tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
//				emoji.ROCKET_ICON+" Split bill with Telegram user(s)",
//				"",
//			),
//		},
//		[]tgbotapi.InlineKeyboardButton{
//			{
//				Text:         emoji.CLIPBOARD_ICON + "New bill manually",
//				CallbackData: "new-bill-manually",
//			},
//		},
//		[]tgbotapi.InlineKeyboardButton{
//			{Text: emoji.RETURN_BACK_ICON + " " + whc.Translate(trans.MESSAGE_TEXT_HELP_BACK_TO_ROOT),
//				CallbackData: "help",
//			},
//		},
//	)
//	return
//}
