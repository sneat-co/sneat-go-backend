package cmds4listusbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"strings"
)

var listCommandPrefixes = []string{
	"/buy", "buy", "купить",
	"do", "/do",
	"watch", "/watch",
}

var listCommand = botsfw.Command{
	Code:     "list",
	Commands: []string{"/buy", "/do", "/watch"},
	Icon:     "🛒",
	InputTypes: []botsfw.WebhookInputType{
		botsfw.WebhookInputText,
		botsfw.WebhookInputCallbackQuery,
	},
	Matcher: func(_ botsfw.Command, context botsfw.WebhookContext) bool {
		input := context.Input()
		if input.InputType() == botsfw.WebhookInputText {
			text := strings.ToLower(strings.TrimSpace(input.(botsfw.WebhookTextMessage).Text()))
			for _, prefix := range listCommandPrefixes {
				if strings.HasPrefix(text, prefix+" ") {
					return true
				}
			}
			return false
		}
		return false
	},
	Action:         listAction,
	CallbackAction: listCallbackAction,
}
