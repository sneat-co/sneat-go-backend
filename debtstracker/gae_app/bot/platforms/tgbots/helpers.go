package tgbots

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"regexp"
	"strings"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
)

func GetTelegramBotApiByBotCode(c context.Context, code string) *tgbotapi.BotAPI {
	if s, ok := _bots.ByCode[code]; ok {
		return tgbotapi.NewBotAPIWithClient(s.Token, dtdal.HttpClient(c))
	} else {
		return nil
	}
}

var reTelegramStartCommandPrefix = regexp.MustCompile(`/start(@\w+)?\s+`)

func ParseStartCommand(whc botsfw.WebhookContext) (startParam string, startParams []string) {
	input := whc.Input()

	switch input := input.(type) {
	case botsfw.WebhookTextMessage:
		startParam = input.Text()
	case botsfw.WebhookReferralMessage:
		startParam = input.RefData()
	default:
		panic("Unknown input type")
	}
	if strings.HasPrefix(startParam, "/start") && startParam != "/start" {
		if loc := reTelegramStartCommandPrefix.FindStringIndex(startParam); len(loc) > 0 {
			startParam = startParam[loc[1]:]
			var utmMedium, utmSource string
			startParams = strings.Split(startParam, "__")
			for _, p := range startParams {
				switch {
				case strings.HasPrefix(p, "l="):
					code5 := p[len("l="):]
					if len(code5) == 5 {
						if err := whc.SetLocale(code5); err != nil {
							panic(fmt.Errorf("failed to set locale: %w", err))
						}
						whc.ChatData().SetPreferredLanguage(code5)
					}
				case strings.HasPrefix(p, "utm_m="):
					utmMedium = p[len("utm_m="):]
				case strings.HasPrefix(p, "utm_s="):
					utmSource = p[len("utm_s="):]
				}
			}
			if utmMedium != "" || utmSource != "" { // TODO: Handle analytics
				logus.Debugf(whc.Context(), "TODO: utm_medium=%v, utm_source=%v", utmMedium, utmSource)
			}
		} else {
			logus.Debugf(whc.Context(), "reTelegramStartCommandPrefix did not match - no start parameters")
		}
		return
	}
	return
}
