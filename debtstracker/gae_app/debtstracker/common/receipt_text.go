package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/i18n"
	"github.com/strongo/log"
	"github.com/strongo/strongoapp"
	"html"
	"html/template"
	"time"
)

type ShowReceiptTo int8

const (
	ShowReceiptToAutodetect ShowReceiptTo = iota
	ShowReceiptToCreator
	ShowReceiptToCounterparty
)

type ReceiptPartyAction int8

const (
	ReceiptPartyGive = iota
	ReceiptPartyGot
)

type receiptTextBuilder struct {
	//whc botsfw.WebhookContext
	translator    i18n.SingleLocaleTranslator
	transfer      models.Transfer
	showReceiptTo ShowReceiptTo
	viewerUserID  string
	partyAction   ReceiptPartyAction
	//
	strongoapp.ExecutionContext
	//
	//showAds        bool
}

func newReceiptTextBuilder(translator i18n.SingleLocaleTranslator, transfer models.Transfer, showReceiptTo ShowReceiptTo) receiptTextBuilder {
	if transfer.ID == "" {
		panic("transferID == 0")
	}
	r := receiptTextBuilder{
		translator:    translator,
		transfer:      transfer,
		showReceiptTo: showReceiptTo,
	}
	switch showReceiptTo {
	case ShowReceiptToCreator:
		r.viewerUserID = transfer.Data.CreatorUserID
	case ShowReceiptToCounterparty:
		r.viewerUserID = transfer.Data.Counterparty().UserID
	default:
		panic(fmt.Sprintf("Unknown showReceiptTo: %v", showReceiptTo))
	}
	if (showReceiptTo == ShowReceiptToCreator && r.transfer.Data.Direction() == models.TransferDirectionCounterparty2User) ||
		(showReceiptTo == ShowReceiptToCounterparty && r.transfer.Data.Direction() == models.TransferDirectionUser2Counterparty) {
		r.partyAction = ReceiptPartyGot
	} else if (showReceiptTo == ShowReceiptToCounterparty && r.transfer.Data.Direction() == models.TransferDirectionCounterparty2User) ||
		(showReceiptTo == ShowReceiptToCreator && r.transfer.Data.Direction() == models.TransferDirectionUser2Counterparty) {
		r.partyAction = ReceiptPartyGive
	} else {
		if showReceiptTo != ShowReceiptToCreator && showReceiptTo != ShowReceiptToCounterparty {
			panic(fmt.Sprintf("Unknown ShowReceiptTo: %v", r.showReceiptTo))
		}
		panic(fmt.Sprintf("Invalid direction (%v) or showReceiptTo (%v)", r.transfer.Data.Direction(), showReceiptTo))
	}
	return r
}

//func (r receiptTextBuilder) validateRequiredParams() {
//}

func (r receiptTextBuilder) receiptCommonFooter(buffer *bytes.Buffer) {
	transfer := r.transfer
	if r.showReceiptTo == ShowReceiptToCreator && transfer.Data.Creator().Note != "" {
		_, _ = buffer.WriteString("\n" + fmt.Sprintf(emoji.MEMO_ICON+" <b>%v</b>: %v", r.translator.Translate(trans.MESSAGE_TEXT_NOTE), html.EscapeString(transfer.Data.Creator().Note)))
	}
	if r.showReceiptTo == ShowReceiptToCounterparty && transfer.Data.Counterparty().Note != "" {
		_, _ = buffer.WriteString("\n" + fmt.Sprintf(emoji.MEMO_ICON+" <b>%v</b>: %v", r.translator.Translate(trans.MESSAGE_TEXT_NOTE), html.EscapeString(transfer.Data.Counterparty().Note)))
	}

	if transfer.Data.Creator().Comment != "" {
		label := r.translator.Translate(trans.MESSAGE_TEXT_COMMENT)
		_, _ = buffer.WriteString("\n" + fmt.Sprintf(emoji.NEWSPAPER_ICON+" <b>%v</b>: %v", label, html.EscapeString(transfer.Data.Creator().Comment)))
	}
	if transfer.Data.Counterparty().Comment != "" {
		label := r.translator.Translate(trans.MESSAGE_TEXT_COMMENT)
		_, _ = buffer.WriteString("\n" + fmt.Sprintf(emoji.NEWSPAPER_ICON+" <b>%v</b>: %v", label, html.EscapeString(transfer.Data.Counterparty().Comment)))
	}

	//if r.counterpartyID > 0 {
	//	if transfer.CreatorNote != "" || transfer.CreatorComment != "" {
	//		buffer.WriteString(common.HORIZONTAL_LINE)
	//	} else {
	//		buffer.WriteString("\n\n")
	//	}
	//
	//	counterpartyBalance, _ := counterparty.Balance()
	//	utmParams := NewUtmParams(whc, UTM_CAMPAIGN_RECEIPT)
	//	if len(counterpartyBalance) == 0 {
	//		counterpartyLink := GetCounterpartyLink(whc.AppUserID(), whc.Locale(), counterparty.Info(counterpartyID, "", ""), utmParams)
	//		switch transfer.Direction {
	//		case TransferDirectionCounterparty2User:
	//			buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_ON_RETURN_COUNTERPARTY_DOES_NOT_OWE_ANYTHING_TO_USER_ANYMORE, counterpartyLink))
	//		case TransferDirectionUser2Counterparty:
	//			buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_ON_RETURN_USER_DOES_NOT_OWE_ANYTHING_TO_COUNTERPARTY_ANYMORE, counterpartyLink))
	//		default:
	//			panic(fmt.Sprintf("Transfer %v has unkown direction: [%v]", tm.transferID, transfer.Direction))
	//		}
	//	} else {
	//		r.addBalance(whc, buffer, counterpartyID, counterparty, counterpartyBalance, utmParams)
	//	}
	//}

	//if r.showAds {
	//	if !strings.HasSuffix(buffer.String(), common.HORIZONTAL_LINE) {
	//		buffer.WriteString(common.HORIZONTAL_LINE)
	//	}
	//	buffer.WriteString(dtb_general.AdSlot(r.whc, UTM_CAMPAIGN_RECEIPT))
	//}
}

