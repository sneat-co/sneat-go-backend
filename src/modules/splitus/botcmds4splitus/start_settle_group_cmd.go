package botcmds4splitus

import (
	"errors"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/debtusbots/profiles/shared_space"
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
	Action: shared_space.NewSpaceAction(func(whc botsfw.WebhookContext, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		return settleGroupAskForCounterpartyAction(whc, space)
	}),
	CallbackAction: shared_space.NewSpaceCallbackAction(func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		return settleGroupAskForCounterpartyAction(whc, space)
	}),
}

func settleGroupStartAction(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
	var space dbo4spaceus.SpaceEntry
	for _, p := range startParams {
		switch {
		case strings.HasPrefix(p, "space="):
			space.ID = p[len("space="):]
		}
	}
	//if space, err = dtdal.Group.GetGroupByID(whc.Context(), nil, space.ContactID); err != nil {
	//	return
	//}
	return settleGroupAskForCounterpartyAction(whc, space)
}

func settleGroupAskForCounterpartyAction(_ botsfw.WebhookContext, _ dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
	err = errors.New("not implemented yet")
	//	isDebtor, isSponsor := false, false
	//
	//	groupMembers := space.Data.GetGroupMembers()
	//
	//	userID := whc.AppUserID()
	//
	//	var userMember briefs4splitus.GroupMemberJson
	//
	//	balanceCurrencies := func(b money.Balance) (currencies []money.CurrencyCode) {
	//		currencies = make([]money.CurrencyCode, 0, len(b))
	//		for currency := range b {
	//			currencies = append(currencies, currency)
	//		}
	//		return
	//	}
	//
	//	for i, m := range groupMembers {
	//		if m.UserID == userID {
	//			for _, v := range m.Balance {
	//				if v > 0 {
	//					if isSponsor = true; isDebtor {
	//						break
	//					}
	//				} else if v < 0 {
	//					if isDebtor = true; isSponsor {
	//						break
	//					}
	//				}
	//			}
	//			userMember = m
	//			groupMembers = groupMembers[:i+copy(groupMembers[i:], groupMembers[i+1:])]
	//			goto userMemberFound
	//		}
	//	}
	//
	//	m.Text = "You are not a member of this space"
	//	logus.Warningf(whc.Context(), m.Text)
	//	return
	//
	//userMemberFound:
	//
	//	if isSponsor && !isDebtor {
	//		groupMembers = filterGroupMembersByBalance(groupMembers, false, balanceCurrencies(userMember.Balance)...)
	//	} else if isDebtor && !isSponsor {
	//		groupMembers = filterGroupMembersByBalance(groupMembers, true, balanceCurrencies(userMember.Balance)...)
	//	}
	//
	//	switch len(groupMembers) {
	//	case 0:
	//		m.Text = "There are no members to settele up with."
	//		//case 1:
	//		//	return settleGroupCounterpartyChosenAction(whc, space, userMember.ContactID)
	//	default:
	//		membersToKeyboard := func() botsgocore.Keyboard {
	//			keyboard := make([][]tgbotapi.InlineKeyboardButton, len(groupMembers))
	//			for i, m := range groupMembers {
	//				keyboard[i] = []tgbotapi.InlineKeyboardButton{
	//					{
	//						Text:         m.Name,
	//						CallbackData: fmt.Sprintf("%v?space=%v&member=%v", SettleGroupCounterpartyChosenCommandCode, space.ContactID, m.ContactID),
	//					},
	//				}
	//			}
	//			return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	//		}
	//
	//		var buf bytes.Buffer
	//		buf.WriteString(whc.Translate(trans.MT_GROUP_LABEL, space.Data.Name) + "\n\n")
	//
	//		switch {
	//		case isSponsor && !isDebtor:
	//			if len(userMember.Balance) == 1 {
	//				for c, v := range userMember.Balance {
	//					buf.WriteString(fmt.Sprintf("You are owed %v %v by this space.\n\n", v, c))
	//				}
	//			}
	//			buf.WriteString("Who from space debtors will pay to you?")
	//		case isDebtor && !isSponsor:
	//			if len(userMember.Balance) == 1 {
	//				for c, v := range userMember.Balance {
	//					buf.WriteString(fmt.Sprintf("You owe %v %v to this space.\n\n", v, c))
	//				}
	//			}
	//			buf.WriteString("Who from space sponsors will collect your debt?")
	//		case isSponsor && isDebtor:
	//			buf.WriteString("Please choose with whom you are going to settle up?")
	//		}
	//		m.Keyboard = membersToKeyboard()
	//		m.Text = buf.String()
	//	}
	//	m.Format = botsfw.MessageFormatHTML

	return
}

