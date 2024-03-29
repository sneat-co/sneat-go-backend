package debtus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_admin"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_invite"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_retention"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_settings"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_transfer"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
)

var botParams = shared_all.BotParams{
	//GetGroupBillCardInlineKeyboard:   getGroupBillCardInlineKeyboard,
	//GetPrivateBillCardInlineKeyboard: getPrivateBillCardInlineKeyboard,
	//DelayUpdateBillCardOnUserJoin:    delayUpdateBillCardOnUserJoin,
	//OnAfterBillCurrencySelected:      getWhoPaidInlineKeyboard,
	//ShowGroupMembers:                 showGroupMembers,
	HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return dtb_general.HelpCommandAction(whc, true)
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
	StartInBotAction: dtb_settings.StartInBotAction,
	SetMainMenu:      dtb_general.SetMainMenuKeyboard,
}

func init() {
	shared_all.AddSharedRoutes(Router, botParams)
}

var textAndContactCommands = []botsfw.Command{ // TODO: Check for Action || CallbackAction and register accordingly.
	//OnboardingAskInviteChannelCommand, // We need it as otherwise do not handle replies.
	//SetPreferredLanguageCommand,
	//OnboardingAskInviteCodeCommand,
	//OnboardingCheckInviteCommand,
	//
	dtb_general.FeedbackCommand,
	dtb_general.FeedbackTextCommand,
	dtb_general.DeleteAllCommand,
	dtb_general.BetaCommand,
	//
	dtb_admin.AdminCommand,
	//
	dtb_settings.SettingsCommand,
	dtb_settings.LoginPinCommand,
	//dtb_settings.OnboardingTellAboutInviteCodeCommand, // We need it as otherwise do not handle replies. Consider incorporate to StartCommand?
	dtb_settings.FixBalanceCommand,
	dtb_settings.ContactsListCommand,
	//
	//dtb_settings.AskCurrencySettingsCommand,
	//
	dtb_general.Login2WebCommand,
	dtb_general.MainMenuCommand,
	dtb_general.ClearCommand,
	dtb_general.AdsCommand,
	//
	dtb_transfer.StartLendingWizardCommand,
	dtb_transfer.StartBorrowingWizardCommand,
	dtb_transfer.StartReturnWizardCommand,
	dtb_transfer.BalanceCommand,
	dtb_transfer.HistoryCommand,
	dtb_transfer.CancelTransferWizardCommand,
	dtb_transfer.ParseTransferCommand,
	dtb_transfer.AskHowMuchHaveBeenReturnedCommand,
	dtb_transfer.SetNextReminderDateCallbackCommand,
	//
	dtb_retention.DeleteUserCommand,
	//
	dtb_invite.InviteCommand,
	dtb_transfer.AskEmailForReceiptCommand,       // TODO: Should it be in dtb_transfer?
	dtb_transfer.AskPhoneNumberForReceiptCommand, // TODO: Should it be in dtb_transfer?
	dtb_invite.CreateMassInviteCommand,
	//
}

var callbackCommands = []botsfw.Command{
	dtb_general.MainMenuCommand,
	dtb_general.PleaseWaitCommand,
	//dtb_invite.InviteCommand,
	//
	dtb_settings.SettingsCommand,
	dtb_settings.ContactsListCommand,
	//
	//dtb_fbm.FbmGetStartedCommand, // TODO: Move command to other package?
	//dtb_fbm.FbmMainMenuCommand,
	//dtb_fbm.FbmDebtsCommand,
	//dtb_fbm.FbmBillsCommand,
	//dtb_fbm.FbmSettingsCommand,
	//
	dtb_invite.CreateMassInviteCommand,
	dtb_invite.AskInviteAddressCallbackCommand,
	//
	dtb_transfer.CreateReceiptIfNoInlineNotificationCommand,
	dtb_transfer.SendReceiptCallbackCommand,
	//dtb_transfer.AcknowledgeReceiptCommand,
	dtb_transfer.ViewReceiptInTelegramCallbackCommand,
	dtb_transfer.ChangeReceiptAnnouncementLangCommand,
	dtb_transfer.ViewReceiptCallbackCommand,
	dtb_transfer.AcknowledgeReceiptCallbackCommand,
	dtb_transfer.TransferHistoryCallbackCommand,
	dtb_transfer.AskForInterestAndCommentCallbackCommand,
	dtb_transfer.BalanceCallbackCommand,
	dtb_transfer.DueReturnsCallbackCommand,
	dtb_transfer.ReturnCallbackCommand,
	dtb_transfer.EnableReminderAgainCallbackCommand,
	dtb_transfer.SetNextReminderDateCallbackCommand,
	//dtb_transfer.CounterpartyNoTelegramCommand,
	dtb_transfer.RemindAgainCallbackCommand,
	//dtb_general.FeedbackCallbackCommand,
	dtb_general.FeedbackCommand,
	dtb_general.CanYouRateCommand,
	dtb_general.FeedbackTextCommand,
	shared_all.AddReferrerCommand,
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
			dtb_invite.ChosenInlineResultCommand,
		},
		botsfw.WebhookInputNewChatMembers: {
			newChatMembersCommand,
		},
	},
	func() string { return "Please report any errors to @DebtsTrackerGroup" },
)