func TextReceiptForTransfer(c context.Context, translator i18n.SingleLocaleTranslator, transfer models.Transfer, showToUserID string, showReceiptTo ShowReceiptTo, utmParams UtmParams) string {
	log.Debugf(c, "TextReceiptForTransfer(transferID=%v, showToUserID=%v, showReceiptTo=%v)", transfer.ID, showToUserID, showReceiptTo)

	if transfer.ID == "" {
		panic("transferID == 0")
	}
	if transfer.Data == nil {
		panic("transferID == 0")
	}

	//transferEntity := transfer.TransferData

	switch showReceiptTo {
	case ShowReceiptToCreator:
		if showToUserID != "" && showToUserID != transfer.Data.CreatorUserID {
			panic("showToUserID != 0 && showToUserID != transferEntity.CreatorUserID")
		}
	case ShowReceiptToCounterparty:
		if showToUserID != "" && transfer.Data.Counterparty().UserID != "" && showToUserID != transfer.Data.Counterparty().UserID {
			panic("showToUserID != 0 && showToUserID != transferEntity.Counterparty().UserID")
		}
	case ShowReceiptToAutodetect:
		switch showToUserID {
		case transfer.Data.CreatorUserID:
			showReceiptTo = ShowReceiptToCreator
		case transfer.Data.Counterparty().UserID:
			showReceiptTo = ShowReceiptToCounterparty
		default:
			if transfer.Data.Counterparty().UserID == "" {
				showReceiptTo = ShowReceiptToCounterparty
			} else {
				panic(fmt.Sprintf("Parameter showToUserID=%v is not related to transferEntity with id=%v", showToUserID, transfer.ID))
			}
		}
	}

	r := newReceiptTextBuilder(translator, transfer, showReceiptTo)

	var buffer bytes.Buffer
	r.WriteReceiptText(&buffer, utmParams)
	r.receiptCommonFooter(&buffer)
	return buffer.String()
}

func (r receiptTextBuilder) getReceiptCounterparty() *models.TransferCounterpartyInfo {
	switch r.showReceiptTo {
	case ShowReceiptToCreator:
		return r.transfer.Data.Counterparty()
	case ShowReceiptToCounterparty:
		return r.transfer.Data.Creator()
	default:
		panic(fmt.Sprintf("Unknown ShowReceiptTo: %v", r.showReceiptTo))
	}
}

//func (r receiptTextBuilder) receiptOnReturn(utmParams UtmParams) string {
//	var messageTextToTranslate string
//	return r.translateAndFormatMessage(messageTextToTranslate, r.transfer.Data.GetAmount(), utmParams)
//}

