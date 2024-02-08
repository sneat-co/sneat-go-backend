package splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/debtstracker-translations/trans"
	"time"

	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

const NEW_CHAT_MEMBERS_COMMAND = "new-chat-members"

var newChatMembersCommand = botsfw.Command{
	Code: NEW_CHAT_MEMBERS_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()

		newMembersMessage := whc.Input().(botsfw.WebhookNewChatMembersMessage)

		newMembers := newMembersMessage.NewChatMembers()

		{ // filter out bots
			j := 0
			for _, member := range newMembers {
				if !member.IsBotUser() {
					newMembers[j] = member
					j += 1
				}
			}
			newMembers = newMembers[:j]
		}

		if len(newMembers) == 0 {
			return
		}

		var newUsers []facade.NewUser

		{ // Get or create related user records
			for _, chatMember := range newMembers {
				tgChatMember := chatMember.(tgbotapi.ChatMember)
				var botUser record.DataWithID[string, botsfwmodels.BotUserData]
				if botUser, err = whc.BotUser(); err != nil && !dal.IsNotFound(err) {
					return
				}
				if !botUser.Record.Exists() {
					botUser.Data = &botsfwmodels.BotUserBaseData{
						BotBaseData: botsfwmodels.BotBaseData{
							DtCreated: time.Now(),
						},
					}
					if err = whc.Tx().Set(c, botUser.Record); err != nil {
						return
					}
				}
				newUsers = append(newUsers, facade.NewUser{
					Name:        tgChatMember.GetFullName(),
					BotUserData: botUser.Data,
					ChatMember:  chatMember,
				})
			}
		}

		var group models.Group
		if group, err = shared_group.GetGroup(whc, nil); err != nil {
			return
		}
		if group, newUsers, err = facade.Group.AddUsersToTheGroupAndOutstandingBills(whc.Context(), group.ID, newUsers); err != nil {
			return
		}

		if len(newUsers) == 0 {
			return
		}

		createWelcomeMsg := func(member botsfw.WebhookActor) botsfw.MessageFromBot {
			m := whc.NewMessageByCode(trans.MESSAGE_TEXT_USER_JOINED_GROUP, member.GetFirstName())
			m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					{
						Text: whc.CommandText(trans.COMMAND_TEXT_SETTING, emoji.SETTINGS_ICON),
						URL:  fmt.Sprintf("https:/t.me/%v?start=group-%v", whc.GetBotCode(), group.ID),
					},
				},
			)

			return m
		}
		m = createWelcomeMsg(newUsers[0].ChatMember)

		if len(newUsers) > 1 {
			responder := whc.Responder()
			c := whc.Context()
			for _, newUser := range newUsers {
				if _, err = responder.SendMessage(c, createWelcomeMsg(newUser.ChatMember), botsfw.BotAPISendMessageOverHTTPS); err != nil {
					return
				}
			}
		}
		return
	},
}
