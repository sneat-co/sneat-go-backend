package common4all

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"net/url"
)

const (
	UTM_MEDIUM_BOT                     = "bot"
	UTM_MEDIUM_SMS                     = "sms"
	UTM_CAMPAIGN_RECEIPT               = "receipt"
	UTM_CAMPAIGN_REMINDER              = "reminder"
	UTM_CAMPAIGN_RECEIPT_DISCARD       = "receipt-discard"
	UTM_CAMPAIGN_ONBOARDING_INVITE     = "oboarding-invite"
	UTM_CAMPAIGN_DEBT_CREATED          = "debt-created"
	UTM_CAMPAIGN_DEBT_RETURNED         = "debt-returned"
	UTM_CAMPAIGN_TRANSFER_SEND_RECEIPT = "transfer-send-receipt"
)

type UtmParams struct {
	Source   string // (referrer: google, citysearch, newsletter4), in our case it's BotID
	Medium   string // In our case it's bot platform e.g. 'telegram', 'fbm', etc.
	Campaign string // Identify where link is placed. For example 'receipt'
}

func NewUtmParams(whc botsfw.WebhookContext, campaign string) UtmParams {
	return FillUtmParams(whc, UtmParams{Campaign: campaign})
}

func UtmSourceFromContext(whc botsfw.WebhookContext) string {
	return whc.GetBotCode()
}

func FillUtmParams(whc botsfw.WebhookContext, utm UtmParams) UtmParams {
	if utm.Source == "" {
		utm.Source = whc.GetBotCode()
	}
	if utm.Medium == "" {
		utm.Medium = whc.BotPlatform().ID()
	}
	return utm
}

func (utm UtmParams) IsEmpty() bool {
	return utm.Source == "" && utm.Medium == "" && utm.Campaign == ""
}

func (utm UtmParams) String() string {
	switch "" {
	case utm.Source:
		panic("utm.Source is not provided")
	case utm.Medium:
		panic("utm.Medium is not provided")
	case utm.Campaign:
		panic("utm.Campaign is not provided")
	}
	return fmt.Sprintf("utm_source=%v&utm_medium=%v&utm_campaign=%v",
		url.QueryEscape(utm.Source), url.QueryEscape(utm.Medium), url.QueryEscape(utm.Campaign))
}

func (utm UtmParams) ShortString() string {
	switch "" {
	case utm.Source:
		panic("utm.Source is not provided")
	case utm.Medium:
		panic("utm.Medium is not provided")
	case utm.Campaign:
		panic("utm.Campaign is not provided")
	}
	return fmt.Sprintf("utm=%v;%v;%v",
		url.QueryEscape(utm.Source), url.QueryEscape(utm.Medium), url.QueryEscape(utm.Campaign))
}
