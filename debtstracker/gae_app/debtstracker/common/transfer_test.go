package common

import (
	"github.com/strongo/i18n"
	"regexp"
	"testing"
)

func TestGetTransferUrlForUser(t *testing.T) {

	var (
		transferUrl string
		re          *regexp.Regexp
		utm         UtmParams
	)
	{
		transferUrl = GetTransferUrlForUser("123", "", i18n.LocaleRuRu, utm)

		re = regexp.MustCompile(`^https://(\w+\.\w+)/transfer\?id=\d+&lang=ru`)
		if !re.MatchString(transferUrl) {
			t.Errorf("Unexpected transfer URL:\n%v", transferUrl)
			//} else {
			//	t.Logf("Transfer URL: %v", transferUrl)
		}
	}

	utm = UtmParams{
		Source:   "S1",
		Medium:   "M1",
		Campaign: "C1",
	}

	{
		transferUrl = GetTransferUrlForUser("123", "", i18n.LocaleRuRu, utm)

		re = regexp.MustCompile(`^https://(\w+\.\w+)/transfer\?id=\d+&lang=ru&utm=S1;M1;C1`)
		if !re.MatchString(transferUrl) {
			t.Errorf("Unexpected transfer URL: %v", transferUrl)
		} else {
			t.Logf("Transfer URL: %v", transferUrl)
		}
	}

	{
		transferUrl = GetTransferUrlForUser("123", "234", i18n.LocaleRuRu, utm)

		re = regexp.MustCompile(`^https://(\w+\.\w+)/transfer\?id=\d+&lang=ru&utm=S1;M1;C1&secret=[\-\.\w]+$`)
		if !re.MatchString(transferUrl) {
			t.Errorf("Unexpected transfer URL:\n%v", transferUrl)
		} else {
			t.Logf("Transfer URL: %v", transferUrl)
		}
	}
}
