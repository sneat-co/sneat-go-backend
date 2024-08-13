package dtb_transfer

import (
	"fmt"
)

const (
	ReceiptActionDoNotSend     = "do-not-send"
	SendReceiptCallbackPath    = "send-receipt"
	SendReceiptByChooseChannel = "select"
	WizardParamTransfer        = "transfer"
	WizardParamReminder        = "reminder"
	WizardParamSpace           = "space"
	WizardParamCounterparty    = "counterparty" // TODO: Decide use this or WizardParamContact
	WizardParamContact         = "contact"      // TODO: Decide use this or WizardParamCounterparty
)

type SendReceipt struct {
	//transferID int
	By string
}

func SendReceiptCallbackData(transferID string, by string) string {
	return fmt.Sprintf("%s?by=%s&transfer=%s", SendReceiptCallbackPath, by, transferID)
}

func SendReceiptUrl(transferID string, by string) string {
	return fmt.Sprintf("https://debtus.app/pwa/send-receipt?by=%s&transfer=%s", by, transferID)
}
