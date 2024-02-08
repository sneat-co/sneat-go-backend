package common

import (
	"bytes"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
	"regexp"
	"testing"

	"context"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func TestWriteReceiptText(t *testing.T) {
	var (
		buffer bytes.Buffer
	)

	//c, done, err := aetest.NewContext()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer done()

	c := context.TODO()

	//logger := &bots.MockLogger{T: t}

	translator := i18n.NewSingleMapTranslator(i18n.LocaleEnUS, i18n.NewMapTranslator(c, trans.TRANS))

	transfer := models.NewTransfer("123", models.NewTransferData(
		"12",
		false,
		money.Amount{Currency: "EUR", Value: 98765},
		&models.TransferCounterpartyInfo{
			ContactID:   "23",
			ContactName: "John Whites",
		},
		&models.TransferCounterpartyInfo{
			UserID:   "12",
			UserName: "Anna Blacks",
		},
	))

	receiptTextBuilder := newReceiptTextBuilder(translator, transfer, ShowReceiptToCounterparty)

	utmParams := UtmParams{
		Source:   "BotIdUnitTest",
		Medium:   telegram.PlatformID,
		Campaign: "unit-test-campaign",
	}

	receiptTextBuilder.WriteReceiptText(&buffer, utmParams)

	re := regexp.MustCompile(`Anna Blacks borrowed from you <b>987.65 EUR</b>.`)
	if matched := re.MatchString(buffer.String()); !matched {
		t.Errorf("Unexpected output:\nOutput:\n%v\nRegex:\n%v", buffer.String(), re.String())
	}
}
