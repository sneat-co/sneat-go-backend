package anybot

import (
	"context"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/strongo/i18n"
)

type SneatAppTgData struct {
	SpaceRef string `json:"spaceRef,omitempty" firestore:"spaceRef,omitempty"`
}

func (v *SneatAppTgData) GetSpaceRef() core4spaceus.SpaceRef {
	return core4spaceus.SpaceRef(v.SpaceRef)
}

func (v *SneatAppTgData) SetSpaceRef(spaceRef core4spaceus.SpaceRef) {
	v.SpaceRef = string(spaceRef)
}

type SneatAppTgChatDbo struct {
	botsfwtgmodels.TgChatBaseData
	SneatAppTgData
}

type SneatAppTgChatEntry = record.DataWithID[string, *SneatAppTgChatDbo]

type SneatAppTgUserDbo struct {
	botsfwtgmodels.TgPlatformUserBaseDbo
	SneatAppTgData
}

func NewProfile(id string, router *botsfw.WebhooksRouter) botsfw.BotProfile {
	return botsfw.NewBotProfile(id, router,
		func() botsfwmodels.BotChatData {
			return new(SneatAppTgChatDbo)
		},
		func() botsfwmodels.PlatformUserData {
			return new(SneatAppTgUserDbo)
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
