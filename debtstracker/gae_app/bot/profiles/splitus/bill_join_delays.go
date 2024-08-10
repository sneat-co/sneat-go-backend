package splitus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/platforms/tgbots"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/delaying"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"strings"
)

func delayUpdateBillCardOnUserJoin(c context.Context, billID string, message string) error {
	if err := delayUpdateBillCards.EnqueueWork(
		c,
		delaying.With(common.QUEUE_BILLS, "update-bill-cards", 0),
		billID,
		message,
	); err != nil {
		logus.Errorf(c, "Failed to queue update of bill cards: %v", err)
	}
	return nil
}

func delayedUpdateBillCards(c context.Context, billID string, footer string) error {
	logus.Debugf(c, "delayedUpdateBillCards(billID=%s)", billID)
	if bill, err := facade2debtus.GetBillByID(c, nil, billID); err != nil {
		return err
	} else {
		for _, tgChatMessageID := range bill.Data.TgChatMessageIDs {
			if err = delayUpdateBillTgChatCard.EnqueueWork(c, delaying.With(common.QUEUE_BILLS, "update-bill-tg-chat-card", 0), billID, tgChatMessageID, footer); err != nil {
				logus.Errorf(c, "Failed to queue updated for %v: %v", tgChatMessageID, err)
				return err
			}
		}
	}
	return nil
}

func delayedUpdateBillTgChartCard(c context.Context, billID string, tgChatMessageID, footer string) error {
	logus.Debugf(c, "delayedUpdateBillTgChartCard(billID=%s, tgChatMessageID=%v)", billID, tgChatMessageID)
	if bill, err := facade2debtus.GetBillByID(c, nil, billID); err != nil {
		return err
	} else {
		ids := strings.Split(tgChatMessageID, "@")
		inlineMessageID, botCode, localeCode5 := ids[0], ids[1], ids[2]
		translator := i18n.NewSingleMapTranslator(i18n.GetLocaleByCode5(localeCode5), i18n.NewMapTranslator(c, trans.TRANS))

		editMessage := tgbotapi.NewEditMessageText(0, 0, inlineMessageID, "")
		editMessage.ParseMode = "HTML"
		editMessage.DisableWebPagePreview = true

		if err := updateInlineBillCardMessage(c, translator, true, editMessage, bill, botCode, footer); err != nil {
			return err
		} else {
			telegramBots := tgbots.Bots(dtdal.HttpAppHost.GetEnvironment(c, nil), nil)
			botSettings, ok := telegramBots.ByCode[botCode]
			if !ok {
				logus.Errorf(c, "No bot settings for bot: "+botCode)
				return nil
			}

			tgApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, dtdal.HttpClient(c))
			if _, err := tgApi.Send(editMessage); err != nil {
				logus.Errorf(c, "Failed to sent message to Telegram: %v", err)
				return err
			}
		}
	}
	return nil
}

func updateInlineBillCardMessage(c context.Context, translator i18n.SingleLocaleTranslator, isGroupChat bool, editedMessage *tgbotapi.EditMessageTextConfig, bill models.Bill, botCode string, footer string) (err error) {
	if bill.ID == "" {
		panic("bill.ID is empty string")
	}
	if bill.Data == nil {
		panic("bill.BillEntity == nil")
	}

	if editedMessage.Text, err = getBillCardMessageText(c, botCode, translator, bill, true, footer); err != nil {
		return
	}
	if isGroupChat {
		editedMessage.ReplyMarkup = getPublicBillCardInlineKeyboard(translator, botCode, bill.ID)
	} else {
		editedMessage.ReplyMarkup = getPrivateBillCardInlineKeyboard(translator, botCode, bill)
	}

	return
}

func getPublicBillCardInlineKeyboard(translator i18n.SingleLocaleTranslator, botCode string, billID string) *tgbotapi.InlineKeyboardMarkup {
	goToBotLink := func(command string) string {
		return fmt.Sprintf("https://t.me/%v?start=%v-%v", botCode, command, billID)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: translator.Translate(trans.BUTTON_TEXT_JOIN),
				URL:  goToBotLink(joinBillCommandCode),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: translator.Translate(trans.BUTTON_TEXT_EDIT_BILL),
				URL:  goToBotLink(editBillCommandCode),
			},
			{
				Text:         translator.Translate(trans.BUTTON_TEXT_DUE, translator.Translate(trans.NOT_SET)),
				CallbackData: billCallbackCommandData(setBillDueDateCommandCode, billID),
			},
		},
	)
	return keyboard
}
