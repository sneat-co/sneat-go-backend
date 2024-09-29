package botcmds4splitus

import (
	"bytes"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/debtusbots/profiles/shared_space"
	"github.com/strongo/i18n"
	"net/url"

	"context"
	"github.com/sneat-co/debtstracker-translations/emoji"
)

const GroupMembersCommandCode = "group-members"

var groupMembersCommand = botsfw.Command{
	Code:     GroupMembersCommandCode,
	Commands: []string{"/members"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		// TODO: implement persisted active space context
		return showGroupMembers(whc, dal4contactus.ContactusSpaceEntry{}, false)
	},
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {

		ctx := whc.Context()

		spaceID := shared_space.GetSpaceIdFromCallbackUrl(callbackUrl)
		contactusSpace := dal4contactus.NewContactusSpaceEntry(spaceID)

		var db dal.DB
		if err = db.Get(ctx, contactusSpace.Record); err != nil {
			return
		}
		return showGroupMembers(whc, contactusSpace, true)
	},
}

func groupMembersCard(
	_ context.Context,
	t i18n.SingleLocaleTranslator,
	contactusSpace dal4contactus.ContactusSpaceEntry,
	selectedMemberID int64,
) (text string, err error) {

	membersCount := contactusSpace.Data.GetContactsCount(const4contactus.SpaceMemberRoleMember)
	var buffer bytes.Buffer
	buffer.WriteString(t.Translate(trans.MESSAGE_TEXT_MEMBERS_CARD_TITLE, membersCount) + "\n\n")

	if membersCount > 0 {
		panic("not implemented")
		//members := contactusSpace.Data.GetGroupMembers()
		//if len(members) == 0 {
		//	msg := fmt.Sprintf("ERROR: contactusSpace.MembersCount:%d != 0 && len(members) == 0", membersCount)
		//	buffer.WriteString("\n" + msg + "\n")
		//	logus.Errorf(c, msg)
		//}
		//
		//splitMode := contactusSpace.Data.GetSplitMode()
		//
		//var totalShares int
		//
		//if splitMode != models4debtus.SplitModeEqually {
		//	totalShares = contactusSpace.Data.TotalShares()
		//}
		//
		//for i, member := range members {
		//	if member.TgUserID == "" {
		//		fmt.Fprintf(&buffer, `  %d. %v`, i+1, member.Name) // TODO: Do a proper padding with 0 on left of #
		//	} else {
		//		fmt.Fprintf(&buffer, `  %d. <a href="tg://user?id=%v">%v</a>`, i+1, member.TgUserID, member.Name)
		//	}
		//	if splitMode != models4debtus.SplitModeEqually {
		//		fmt.Fprintf(&buffer, " (%d%%)", decimal.Decimal64p2(member.Shares*100/totalShares))
		//	}
		//	fmt.Fprintln(&buffer)
		//}
	}

	buffer.WriteString("\n" + t.Translate(trans.MESSAGE_TEXT_MEMBERS_CARD_FOOTER))

	return buffer.String(), nil
}

func showGroupMembers(whc botsfw.WebhookContext, contactusSpace dal4contactus.ContactusSpaceEntry, isEdit bool) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	if m.Text, err = groupMembersCard(ctx, whc, contactusSpace, 0); err != nil {
		return
	}

	m.Format = botsfw.MessageFormatHTML
	tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.Translate(trans.BUTTON_TEXT_JOIN),
				CallbackData: joinSpaceCommandCode,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
				emoji.CONTACTS_ICON+" "+whc.Translate(trans.COMMAND_TEXT_INVITE_MEMBER),
				shared_space.SpaceCallbackCommandData(joinSpaceCommandCode, contactusSpace.ID),
			),
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.CommandText(trans.COMMAND_TEXT_SETTING, emoji.SETTINGS_ICON),
				CallbackData: shared_space.SpaceCallbackCommandData(cmds4anybot.SettingsCommandCode, contactusSpace.ID),
			},
		},
	)
	m.Keyboard = tgKeyboard
	m.IsEdit = isEdit
	return
}
