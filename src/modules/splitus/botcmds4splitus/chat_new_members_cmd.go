package botcmds4splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"time"

	"github.com/sneat-co/debtstracker-translations/emoji"
)

const NewChatMembersCommandCode = "new-chat-members"

var newChatMembersCommand = botsfw.Command{
	Code: NewChatMembersCommandCode,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()

		newMembersMessage := whc.Input().(botsfw.WebhookNewChatMembersMessage)

		newMembers := newMembersMessage.NewChatMembers()

		{ // filter out botscore
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

		var newUsers []facade4splitus.NewUser

		{ // Get or create related user records
			for _, chatMember := range newMembers {
				tgChatMember := chatMember.(tgbotapi.ChatMember)
				var botUser record.DataWithID[string, botsfwmodels.PlatformUserData]
				if botUser, err = whc.BotUser(); err != nil && !dal.IsNotFound(err) {
					return
				}
				if !botUser.Record.Exists() {
					botUser.Data = &botsfwmodels.PlatformUserBaseDbo{
						BotBaseData: botsfwmodels.BotBaseData{
							DtCreated: time.Now(),
						},
					}
					if err = whc.Tx().Set(ctx, botUser.Record); err != nil {
						return
					}
				}
				newUsers = append(newUsers, facade4splitus.NewUser{
					Name:             tgChatMember.GetFullName(),
					PlatformUserData: botUser.Data,
					ChatMember:       chatMember,
				})
			}
		}

		var splitusSpace models4splitus.SplitusSpaceEntry
		//if splitusSpace, err = shared_space.GetSpaceEntryByCallbackUrl(whc, nil); err != nil {
		//	return
		//}
		if splitusSpace, newUsers, err = facade4splitus.AddUsersToTheGroupAndOutstandingBills(whc.Context(), splitusSpace.ID, newUsers); err != nil {
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
						URL:  fmt.Sprintf("https:/t.me/%v?start=splitusSpace-%v", whc.GetBotCode(), splitusSpace.ID),
					},
				},
			)

			return m
		}
		m = createWelcomeMsg(newUsers[0].ChatMember)

		if len(newUsers) > 1 {
			responder := whc.Responder()
			ctx := whc.Context()
			for _, newUser := range newUsers {
				if _, err = responder.SendMessage(ctx, createWelcomeMsg(newUser.ChatMember), botsfw.BotAPISendMessageOverHTTPS); err != nil {
					return
				}
			}
		}
		return
	},
}
