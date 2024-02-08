package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
	"time"

	"context"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
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
func (m BalanceMessageBuilder) ByContact(c context.Context, linker common.Linker, userContactJsons []models.UserContactJson) string {
	var buffer bytes.Buffer
	translator := m.translator

	getContactName := func(userContactJson models.UserContactJson) string {
		return fmt.Sprintf(`<a href="%v">%v</a>`, linker.UrlToContact(userContactJson.ID), html.EscapeString(userContactJson.Name))
	}

	writeBalanceRow := func(userContactJson models.UserContactJson, b money.Balance, msg string) {
		if len(b) > 0 {
			amounts := b.CommaSeparatedUnsignedWithSymbols(simpleTranslator{t: translator})
			msg = m.translator.Translate(msg)
			name := getContactName(userContactJson)
			buffer.WriteString(fmt.Sprintf(msg, name, amounts) + "\n")
		}
	}

	writeBalanceErrorRow := func(userContactJson models.UserContactJson, err error) {
		buffer.WriteString(getContactName(userContactJson))
		buffer.WriteString(" - " + emoji.ERROR_ICON + " ERROR: " + err.Error() + "\n")
	}

	var (
		counterpartiesWithZeroBalance      bytes.Buffer
		counterpartiesWithZeroBalanceCount int
	)

	now := time.Now()

	for _, userContactJson := range userContactJsons {
		counterpartyBalanceWithInterest, err := userContactJson.BalanceWithInterest(c, now)
		if err != nil {
			log.Errorf(c, "Failed to get userContactJson balance with interest for contact %v: %v", userContactJson.ID, err)
			writeBalanceErrorRow(userContactJson, err)
			continue
		}
		//counterpartyBalance := userContactJson.Balance()
		//log.Debugf(c, "counterpartyBalanceWithInterest: %v\ncounterpartyBalance: %v", counterpartyBalanceWithInterest, counterpartyBalance)
		if counterpartyBalanceWithInterest.IsZero() {
			counterpartiesWithZeroBalanceCount += 1
			counterpartiesWithZeroBalance.WriteString(userContactJson.ID)
			counterpartiesWithZeroBalance.WriteString(", ")
			continue
		}
		writeBalanceRow(userContactJson, counterpartyBalanceWithInterest.OnlyPositive(), trans.MESSAGE_TEXT_BALANCE_SINGLE_CURRENCY_COUNTERPARTY_DEBT_TO_USER)
		writeBalanceRow(userContactJson, counterpartyBalanceWithInterest.OnlyNegative(), trans.MESSAGE_TEXT_BALANCE_SINGLE_CURRENCY_COUNTERPARTY_DEBT_BY_USER)
	}
	//if counterpartiesWithZeroBalanceCount > 0 {
	//	log.Debugf(c, "There are %d userContactJsons with zero balance: %v", counterpartiesWithZeroBalanceCount, strings.TrimRight(counterpartiesWithZeroBalance.String(), ", "))
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
	return "\n" + header + common.HORIZONTAL_LINE + balanceMessageBuilder.ByCurrency(false, b) + common.HORIZONTAL_LINE
}
