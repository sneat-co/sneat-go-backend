package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/strongo/strongoapp"
)

const DELETE_ALL_COMMAND = "delete-all"

var DeleteAllCommand = botsfw.Command{
	Code:     DELETE_ALL_COMMAND,
	Icon:     emoji.MAIN_MENU_ICON,
	Commands: []string{"/deleteall"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		botSettings := whc.GetBotSettings()
		if botSettings.Env != strongoapp.LocalHostEnv && botSettings.Env != "dev" {
			return whc.NewMessage(fmt.Sprintf("This command supported just in development, got botSettings.Env: %v", botSettings.Env)), nil
		} else if botSettings.Env == "prod" {
			return whc.NewMessage("This command supported production environment"), nil
		}

		// We create a success message ahead of actual operation as keyboard creation will fail once user deleted.
		m = whc.NewMessage("Deleted all records")
		SetMainMenuKeyboard(whc, &m)

		var chatID string
		if chatID, err = whc.Input().BotChatID(); err != nil {
			return
		}

		if err = dtdal.Admin.DeleteAll(whc.Context(), botSettings.Code, chatID); err != nil {
			return
		}

		return
	},
}
