package common4debtus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/strongo/i18n"
	"regexp"
	"testing"
)

func TestGetTransferUrlForUser(t *testing.T) {

	backupIssueBotToken := token4auth.IssueBotToken
	defer func() {
		token4auth.IssueBotToken = backupIssueBotToken
	}()
	token4auth.IssueBotToken = func(ctx context.Context, userID, createdOnPlatform, createdOnID string) (string, error) {
		return "SECRET_TOKEN", nil
	}

	var (
		transferUrl string
		re          *regexp.Regexp
		utm         common4all.UtmParams
	)
	{
		transferUrl = GetTransferUrlForUser(context.Background(), "123", "", i18n.LocaleRuRu, utm)

		re = regexp.MustCompile(`^https://(\w+\.\w+)/transfer\?id=\d+&lang=ru`)
		if !re.MatchString(transferUrl) {
			t.Errorf("Unexpected transfer URL:\n%v", transferUrl)
			//} else {
			//	t.Logf("Transfer URL: %v", transferUrl)
		}
	}

	utm = common4all.UtmParams{
		Source:   "S1",
		Medium:   "M1",
		Campaign: "C1",
	}

	{
		transferUrl = GetTransferUrlForUser(context.Background(), "123", "", i18n.LocaleRuRu, utm)

		re = regexp.MustCompile(`^https://(\w+\.\w+)/transfer\?id=\d+&lang=ru&utm=S1;M1;C1`)
		if !re.MatchString(transferUrl) {
			t.Errorf("Unexpected transfer URL: %v", transferUrl)
		} else {
			t.Logf("Transfer URL: %v", transferUrl)
		}
	}

	{
		transferUrl = GetTransferUrlForUser(context.Background(), "123", "234", i18n.LocaleRuRu, utm)

		re = regexp.MustCompile(`^https://(\w+\.\w+)/transfer\?id=\d+&lang=ru&utm=S1;M1;C1&secret=[\-\.\w]+$`)
		if !re.MatchString(transferUrl) {
			t.Errorf("Unexpected transfer URL for user:\n%v", transferUrl)
		} else {
			t.Logf("Transfer URL: %v", transferUrl)
		}
	}
}
