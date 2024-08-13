package botcmds4splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
	"net/url"
)

var billSplitModesListCommand = billCallbackCommand("split-modes",
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "billSplitModesListCommand.CallbackAction()")
		var mt string
		if mt, err = getBillCardMessageText(c, whc.GetBotCode(), whc, bill, true, ""); err != nil {
			return
		}
		if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
			return
		}
		callbackData := fmt.Sprintf("split-mode?bill=%v&mode=", bill.ID)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.SPLIT_MODE_EQUALLY),
					CallbackData: callbackData + string(models4splitus.SplitModeEqually),
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.SPLIT_MODE_PERCENTAGE),
					CallbackData: callbackData + string(models4splitus.SplitModePercentage),
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.SPLIT_MODE_SHARES),
					CallbackData: callbackData + string(models4splitus.SplitModeShare),
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.SPLIT_MODE_EXACT_AMOUNT),
					CallbackData: callbackData + string(models4splitus.SplitModeExactAmount),
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.BUTTON_TEXT_CANCEL),
					CallbackData: billCardCallbackCommandData(bill.ID),
				},
			},
		)
		m.Keyboard = keyboard
		return
	},
)
