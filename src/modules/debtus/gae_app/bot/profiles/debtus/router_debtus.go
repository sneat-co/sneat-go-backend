package debtus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/cmd/dtb_admin"
	dtb_general2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/cmd/dtb_general"
	dtb_invite2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/cmd/dtb_invite"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/cmd/dtb_retention"
	dtb_settings2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/cmd/dtb_settings"
	dtb_transfer2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/cmd/dtb_transfer"
	shared_all2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/shared_all"
)

var botParams = shared_all2.BotParams{
	//GetGroupBillCardInlineKeyboard:   getGroupBillCardInlineKeyboard,
	//GetPrivateBillCardInlineKeyboard: getPrivateBillCardInlineKeyboard,
	//DelayUpdateBillCardOnUserJoin:    delayUpdateBillCardOnUserJoin,
	//OnAfterBillCurrencySelected:      getWhoPaidInlineKeyboard,
	//ShowGroupMembers:                 showGroupMembers,
	HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return dtb_general2.HelpCommandAction(whc, true)
	},
	//InGroupWelcomeMessage: func(whc botsfw.WebhookContext, group models.Group) (m botsfw.MessageFromBot, err error) {
	//	m, err = shared_all.GroupSettingsAction(whc, group, false)
	//	if err != nil {
	//		return
	//	}
	//	if _, err = whc.Responder().SendMessage(whc.Context(), m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
	//		return
	//	}
	//
	//	return whc.NewEditMessage(whc.Translate(trans.MESSAGE_TEXT_HI)+
	//		"\n\n"+ whc.Translate(trans.SPLITUS_TEXT_HI_IN_GROUP)+
	//		"\n\n"+ whc.Translate(trans.SPLITUS_TEXT_ABOUT_ME_AND_CO),
	//		bots.MessageFormatHTML)
	//},
	InBotWelcomeMessage: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Hi there"
		m.Format = botsfw.MessageFormatHTML
		//m.IsEdit = true
		return
	},
	//
	//
	//
	StartInBotAction: dtb_settings2.StartInBotAction,
	SetMainMenu:      dtb_general2.SetMainMenuKeyboard,
}

func init() {
	shared_all2.AddSharedRoutes(Router, botParams)
}

var textAndContactCommands = []botsfw.Command{ // TODO: Check for Action || CallbackAction and register accordingly.
	//OnboardingAskInviteChannelCommand, // We need it as otherwise do not handle replies.
	//SetPreferredLanguageCommand,
	//OnboardingAskInviteCodeCommand,
	//OnboardingCheckInviteCommand,
	//
	dtb_general2.FeedbackCommand,
	dtb_general2.FeedbackTextCommand,
	dtb_general2.DeleteAllCommand,
	dtb_general2.BetaCommand,
	//
	dtb_admin.AdminCommand,
	//
	dtb_settings2.SettingsCommand,
	dtb_settings2.LoginPinCommand,
	//dtb_settings.OnboardingTellAboutInviteCodeCommand, // We need it as otherwise do not handle replies. Consider incorporate to StartCommand?
	dtb_settings2.FixBalanceCommand,
	dtb_settings2.ContactsListCommand,
	//
	//dtb_settings.AskCurrencySettingsCommand,
	//
	dtb_general2.Login2WebCommand,
	dtb_general2.MainMenuCommand,
	dtb_general2.ClearCommand,
	dtb_general2.AdsCommand,
	//
	dtb_transfer2.StartLendingWizardCommand,
	dtb_transfer2.StartBorrowingWizardCommand,
	dtb_transfer2.StartReturnWizardCommand,
	dtb_transfer2.BalanceCommand,
	dtb_transfer2.HistoryCommand,
	dtb_transfer2.CancelTransferWizardCommand,
	dtb_transfer2.ParseTransferCommand,
	dtb_transfer2.AskHowMuchHaveBeenReturnedCommand,
	dtb_transfer2.SetNextReminderDateCallbackCommand,
	//
	dtb_retention.DeleteUserCommand,
	//
	dtb_invite2.InviteCommand,
	dtb_transfer2.AskEmailForReceiptCommand,       // TODO: Should it be in dtb_transfer?
	dtb_transfer2.AskPhoneNumberForReceiptCommand, // TODO: Should it be in dtb_transfer?
	dtb_invite2.CreateMassInviteCommand,
	//
}

var callbackCommands = []botsfw.Command{
	dtb_general2.MainMenuCommand,
	dtb_general2.PleaseWaitCommand,
	//dtb_invite.InviteCommand,
	//
	dtb_settings2.SettingsCommand,
	dtb_settings2.ContactsListCommand,
	//
	//dtb_fbm.FbmGetStartedCommand, // TODO: Move command to other package?
	//dtb_fbm.FbmMainMenuCommand,
	//dtb_fbm.FbmDebtsCommand,
	//dtb_fbm.FbmBillsCommand,
	//dtb_fbm.FbmSettingsCommand,
	//
	dtb_invite2.CreateMassInviteCommand,
	dtb_invite2.AskInviteAddressCallbackCommand,
	//
	dtb_transfer2.CreateReceiptIfNoInlineNotificationCommand,
	dtb_transfer2.SendReceiptCallbackCommand,
	//dtb_transfer.AcknowledgeReceiptCommand,
	dtb_transfer2.ViewReceiptInTelegramCallbackCommand,
	dtb_transfer2.ChangeReceiptAnnouncementLangCommand,
	dtb_transfer2.ViewReceiptCallbackCommand,
	dtb_transfer2.AcknowledgeReceiptCallbackCommand,
	dtb_transfer2.TransferHistoryCallbackCommand,
	dtb_transfer2.AskForInterestAndCommentCallbackCommand,
	dtb_transfer2.BalanceCallbackCommand,
	dtb_transfer2.DueReturnsCallbackCommand,
	dtb_transfer2.ReturnCallbackCommand,
	dtb_transfer2.EnableReminderAgainCallbackCommand,
	dtb_transfer2.SetNextReminderDateCallbackCommand,
	//dtb_transfer.CounterpartyNoTelegramCommand,
	dtb_transfer2.RemindAgainCallbackCommand,
	//dtb_general.FeedbackCallbackCommand,
	dtb_general2.FeedbackCommand,
	dtb_general2.CanYouRateCommand,
	dtb_general2.FeedbackTextCommand,
	shared_all2.AddReferrerCommand,
}

var Router = botsfw.NewWebhookRouter(
	map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText:          textAndContactCommands,
		botsfw.WebhookInputContact:       textAndContactCommands,
		botsfw.WebhookInputCallbackQuery: callbackCommands,
		//
		botsfw.WebhookInputInlineQuery: {
			InlineQueryCommand,
		},
		botsfw.WebhookInputChosenInlineResult: {
			dtb_invite2.ChosenInlineResultCommand,
		},
		botsfw.WebhookInputNewChatMembers: {
			newChatMembersCommand,
		},
	},
	func() string { return "Please report any errors to @DebtsTrackerGroup" },
)
