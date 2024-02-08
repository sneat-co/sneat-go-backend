package splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/i18n"
)

func getGroupBillCardInlineKeyboard(translator i18n.SingleLocaleTranslator, bill models.Bill) *tgbotapi.InlineKeyboardMarkup {
	//	//{{Text: "I paid for the bill alone", CallbackData: joinBillCallbackPrefix + "&i=paid-alone"}},
	//	//{{Text:"I paid part of this bill",CallbackData:  joinBillCallbackPrefix + "&i=paid-part"}},
	//	//{{Text: "I owe for this bill", CallbackData: joinBillCallbackPrefix + "&i=owe"}},
	//	//{{Text: "I don't share this bill", CallbackData: billCallbackCommandData(leaveBillCommandCode, bill.ID)}},
	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				{
					Text:         translator.Translate(trans.BUTTON_TEXT_MANAGE_MEMBERS),
					CallbackData: GetBillMembersCallbackData(bill.ID),
				},
			},
			{
				{
					Text:         translator.Translate(trans.BUTTON_TEXT_SPLIT_MODE, translator.Translate(string(bill.Data.SplitMode))),
					CallbackData: billCallbackCommandData(billSharesCommandCode, bill.ID),
				},
			},
			{
				{
					Text:         emoji.GREEN_CHECKBOX + " Finalize bill",
					CallbackData: billCallbackCommandData(finalizeBillCommandCode, bill.ID),
				},
				{
					Text:         emoji.CROSS_MARK + " Delete",
					CallbackData: billCallbackCommandData(deleteBillCommandCode, bill.ID),
				},
			},
		},
	}
}

func getPrivateBillCardInlineKeyboard(translator i18n.SingleLocaleTranslator, botCode string, bill models.Bill) *tgbotapi.InlineKeyboardMarkup {
	callbackData := fmt.Sprintf("split-mode?bill=%v&mode=", bill.ID)
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         translator.Translate(trans.BUTTON_TEXT_MANAGE_MEMBERS),
				CallbackData: GetBillMembersCallbackData(bill.ID),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         translator.Translate(trans.BUTTON_TEXT_CHANGE_BILL_PAYER),
				CallbackData: fmt.Sprintf(CHANGE_BILL_PAYER_COMMAND+"?bill=%v", bill.ID)},
			{
				Text:         translator.Translate(trans.BUTTON_TEXT_SPLIT_MODE, translator.Translate(string(bill.Data.SplitMode))),
				CallbackData: fmt.Sprintf("split-modes?bill=%v", bill.ID),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         translator.Translate("💯 Change total"),
				CallbackData: billCallbackCommandData(CHANGE_BILL_TOTAL_COMMAND, bill.ID),
			},
			{
				Text:         translator.Translate("✍ Adjust per person"),
				CallbackData: callbackData + string(models.SplitModePercentage),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         translator.Translate("📝 Комментарий"),
				CallbackData: billCallbackCommandData(ADD_BILL_COMMENT_COMMAND, bill.ID),
			},
			{
				Text:         translator.Translate(trans.BUTTON_TEXT_FINALIZE_BILL),
				CallbackData: billCallbackCommandData(CLOSE_BILL_COMMAND, bill.ID),
			},
		},
	)
}
