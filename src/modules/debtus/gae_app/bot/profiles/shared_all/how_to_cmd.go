package shared_all

import "github.com/bots-go-framework/bots-fw/botsfw"

const howToCommandCode = "how-to"

var howToCommand = botsfw.Command{
	Code: howToCommandCode,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "<b>How To</b> - not implemented yet"
		return
	},
}
