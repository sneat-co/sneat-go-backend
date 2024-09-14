package cmds4sneatbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
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
	return spaceAction(whc, core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, ""))
}

func sneatBotWelcomeMessage(whc botsfw.WebhookContext) (text string, err error) {
	text = whc.Translate(trans.SNEATBOT_MSG_TXT_START)

	sender := whc.Input().GetSender()

	name := sender.GetFirstName()
	if name == "" {
		if name = sender.GetLastName(); name == "" {
			if name = sender.GetUserName(); name == "" {
				if name = whc.GetBotUserID(); name == "" {
					if name = whc.AppUserID(); name == "" {
						name = "stranger"
					}
				}
			}
		}
	}
	text = fmt.Sprintf(text, name)
	return
}
