package bots

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/bots/listusbot"
	"github.com/sneat-co/sneat-go-backend/src/bots/sneatbot"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/debtusbotconst"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/platforms/debtustgbots"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/i18n"
	"os"
	"strings"
)

var _bots botsfw.SettingsBy

var GetDb = facade.GetDatabase

func telegramBots(environment string) botsfw.SettingsBy {
	if _bots.ByCode != nil {
		return _bots
	}

	getAppUser := func(ctx context.Context, tx dal.ReadSession, botID string, appUserID string) (appUser record.DataWithID[string, botsfwmodels.AppUserData], err error) {
		return
	}

	errFooterText := func() string {
		return "Please report any issues to @trakhimenok"
	}

	switch environment {
	case botsfw.EnvProduction:
		_bots = botsfw.NewBotSettingsBy( // TODO: Get bot tokens from environment variables
			telegram.NewTelegramBot(environment, sneatbot.GetProfile(errFooterText), "SneatBot", "", "", "", "", i18n.LocaleEnUS, GetDb, getAppUser),
			telegram.NewTelegramBot(environment, listusbot.GetProfile(errFooterText), "ListusBot", "", "", "", "", i18n.LocaleEnUS, GetDb, getAppUser),
		)
	case botsfw.EnvLocal:
		sneatTgDevBots := os.Getenv("SNEAT_TG_DEV_BOTS")
		if sneatTgDevBots == "" {
			panic("Environment variable SNEAT_TG_DEV_BOTS is not set")
		}

		botIDs := strings.Split(sneatTgDevBots, ",")

		botSettings := make([]botsfw.BotSettings, 0, len(botIDs))

		for i, bot := range botIDs {
			sneatDevBotVals := strings.Split(bot, ":")
			if len(sneatDevBotVals) != 2 {
				panic(fmt.Sprintf("Invalid SNEAT_TG_DEV_BOT (should be in format of 'id:profileID'): %s", sneatTgDevBots))
			}
			botID := sneatDevBotVals[0]
			if botID == "" {
				panic(fmt.Sprintf("Invalid SNEAT_TG_DEV_BOTS[%d] (should be in format of 'id:profileID'): ", i) + sneatTgDevBots)
			}
			profileID := sneatDevBotVals[1]
			if profileID == "" {
				panic(fmt.Sprintf("Invalid SNEAT_TG_DEV_BOTS[%d] (should be in format of 'id:profileID'): ", i) + sneatTgDevBots)
			}
			var profile botsfw.BotProfile
			switch profileID {
			case sneatbot.ProfileID:
				profile = sneatbot.GetProfile(errFooterText)
			case debtusbotconst.DebtusBotProfileID:
				profile = debtustgbots.GetDebtusBotProfile(errFooterText)
			case listusbot.ProfileID:
				profile = listusbot.GetProfile(errFooterText)
			default:
				panic(fmt.Sprintf("Unsupported profileID: %s", profileID))
			}
			botSettings = append(botSettings, telegram.NewTelegramBot(botsfw.EnvLocal, profile, botID, "", "", "", "", i18n.LocaleEnUS, GetDb, getAppUser))
		}
		_bots = botsfw.NewBotSettingsBy(botSettings...)
	default:
		panic(fmt.Sprintf("Unsupported environment: '%s'", environment))
	}

	return _bots
}
