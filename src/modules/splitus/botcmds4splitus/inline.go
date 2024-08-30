package botcmds4splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"html"
	"net/url"
	"regexp"
	"strings"

	"errors"
	"github.com/bots-go-framework/bots-fw-telegram"
)

var reInlineQueryNewBill = regexp.MustCompile(`^\s*(\d+(?:\.\d*)?)([^\s]*)\s+(.+?)\s*$`)

var inlineQueryCommand = botsfw.Command{
	Code:       "inline-query",
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputInlineQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		whc.Input().LogRequest()
		ctx := whc.Context()
		if tgInput, ok := whc.Input().(telegram.TgWebhookInput); ok {
			update := tgInput.TgUpdate()

			if appUserData, err := whc.AppUserData(); err != nil {
				return m, err
			} else if preferredLocale := appUserData.BotsFwAdapter().GetPreferredLocale(); preferredLocale != "" {
				logus.Debugf(ctx, "User has preferring locale")
				_ = whc.SetLocale(preferredLocale)
			} else if tgLang := update.InlineQuery.From.LanguageCode; len(tgLang) >= 2 {
				switch strings.ToLower(tgLang[:2]) {
				case "ru":
					logus.Debugf(ctx, "Telegram client has known language code")
					if err = whc.SetLocale(i18n.LocaleRuRu.Code5); err != nil {
						return m, err
					}
				}
			}
		}
		inlineQuery := whc.Input().(botinput.WebhookInlineQuery)
		query := strings.TrimSpace(inlineQuery.GetQuery())
		logus.Debugf(ctx, "inlineQueryCommand.Action(query=%v)", query)
		switch {
		case query == "":
			return inlineEmptyQuery(whc, inlineQuery)
		case strings.HasPrefix(query, joinSpaceCommandCode+"?id="):
			return inlineQueryJoinGroup(whc, query)
		default:
			if reMatches := reInlineQueryNewBill.FindStringSubmatch(query); reMatches != nil {
				return inlineQueryNewBill(whc, reMatches[1], reMatches[2], reMatches[3])
			}
			logus.Debugf(ctx, "Inline query not matched to any action: [%v]", query)
		}

		return
	},
}

func inlineEmptyQuery(whc botsfw.WebhookContext, inlineQuery botinput.WebhookInlineQuery) (m botsfw.MessageFromBot, err error) {
	logus.Debugf(whc.Context(), "InlineEmptyQuery()")
	m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
		InlineQueryID:     inlineQuery.GetInlineQueryID(),
		CacheTime:         60,
		SwitchPMText:      "Help: How to use this bot?",
		SwitchPMParameter: "help_inline",
	})
	return
}

func inlineQueryJoinGroup(whc botsfw.WebhookContext, query string) (m botsfw.MessageFromBot, err error) {
	err = errors.New("inlineQueryJoinGroup is not implemented")
	//ctx := whc.Context()
	//
	//inlineQuery := whc.Input().(botsfw.WebhookInlineQuery)
	//
	//var group models4splitus.GroupEntry
	//if group.ContactID = query[len(joinSpaceCommandCode+"?id="):]; group.ContactID == "" {
	//	err = errors.New("Missing group ContactID")
	//	return
	//}
	//if group, err = dtdal.Group.GetGroupByID(ctx, nil, group.ContactID); err != nil {
	//	return
	//}
	//
	//joinBillInlineResult := tgbotapi.InlineQueryResultArticle{
	//	ContactID:          query,
	//	Type:        "article",
	//	Title:       "Send invite for joining",
	//	Description: "group: " + group.Data.Name,
	//	InputMessageContent: tgbotapi.InputTextMessageContent{
	//		Text:      fmt.Sprintf("I'm inviting you to join <b>bills sharing</b> group @%v.", whc.GetBotCode()),
	//		ParseMode: "HTML",
	//	},
	//	ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
	//		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
	//			{
	//				{
	//					Text:         "Join",
	//					CallbackData: query,
	//				},
	//			},
	//		},
	//	},
	//}
	//
	//m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
	//	InlineQueryID: inlineQuery.GetInlineQueryID(),
	//	CacheTime:     60,
	//	Results: []interface{}{
	//		joinBillInlineResult,
	//	},
	//})
	return
}

func inlineQueryNewBill(whc botsfw.WebhookContext, amountNum, amountCurr, billName string) (m botsfw.MessageFromBot, err error) {
	if len(amountCurr) == 3 {
		amountCurr = strings.ToUpper(amountCurr)
	}

	m.Text = fmt.Sprintf("Amount: %v %v, BillEntry name: %v", amountNum, amountCurr, billName)

	inlineQuery := whc.Input().(botinput.WebhookInlineQuery)

	params := fmt.Sprintf("amount=%v&lang=%v", url.QueryEscape(amountNum+amountCurr), whc.Locale().Code5)

	resultID := "bill?" + params

	newBillInlineResult := tgbotapi.InlineQueryResultArticle{
		ID:          resultID,
		Type:        "article",
		Title:       fmt.Sprintf("%v: %v", whc.Translate(trans.COMMAND_TEXT_NEW_BILL), billName),
		Description: fmt.Sprintf("%v: %v %v", whc.Translate(trans.HTML_AMOUNT), amountNum, amountCurr),
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text: fmt.Sprintf("<b>%v</b>: %v %v - %v",
				whc.Translate(trans.MESAGE_TEXT_CREATING_BILL),
				html.EscapeString(amountNum),
				html.EscapeString(amountCurr),
				html.EscapeString(billName),
			),
			ParseMode: "HTML",
		},
		ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					{
						Text:         whc.Translate(trans.MESSAGE_TEXT_PLEASE_WAIT),
						CallbackData: "creating-bill?" + params,
					},
				},
			},
		},
	}

	m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.GetInlineQueryID(),
		CacheTime:     60,
		Results: []interface{}{
			newBillInlineResult,
		},
	})

	return
}
