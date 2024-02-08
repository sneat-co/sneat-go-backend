package dtb_transfer

import (
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp"

	//"fmt"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
)

func getTestMocks(t *testing.T, locale i18n.Locale) BalanceMessageBuilder {
	translator := i18n.NewMapTranslator(context.TODO(), trans.TRANS)
	singleLocaleTranslator := i18n.NewSingleMapTranslator(locale, translator)
	return NewBalanceMessageBuilder(singleLocaleTranslator)
}

func enMock(t *testing.T) BalanceMessageBuilder { return getTestMocks(t, i18n.LocaleEnUS) }
func ruMock(t *testing.T) BalanceMessageBuilder { return getTestMocks(t, i18n.LocaleRuRu) }

var (
	ruLinker = common.NewLinker(strongoapp.LocalHostEnv, "123", i18n.LocaleRuRu.Code5, "unit-Test")
	enLinker = common.NewLinker(strongoapp.LocalHostEnv, "123", i18n.LocaleEnUS.Code5, "unit-Test")
)

//type testBalanceDataProvider struct {
//}

func TestBalanceMessageSingleCounterparty(t *testing.T) {
	balanceJson := json.RawMessage(`{"USD": 10}`)
	counterparties := []models.UserContactJson{
		{
			ID:     "1",
			Name:   "John Doe",
			Status: "active",
			//UserID: 1,
			BalanceJson: &balanceJson,
		},
	}

	c := context.TODO()
	expectedEn := `<a href="https://debtstracker.local/contact?id=1&lang=en-US&secret=SECRET">John Doe</a>`
	expectedRu := `<a href="https://debtstracker.local/contact?id=1&lang=ru-RU&secret=SECRET">John Doe</a>`

	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 10 USD", expectedEn), enMock(t).ByContact(c, enLinker, counterparties))

	assert(t, i18n.LocaleRuRu, 0, fmt.Sprintf("%v - долг вам 10 USD", expectedRu), ruMock(t).ByContact(c, ruLinker, counterparties))

	balanceJson = json.RawMessage(`{"USD": -10}`)
	counterparties[0].BalanceJson = &balanceJson
	assert(t, i18n.LocaleRuRu, 0, fmt.Sprintf("%v - вы должны 10 USD", expectedRu), ruMock(t).ByContact(c, ruLinker, counterparties))

	balanceJson = json.RawMessage(`{"USD": 10, "EUR": 20}`)
	counterparties[0].BalanceJson = &balanceJson
	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 20 EUR and 10 USD", expectedEn), enMock(t).ByContact(c, enLinker, counterparties))

	balanceJson = json.RawMessage(`{"USD": 10, "EUR": 20, "RUB": 15}`)
	counterparties[0].BalanceJson = &balanceJson
	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 20 EUR, 15 RUB and 10 USD", expectedEn), enMock(t).ByContact(c, enLinker, counterparties))

}

func TestBalanceMessageTwoCounterparties(t *testing.T) {
	john := models.UserContactJson{
		ID:   "1",
		Name: "Johnny The Doe",
	}

	jack := models.UserContactJson{
		ID:   "2",
		Name: "Jacky Dark Brown",
	}

	c := context.TODO()

	johnLink := fmt.Sprintf(`<a href="https://debtstracker.local/contact?id=1&lang=en-US&secret=SECRET">%v</a>`, john.Name)
	jackLink := fmt.Sprintf(`<a href="https://debtstracker.local/contact?id=2&lang=en-US&secret=SECRET">%v</a>`, jack.Name)

	var johnBalance, jackBalance json.RawMessage
	johnBalance = json.RawMessage(`{"USD": 10}`)
	john.BalanceJson = &johnBalance
	jackBalance = json.RawMessage(`{"USD": 15}`)
	jack.BalanceJson = &jackBalance
	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 10 USD\n%v - owes you 15 USD", johnLink, jackLink), enMock(t).ByContact(c, enLinker, []models.UserContactJson{john, jack}))

	johnBalance = json.RawMessage(`{"USD": 10, "EUR": 20}`)
	john.BalanceJson = &johnBalance
	jackBalance = json.RawMessage(`{"USD": 40, "EUR": 15}`)
	jack.BalanceJson = &jackBalance
	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 20 EUR and 10 USD\n%v - owes you 40 USD and 15 EUR", johnLink, jackLink), enMock(t).ByContact(c, enLinker, []models.UserContactJson{john, jack}))

	johnBalance = json.RawMessage(`{"USD": 10, "EUR": 20, "RUB": 100}`)
	john.BalanceJson = &johnBalance
	jackBalance = json.RawMessage(`{"USD": 40, "EUR": 15}`)
	jack.BalanceJson = &jackBalance
	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 100 RUB, 20 EUR and 10 USD\n%v - owes you 40 USD and 15 EUR", johnLink, jackLink), enMock(t).ByContact(c, enLinker, []models.UserContactJson{john, jack}))

	johnBalance = json.RawMessage(`{"USD": -10}`)
	john.BalanceJson = &johnBalance
	jackBalance = json.RawMessage(`{"USD": -15}`)
	jack.BalanceJson = &jackBalance
	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - you owe 10 USD\n%v - you owe 15 USD", johnLink, jackLink), enMock(t).ByContact(c, enLinker, []models.UserContactJson{john, jack}))

	johnBalance = json.RawMessage(`{"USD": -10}`)
	john.BalanceJson = &johnBalance
	jackBalance = json.RawMessage(`{"USD": 15}`)
	jack.BalanceJson = &jackBalance
	assert(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - you owe 10 USD\n%v - owes you 15 USD", johnLink, jackLink), enMock(t).ByContact(c, enLinker, []models.UserContactJson{john, jack}))

}

func TestBalanceMessageBuilder_ByCurrency(t *testing.T) {
	balance := money.Balance{
		money.CurrencyUSD: decimal.NewDecimal64p2(10, 0),
		money.CurrencyRUB: decimal.NewDecimal64p2(50, 0),
		money.CurrencyEUR: decimal.NewDecimal64p2(15, 0),
	}
	assert(t, i18n.LocaleRuRu, 0, "<b>Всего</b>\nВам должны 50 RUB, 15 EUR и 10 USD", ruMock(t).ByCurrency(true, balance))
}

var reCleanSecret = regexp.MustCompile(`secret.+?"`)

func assert(t *testing.T, locale i18n.Locale, warningsCount int, expected, actual string) {
	actual = reCleanSecret.ReplaceAllString(actual, `secret=SECRET"`)
	if actual != expected {
		t.Errorf("Unexpected output for locale %v:\nExpected:\n%v\nActual:\n%v", locale.Code5, expected, actual)
	}
	//if len(log.Warnings) != warningsCount {
	//	t.Errorf("Unexpected warnings count: %v", log.Warnings)
	//}
}
