package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"net/url"
)

const groupCommandCode = "group"

var groupCommand = botsfw.NewCallbackCommand(groupCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		panic("implement me")
		//// we can't use GroupCallbackCommand as we have parameter id=[first|last|<id>]
		//ctx := whc.Context()
		//logus.Debugf(ctx, "groupCommand.CallbackAction()")
		//
		//var user botsfwmodels.AppUserData
		//if user, err = whc.AppUserData(); err != nil {
		//	return
		//}
		//appUserEntity := user.(*models.DebutsAppUserDataOBSOLETE) // TODO: Create shortcut function
		//
		//groups := appUserEntity.ActiveGroups()
		//
		//if len(groups) == 0 {
		//	return groupsAction(whc, true, 0)
		//}
		//
		//query := callbackUrl.Query()
		//
		//id := query.Get("id")
		//
		//var (
		//	i             int
		//	userGroupJson models.UserGroupJson
		//)
		//switch id {
		//case "first":
		//	i = 0
		//case "last":
		//	i = len(groups) - 1
		//default:
		//	userGroupJson.ContactID = id
		//	for j, g := range groups {
		//		if g.ContactID == userGroupJson.ContactID {
		//			i = j
		//		}
		//	}
		//}
		//userGroupJson = groups[i]
		//
		//do := query.Get("do")
		//switch do {
		//case "leave":
		//	if _, _, err = facade4debtus.Group.LeaveGroup(c, userGroupJson.ContactID, whc.AppUserID()); err != nil {
		//		if err == facade4debtus.ErrAttemptToLeaveUnsettledGroup {
		//			err = nil
		//			m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{Text: "Please settle group debts before leaving it."})
		//		}
		//		return
		//	}
		//	return groupsAction(whc, true, 0)
		//}
		//
		//var group models.Group
		//
		//if group, err = dtdal.Group.GetGroupByID(c, nil, userGroupJson.ContactID); err != nil {
		//	return
		//}
		//
		//buf := new(bytes.Buffer)
		//
		//_, _ = fmt.Fprintf(buf, "<b>Group #%d</b>: %v", i+1, userGroupJson.Name)
		//var groupMemberJson models.GroupMemberJson
		//if groupMemberJson, err = group.Data.GetGroupMemberByUserID(whc.AppUserID()); err != nil {
		//	return
		//}
		//writeBalanceSide := func(title string, sign decimal.Decimal64p2, b money.Balance) {
		//	if len(b) > 0 {
		//		fmt.Fprintf(buf, "\n<b>%v</b>: ", title)
		//		if len(b) == 1 {
		//			for currency, value := range b {
		//				fmt.Fprintf(buf, "%v %v", sign*value, currency)
		//			}
		//		} else {
		//			for currency, value := range b {
		//				fmt.Fprintf(buf, "\n%v %v", sign*value, currency)
		//			}
		//		}
		//	}
		//}
		//writeBalanceSide("Owed to me", +1, groupMemberJson.Balance.OnlyPositive())
		//writeBalanceSide("I owe", -1, groupMemberJson.Balance.OnlyNegative())
		//fmt.Fprintf(buf, "\n<b>Members</b>: %v", group.Data.MembersCount)
		//
		//m.Text = buf.String()
		//
		//m.IsEdit = true
		//m.Format = botsfw.MessageFormatHTML
		//tgKeyboard := tgbotapi.NewInlineKeyboardMarkup(groupsNavButtons(whc, groups, userGroupJson.ContactID))
		//tgKeyboard.InlineKeyboard = append(tgKeyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		//	{
		//		Text:         whc.Translate("Leave group"),
		//		CallbackData: CallbackLink.ToGroup(groups[len(groups)-1].ContactID, true) + "&do=leave",
		//	},
		//})
		//m.Keyboard = tgKeyboard
		//return
	},
)

//func groupsNavButtons(translator i18n.SingleLocaleTranslator, groups []models.UserGroupJson, currentGroupID string) []tgbotapi.InlineKeyboardButton {
//	var currentGroupIndex = -1
//	if currentGroupID != "" {
//
//		for i, group := range groups {
//			if group.ContactID == currentGroupID {
//				currentGroupIndex = i
//				break
//			}
//		}
//	}
//	buttons := []tgbotapi.InlineKeyboardButton{}
//	if len(groups) > 0 || currentGroupIndex < 0 {
//		switch currentGroupIndex {
//		case -1:
//			buttons = append(buttons, tgbotapi.InlineKeyboardButton{
//				Text:         "⬅️",
//				CallbackData: CallbackLink.ToGroup(groups[len(groups)-1].ContactID, true),
//			})
//		case 0:
//			buttons = append(buttons, tgbotapi.InlineKeyboardButton{
//				Text:         "⬅️",
//				CallbackData: groupsCommandCode + "?edit=1",
//			})
//		default:
//			buttons = append(buttons, tgbotapi.InlineKeyboardButton{
//				Text:         "⬅️",
//				CallbackData: CallbackLink.ToGroup(groups[currentGroupIndex-1].ContactID, true),
//			})
//		}
//	}
//	if currentGroupID != "" {
//		buttons = append(buttons, tgbotapi.InlineKeyboardButton{
//			Text:         translator.Translate(trans.COMMAND_TEXT_GROUPS),
//			CallbackData: groupsCommandCode + "?edit=1",
//		})
//
//	}
//	if len(groups) > 0 || currentGroupIndex < 0 {
//		switch currentGroupIndex {
//		case -1:
//			buttons = append(buttons, tgbotapi.InlineKeyboardButton{
//				Text:         "➡️",
//				CallbackData: CallbackLink.ToGroup(groups[0].ContactID, true),
//			})
//		case len(groups) - 1:
//			buttons = append(buttons, tgbotapi.InlineKeyboardButton{
//				Text:         "➡️",
//				CallbackData: groupsCommandCode + "?edit=1",
//			})
//		default:
//			buttons = append(buttons, tgbotapi.InlineKeyboardButton{
//				Text:         "➡️",
//				CallbackData: CallbackLink.ToGroup(groups[currentGroupIndex+1].ContactID, true),
//			})
//		}
//	}
//	return buttons
//}
