package splitus

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
	"net/url"

	"context"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/log"
)

const GroupMembersCommandCode = "group-members"

var groupMembersCommand = botsfw.Command{
	Code:     GroupMembersCommandCode,
	Commands: []string{"/members"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return showGroupMembers(whc, models.Group{}, false)
	},
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		var group models.Group
		if group, err = shared_group.GetGroup(whc, callbackUrl); err != nil {
			err = nil
			return
		}
		return showGroupMembers(whc, group, true)
	},
}

func groupMembersCard(
	c context.Context,
	t i18n.SingleLocaleTranslator,
	group models.Group,
	selectedMemberID int64,
) (text string, err error) {
	var buffer bytes.Buffer
	buffer.WriteString(t.Translate(trans.MESSAGE_TEXT_MEMBERS_CARD_TITLE, group.Data.MembersCount) + "\n\n")

	if group.Data == nil {
		if group, err = dtdal.Group.GetGroupByID(c, nil, group.ID); err != nil {
			return
		}
	}

	if group.Data.MembersCount > 0 {
		members := group.Data.GetGroupMembers()
		if len(members) == 0 {
			msg := fmt.Sprintf("ERROR: group.MembersCount:%d != 0 && len(members) == 0", group.Data.MembersCount)
			buffer.WriteString("\n" + msg + "\n")
			log.Errorf(c, msg)
		}

		splitMode := group.Data.GetSplitMode()

		var totalShares int

		if splitMode != models.SplitModeEqually {
			totalShares = group.Data.TotalShares()
		}

		for i, member := range members {
			if member.TgUserID == "" {
				fmt.Fprintf(&buffer, `  %d. %v`, i+1, member.Name) // TODO: Do a proper padding with 0 on left of #
			} else {
				fmt.Fprintf(&buffer, `  %d. <a href="tg://user?id=%v">%v</a>`, i+1, member.TgUserID, member.Name)
			}
			if splitMode != models.SplitModeEqually {
				fmt.Fprintf(&buffer, " (%d%%)", decimal.Decimal64p2(member.Shares*100/totalShares))
			}
			fmt.Fprintln(&buffer)
		}
	}

	buffer.WriteString("\n" + t.Translate(trans.MESSAGE_TEXT_MEMBERS_CARD_FOOTER))

	return buffer.String(), nil
}

func showGroupMembers(whc botsfw.WebhookContext, group models.Group, isEdit bool) (m botsfw.MessageFromBot, err error) {

	if group.Data == nil {
		if group, err = shared_group.GetGroup(whc, nil); err != nil {
			return
		}
	}

	c := whc.Context()

	if m.Text, err = groupMembersCard(c, whc, group, 0); err != nil {
		return
	}

	m.Format = botsfw.MessageFormatHTML
	tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.Translate(trans.BUTTON_TEXT_JOIN),
				CallbackData: joinGroupCommanCode,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
				emoji.CONTACTS_ICON+" "+whc.Translate(trans.COMMAND_TEXT_INVITE_MEMBER),
				shared_group.GroupCallbackCommandData(joinGroupCommanCode, group.ID),
			),
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.CommandText(trans.COMMAND_TEXT_SETTING, emoji.SETTINGS_ICON),
				CallbackData: shared_group.GroupCallbackCommandData(shared_all.SettingsCommandCode, group.ID),
			},
		},
	)
	m.Keyboard = tgKeyboard
	m.IsEdit = isEdit
	return
}