var settleGroupCounterpartyChosenCommand = shared_space.SpaceCallbackCommand(
	SettleGroupCounterpartyChosenCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		return settleGroupCounterpartyChosenAction(whc, space, callbackUrl.Query().Get("member"))
	},
)

func settleGroupCounterpartyChosenAction(whc botsfw.WebhookContext, space dbo4spaceus.SpaceEntry, memberID string) (m botsfw.MessageFromBot, err error) {
	err = errors.New("not implemented yet")
	return
	//var userMember, counterpartyMember briefs4splitus.GroupMemberJson
	//userID := whc.AppUserID()
	//for _, m := range space.Data.GetGroupMembers() {
	//	if m.UserID == userID {
	//		userMember = m
	//		if counterpartyMember.ContactID != "" {
	//			break
	//		}
	//	} else if m.ContactID == memberID {
	//		counterpartyMember = m
	//		if userMember.ContactID != "" {
	//			break
	//		}
	//	}
	//}
	//m.IsEdit = whc.InputType() == botsfw.WebhookInputCallbackQuery
	//
	//if userMember.ContactID == "" {
	//	m.Text = "You are not a member of this space."
	//	return
	//} else if counterpartyMember.ContactID == "" {
	//	m.Text = "Selected member has left this space."
	//	return
	//}
	//m.Text = fmt.Sprintf("Have you returned this debt to %v already or you will return it?", counterpartyMember.Name)
	//
	//m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
	//	[]tgbotapi.InlineKeyboardButton{
	//		{
	//			Text:         "I have returned this debt",
	//			CallbackData: fmt.Sprintf("%v?debt=returned&space=%v&member=%v", SettleGroupCounterpartyConfirmedCommandCode, space.ContactID, memberID),
	//		},
	//	},
	//	[]tgbotapi.InlineKeyboardButton{
	//		{
	//			Text:         "I will returned this debt",
	//			CallbackData: fmt.Sprintf("%v?debt=will-return&space=%v&member=%v", SettleGroupCounterpartyConfirmedCommandCode, space.ContactID, memberID),
	//		},
	//	},
	//)
	//logus.Debugf(whc.Context(), "counterpartyMember: %v", counterpartyMember)
	//return
}

var settleGroupCounterpartyConfirmedCommand = shared_space.SpaceCallbackCommand(
	SettleGroupCounterpartyConfirmedCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		q := callbackUrl.Query()
		currency := "RUB" // q.Get("currency")
		return settleGroupCounterpartyConfirmedAction(whc, space, q.Get("member"), money.CurrencyCode(currency))
	},
)

func settleGroupCounterpartyConfirmedAction(whc botsfw.WebhookContext, space dbo4spaceus.SpaceEntry, memberID string, currency money.CurrencyCode) (m botsfw.MessageFromBot, err error) {
	err = errors.New("not implemented yet")
	return
	//var userMember, counterpartyMember briefs4splitus.GroupMemberJson
	//
	//if counterpartyMember, err = space.Data.GetGroupMemberByID(memberID); err != nil {
	//	return
	//}
	//
	//userID := whc.AppUserID()
	//
	//for _, m := range space.Data.GetGroupMembers() {
	//	if m.UserID == userID {
	//		userMember = m
	//		break
	//	}
	//}
	//
	//var debtorID, sponsorID string
	//
	//userBalance := userMember.Balance[currency]
	//counterpartyBalance := counterpartyMember.Balance[currency]
	//
	//if userBalance > 0 && counterpartyBalance < 0 {
	//	debtorID = counterpartyMember.ContactID
	//	sponsorID = userMember.ContactID
	//} else if userBalance < 0 && counterpartyBalance > 0 {
	//	debtorID = userMember.ContactID
	//	sponsorID = counterpartyMember.ContactID
	//} else {
	//	err = errors.New("Balance changed")
	//	return
	//}
	//
	//if err = gaedal.Settle2members(whc.Context(), space.ContactID, debtorID, sponsorID, currency, 700); err != nil {
	//	return
	//}
	//
	//m.Text = "Settled up"
	//m.IsEdit = true
	//logus.Debugf(whc.Context(), "counterpartyMember: %v", counterpartyMember)
	//return
}

//func filterGroupMembersByBalance(members []briefs4splitus.GroupMemberJson, positive bool, currencies ...money.CurrencyCode) (result []briefs4splitus.GroupMemberJson) {
//	result = make([]briefs4splitus.GroupMemberJson, 0, len(members))
//	for _, m := range members {
//		for c, v := range m.Balance {
//			if (positive && v > 0) || (!positive && v < 0) {
//				for _, currency := range currencies {
//					if c == currency {
//						result = append(result, m)
//						break
//					}
//				}
//			}
//		}
//	}
//	return
//}
