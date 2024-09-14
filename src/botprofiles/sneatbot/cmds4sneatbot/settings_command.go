package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/cmds4anybot"
	"net/url"
)

var settingsCommand = cmds4anybot.SettingsCommandTemplate

func init() {
	settingsCommand.Action = func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return cmds4anybot.SettingsMainAction(whc)
	}
	settingsCommand.CallbackAction = func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return cmds4anybot.SettingsMainAction(whc)
	}
}
