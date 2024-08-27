package botcmds4splitus

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"net/url"

	"context"
	"errors"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/strongo/decimal"
)

const billCardCommandCode = "bill-card"

var billCardCommand = botsfw.Command{
	Code: billCardCommandCode,
	CallbackAction: billCallbackAction(func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		if m.Text, err = getBillCardMessageText(ctx, whc.GetBotCode(), whc, bill, false, ""); err != nil {
			return
		}
		m.Format = botsfw.MessageFormatHTML
		m.Keyboard = getGroupBillCardInlineKeyboard(whc, bill)
		return
	}),
}

func startBillAction(whc botsfw.WebhookContext, billParam string) (m botsfw.MessageFromBot, err error) {
	var bill models4splitus.BillEntry
	if bill.ID = billParam[len("bill-"):]; bill.ID == "" {
		return m, errors.New("Invalid bill parameter")
	}
	if bill, err = facade4splitus.GetBillByID(whc.Context(), nil, bill.ID); err != nil {
		return
	}
	return ShowBillCard(whc, false, bill, "")
}

func billCardCallbackCommandData(billID string) string {
	return billCallbackCommandData(billCardCommandCode, billID)
}

const billMembersCommandCode = "bill-members"

func billCallbackCommandData(command string, billID string) string {
	return command + "?bill=" + billID
}

var billMembersCommand = billCallbackCommand(billMembersCommandCode,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		var buffer bytes.Buffer
		if err = writeBillCardTitle(whc.Context(), bill, whc.GetBotCode(), &buffer, whc); err != nil {
			return
		}
		buffer.WriteString("\n\n")
		writeBillMembersList(whc.Context(), &buffer, whc, bill, "")
		m.Text = buffer.String()
		m.Format = botsfw.MessageFormatHTML

		m.Keyboard = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					{
						Text:         whc.Translate(trans.BUTTON_TEXT_JOIN),
						CallbackData: billCallbackCommandData(joinBillCommandCode, bill.ID),
					},
				},
				{
					{
						Text:         whc.Translate(trans.COMMAND_TEXT_INVITE_MEMBER),
						CallbackData: billCallbackCommandData(INVITE_BILL_MEMBER_COMMAND, bill.ID),
					},
				},
				{
					{
						Text:         whc.Translate(emoji.RETURN_BACK_ICON),
						CallbackData: billCardCallbackCommandData(bill.ID),
					},
				},
			},
		}
		return
	},
)

func writeBillMembersList(
	ctx context.Context,
	buffer *bytes.Buffer,
	translator i18n.SingleLocaleTranslator,
	bill models4splitus.BillEntry,
	selectedMemberID string,
) {
	billCurrency := money.CurrencyCode(bill.Data.Currency)
	type MemberRowParams struct {
		N          int
		MemberName string
		Percent    decimal.Decimal64p2
		Owes       money.Amount
		Paid       money.Amount
	}
	billMembers := bill.Data.GetBillMembers()

	totalShares := 0

	for _, member := range billMembers {
		totalShares += member.Shares
	}

	for i, member := range bill.Data.GetBillMembers() {
		templateParams := MemberRowParams{
			N:          i + 1,
			MemberName: member.Name,
			Owes:       money.NewAmount(billCurrency, member.Owes),
			Paid:       money.NewAmount(billCurrency, member.Paid),
		}
		if totalShares == 0 {
			templateParams.Percent = decimal.Decimal64p2(1 * 100 / len(billMembers))
		} else {
			templateParams.Percent = decimal.Decimal64p2(member.Shares * 100 * 100 / totalShares)
		}

		var (
			templateName string
			err          error
		)
		if member.Paid == bill.Data.AmountTotal {
			buffer.WriteString("<b>")
		}
		if err = common4debtus.HtmlTemplates.RenderTemplate(ctx, buffer, translator, trans.MESSAGE_TEXT_BILL_CARD_MEMBER_TITLE, templateParams); err != nil {
			logus.Errorf(ctx, "Failed to render template")
			return
		}
		if member.Paid == bill.Data.AmountTotal {
			buffer.WriteString("</b>")
		}

		if selectedMemberID == "" {
			switch {
			case member.Owes > 0 && member.Paid > 0:
				templateName = trans.MESSAGE_TEXT_BILL_CARD_MEMBERS_ROW_PART_PAID
			case member.Owes > 0:
				templateName = trans.MESSAGE_TEXT_BILL_CARD_MEMBERS_ROW_OWES
			case member.Paid > 0:
				templateName = trans.MESSAGE_TEXT_BILL_CARD_MEMBERS_ROW_PAID
			default:
				templateName = trans.MESSAGE_TEXT_BILL_CARD_MEMBERS_ROW
			}
		} else {
			templateName = trans.MESSAGE_TEXT_BILL_CARD_MEMBERS_ROW
		}

		logus.Debugf(ctx, "Will render template")
		buffer.WriteString(" ")
		if err = common4debtus.HtmlTemplates.RenderTemplate(ctx, buffer, translator, templateName, templateParams); err != nil {
			logus.Errorf(ctx, "Failed to render template")
			return
		}
		buffer.WriteString("\n\n")
	}
}

