package common

import (
	"github.com/strongo/i18n"
	"regexp"
	"testing"
)

func TestGetCounterpartyUrl(t *testing.T) {
	var (
		utm UtmParams
	)
	counterpartyUrl := GetCounterpartyUrl("123", "", i18n.LocaleRuRu, utm)

	re := regexp.MustCompile(`^https://debtstracker\.io/counterparty\?id=\d+&lang=\w{2}$`)
	if !re.MatchString(counterpartyUrl) {
		t.Errorf("Unexpected counterpart URL:\n%v", counterpartyUrl)
	}
}
