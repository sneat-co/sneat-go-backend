package debtustgbots

import (
	"context"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
)

func GetTelegramBotApiByBotCode(ctx context.Context, code string) *tgbotapi.BotAPI {
	if s, ok := _bots.ByCode[code]; ok {
		return tgbotapi.NewBotAPIWithClient(s.Token, dtdal.HttpClient(ctx))
	} else {
		return nil
	}
}
