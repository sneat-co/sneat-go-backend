package dtb_inline

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/decimal"
	"github.com/strongo/logus"
	"html"
	"net/url"
	"regexp"
	"strings"
)

var ReInlineQueryAmount = regexp.MustCompile(`^\s*(\d+(?:\.\d*)?)\s*((?:\b|\B).+?)?\s*$`)

func InlineNewRecord(whc botsfw.WebhookContext, amountMatches []string) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	logus.Debugf(ctx, "InlineNewRecord()")

	inlineQuery := whc.Input().(botinput.WebhookInlineQuery)
	var (
		amountValue    decimal.Decimal64p2
		amountCurrency money.CurrencyCode
	)
	if amountValue, err = decimal.ParseDecimal64p2(strings.TrimRight(amountMatches[1], ".")); err != nil {
		return
	}
	currencyCode := strings.TrimRight(amountMatches[2], ".,;()[]{} ")
	logus.Debugf(ctx, "currencyCode: %v", currencyCode)
	if currencyCode != "" {
		if len(currencyCode) > 20 {
			currencyCode = currencyCode[:20]
		}
		ccLow := strings.ToLower(currencyCode)
		if ccLow == money.CurrencySymbolRUR || ccLow == "р" || ccLow == "руб" || ccLow == "рубля" || ccLow == "рублей" || ccLow == "rub" || ccLow == "rubles" || ccLow == "ruble" || ccLow == "rubley" {
			amountCurrency = money.CurrencySymbolRUR
		} else if ccLow == "eur" || ccLow == "euro" || ccLow == money.CurrencySymbolEUR {
			amountCurrency = money.CurrencyEUR
		} else if ccLow == "гривна" || ccLow == "гривен" || ccLow == "г" || ccLow == money.CurrencySymbolUAH {
			amountCurrency = money.CurrencyUAH
		} else if ccLow == "тенге" || ccLow == "теңге" || ccLow == "т" || ccLow == money.CurrencySymbolKZT {
			amountCurrency = money.CurrencyKZT
		} else {
			amountCurrency = money.CurrencyCode(currencyCode)
		}
	} else {
		amountCurrency = money.CurrencyUSD
	}

	amountText := html.EscapeString(money.NewAmount(amountCurrency, amountValue).String())

	newBillCallbackData := fmt.Sprintf("new-bill?v=%v&ctx=%v", amountMatches[1], url.QueryEscape(string(amountCurrency)))
	m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.GetInlineQueryID(),
		Results: []interface{}{
			tgbotapi.InlineQueryResultArticle{
				ID:          "SplitBill_" + whc.Locale().Code5,
				Type:        "article",
				Title:       "🛒 " + whc.Translate(trans.ARTICLE_TITLE_SPLIT_BILL),
				Description: whc.Translate(trans.ARTICLE_SUBTITLE_SPLIT_BILL, amountText),
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:                  whc.Translate(trans.MESSAGE_TEXT_BILL_HEADER, amountText),
					ParseMode:             "HTML",
					DisableWebPagePreview: true,
				},
				ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
						{
							{Text: whc.Translate(trans.COMMAND_TEXT_I_PAID), CallbackData: newBillCallbackData + "&i=paid"},
							{Text: whc.Translate(trans.COMMAND_TEXT_I_OWE), CallbackData: newBillCallbackData + "&i=owe"},
						},
					},
				},
			},
			tgbotapi.InlineQueryResultArticle{
				ID:          "NewDebt_" + whc.Locale().Code5,
				Type:        "article",
				Title:       "💵 " + whc.Translate(trans.ARTICLE_NEW_DEBT_TITLE),
				Description: whc.Translate(trans.ARTICLE_NEW_DEBT_SUBTITLE, amountText),
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:                  whc.Translate(trans.MESSAGE_TEXT_NEW_DEBT_HEADER, amountText),
					ParseMode:             "HTML",
					DisableWebPagePreview: true,
				},
				ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
						{
							{Text: whc.Translate(trans.COMMAND_TEXT_I_OWE), CallbackData: "i-owed?debt=new"},
							{Text: whc.Translate(trans.COMMAND_TEXT_OWED_TO_ME), CallbackData: "owed2me?debt=new"},
						},
					},
				},
			},
		},
	})
	return m, err
}
