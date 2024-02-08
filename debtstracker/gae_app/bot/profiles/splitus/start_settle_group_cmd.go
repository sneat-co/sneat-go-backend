package splitus

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	botsgocore "github.com/bots-go-framework/bots-go-core"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal/gaedal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
	"net/url"
	"strings"
)

const (
	SettleGroupAskForCounterpartyCommandCode    = "sttl-grp"
	SettleGroupCounterpartyChosenCommandCode    = "sttl-grp-cp-chsn"
	SettleGroupCounterpartyConfirmedCommandCode = "sttl-grp-cp-cnfrmd"
)

var settleGroupAskForCounterpartyCommand = botsfw.Command{
	Code: SettleGroupAskForCounterpartyCommandCode,
	Action: shared_group.NewGroupAction(func(whc botsfw.WebhookContext, group models.Group) (m botsfw.MessageFromBot, err error) {
		return settleGroupAskForCounterpartyAction(whc, group)
	}),
	CallbackAction: shared_group.NewGroupCallbackAction(func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.Group) (m botsfw.MessageFromBot, err error) {
		return settleGroupAskForCounterpartyAction(whc, group)
	}),
}

func settleGroupStartAction(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
	var group models.Group
	for _, p := range startParams {
		switch {
		case strings.HasPrefix(p, "group="):
			group.ID = p[len("group="):]
		}
	}
	if group, err = dtdal.Group.GetGroupByID(whc.Context(), nil, group.ID); err != nil {
		return
	}
	return settleGroupAskForCounterpartyAction(whc, group)
}

func settleGroupAskForCounterpartyAction(whc botsfw.WebhookContext, group models.Group) (m botsfw.MessageFromBot, err error) {
	isDebtor, isSponsor := false, false

	groupMembers := group.Data.GetGroupMembers()

	userID := whc.AppUserID()

	var userMember models.GroupMemberJson

	balanceCurrencies := func(b money.Balance) (currencies []money.CurrencyCode) {
		currencies = make([]money.CurrencyCode, 0, len(b))
		for currency := range b {
			currencies = append(currencies, currency)
		}
		return
	}

	for i, m := range groupMembers {
		if m.UserID == userID {
			for _, v := range m.Balance {
				if v > 0 {
					if isSponsor = true; isDebtor {
						break
					}
				} else if v < 0 {
					if isDebtor = true; isSponsor {
						break
					}
				}
			}
			userMember = m
			groupMembers = groupMembers[:i+copy(groupMembers[i:], groupMembers[i+1:])]
			goto userMemberFound
		}
	}

	m.Text = "You are not a member of this group"
	log.Warningf(whc.Context(), m.Text)
	return

userMemberFound:

	if isSponsor && !isDebtor {
		groupMembers = filterGroupMembersByBalance(groupMembers, false, balanceCurrencies(userMember.Balance)...)
	} else if isDebtor && !isSponsor {
		groupMembers = filterGroupMembersByBalance(groupMembers, true, balanceCurrencies(userMember.Balance)...)
	}

	switch len(groupMembers) {
	case 0:
		m.Text = "There are no members to settele up with."
		//case 1:
		//	return settleGroupCounterpartyChosenAction(whc, group, userMember.ID)
	default:
		membersToKeyboard := func() botsgocore.Keyboard {
			keyboard := make([][]tgbotapi.InlineKeyboardButton, len(groupMembers))
			for i, m := range groupMembers {
				keyboard[i] = []tgbotapi.InlineKeyboardButton{
					{
						Text:         m.Name,
						CallbackData: fmt.Sprintf("%v?group=%v&member=%v", SettleGroupCounterpartyChosenCommandCode, group.ID, m.ID),
					},
				}
			}
			return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
		}

		var buf bytes.Buffer
		buf.WriteString(whc.Translate(trans.MT_GROUP_LABEL, group.Data.Name) + "\n\n")

		switch {
		case isSponsor && !isDebtor:
			if len(userMember.Balance) == 1 {
				for c, v := range userMember.Balance {
					buf.WriteString(fmt.Sprintf("You are owed %v %v by this group.\n\n", v, c))
				}
			}
			buf.WriteString("Who from group debtors will pay to you?")
		case isDebtor && !isSponsor:
			if len(userMember.Balance) == 1 {
				for c, v := range userMember.Balance {
					buf.WriteString(fmt.Sprintf("You owe %v %v to this group.\n\n", v, c))
				}
			}
			buf.WriteString("Who from group sponsors will collect your debt?")
		case isSponsor && isDebtor:
			buf.WriteString("Please choose with whom you are going to settle up?")
		}
		m.Keyboard = membersToKeyboard()
		m.Text = buf.String()
	}
	m.Format = botsfw.MessageFormatHTML

	return
}

