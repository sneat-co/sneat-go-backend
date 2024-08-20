package dtb_invite

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
)

const INVITE_COMMAND = "invite"

var InviteCommand = botsfw.Command{
	Code:     INVITE_COMMAND,
	Commands: []string{dtb_general.INVITES_SHOT_COMMAND, "/Пригласить_друга", "/invite"},
	Replies: []botsfw.Command{
		AskInviteAddressTelegramCommand,
		AskInviteAddressEmailCommand,
		AskInviteAddressSmsCommand,
	},
	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
		m := whc.NewMessageByCode(trans.MESSAGE_TEXT_ABOUT_INVITES)
		m.Keyboard = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(AskInviteAddressTelegramCommand.DefaultTitle(whc), "/invite"),
				},
				{
					{
						Text:         AskInviteAddressSmsCommand.DefaultTitle(whc),
						CallbackData: "invite?by=sms",
					},
					{
						Text:         AskInviteAddressEmailCommand.DefaultTitle(whc),
						CallbackData: "invite?by=email",
					},
				},
			},
		}
		whc.ChatData().SetAwaitingReplyTo(INVITE_COMMAND)
		return m, nil
	},
}
