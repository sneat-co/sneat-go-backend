package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"net/url"
)

const ACKNOWLEDGE_RECEIPT_CALLBACK_COMMAND = "ack-receipt"

var AcknowledgeReceiptCallbackCommand = botsfw.NewCallbackCommand(ACKNOWLEDGE_RECEIPT_CALLBACK_COMMAND, func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	query := callbackUrl.Query()
	receiptID := query.Get("id")
	if receiptID == "" {
		return m, fmt.Errorf("receiptID is empty")
	}

	return AcknowledgeReceipt(whc, receiptID, query.Get("do"))
})
