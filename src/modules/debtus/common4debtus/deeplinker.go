package common4debtus

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/strongo/strongoapp"
)

type deeplink struct {
}

func (deeplink) AppHashPathToReceipt(receiptID string) string {
	return fmt.Sprintf("receipt=%s", receiptID)
}

var Deeplink = deeplink{}

type Linker struct {
	userID string
	locale string
	issuer string
	host   string
}

func NewLinker(environment string, userID string, locale, issuer string) Linker {
	return Linker{
		userID: userID,
		locale: locale,
		issuer: issuer,
		host:   host(environment),
	}
}

func NewLinkerFromWhc(whc botsfw.WebhookContext) Linker {
	return NewLinker(whc.Environment(), whc.AppUserID(), whc.Locale().SiteCode(), formatIssuer(whc.BotPlatform().ID(), whc.GetBotCode()))
}

func host(environment string) string {
	switch environment {
	case "prod":
		return "debtusbot.io"
	case strongoapp.LocalHostEnv:
		return "debtusbot.local"
	case "dev":
		return "debtusbot-dev1.appspot.com"
	}
	panic(fmt.Sprintf("Unknown environment: %v", environment))
}

func (l Linker) UrlToContact(contactID string) string {
	return l.url("/contact", fmt.Sprintf("?id=%s", contactID), "")
}

func formatIssuer(botPlatform, botID string) string {
	return fmt.Sprintf("%s:%s", botPlatform, botID)
}

func (l Linker) url(path, query, hash string) string {
	var buffer bytes.Buffer
	buffer.WriteString("https://" + l.host + path + query)
	if hash != "" {
		buffer.WriteString(hash)
	}
	if query != "" || hash != "" {
		buffer.WriteString("&")
	}
	isAdmin := false // TODO: How to get isAdmin?
	token := auth.IssueToken(l.userID, l.issuer, isAdmin)
	buffer.WriteString("lang=" + l.locale)
	buffer.WriteString("&secret=" + token)
	return buffer.String()
}

func (l Linker) ToMainScreen(_ botsfw.WebhookContext) string {
	return l.url("/app/", "", "#")
}
