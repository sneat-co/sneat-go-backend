package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/strongo/i18n"
)

var botParams = cmds4anybot.BotParams{
	StartInGroupAction: startInGroupAction,
	StartInBotAction:   startInBotAction,
	HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "HelpCommandAction is not implemented yet"
		return
	},
	//GetGroupBillCardInlineKeyboard:   getGroupBillCardInlineKeyboard,
	//GetPrivateBillCardInlineKeyboard: getPrivateBillCardInlineKeyboard,
	//DelayUpdateBillCardOnUserJoin:    delayUpdateBillCardOnUserJoin,
	//OnAfterBillCurrencySelected:      getWhoPaidInlineKeyboard,
	//ShowGroupMembers:                 showGroupMembers,
	GetWelcomeMessageText: func(whc botsfw.WebhookContext) (text string, err error) {
		var user dbo4userus.UserEntry
		if user, err = cmds4anybot.GetUser(whc); err != nil {
			return
		}
		text = whc.Translate(
			trans.MESSAGE_TEXT_HI_USERNAME, user.Data.Names.FirstName) + " " + whc.Translate(trans.SPLITUS_TEXT_HI) +
			"\n\n" + whc.Translate(trans.SPLITUS_TEXT_ABOUT_ME_AND_CO) +
			"\n\n" + whc.Translate(trans.SPLITUS_TG_COMMANDS)

		return
	},
	//
	//
	//
	SetMainMenu: func(whc botsfw.WebhookContext, messageText string, showHint bool) (m botsfw.MessageFromBot, err error) {
		setMainMenu(whc, &m)
		return
	},
}

var Router = botsfw.NewWebhookRouter(
	func() string { return "Please report any errors to @SplitusGroup" },
)

func init() {
	cmds4anybot.AddSharedCommands(Router, botParams)
	commandsByType := map[botinput.WebhookInputType][]botsfw.Command{
		// TODO: Move input types inside commands and register as slice
		botinput.WebhookInputText: {
			EditedBillCardHookCommand,
			billsCommand,
			groupBalanceCommand,
			menuCommand,
			setBillDueDateCommand,
			groupsCommand,
			settingsCommand,
			settleBillsCommand,
			outstandingBalanceCommand,
		},
		botinput.WebhookInputCallbackQuery: {
			joinBillCommand,
			closeBillCommand,
			editBillCommand,
			newBillCommand,
			groupBalanceCommand,
			billsCommand,
			billSharesCommand,
			billSplitModesListCommand,
			finalizeBillCommand,
			deleteBillCommand,
			restoreBillCommand,
			billChangeSplitModeCommand,
			changeBillPayerCommand,
			spaceSplitCommand,
			joinSpaceCommand,
			//billCardCommand,
			setBillCurrencyCommand,
			groupCommand,
			leaveGroupCommand,
			billCardCommand,
			billMembersCommand,
			inviteToBillCommand,
			setBillDueDateCommand,
			changeBillTotalCommand,
			addBillComment,
			groupMembersCommand,
			groupSettingsSetCurrencyCommand(),
			groupsCommand,
			settingsCommand,
			groupSettingsChooseCurrencyCommand,
			settleGroupAskForCounterpartyCommand,
			settleGroupCounterpartyChosenCommand,
			settleGroupCounterpartyConfirmedCommand,
			settleBillsCommand,
		},
		botinput.WebhookInputInlineQuery: {
			inlineQueryCommand,
		},
		botinput.WebhookInputChosenInlineResult: {
			chosenInlineResultCommand,
		},
		botinput.WebhookInputNewChatMembers: {
			newChatMembersCommand,
		},
	}
	Router.AddCommandsGroupedByType(commandsByType)
}

func getWhoPaidInlineKeyboard(translator i18n.SingleLocaleTranslator, billID string) *tgbotapi.InlineKeyboardMarkup {
	callbackDataPrefix := billCallbackCommandData(joinBillCommandCode, billID)
	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "✋ " + translator.Translate(trans.BUTTON_TEXT_I_PAID_FOR_THE_BILL), CallbackData: callbackDataPrefix + "&i=paid"}},
			{{Text: "🙏 " + translator.Translate(trans.BUTTON_TEXT_I_OWE_FOR_THE_BILL), CallbackData: callbackDataPrefix + "&i=owe"}},
			{{Text: "🚫 " + translator.Translate(trans.BUTTON_TEXT_I_DO_NOT_SHARE_THIS_BILL), CallbackData: billCallbackCommandData(leaveBillCommandCode, billID)}},
		},
	}
}
