package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/anybot/cmds4anybot"
	"github.com/strongo/logus"
	"strings"
)

func startInGroupAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	panic("implement me")
	//ctx := whc.Context()
	//logus.Debugf(ctx, "splitus.startInGroupAction()")
	//var group models.Group
	//if group, err = shared_space.GetGroup(whc, nil); err != nil {
	//	return
	//}
	////var appUserData botsfwmodels.AppUserData
	////if appUserData, err = whc.AppUserData(); err != nil {
	////	return
	////}
	//
	////appUser := appUserData.(*models.DebutsAppUserDataOBSOLETE)
	//
	//var botUser record.DataWithID[string, botsfwmodels.BotUserData]
	//
	//if botUser, err = whc.BotUser(); err != nil && !dal.IsNotFound(err) {
	//	return
	//}
	//
	//if group, _, err = facade4debtus.Group.AddUsersToTheGroupAndOutstandingBills(c, group.ContactID, []facade4debtus.NewUser{
	//	{
	//		//Name:        appUserData.FullName(),
	//		BotUserData: botUser.Data,
	//		ChatMember:  whc.Input().GetSender(),
	//	},
	//}); err != nil {
	//	err = fmt.Errorf("%w: failed to add user to the group", err)
	//	return
	//}
	//m.Text = whc.Translate(trans.MESSAGE_TEXT_HI) +
	//	"\n\n" + whc.Translate(trans.SPLITUS_TEXT_HI_IN_GROUP) +
	//	"\n\n<b>" + whc.Translate(trans.MESSAGE_TEXT_ASK_PRIMARY_CURRENCY_FOR_GROUP) + "</b>"
	//
	//m.Format = botsfw.MessageFormatHTML
	//m.Keyboard = currenciesInlineKeyboard(
	//	GroupSettingsSetCurrencyCommandCode+"?start=y&group="+group.ContactID,
	//	[]tgbotapi.InlineKeyboardButton{
	//		{
	//			Text: whc.Translate(trans.BT_OTHER_CURRENCY),
	//			URL:  fmt.Sprintf("https://t.me/%v?start=%v__group=%v", whc.GetBotCode(), GroupSettingsChooseCurrencyCommandCode, group.ContactID),
	//		},
	//	},
	//)
	//return
}

func startInBotAction(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
	logus.Debugf(whc.Context(), "splitus.startInBotAction() => startParams: %v", startParams)
	if len(startParams) > 0 {
		switch {
		case strings.HasPrefix(startParams[0], "bill-"):
			return startBillAction(whc, startParams[0])
		case startParams[0] == SettleGroupAskForCounterpartyCommandCode:
			return settleGroupStartAction(whc, startParams[1:])
		}
	}
	err = cmds4anybot.ErrUnknownStartParam
	return
}