const INVITE_BILL_MEMBER_COMMAND = "invite2bill"

const INLINE_COMMAND_JOIN = "join"

var inviteToBillCommand = billCallbackCommand(INVITE_BILL_MEMBER_COMMAND,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		m.Keyboard = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
						"via Telegram",
						INLINE_COMMAND_JOIN+"?bill="+bill.ID,
					),
				},
				{
					{
						Text:         whc.Translate(emoji.RETURN_BACK_ICON),
						CallbackData: billCardCallbackCommandData(bill.ID),
					},
				},
			},
		}
		return
	},
)

func ShowBillCard(whc botsfw.WebhookContext, isEdit bool, bill models4splitus.BillEntry, footer string) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	m = whc.NewMessage("")
	m.IsEdit = isEdit
	if m.Text, err = getBillCardMessageText(ctx, whc.GetBotCode(), whc, bill, true, footer); err != nil {
		return
	}
	var isInGroup bool
	if isInGroup, err = whc.IsInGroup(); err != nil {
		return
	} else if isInGroup || whc.Chat() == nil {
		m.Keyboard = getGroupBillCardInlineKeyboard(whc, bill)
	} else {
		m.Keyboard = getPrivateBillCardInlineKeyboard(whc, whc.GetBotCode(), bill)
	}
	return
}

func writeBillCardTitle(ctx context.Context, bill models4splitus.BillEntry, botID string, buffer *bytes.Buffer, translator i18n.SingleLocaleTranslator) error {
	var amount interface{}
	if bill.Data.Currency == "" {
		amount = bill.Data.AmountTotal
	} else {
		amount = bill.Data.TotalAmount()
	}
	titleWithLink := fmt.Sprintf(`<a href="https://t.me/%v?start=bill-%v">%v</a>`, botID, bill.ID, bill.Data.Name)
	logus.Debugf(ctx, "titleWithLink: %v", titleWithLink)
	header := translator.Translate(trans.MESSAGE_TEXT_BILL_CARD_HEADER, amount, titleWithLink)
	logus.Debugf(ctx, "header: %v", header)
	if _, err := buffer.WriteString(header); err != nil {
		logus.Errorf(ctx, "Failed to write bill header")
		return err
	}
	return nil
}

func getBillCardMessageText(ctx context.Context, botID string, translator i18n.SingleLocaleTranslator, bill models4splitus.BillEntry, showMembers bool, footer string) (string, error) {
	logus.Debugf(ctx, "getBillCardMessageText() => bill.BillDbo: %v", bill.Data)

	var buffer bytes.Buffer
	logus.Debugf(ctx, "Will write bill header...")

	if err := writeBillCardTitle(ctx, bill, botID, &buffer, translator); err != nil {
		return "", err
	}
	//buffer.WriteString("\n" + strings.Repeat("―", 15))

	buffer.WriteString("\n" + translator.Translate(trans.MT_TEXT_MEMBERS_COUNT, len(bill.Data.Members)))

	if showMembers {
		//buffer.WriteString("\n")
		//buffer.WriteString(translator.Translate(trans.MESSAGE_TEXT_SPLIT_LABEL_WITH_VALUE, translator.Translate(string(bill.SplitMode))))
		//if bill.Status != models.BillStatusOutstanding {
		//	buffer.WriteString(", " + translator.Translate(trans.MESSAGE_TEXT_STATUS, bill.Status))
		//}
		//buffer.WriteString(fmt.Sprintf("\n\n<b>%v</b> (%d)\n\n", translator.Translate(trans.MESSAGE_TEXT_MEMBERS_TITLE), bill.MembersCount))
		buffer.WriteString("\n\n")
		writeBillMembersList(ctx, &buffer, translator, bill, "")
	}

	if footer != "" {
		if !showMembers || len(bill.Data.Members) == 0 {
			buffer.WriteString("\n\n")
		}
		buffer.WriteString(footer)
	}
	logus.Debugf(ctx, "getBillCardMessageText() completed")
	return buffer.String(), nil
}
