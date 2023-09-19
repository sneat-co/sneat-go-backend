package bots

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo2firestore"
	strongo "github.com/strongo/app"
	"github.com/strongo/i18n"
	"os"
	"strings"
)

var _bots botsfw.SettingsBy

func telegramBots(environment strongo.Environment) botsfw.SettingsBy {
	if _bots.ByCode != nil {
		return _bots
	}

	var getDb = func(c context.Context) (dal.DB, error) {
		fsClient, err := firestore.NewClient(c, "demo-local-sneat-app")
		if err != nil {
			return nil, err
		}
		return dalgo2firestore.NewDatabase("sneat", fsClient), nil
	}

	getAppUser := func(ctx context.Context, tx dal.ReadSession, botID string, appUserID string) (appUser record.DataWithID[string, botsfwmodels.AppUserData], err error) {
		return
	}

	switch environment {
	case strongo.EnvProduction:
		_bots = botsfw.NewBotSettingsBy(
			telegram.NewTelegramBot(strongo.EnvProduction, SneatBotProfile, "SneatBot", "", "", "", "", i18n.LocaleEnUS, getDb, getAppUser),
			//telegram.NewTelegramBot(strongo.EnvProduction, bot4listus.ListusBotProfile, "ListusBot", "", "", "", "", i18n.LocaleEnUS, getDb, getAppUser),
		)
	case strongo.EnvLocal:
		sneatTgDevBot := os.Getenv("SNEAT_TG_DEV_BOTS")
		if sneatTgDevBot == "" {
			panic("Environment variable SNEAT_TG_DEV_BOTS is not set")
		}

		bots := strings.Split(sneatTgDevBot, ",")

		botSettings := make([]botsfw.BotSettings, 0, len(bots))

		for i, bot := range bots {
			sneatDevBotVals := strings.Split(bot, ":")
			if len(sneatDevBotVals) != 2 {
				panic(fmt.Sprintf("Invalid SNEAT_TG_DEV_BOT (should be in format of 'id:profileID'): %s", sneatTgDevBot))
			}
			botID := sneatDevBotVals[0]
			if botID == "" {
				panic(fmt.Sprintf("Invalid SNEAT_TG_DEV_BOTS[%d] (should be in format of 'id:profileID'): ", i) + sneatTgDevBot)
			}
			profileID := sneatDevBotVals[1]
			if profileID == "" {
				panic(fmt.Sprintf("Invalid SNEAT_TG_DEV_BOTS[%d] (should be in format of 'id:profileID'): ", i) + sneatTgDevBot)
			}
			var profile botsfw.BotProfile
			switch profileID {
			case "sneat":
				profile = SneatBotProfile
			//case "listus":
			//profile = bot4listus.ListusBotProfile
			default:
				panic(fmt.Sprintf("Unsupported profileID: %s", profileID))
			}
			botSettings = append(botSettings, telegram.NewTelegramBot(strongo.EnvLocal, profile, botID, "", "", "", "", i18n.LocaleEnUS, getDb, getAppUser))
		}
		_bots = botsfw.NewBotSettingsBy(botSettings...)
	default:
		panic(fmt.Sprintf("Unsupported environment: %d='%s'", environment, strongo.EnvironmentNames[environment]))
	}

	return _bots
}
