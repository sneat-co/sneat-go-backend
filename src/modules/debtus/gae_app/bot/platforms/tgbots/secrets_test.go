package tgbots

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp"
	"testing"
)

func TestGetBotSettingsByLang(t *testing.T) {
	t.Skip("TODO: fix this test to run on CI")
	verify := func(profile, locale, code string) {
		botSettings, err := GetBotSettingsByLang(strongoapp.LocalHostEnv, bot.ProfileDebtus, locale)
		if err != nil {
			t.Fatal(err)
		}
		if botSettings.Code != code {
			t.Error(code + " not found in settings, got: " + botSettings.Code)
		}
	}
	verify(bot.ProfileDebtus, i18n.LocalCodeRuRu, "DebtsTrackerRuBot")
	verify(bot.ProfileDebtus, i18n.LocaleCodeEnUS, "DebtsTrackerBot")
}
