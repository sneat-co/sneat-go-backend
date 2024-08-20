package dtb_settings

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_all"
	"net/url"
)

var SettingsCommand = shared_all.SettingsCommandTemplate

func init() {
	SettingsCommand.Action = func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return shared_all.SettingsMainAction(whc)
	}
	SettingsCommand.CallbackAction = func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return shared_all.SettingsMainAction(whc)
	}
}
