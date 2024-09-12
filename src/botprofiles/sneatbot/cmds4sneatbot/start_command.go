package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"strings"
)

//var startCommand = botsfw.Command{
//	Code:     "start",
//	Commands: []string{"/start"},
//	InputTypes: []botinput.WebhookInputType{
//		botinput.WebhookInputText,
//		botinput.WebhookInputCallbackQuery,
//		botinput.WebhookInputInlineQuery,
//	},
//	Action: sneatBotStartAction,
//}

func startActionWithStartParams(whc botsfw.WebhookContext, _ []string) (m botsfw.MessageFromBot, err error) {
	return sneatBotStartAction(whc)
}

func sneatBotStartAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	return spaceAction(whc, core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, ""))
}

func sneatBotWelcomeMessage(_ botsfw.WebhookContext) (text string, err error) {
	msg := make([]string, 0)
	msg = append(msg, "Hello, stranger! I'm a @SneatBot.")
	msg = append(msg, "I can help you to manage your day-to-day family life.")
	msg = append(msg, "Or you can create a space to manage your group/team/community.")
	return strings.Join(msg, "\n\n"), err
}
