package dtb_retention

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/cmd/dtb_general"
)

var DeleteUserCommand = botsfw.Command{
	Code:     "delete-user",
	Icon:     emoji.NO_ENTRY_SIGN_ICON,
	Commands: []string{"/deleteuser"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		err = botsfw.SetAccessGranted(whc, false)
		if err != nil {
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_FAILED_TO_DELETE_USER, err)
			dtb_general.SetMainMenuKeyboard(whc, &m)
			return m, nil
		} else {
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_USER_DELETED)
			return m, nil
		}
	},
}
