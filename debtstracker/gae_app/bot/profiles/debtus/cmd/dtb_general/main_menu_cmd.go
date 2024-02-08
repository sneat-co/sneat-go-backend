package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/log"
	"net/url"
)

const MAIN_MENU_COMMAND = "main-menu"

var MainMenuCommand = botsfw.Command{
	Code:     MAIN_MENU_COMMAND,
	Icon:     emoji.MAIN_MENU_ICON,
	Commands: trans.Commands(trans.COMMAND_MENU, emoji.MAIN_MENU_ICON),
	Title:    trans.COMMAND_TEXT_MAIN_MENU_TITLE,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return MainMenuAction(whc, "", true)
	},
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return MainMenuAction(whc, "", true)
	},
}

func MainMenuAction(whc botsfw.WebhookContext, messageText string, showHint bool) (m botsfw.MessageFromBot, err error) {
	if messageText == "" {
		//if whc.BotPlatform().ID() != fbm.PlatformID {
		if showHint {
			messageText = fmt.Sprintf("%v\n\n%v", whc.Translate(trans.MESSAGE_TEXT_WHATS_NEXT), whc.Translate(trans.MESSAGE_TEXT_DEBTUS_COMMANDS))
		} else {
			messageText = whc.Translate(trans.MESSAGE_TEXT_WHATS_NEXT)
		}
		//}
	}
	log.Infof(whc.Context(), "MainMenuCommand.Action()")
	whc.ChatData().SetAwaitingReplyTo("")
	m = whc.NewMessage(messageText)
	SetMainMenuKeyboard(whc, &m)
	return m, nil
}
