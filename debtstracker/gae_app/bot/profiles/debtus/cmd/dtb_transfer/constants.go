package dtb_transfer

import (
	"fmt"
)

const (
	RECEIPT_ACTION__DO_NOT_SEND    = "do-not-send"
	SEND_RECEIPT_CALLBACK_PATH     = "send-receipt"
	SEND_RECEIPT_BY_CHOOSE_CHANNEL = "select"
	WIZARD_PARAM_TRANSFER          = "transfer"
	WIZARD_PARAM_REMINDER          = "reminder"
	WIZARD_PARAM_COUNTERPARTY      = "counterparty" // TODO: Decide use this or WIZARD_PARAM_CONTACT
	WIZARD_PARAM_CONTACT           = "contact"      // TODO: Decide use this or WIZARD_PARAM_COUNTERPARTY
)

type SendReceipt struct {
	//transferID int
	By string
}

func SendReceiptCallbackData(transferID string, by string) string {
	return fmt.Sprintf("%s?by=%s&transfer=%s", SEND_RECEIPT_CALLBACK_PATH, by, transferID)
}

func SendReceiptUrl(transferID string, by string) string {
	return fmt.Sprintf("https://debtus.app/pwa/send-receipt?by=%s&transfer=%s", by, transferID)
}
