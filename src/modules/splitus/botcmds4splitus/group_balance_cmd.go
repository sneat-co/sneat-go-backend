package botcmds4splitus

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_space"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/url"
)

const GroupBalanceCommandCode = "group-balance"

var groupBalanceCommand = botsfw.Command{
	Code:     GroupBalanceCommandCode,
	Commands: []string{"/balance"},
	Action:   shared_space.NewSplitusSpaceAction(groupBalanceAction),
	CallbackAction: shared_space.NewSplitusSpaceCallbackAction(func(whc botsfw.WebhookContext, callbackUrl *url.URL, splitusSpace models4splitus.SplitusSpaceEntry) (m botsfw.MessageFromBot, err error) {
		return groupBalanceAction(whc, splitusSpace)
	}),
}

func groupBalanceAction(whc botsfw.WebhookContext, splitusSpace models4splitus.SplitusSpaceEntry) (m botsfw.MessageFromBot, err error) {
	var buf bytes.Buffer
	writeMembers := func(members []briefs4splitus.SpaceSplitMember) {
		for i, m := range members {
			_, _ = fmt.Fprintf(&buf, " %d. %v:", i+1, m.Name)
			for currency, amount := range m.Balance {
				if amount < 0 {
					amount *= -1
				}
				fmt.Fprintf(&buf, " %v %v,", amount, currency)
			}
			buf.Truncate(buf.Len() - 1)
			buf.WriteString("\n")
		}
	}
	groupMembers := splitusSpace.Data.GetGroupMembers()
	sponsors, debtors := getGroupSponsorsAndDebtors(groupMembers)

	ctx := whc.Context()

	spaceID := splitusSpace.Key.Parent().ID.(string)
	user := facade.NewUserContext(whc.AppUserID())
	var space dbo4spaceus.SpaceEntry
	if space, err = facade4spaceus.GetSpace(ctx, user, spaceID); err != nil {
		return
	}

	buf.WriteString(whc.Translate(trans.MT_GROUP_LABEL, space.Data.Title))
	buf.WriteString("\n")

	buf.WriteString("\n")
	buf.WriteString(whc.Translate(trans.MT_SPONSORS_HEADER))
	buf.WriteString("\n")
	writeMembers(sponsors)

	buf.WriteString("\n")
	buf.WriteString(whc.Translate(trans.MT_DEBTORS_HEADER))
	buf.WriteString("\n")
	writeMembers(debtors)

	m.Text = buf.String()
	m.Format = botsfw.MessageFormatHTML
	m.IsEdit = whc.Input().InputType() == botinput.WebhookInputCallbackQuery

	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "Settle up",
				URL:  cmds4anybot.StartBotLink(whc.GetBotCode(), SettleGroupAskForCounterpartyCommandCode, "splitusSpace="+splitusSpace.ID),
			},
		},
	)
	return
}

func getGroupSponsorsAndDebtors(members []briefs4splitus.SpaceSplitMember, excludeMemberIDs ...string) (sponsors, debtors []briefs4splitus.SpaceSplitMember) {
	sponsors = make([]briefs4splitus.SpaceSplitMember, 0, len(members))
	debtors = make([]briefs4splitus.SpaceSplitMember, 0, len(members))

	for _, m := range members {
		for _, id := range excludeMemberIDs {
			if m.ID == id {
				continue
			}
		}
		for _, v := range m.Balance {
			if v > 0 {
				sponsors = append(sponsors, m)
			} else if v < 0 {
				debtors = append(debtors, m)
			}
		}
	}
	return
}

//func removeGroupMemberByID(members []models.SpaceSplitMember, excludeMemberID string) ([]models.SpaceSplitMember) {
//	for i, m := range members {
//		if m.ContactID == excludeMemberID {
//			return append(members[:i], members[i+1:]...)
//		}
//	}
//	return models.SpaceSplitMember{}, members
//}
