package common4debtus

import (
	"bytes"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/strongo/i18n"
	"io"
)

func GetCounterpartyUrl(counterpartyID string, currentUserID string, locale i18n.Locale, utmParams UtmParams) string {
	var buffer bytes.Buffer
	WriteCounterpartyUrl(&buffer, counterpartyID, currentUserID, locale, utmParams)
	return buffer.String()
}

func WriteCounterpartyUrl(writer io.Writer, counterpartyID string, currentUserID string, locale i18n.Locale, utmParams UtmParams) {
	host := GetWebsiteHost(utmParams.Source)
	_, _ = writer.Write([]byte(fmt.Sprintf("https://%v/counterparty?id=%v&lang=%v", host, counterpartyID, locale.SiteCode())))
	// TODO: Commented due to Telegram issue with too long URL
	if !utmParams.IsEmpty() {
		_, _ = writer.Write([]byte(fmt.Sprintf("&%v", utmParams.ShortString())))
	}
	if currentUserID != "" && currentUserID != "0" {
		token := token4auth.IssueBotToken(currentUserID, utmParams.Medium, utmParams.Source)
		_, _ = writer.Write([]byte(fmt.Sprintf("&secret=%v", token)))
	}
}
