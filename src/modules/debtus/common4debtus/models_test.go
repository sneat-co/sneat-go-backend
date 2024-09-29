package common4debtus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/strongo/i18n"
	"regexp"
	"testing"
)

func TestGetCounterpartyUrl(t *testing.T) {
	var (
		utm common4all.UtmParams
	)
	counterpartyUrl, _ := GetCounterpartyUrl(context.Background(), "123", "", i18n.LocaleRuRu, utm)

	re := regexp.MustCompile(`^https://debtusbot\.io/counterparty\?id=\d+&lang=\w{2}$`)
	if !re.MatchString(counterpartyUrl) {
		t.Errorf("Unexpected counterpart URL:\n%v", counterpartyUrl)
	}
}
