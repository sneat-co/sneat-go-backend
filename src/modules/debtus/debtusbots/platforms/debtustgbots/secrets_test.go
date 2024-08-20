package debtustgbots

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/debtusbotconst"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp"
	"testing"
)

func TestGetBotSettingsByLang(t *testing.T) {
	t.Skip("TODO: fix this test to run on CI")
	verify := func(profile, locale, code string) {
		botSettings, err := GetBotSettingsByLang(strongoapp.LocalHostEnv, debtusbotconst.DebtusBotProfileID, locale)
		if err != nil {
			t.Fatal(err)
		}
		if botSettings.Code != code {
			t.Error(code + " not found in settings, got: " + botSettings.Code)
		}
	}
	verify(debtusbotconst.DebtusBotProfileID, i18n.LocalCodeRuRu, "DebtsTrackerRuBot")
	verify(debtusbotconst.DebtusBotProfileID, i18n.LocaleCodeEnUS, "DebtsTrackerBot")
}
