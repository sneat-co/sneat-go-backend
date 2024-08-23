package common4debtus

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/i18n"
	"io"
	"strconv"
	"strings"

	"html/template"
)

func GetBalanceUrlForUser(userID int64, locale i18n.Locale, createdOnPlatform, createdOnID string) string {
	return getUrlForUser(userID, locale, "debts", createdOnPlatform, createdOnID)
}

func GetHistoryUrlForUser(userID int64, locale i18n.Locale, createdOnPlatform, createdOnID string) string {
	return getUrlForUser(userID, locale, "history", createdOnPlatform, createdOnID)
}

func getUrlForUser(userID int64, locale i18n.Locale, page, createdOnPlatform, createdOnID string) string {
	token := token4auth.IssueBotToken(strconv.FormatInt(userID, 10), createdOnPlatform, createdOnID)
	host := GetWebsiteHost(createdOnID)
	url := fmt.Sprintf("https://%v/app/#", host)
	switch page {
	case "history":
		url += "/history/"
	case "debts":
		url += "/debts/"
	default:
		url += "page=" + page
	}
	return url + fmt.Sprintf("&lang=%v&secret=%v", locale.SiteCode(), token)
}

func GetTransferUrlForUser(transferID string, userID string, locale i18n.Locale, utmParams UtmParams) string {
	var buffer bytes.Buffer
	WriteTransferUrlForUser(&buffer, transferID, userID, locale, utmParams)
	return buffer.String()
}

func WriteTransferUrlForUser(writer io.Writer, transferID string, userID string, locale i18n.Locale, utmParams UtmParams) {
	host := GetWebsiteHost(utmParams.Source)
	_, _ = writer.Write([]byte(fmt.Sprintf(
		"https://%v/transfer?id=%v&lang=%v",
		host, transferID, locale.SiteCode(),
	)))
	if !utmParams.IsEmpty() {
		_, _ = writer.Write([]byte(fmt.Sprintf("&%v", utmParams.ShortString())))
	}
	if userID != "" {
		token := token4auth.IssueBotToken(userID, utmParams.Medium, utmParams.Source)
		_, _ = writer.Write([]byte(fmt.Sprintf("&secret=%v", token)))
	}
}

func GetChooseCurrencyUrlForUser(userID string, locale i18n.Locale, createdOnPlatform, createdOnID, contextData string) string {
	token := token4auth.IssueBotToken(userID, createdOnPlatform, createdOnID)
	host := GetWebsiteHost(createdOnID)
	return fmt.Sprintf(
		"https://%v/app/#/choose-currency?lang=%v&%v&secret=%v",
		host, locale.SiteCode(), contextData, token,
	)
}

func GetWebsiteHost(createdOnID string) string {
	createdOnID = strings.ToLower(createdOnID)
	if strings.Contains(createdOnID, "dev") {
		return "debtusbot-dev1.appspot.com"
	} else if strings.Contains(createdOnID, ".local") {
		return "debtusbot.local"
	} else {
		return "debtusbot.io"
	}
}

func GetPathAndQueryForInvite(inviteCode string, utmParams UtmParams) string {
	return fmt.Sprintf("ack?invite=%v#%v", template.URLQueryEscaper(inviteCode), utmParams)
}

func GetNewDebtPageUrl(whc botsfw.WebhookContext, direction models4debtus.TransferDirection, utmCampaign string) string {
	botID := whc.GetBotCode()
	botPlatform := whc.BotPlatform().ID()
	token := token4auth.IssueToken(whc.AppUserID(), token4auth.GetBotIssuer(botPlatform, botID))
	host := GetWebsiteHost(botID)
	// utmParams := NewUtmParams(whc, utmCampaign)
	return fmt.Sprintf(
		"https://%v/open/new-debt?d=%v&lang=%v&secret=%v",
		host, direction, whc.Locale().SiteCode(), token, // utmParams, - commented out as: Viber response.Status=3: keyboard is not valid. is too long (length: 274, maximum allowed: 250)]
	)
}
