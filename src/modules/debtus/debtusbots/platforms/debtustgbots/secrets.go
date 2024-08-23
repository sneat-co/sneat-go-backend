package debtustgbots

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/collectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/botcmds4splitus"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp"
	"slices"
)

var _bots botsfw.SettingsBy

const DefaultLocale = i18n.LocaleCodeEnUS

//var DebtusBotToken = ""

func newTelegramBot(
	mode string,
	botProfile botsfw.BotProfile,
	code, gaToken string,
	locale i18n.Locale,
) botsfw.BotSettings {
	return telegram.NewTelegramBot(mode, botProfile, code, "", "", "", gaToken, i18n.LocaleEnUS, nil, nil)
}

func Bots(environment string) botsfw.SettingsBy { //TODO: Consider to do pre-deployment replace

	errFooterText := func() string {
		return "error footer"
	}

	debtusBotProfile := GetDebtusBotProfile(errFooterText)
	splitusBotProfile := botsfw.NewBotProfile("splitus", &botcmds4splitus.Router, newBotChatData, newBotUserData, newAppUserData, getAppUserByID, i18n.LocaleEnUS, nil)
	collectusBotProfile := botsfw.NewBotProfile("collectus", &collectus.Router, newBotChatData, newBotUserData, newAppUserData, getAppUserByID, i18n.LocaleEnUS, nil)

	const prod = "prod"

	if len(_bots.ByCode) == 0 {
		//logus.Debugf(c, "Bots() => hostname:%v, environment:%s:%s", hostname, environment, strongoapp.EnvironmentNames[environment])
		switch environment {
		case prod:
			_bots = botsfw.NewBotSettingsBy( // Production bots
				newTelegramBot(prod, debtusBotProfile, "DebtusBot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
				newTelegramBot(prod, splitusBotProfile, "SplitusBot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
				newTelegramBot(prod, collectusBotProfile, "CollectusBot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerBot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerRuBot", common4debtus.GA_TRACKING_ID, i18n.LocaleRuRu),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerFaBot", common4debtus.GA_TRACKING_ID, i18n.LocalesByCode5[i18n.LocaleCodeFaIR]),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerItBot", common4debtus.GA_TRACKING_ID, i18n.LocaleItIt),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerFrBot", common4debtus.GA_TRACKING_ID, i18n.LocaleFrFr),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerDeBot", common4debtus.GA_TRACKING_ID, i18n.LocaleDeDe),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerPLbot", common4debtus.GA_TRACKING_ID, i18n.LocalePlPl),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerPtBot", common4debtus.GA_TRACKING_ID, i18n.LocalePtBr),
				//newTelegramBot(prod, debtusBotProfile, "DebtsTrackerEsBot", common4debtus.GA_TRACKING_ID, i18n.LocalePtBr),
			)
		case "dev":
			_bots = botsfw.NewBotSettingsBy( // Development bots
				newTelegramBot("dev", debtusBotProfile, "DebtsTrackerDev1Bot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
				newTelegramBot("dev", debtusBotProfile, "DebtsTrackerDev1RuBot", common4debtus.GA_TRACKING_ID, i18n.LocaleRuRu),
				//telegram.NewTelegramBot(strongoapp.EnvDevTest, bot.ProfileDebtus, "DebtsTrackerDev2RuBot", "360514041:AAFXuT0STHBD9cOn1SFmKzTYDmalP0Rz-7M", "", "", anybot.GA_TRACKING_ID, i18n.LocalesByCode5[i18n.LocalCodeRuRu]),
			)
		case "staging":
			_bots = botsfw.NewBotSettingsBy( // Staging bots
				newTelegramBot("staging", debtusBotProfile, "DebtsTrackerSt1Bot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
			)
		case strongoapp.LocalHostEnv:
			_bots = botsfw.NewBotSettingsBy( // Staging bots
				newTelegramBot(strongoapp.LocalHostEnv, debtusBotProfile, "DebtsTrackerLocalBot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
				newTelegramBot(strongoapp.LocalHostEnv, splitusBotProfile, "SplitusLocalBot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
				newTelegramBot(strongoapp.LocalHostEnv, collectusBotProfile, "CollectusLocalBot", common4debtus.GA_TRACKING_ID, i18n.LocaleEnUS),
			)
		case strongoapp.UnknownEnv:
			// Pass for unit tests?
		default:
			panic(fmt.Sprintf("Unknown environment => %s", environment))
		}
	}
	return _bots
}

func GetBotSettingsByLang(environment string, profile, lang string) (botSettings *botsfw.BotSettings, err error) {
	botSettingsBy := Bots(environment)
	if profileBots, ok := botSettingsBy.ByProfile[profile]; !ok {
		err = fmt.Errorf("no bot settings for profileID=%s", profile)
		return
	} else {
		locales := make([]string, 0, len(profileBots))
		getBotSettingsByLocale := func(locale string) *botsfw.BotSettings {
			for _, bs := range profileBots {
				if bs.Locale.Code5 == lang {
					return bs
				}
				if slices.Contains(locales, bs.Locale.Code5) {
					locales = append(locales, bs.Locale.Code5)
				}
			}
			return nil
		}
		if bs := getBotSettingsByLocale(lang); bs != nil {
			return bs, nil
		}
		if bs := getBotSettingsByLocale(DefaultLocale); bs != nil {
			return bs, nil
		}
		return nil, fmt.Errorf("no bot setting for both %s & %s locales", lang, DefaultLocale)
	}
}
