package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/decimal"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"sort"
	"time"

	"context"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"golang.org/x/net/html"
)

type BalanceMessageBuilder struct {
	translator i18n.SingleLocaleTranslator
	NeedsTotal bool
}

func NewBalanceMessageBuilder(translator i18n.SingleLocaleTranslator) BalanceMessageBuilder {
	return BalanceMessageBuilder{translator: translator}
}

type simpleTranslator struct {
	t i18n.SingleLocaleTranslator
}

func (t simpleTranslator) Translate(key string) string {
	return t.t.Translate(key)
}

func (m BalanceMessageBuilder) ByContact(
	ctx context.Context,
	linker common4debtus.Linker,
	contactBriefs map[string]*briefs4contactus.ContactBrief,
	debtusContactBriefs map[string]*models4debtus.DebtusContactBrief,
) string {
	var buffer bytes.Buffer
	translator := m.translator

	getContactName := func(contactID string, contactBrief briefs4contactus.ContactBrief) string {
		return fmt.Sprintf(`<a href="%v">%v</a>`, linker.UrlToContact(contactID), html.EscapeString(contactBrief.Names.GetFullName()))
	}

	writeBalanceRow := func(contactID string, contactBrief briefs4contactus.ContactBrief, b money.Balance, msg string) {
		if len(b) > 0 {
			amounts := b.CommaSeparatedUnsignedWithSymbols(simpleTranslator{t: translator})
			msg = m.translator.Translate(msg)
			name := getContactName(contactID, contactBrief)
			buffer.WriteString(fmt.Sprintf(msg, name, amounts) + "\n")
		}
	}

	writeBalanceErrorRow := func(contactID string, contactBrief briefs4contactus.ContactBrief, err error) {
		buffer.WriteString(getContactName(contactID, contactBrief))
		buffer.WriteString(" - " + emoji.ERROR_ICON + " ERROR: " + err.Error() + "\n")
	}

	var (
		counterpartiesWithZeroBalance      bytes.Buffer
		counterpartiesWithZeroBalanceCount int
	)

	now := time.Now()

	sortedContactBriefs := make([]models4debtus.DebtusContactBriefWithContactID, 0, len(debtusContactBriefs))

	for contactID, debtusContactBrief := range debtusContactBriefs {
		counterpartyBalanceWithInterest, err := debtusContactBrief.BalanceWithInterest(ctx, now)
		if err != nil {
			logus.Errorf(ctx, "Failed to get debtusContactBrief balance with interest for contact %v: %v", contactID, err)
			contactBrief := contactBriefs[contactID]
			writeBalanceErrorRow(contactID, *contactBrief, err)
			continue
		}
		//counterpartyBalance := debtusContactBrief.Balance()
		//logus.Debugf(ctx, "counterpartyBalanceWithInterest: %v\ncounterpartyBalance: %v", counterpartyBalanceWithInterest, counterpartyBalance)
		if counterpartyBalanceWithInterest.IsZero() {
			counterpartiesWithZeroBalanceCount += 1
			counterpartiesWithZeroBalance.WriteString(contactID)
			counterpartiesWithZeroBalance.WriteString(", ")
			continue
		}
		debtusContactBriefWithContactID := models4debtus.DebtusContactBriefWithContactID{ContactID: contactID, DebtusContactBrief: *debtusContactBrief}
		debtusContactBriefWithContactID.Balance = counterpartyBalanceWithInterest
		sortedContactBriefs = append(sortedContactBriefs, debtusContactBriefWithContactID)
	}

	sort.Slice(sortedContactBriefs, func(i, j int) bool {
		b1 := sortedContactBriefs[i].DebtusContactBrief.Balance
		b2 := sortedContactBriefs[j].DebtusContactBrief.Balance
		var a1, a2 decimal.Decimal64p2
		for _, v := range b1 {
			a1 += v
		}
		for _, v := range b2 {
			a2 += v
		}
		return a1 > a2
	})
	for _, debtusContactBriefWithContactID := range sortedContactBriefs {
		contactID := debtusContactBriefWithContactID.ContactID
		debtusContactBrief := debtusContactBriefWithContactID.DebtusContactBrief
		contactBrief := contactBriefs[contactID]
		writeBalanceRow(contactID, *contactBrief, debtusContactBrief.Balance.OnlyPositive(), trans.MESSAGE_TEXT_BALANCE_SINGLE_CURRENCY_COUNTERPARTY_DEBT_TO_USER)
		writeBalanceRow(contactID, *contactBrief, debtusContactBrief.Balance.OnlyNegative(), trans.MESSAGE_TEXT_BALANCE_SINGLE_CURRENCY_COUNTERPARTY_DEBT_BY_USER)

	}

	//if counterpartiesWithZeroBalanceCount > 0 {
	//	logus.Debugf(ctx, "There are %d debtusContactBriefs with zero balance: %v", counterpartiesWithZeroBalanceCount, strings.TrimRight(counterpartiesWithZeroBalance.String(), ", "))
	//}
	if l := buffer.Len() - 1; l > 0 {
		buffer.Truncate(l)
	}
	return buffer.String()
}

func (m BalanceMessageBuilder) ByCurrency(isTotal bool, balance money.Balance) string {
	var buffer bytes.Buffer
	translator := m.translator
	if isTotal {
		buffer.WriteString("<b>" + translator.Translate(trans.MESSAGE_TEXT_BALANCE_CURRENCY_TOTAL_INTRO) + "</b>\n")
	}
	debtByUser := balance.OnlyNegative()
	debtToUser := balance.OnlyPositive()
	commaSeparatedAmounts := func(prefix string, owed money.Balance) {
		if !owed.IsZero() {
			buffer.WriteString(fmt.Sprintf(translator.Translate(prefix), owed.CommaSeparatedUnsignedWithSymbols(simpleTranslator{t: translator})) + "\n")
		}
	}
	commaSeparatedAmounts(trans.MESSAGE_TEXT_BALANCE_CURRENCY_ROW_DEBT_BY_USER, debtByUser)
	commaSeparatedAmounts(trans.MESSAGE_TEXT_BALANCE_CURRENCY_ROW_DEBT_TO_USER, debtToUser)

	if l := buffer.Len() - 1; l > 0 {
		buffer.Truncate(l)
	}
	return buffer.String()
}

func BalanceForCounterpartyWithHeader(counterpartyLink string, b money.Balance, translator i18n.SingleLocaleTranslator) string {
	balanceMessageBuilder := NewBalanceMessageBuilder(translator)
	header := fmt.Sprintf("<b>%v</b>: %v", translator.Translate(trans.MESSAGE_TEXT_BALANCE_HEADER), counterpartyLink)
	return "\n" + header + common4debtus.HORIZONTAL_LINE + balanceMessageBuilder.ByCurrency(false, b) + common4debtus.HORIZONTAL_LINE
}
