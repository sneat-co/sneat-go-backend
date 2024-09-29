package common4debtus

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/strongo/i18n"
	"io"
)

func GetCounterpartyUrl(ctx context.Context, counterpartyID string, currentUserID string, locale i18n.Locale, utmParams common4all.UtmParams) (string, error) {
	var buffer bytes.Buffer
	if err := WriteCounterpartyUrl(ctx, &buffer, counterpartyID, currentUserID, locale, utmParams); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func WriteCounterpartyUrl(
	ctx context.Context, writer io.Writer, counterpartyID string, currentUserID string, locale i18n.Locale, utmParams common4all.UtmParams,
) (err error) {
	host := GetWebsiteHost(utmParams.Source)
	_, _ = writer.Write([]byte(fmt.Sprintf("https://%v/counterparty?id=%v&lang=%v", host, counterpartyID, locale.SiteCode())))
	// TODO: Commented due to Telegram issue with too long URL
	if !utmParams.IsEmpty() {
		_, _ = writer.Write([]byte(fmt.Sprintf("&%v", utmParams.ShortString())))
	}
	if currentUserID != "" && currentUserID != "0" {
		var token string

		if token, err = token4auth.IssueBotToken(ctx, currentUserID, utmParams.Medium, utmParams.Source); err != nil {
			return
		}
		_, _ = writer.Write([]byte(fmt.Sprintf("&secret=%v", token)))
	}
	return err
}
