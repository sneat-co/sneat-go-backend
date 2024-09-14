package dtb_general

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/logus"
)

const ClearCommandCode = "clear"

var ClearCommand = botsfw.Command{
	Code:     ClearCommandCode,
	Commands: trans.Commands(trans.COMMAND_CLEAR),
	//Title:    trans.COMMAND_TEXT_MAIN_MENU_TITLE,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		logus.Warningf(whc.Context(), "User called /clear command (not implemented yet)")
		return MainMenuAction(whc, whc.Translate(trans.MESSAGE_TEXT_NOT_IMPLEMENTED_YET), false)
	},
}
