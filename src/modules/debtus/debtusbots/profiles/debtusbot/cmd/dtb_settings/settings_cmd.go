package dtb_settings

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/anybot/cmds4anybot"
	"net/url"
)

var SettingsCommand = cmds4anybot.SettingsCommandTemplate

func init() {
	SettingsCommand.Action = func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return cmds4anybot.SettingsMainAction(whc)
	}
	SettingsCommand.CallbackAction = func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return cmds4anybot.SettingsMainAction(whc)
	}
}
