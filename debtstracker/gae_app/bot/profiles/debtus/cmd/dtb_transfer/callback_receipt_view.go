package dtb_transfer

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/delaying"
	"github.com/strongo/log"
	"net/url"
	"strings"
)

var ViewReceiptInTelegramCallbackCommand = botsfw.NewCallbackCommand(
	VIEW_RECEIPT_IN_TELEGRAM_COMMAND,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		log.Debugf(c, "ViewReceiptInTelegramCallbackCommand.CallbackAction()")
		query := callbackUrl.Query()

		receiptID := query.Get("id")
		receipt, err := dtdal.Receipt.GetReceiptByID(c, nil, receiptID)
		if err != nil {
			return m, err
		}
		currentUserID := whc.AppUserID()
		if receipt.Data.CreatorUserID != currentUserID {
			if err = linkUsersByReceiptNowOrDelay(c, receipt, currentUserID); err != nil {
				log.Errorf(c, err.Error())
				err = nil // We still can create link to receipt, so log error and continue
			}
		}
		localeCode5 := query.Get("locale")
		if len(localeCode5) != 5 {
			return m, errors.New("len(localeCode5) != 5")
		}

		callbackAnswer := tgbotapi.NewCallbackWithURL(
			GetUrlForReceiptInTelegram(whc.GetBotCode(), receiptID, localeCode5),
			//common.GetReceiptUrlForUser(
			//	receiptID,
			//	whc.AppUserID(),
			//	whc.BotPlatform().ID(),
			//	whc.GetBotCode(),
			//) + "&lang=" + localeCode5,
		)
		m.BotMessage = telegram.CallbackAnswer(callbackAnswer)
		// TODO: https://core.telegram.org/bots/api#answercallbackquery, show_alert = true
		return
	},
)

const delayLinkUserByReceiptKeyName = "delayLinkUserByReceipt"

func DelayLinkUsersByReceipt(c context.Context, receiptID, invitedUserID string) (err error) {
	return delayLinkUserByReceipt.EnqueueWork(c, delaying.With(common.QUEUE_RECEIPTS, delayLinkUserByReceiptKeyName, 0), receiptID, invitedUserID)
}

func delayedLinkUsersByReceipt(c context.Context, receiptID, invitedUserID string) error {
	log.Debugf(c, "delayedLinkUsersByReceipt(receiptID=%v, invitedUserID=%v)", receiptID, invitedUserID)
	receipt, err := dtdal.Receipt.GetReceiptByID(c, nil, receiptID)
	if err != nil {
		if dal.IsNotFound(err) {
			log.Errorf(c, err.Error())
			err = nil
		}
		return err
	}
	return linkUsersByReceipt(c, receipt, invitedUserID)
}

func linkUsersByReceiptNowOrDelay(c context.Context, receipt models.Receipt, invitedUserID string) (err error) {
	if err = linkUsersByReceipt(c, receipt, invitedUserID); err != nil {
		err = fmt.Errorf("failed to link users by receipt: %w", err)
		if strings.Contains(err.Error(), "concurrent transaction") {
			log.Warningf(c, err.Error())
			if err = DelayLinkUsersByReceipt(c, receipt.ID, invitedUserID); err != nil {
				err = fmt.Errorf("failed to delay linking users by receipt: %w", err)
			}
		}
	}
	return
}

func linkUsersByReceipt(c context.Context, receipt models.Receipt, invitedUserID string) (err error) {
	if receipt.Data.CounterpartyUserID == "" {
		linker := facade.NewReceiptUsersLinker(nil) // TODO: Link users
		if _, err = linker.LinkReceiptUsers(c, receipt.ID, invitedUserID); err != nil {
			return err
		}
	} else if receipt.Data.CounterpartyUserID != invitedUserID {
		// TODO: Should we allow to see receipt but block from changing it?
		log.Warningf(c, `Security issue: receipt.CreatorUserID != currentUserID && receipt.CounterpartyUserID != currentUserID
	currentUserID: %d
	receipt.CreatorUserID: %d
	receipt.CounterpartyUserID: %d
				`, invitedUserID, receipt.Data.CreatorUserID, receipt.Data.CounterpartyUserID)
		//} else {
		// receipt.CounterpartyUserID == currentUserID - we are fine
	}
	return nil
}
