package anybot

import (
	"context"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/i18n"
)

func NewProfile(id string, router *botsfw.WebhooksRouter) botsfw.BotProfile {
	return botsfw.NewBotProfile(id, router,
		func() botsfwmodels.BotChatData {
			return new(botsfwtgmodels.TgChatBaseData)
		},
		func() botsfwmodels.PlatformUserData {
			return new(botsfwtgmodels.TgPlatformUserBaseDbo)
		},
		func() botsfwmodels.AppUserData {
			return new(dbo4userus.UserDbo)
		},
		func(ctx context.Context, tx dal.ReadSession, botID string, appUserID string) (appUser record.DataWithID[string, botsfwmodels.AppUserData], err error) {
			return
		},
		i18n.LocaleEnUS,
		[]i18n.Locale{i18n.LocaleEnUS, i18n.LocaleRuRu},
	)
}
