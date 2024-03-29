package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_inline"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
)

func showReceiptAnnouncement(whc botsfw.WebhookContext, receiptID string, creatorName string) (m botsfw.MessageFromBot, err error) {
	var inlineMessageID string
	switch input := whc.Input().(type) {
	case botsfw.WebhookChosenInlineResult:
		inlineMessageID = input.GetInlineMessageID()
	case botsfw.WebhookCallbackQuery:
		inlineMessageID = input.GetInlineMessageID()
	default:
		return m, fmt.Errorf("showReceiptAnnouncement: Unsupported InputType=%T", input)
	}

	c := whc.Context()

	receipt, err := dtdal.Receipt.GetReceiptByID(c, nil, receiptID)
	if err != nil {
		return m, err
	}
	if creatorName == "" {
		user, err := facade.User.GetUserByID(c, nil, receipt.Data.CreatorUserID)
		if err != nil {
			return m, err
		}
		creatorName = user.Data.FullName()
	}

	messageText := getInlineReceiptMessageText(whc, whc.GetBotCode(), whc.Locale().Code5, creatorName, receiptID)
	m, err = whc.NewEditMessage(messageText, botsfw.MessageFormatHTML)
	m.EditMessageUID = telegram.NewInlineMessageUID(inlineMessageID)
	m.DisableWebPagePreview = true
	kbRows := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData(
				whc.Translate(trans.COMMAND_TEXT_VIEW_RECEIPT_DETAILS),
				fmt.Sprintf("%s?id=%s&locale=%s",
					VIEW_RECEIPT_IN_TELEGRAM_COMMAND, receiptID, whc.Locale().Code5,
				),
			),
		},
	}
	kbRows = append(kbRows, dtb_inline.GetChooseLangInlineKeyboard(
		fmt.Sprintf("%s?id=%s", CHANGE_RECEIPT_LANG_COMMAND, receiptID)+"&locale=%v", // Intentionally &locale separate
		whc.Locale().Code5,
	)...)
	m.Keyboard = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: kbRows,
	}
	return
}

const VIEW_RECEIPT_IN_TELEGRAM_COMMAND = "tg-view-receipt"

func GetUrlForReceiptInTelegram(botCode string, receiptID string, localeCode5 string) string {
	return fmt.Sprintf("https://t.me/%v?start=receipt-%v-view_%v", botCode, receiptID, localeCode5)
}
