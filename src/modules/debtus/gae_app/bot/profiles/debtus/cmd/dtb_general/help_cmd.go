package dtb_general

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
)

func HelpCommandAction(whc botsfw.WebhookContext, showFeedbackButton bool) (m botsfw.MessageFromBot, err error) {
	keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: emoji.PUBLIC_LOUDSPEAKER + " " + whc.Translate(trans.COMMAND_TEXT_OPEN_USER_REPORT),
				URL:  getUserReportUrl(whc, ""),
			},
		},
		[]tgbotapi.InlineKeyboardButton{btnSubmitBug(whc, getUserReportUrl(whc, "bug"))},
		[]tgbotapi.InlineKeyboardButton{btnSubmitIdea(whc, getUserReportUrl(whc, "idea"))},
	)
	if showFeedbackButton {
		keyboardMarkup.InlineKeyboard = append(
			keyboardMarkup.InlineKeyboard,
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         emoji.STAR_ICON + " " + whc.Translate(trans.COMMAND_TEXT_ASK_FOR_FEEDBACK),
					CallbackData: FEEDBACK_COMMAND,
				},
			})
	}
	if showFeedbackButton {
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_HELP)
		m.Keyboard = keyboardMarkup
	} else {
		if m, err = whc.NewEditMessage("", botsfw.MessageFormatText); err != nil {
			return
		}
		m.Keyboard = keyboardMarkup
	}

	return m, err
}

func getUserReportUrl(t i18n.SingleLocaleTranslator, submit string) string {
	switch t.Locale().Code5 {
	case i18n.LocalCodeRuRu:
		switch submit {
		case "idea":
			return "https://goo.gl/dAKHFC"
		case "bug":
			return "https://goo.gl/jQM2K5"
		case "":
			return "https://goo.gl/Vge31X"
		default:
			panic("Parameter 'submit' should be either 'idea' or 'bug'")
		}
	default:
		switch submit {
		case "idea":
			return "https://goo.gl/sl09Wr"
		case "bug":
			return "https://goo.gl/x5H6Fn"
		case "":
			return "https://goo.gl/3tB0FG"
		default:
			panic("Parameter 'submit' should be either 'idea' or 'bug'")
		}
	}
}

func btnSubmitIdea(whc botsfw.WebhookContext, url string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text: emoji.BULB_ICON + " " + whc.Translate(trans.COMMAND_TEXT_SUBMIT_AN_IDEA),
		URL:  url,
	}
}

func btnSubmitBug(whc botsfw.WebhookContext, url string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text: emoji.ERROR_ICON + " " + whc.Translate(trans.COMMAND_TEXT_REPORT_A_BUG),
		URL:  url,
	}
}

const ADS_COMMAND = "ads"

var AdsCommand = botsfw.Command{
	Code:     ADS_COMMAND,
	Icon:     emoji.NEWSPAPER_ICON,
	Commands: []string{emoji.NEWSPAPER_ICON, "/ads", "/реклама"},
	Title:    trans.COMMAND_TEXT_HELP,
	Titles:   map[string]string{botsfw.ShortTitle: ""},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		chatData := whc.ChatData()

		yesOption := emoji.PHONE_ICON + " " + whc.Translate(trans.COMMAND_TEXT_SUBSCRIBE_TO_APP)
		noOption := whc.Translate(trans.COMMAND_TEXT_I_AM_FINE_WITH_BOT)
		if chatData.GetAwaitingReplyTo() == "" {
			chatData.SetAwaitingReplyTo(ADS_COMMAND)
			m = whc.NewMessage(emoji.NEWSPAPER_ICON + " " + whc.Translate(trans.MESSAGE_TEXT_YOUR_ABOUT_ADS))
			m.DisableWebPagePreview = true
			m.Keyboard = tgbotapi.NewReplyKeyboard(
				[]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(yesOption)},
				[]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(noOption)},
				[]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(MainMenuCommand.DefaultTitle(whc))},
			)
		} else {
			switch whc.Input().(botsfw.WebhookTextMessage).Text() {
			case yesOption:
				m = whc.NewMessageByCode(trans.MESSAGE_TEXT_SUBSCRIBED_TO_APP)
				SetMainMenuKeyboard(whc, &m)
				chatData.SetAwaitingReplyTo("")
			case noOption:
				m = whc.NewMessageByCode(trans.MESSAGE_TEXT_NOT_INTERESTED_IN_APP)
				SetMainMenuKeyboard(whc, &m)
				chatData.SetAwaitingReplyTo("")
			default:
				m = whc.NewMessageByCode(trans.MESSAGE_TEXT_PLEASE_CHOOSE_FROM_OPTIONS_PROVIDED)
			}
		}
		return m, err
	},
}