func (r receiptTextBuilder) WriteReceiptText(buffer *bytes.Buffer, utmParams UtmParams) {
	var messageTextToTranslate string
	if r.transfer.Data.IsReturn {
		switch r.partyAction {
		case ReceiptPartyGive:
			messageTextToTranslate = trans.MESSAGE_TEXT_RECEIPT_RETURN_FROM_USER
		case ReceiptPartyGot:
			messageTextToTranslate = trans.MESSAGE_TEXT_RECEIPT_RETURN_TO_USER
		default:
			panic(fmt.Sprintf("Unknown partyAction: %v", r.partyAction))
		}
	} else {
		switch r.partyAction {
		case ReceiptPartyGive:
			messageTextToTranslate = trans.MESSAGE_TEXT_RECEIPT_NEW_DEBT_FROM_USER
		case ReceiptPartyGot:
			messageTextToTranslate = trans.MESSAGE_TEXT_RECEIPT_NEW_DEBT_TO_USER
		default:
			panic(fmt.Sprintf("Unknown partyAction: %v", r.partyAction))
		}
	}

	buffer.WriteString(r.translateAndFormatMessage(messageTextToTranslate, r.transfer.Data.GetAmount(), utmParams))

	if r.transfer.Data.HasInterest() {
		buffer.WriteString("\n")
		WriteTransferInterest(buffer, r.transfer, r.translator)
	}

	if !r.transfer.Data.DtDueOn.IsZero() {
		buffer.WriteString("\n" + emoji.ALARM_CLOCK_ICON + " " + fmt.Sprintf(r.translator.Translate(trans.MESSAGE_TEXT_DUE_ON), r.transfer.Data.DtDueOn.Format("2006-01-02 15:04")))
	}

	if amountReturned := r.transfer.Data.AmountReturned(); amountReturned > 0 && amountReturned != r.transfer.Data.AmountInCents {
		buffer.WriteString("\n" + r.translateAndFormatMessage(trans.MESSAGE_TEXT_RECEIPT_ALREADY_RETURNED_AMOUNT, r.transfer.Data.GetReturnedAmount(), utmParams))
	}

	if outstandingAmount := r.transfer.Data.GetOutstandingAmount(time.Now()); outstandingAmount.Value > 0 && outstandingAmount.Value != r.transfer.Data.AmountInCents {
		buffer.WriteString("\n" + r.translateAndFormatMessage(trans.MESSAGE_TEXT_RECEIPT_OUTSTANDING_AMOUNT, outstandingAmount, utmParams))
	}
}

func WriteTransferInterest(buffer *bytes.Buffer, transfer models.Transfer, translator i18n.SingleLocaleTranslator) {
	buffer.WriteString(translator.Translate(trans.MESSAGE_TEXT_INTEREST, transfer.Data.InterestPercent, days(translator, int(transfer.Data.InterestPeriod))))
	if transfer.Data.InterestMinimumPeriod > 1 {
		buffer.WriteString(", " + translator.Translate(trans.MESSAGE_TEXT_INTEREST_MIN_PERIOD, days(translator, transfer.Data.InterestMinimumPeriod)))
	}
}

func days(t i18n.SingleLocaleTranslator, d int) string {
	var messageTextToTranslate string
	if d == 1 {
		messageTextToTranslate = trans.DAY
	} else if d <= 4 {
		messageTextToTranslate = trans.DAYS_234
	} else {
		messageTextToTranslate = trans.DAYS
	}
	return t.Translate(messageTextToTranslate, d)
}

func (r receiptTextBuilder) translateAndFormatMessage(messageTextToTranslate string, amount money.Amount, utmParams UtmParams) string {
	userID := r.viewerUserID

	counterpartyInfo := r.getReceiptCounterparty()

	var counterpartyText string
	{
		// TODO: Disabled URL due to issue with Telegram parser
		if userID == "" || utmParams.Medium == UTM_MEDIUM_SMS || utmParams.Medium == telegram.PlatformID {
			counterpartyText = counterpartyInfo.Name()
		} else {
			counterpartyUrl := GetCounterpartyUrl(counterpartyInfo.ContactID, userID, r.translator.Locale(), utmParams)
			counterpartyText = fmt.Sprintf(`<a href="%v"><b>%v</b></a>`, counterpartyUrl, html.EscapeString(counterpartyInfo.Name()))
		}
		// TODO: Add a @counterparty Telegram nickname if sending receipt to Telegram channel
	}

	//var amountText string
	//{
	//	transferUrl := GetTransferUrlForUser(r.transfer.ID, userID, r.Locale(), utmParams)
	//	if utmParams.Medium == UTM_MEDIUM_SMS {
	//		amountText = fmt.Sprintf(`%v - %v`, r.transfer.GetAmount(), transferUrl)
	//	} else {
	//		amountText = fmt.Sprintf(`<a href="%v">%v</a>`, transferUrl, r.transfer.GetAmount())
	//	}
	//}
	amountText := fmt.Sprintf("<b>%v</b>", amount)

	return r.translator.Translate(messageTextToTranslate, map[string]interface{}{
		"Counterparty": template.HTML(counterpartyText),
		"Amount":       template.HTML(amountText),
	})
}

//func (r receiptBuilder) addBalance(buffer *bytes.Buffer, counterpartyBalance Balance, utmParams UtmParams) string {
//	if counterpartyID == 0 {
//		return ""
//	}
//	showBalanceMessage := true
//	transfer := tm.transfer
//	if len(counterpartyBalance) == 1 {
//		transferAmount := transfer.GetAmount()
//		if singleCurrencyVal, ok := counterpartyBalance[transferAmount.Currency]; !ok || singleCurrencyVal != transfer.Amount {
//			showBalanceMessage = true
//		}
//	}
//	if showBalanceMessage {
//		counterpartyLink := GetCounterpartyLink(whc.AppUserID(), whc.Locale(), counterparty.Info(counterpartyID, "", ""), utmParams)
//		buffer.WriteString(BalanceForCounterpartyWithHeader(counterpartyLink, counterpartyBalance, tm.executionContext.Logger(), tm.executionContext))
//	}
//	return buffer.String()
//}
