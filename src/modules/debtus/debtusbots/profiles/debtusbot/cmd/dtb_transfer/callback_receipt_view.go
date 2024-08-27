package dtb_transfer

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"net/url"
	"strings"
)

var ViewReceiptInTelegramCallbackCommand = botsfw.NewCallbackCommand(
	VIEW_RECEIPT_IN_TELEGRAM_COMMAND,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		logus.Debugf(ctx, "ViewReceiptInTelegramCallbackCommand.CallbackAction()")
		query := callbackUrl.Query()

		receiptID := query.Get("id")
		receipt, err := dtdal.Receipt.GetReceiptByID(ctx, nil, receiptID)
		if err != nil {
			return m, err
		}
		currentUserID := whc.AppUserID()
		if receipt.Data.CreatorUserID != currentUserID {
			if err = linkUsersByReceiptNowOrDelay(ctx, receipt, currentUserID); err != nil {
				logus.Errorf(ctx, err.Error())
				err = nil // We still can create link to receipt, so log error and continue
			}
		}
		localeCode5 := query.Get("locale")
		if len(localeCode5) != 5 {
			return m, errors.New("len(localeCode5) != 5")
		}

		callbackAnswer := tgbotapi.NewCallbackWithURL(
			GetUrlForReceiptInTelegram(whc.GetBotCode(), receiptID, localeCode5),
			//anybot.GetReceiptUrlForUser(
			//	receiptID,
			//	whc.AppUserID(),
			//	whc.BotPlatform().ContactID(),
			//	whc.GetBotCode(),
			//) + "&lang=" + localeCode5,
		)
		m.BotMessage = telegram.CallbackAnswer(callbackAnswer)
		// TODO: https://core.telegram.org/bots/api#answercallbackquery, show_alert = true
		return
	},
)

const delayLinkUserByReceiptKeyName = "delayLinkUserByReceipt"

func DelayLinkUsersByReceipt(ctx context.Context, receiptID, invitedUserID string) (err error) {
	return delayLinkUserByReceipt.EnqueueWork(ctx, delaying.With(const4debtus.QueueReceipts, delayLinkUserByReceiptKeyName, 0), receiptID, invitedUserID)
}

func delayedLinkUsersByReceipt(ctx context.Context, receiptID, invitedUserID string) error {
	logus.Debugf(ctx, "delayedLinkUsersByReceipt(receiptID=%v, invitedUserID=%v)", receiptID, invitedUserID)
	receipt, err := dtdal.Receipt.GetReceiptByID(ctx, nil, receiptID)
	if err != nil {
		if dal.IsNotFound(err) {
			logus.Errorf(ctx, err.Error())
			err = nil
		}
		return err
	}
	return linkUsersByReceipt(ctx, receipt, invitedUserID)
}

func linkUsersByReceiptNowOrDelay(ctx context.Context, receipt models4debtus.ReceiptEntry, invitedUserID string) (err error) {
	if err = linkUsersByReceipt(ctx, receipt, invitedUserID); err != nil {
		err = fmt.Errorf("failed to link users by receipt: %w", err)
		if strings.Contains(err.Error(), "concurrent transaction") {
			logus.Warningf(ctx, err.Error())
			if err = DelayLinkUsersByReceipt(ctx, receipt.ID, invitedUserID); err != nil {
				err = fmt.Errorf("failed to delay linking users by receipt: %w", err)
			}
		}
	}
	return
}

func linkUsersByReceipt(ctx context.Context, receipt models4debtus.ReceiptEntry, invitedUserID string) (err error) {
	if receipt.Data.CounterpartyUserID == "" {
		linker := facade4debtus.NewReceiptUsersLinker(nil) // TODO: Link users
		if _, err = linker.LinkReceiptUsers(ctx, receipt.ID, invitedUserID); err != nil {
			return err
		}
	} else if receipt.Data.CounterpartyUserID != invitedUserID {
		// TODO: Should we allow to see receipt but block from changing it?
		logus.Warningf(ctx, `Security issue: receipt.CreatorUserID != currentUserID && receipt.CounterpartyUserID != currentUserID
	currentUserID: %s
	receipt.CreatorUserID: %s
	receipt.CounterpartyUserID: %s
				`, invitedUserID, receipt.Data.CreatorUserID, receipt.Data.CounterpartyUserID)
		//} else {
		// receipt.CounterpartyUserID == currentUserID - we are fine
	}
	return nil
}
