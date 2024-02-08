package viberbots

import (
	"context"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var _bots botsfw.SettingsBy

func Bots(c context.Context) botsfw.SettingsBy { //TODO: Consider to do pre-deployment replace
	//if len(_bots.ByCode) == 0 {
	//	host := appengine.DefaultVersionHostname(c)
	//	//if host == "" || strings.Contains(host, "dev") {
	//		//_bots = botsfw.NewBotSettingsBy(nil,
	//		//	// Development bot
	//		//	viber.NewViberBot(strongoapp.EnvDevTest, bot.ProfileDebtus, "DebtsTrackerDev", "451be8dd024fbbc7-4fb4285be8dbb24e-1b2d99610f798855", "", i18n.LocalesByCode5[i18n.LocaleCodeEnUS]),
	//		//)
	//		//} else if strings.Contains(host, "st1") {
	//		//_bots = botsfw.NewBotSettingsBy(
	//		//	// Staging bots
	//		//)
	//		//} else if strings.HasPrefix(host, "debtstracker-io.") {
	//		//_bots = botsfw.NewBotSettingsBy(nil,
	//		//	// Production bot
	//		//	viber.NewViberBot(strongoapp.EnvProduction, bot.ProfileDebtus, "DebtsTracker", "xxxx-xxx-xxxx", common.GA_TRACKING_ID, i18n.LocalesByCode5[i18n.LocaleCodeEnUS]),
	//		//)
	//	}
	//}
	return _bots
}

// TODO: Decouple to common lib
//func GetBotSettingsByLang(c context.Context, lang string) (botsfw.BotSettings, error) {
//	botSettingsBy := Bots(c)
//	langLen := len(lang)
//	if langLen == 2 {
//		lang = fmt.Sprintf("%v-%v", strings.ToLower(lang), strings.ToUpper(lang))
//	} else if langLen != 5 {
//		return botsfw.BotSettings{}, fmt.Errorf("Invalid length of lang parameter: %v, %v", langLen, lang)
//	}
//	if botSettings, ok := botSettingsBy.Locale[lang]; ok {
//		return botSettings, nil
//	} else if lang != DEFAULT_LOCALE {
//		if botSettings, ok = botSettingsBy.Locale[DEFAULT_LOCALE]; ok {
//			return botSettings, nil
//		}
//	}
//	return botsfw.BotSettings{}, fmt.Errorf("No bot setting for both %v & %v locales.", lang, DEFAULT_LOCALE)
//}
