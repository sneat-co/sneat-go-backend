package cmds4listusbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"strings"
)

var listCommandPrefixes = []string{
	"/buy", "buy", "купить",
	"do", "/do",
	"watch", "/watch",
}

var listusListCommand = botsfw.Command{
	Code:     "list",
	Commands: []string{"/buy", "/do", "/watch"},
	Icon:     "🛒",
	InputTypes: []botinput.WebhookInputType{
		botinput.WebhookInputText,
		botinput.WebhookInputCallbackQuery,
	},
	Matcher: func(_ botsfw.Command, context botsfw.WebhookContext) bool {
		input := context.Input()
		if input.InputType() == botinput.WebhookInputText {
			text := strings.ToLower(strings.TrimSpace(input.(botinput.WebhookTextMessage).Text()))
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