var settleGroupCounterpartyChosenCommand = shared_group.GroupCallbackCommand(
	SettleGroupCounterpartyChosenCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.Group) (m botsfw.MessageFromBot, err error) {
		return settleGroupCounterpartyChosenAction(whc, group, callbackUrl.Query().Get("member"))
	},
)

func settleGroupCounterpartyChosenAction(whc botsfw.WebhookContext, group models.Group, memberID string) (m botsfw.MessageFromBot, err error) {

	var userMember, counterpartyMember models.GroupMemberJson
	userID := whc.AppUserID()
	for _, m := range group.Data.GetGroupMembers() {
		if m.UserID == userID {
			userMember = m
			if counterpartyMember.ID != "" {
				break
			}
		} else if m.ID == memberID {
			counterpartyMember = m
			if userMember.ID != "" {
				break
			}
		}
	}
	m.IsEdit = whc.InputType() == botsfw.WebhookInputCallbackQuery

	if userMember.ID == "" {
		m.Text = "You are not a member of this group."
		return
	} else if counterpartyMember.ID == "" {
		m.Text = "Selected member has left this group."
		return
	}
	m.Text = fmt.Sprintf("Have you returned this debt to %v already or you will return it?", counterpartyMember.Name)

	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "I have returned this debt",
				CallbackData: fmt.Sprintf("%v?debt=returned&group=%v&member=%v", SettleGroupCounterpartyConfirmedCommandCode, group.ID, memberID),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "I will returned this debt",
				CallbackData: fmt.Sprintf("%v?debt=will-return&group=%v&member=%v", SettleGroupCounterpartyConfirmedCommandCode, group.ID, memberID),
			},
		},
	)
	log.Debugf(whc.Context(), "counterpartyMember: %v", counterpartyMember)
	return
}

var settleGroupCounterpartyConfirmedCommand = shared_group.GroupCallbackCommand(
	SettleGroupCounterpartyConfirmedCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.Group) (m botsfw.MessageFromBot, err error) {
		q := callbackUrl.Query()
		currency := "RUB" // q.Get("currency")
		return settleGroupCounterpartyConfirmedAction(whc, group, q.Get("member"), money.CurrencyCode(currency))
	},
)

func settleGroupCounterpartyConfirmedAction(whc botsfw.WebhookContext, group models.Group, memberID string, currency money.CurrencyCode) (m botsfw.MessageFromBot, err error) {

	var userMember, counterpartyMember models.GroupMemberJson

	if counterpartyMember, err = group.Data.GetGroupMemberByID(memberID); err != nil {
		return
	}

	userID := whc.AppUserID()

	for _, m := range group.Data.GetGroupMembers() {
		if m.UserID == userID {
			userMember = m
			break
		}
	}

	var debtorID, sponsorID string

	userBalance := userMember.Balance[currency]
	counterpartyBalance := counterpartyMember.Balance[currency]

	if userBalance > 0 && counterpartyBalance < 0 {
		debtorID = counterpartyMember.ID
		sponsorID = userMember.ID
	} else if userBalance < 0 && counterpartyBalance > 0 {
		debtorID = userMember.ID
		sponsorID = counterpartyMember.ID
	} else {
		err = errors.New("Balance changed")
		return
	}

	if err = gaedal.Settle2members(whc.Context(), group.ID, debtorID, sponsorID, currency, 700); err != nil {
		return
	}

	m.Text = "Settled up"
	m.IsEdit = true
	log.Debugf(whc.Context(), "counterpartyMember: %v", counterpartyMember)
	return
}

func filterGroupMembersByBalance(members []models.GroupMemberJson, positive bool, currencies ...money.CurrencyCode) (result []models.GroupMemberJson) {
	result = make([]models.GroupMemberJson, 0, len(members))
	for _, m := range members {
		for c, v := range m.Balance {
			if (positive && v > 0) || (!positive && v < 0) {
				for _, currency := range currencies {
					if c == currency {
						result = append(result, m)
						break
					}
				}
			}
		}
	}
	return
}
