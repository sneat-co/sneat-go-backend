package dtb_transfer

import (
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp"
	"github.com/strongo/strongoapp/person"

	"fmt"
	"regexp"
	"testing"

	"context"
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
	ruLinker = common4debtus.NewLinker(strongoapp.LocalHostEnv, "123", i18n.LocaleRuRu.Code5, "unit-Test")
	enLinker = common4debtus.NewLinker(strongoapp.LocalHostEnv, "123", i18n.LocaleEnUS.Code5, "unit-Test")
)

//type testBalanceDataProvider struct {
//}

func TestBalanceMessageSingleCounterparty(t *testing.T) {
	const johnID = "john1"
	debtusContacts := map[string]*models4debtus.DebtusContactBrief{
		johnID: {
			Status:  "active",
			Balance: money.Balance{"USD": 1000},
		},
	}

	contacts := map[string]*briefs4contactus.ContactBrief{
		johnID: {
			Names: &person.NameFields{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
	}

	c := context.TODO()
	expectedEn := `<a href="https://debtus.local/contact?id=john1&lang=en-US&secret=SECRET">John Doe</a>`
	expectedRu := `<a href="https://debtus.local/contact?id=john1&lang=ru-RU&secret=SECRET">John Doe</a>`

	var actual string

	mockEn := enMock(t)
	mockRu := ruMock(t)

	actual = mockEn.ByContact(c, enLinker, contacts, debtusContacts)
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 10 USD", expectedEn), actual)

	actual = mockRu.ByContact(c, ruLinker, contacts, debtusContacts)
	assertText(t, i18n.LocaleRuRu, 0, fmt.Sprintf("%v - долг вам 10 USD", expectedRu), actual)

	debtusContacts[johnID].Balance = money.Balance{"USD": -1000}
	actual = mockRu.ByContact(c, ruLinker, contacts, debtusContacts)
	assertText(t, i18n.LocaleRuRu, 0, fmt.Sprintf("%v - вы должны 10 USD", expectedRu), actual)

	debtusContacts[johnID].Balance = money.Balance{"USD": 1000, "EUR": 2000}
	actual = mockEn.ByContact(c, enLinker, contacts, debtusContacts)
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 20 EUR and 10 USD", expectedEn), actual)

	debtusContacts[johnID].Balance = money.Balance{"USD": 1000, "EUR": 2000, "RUB": 1500}
	actual = mockEn.ByContact(c, enLinker, contacts, debtusContacts)
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 20 EUR, 15 RUB and 10 USD", expectedEn), actual)

}

func TestBalanceMessageTwoCounterparties(t *testing.T) {
	const johnID = "john1"
	const jackID = "jack2"

	//const spaceID = "s3"

	johnContact := briefs4contactus.ContactBrief{
		Names: &person.NameFields{
			FirstName: "John",
		},
	}
	var johnDebtusContactBrief, jackDebtusContactBrief models4debtus.DebtusContactBrief

	jackContact := briefs4contactus.ContactBrief{
		Names: &person.NameFields{
			FirstName: "Jack",
		},
	}

	contacBriefs := map[string]*briefs4contactus.ContactBrief{
		johnID: &johnContact,
		jackID: &jackContact,
	}

	c := context.TODO()

	johnLink := fmt.Sprintf(`<a href="https://debtus.local/contact?id=john1&lang=en-US&secret=SECRET">%v</a>`, johnContact.Names.GetFullName())
	jackLink := fmt.Sprintf(`<a href="https://debtus.local/contact?id=jack2&lang=en-US&secret=SECRET">%v</a>`, jackContact.Names.GetFullName())

	johnDebtusContactBrief.Balance = money.Balance{"USD": 1000}
	jackDebtusContactBrief.Balance = money.Balance{"USD": 1500}
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 15 USD\n%v - owes you 10 USD", jackLink, johnLink),
		enMock(t).ByContact(c, enLinker,
			contacBriefs,
			map[string]*models4debtus.DebtusContactBrief{
				johnID: &johnDebtusContactBrief,
				jackID: &jackDebtusContactBrief,
			}))

	johnDebtusContactBrief.Balance = money.Balance{"USD": 1000, "EUR": 2000}
	jackDebtusContactBrief.Balance = money.Balance{"USD": 4000, "EUR": 1500}
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 40 USD and 15 EUR\n%v - owes you 20 EUR and 10 USD", jackLink, johnLink),
		enMock(t).ByContact(c, enLinker, contacBriefs,
			map[string]*models4debtus.DebtusContactBrief{
				johnID: &johnDebtusContactBrief,
				jackID: &jackDebtusContactBrief,
			}))

	johnDebtusContactBrief.Balance = money.Balance{"USD": 1000, "EUR": 2000, "RUB": 10000}
	jackDebtusContactBrief.Balance = money.Balance{"USD": 4000, "EUR": 1500}
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 100 RUB, 20 EUR and 10 USD\n%v - owes you 40 USD and 15 EUR", johnLink, jackLink),
		enMock(t).ByContact(c, enLinker,
			contacBriefs,
			map[string]*models4debtus.DebtusContactBrief{
				johnID: &johnDebtusContactBrief,
				jackID: &jackDebtusContactBrief,
			}))

	johnDebtusContactBrief.Balance = money.Balance{"USD": -1000}
	jackDebtusContactBrief.Balance = money.Balance{"USD": -1500}
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - you owe 10 USD\n%v - you owe 15 USD", johnLink, jackLink),
		enMock(t).ByContact(c, enLinker,
			contacBriefs,
			map[string]*models4debtus.DebtusContactBrief{
				johnID: &johnDebtusContactBrief,
				jackID: &jackDebtusContactBrief,
			}))

	johnDebtusContactBrief.Balance = money.Balance{"USD": -1000}
	jackDebtusContactBrief.Balance = money.Balance{"USD": 1500}
	assertText(t, i18n.LocaleEnUS, 0, fmt.Sprintf("%v - owes you 15 USD\n%v - you owe 10 USD", jackLink, johnLink),
		enMock(t).ByContact(c, enLinker, contacBriefs, map[string]*models4debtus.DebtusContactBrief{
			johnID: &johnDebtusContactBrief,
			jackID: &jackDebtusContactBrief,
		}))

}

func TestBalanceMessageBuilder_ByCurrency(t *testing.T) {
	balance := money.Balance{
		money.CurrencyUSD: decimal.NewDecimal64p2(10, 0),
		money.CurrencyRUB: decimal.NewDecimal64p2(50, 0),
		money.CurrencyEUR: decimal.NewDecimal64p2(15, 0),
	}
	assertText(t, i18n.LocaleRuRu, 0, "<b>Всего</b>\nВам должны 50 RUB, 15 EUR и 10 USD", ruMock(t).ByCurrency(true, balance))
}

var reCleanSecret = regexp.MustCompile(`secret.+?"`)

func assertText(t *testing.T, locale i18n.Locale, warningsCount int, expected, actual string) {
	t.Helper()
	actual = reCleanSecret.ReplaceAllString(actual, `secret=SECRET"`)
	if actual != expected {
		t.Errorf("Unexpected output for locale %v:\nExpected:\n%v\nActual:\n%v", locale.Code5, expected, actual)
	}
	//if len(logus.Warnings) != warningsCount {
	//	t.Errorf("Unexpected warnings count: %v", logus.Warnings)
	//}
}
