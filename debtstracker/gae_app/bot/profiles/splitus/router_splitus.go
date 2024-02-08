package splitus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/i18n"
)

var botParams = shared_all.BotParams{
	StartInGroupAction: startInGroupAction,
	StartInBotAction:   startInBotAction,
	//GetGroupBillCardInlineKeyboard:   getGroupBillCardInlineKeyboard,
	//GetPrivateBillCardInlineKeyboard: getPrivateBillCardInlineKeyboard,
	//DelayUpdateBillCardOnUserJoin:    delayUpdateBillCardOnUserJoin,
	//OnAfterBillCurrencySelected:      getWhoPaidInlineKeyboard,
	//ShowGroupMembers:                 showGroupMembers,
	InBotWelcomeMessage: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		var user *models.DebutsAppUserDataOBSOLETE
		if user, err = shared_all.GetUser(whc); err != nil {
			return
		}
		m.Text = whc.Translate(
			trans.MESSAGE_TEXT_HI_USERNAME, user.FirstName) + " " + whc.Translate(trans.SPLITUS_TEXT_HI) +
			"\n\n" + whc.Translate(trans.SPLITUS_TEXT_ABOUT_ME_AND_CO) +
			"\n\n" + whc.Translate(trans.SPLITUS_TG_COMMANDS)
		m.Format = botsfw.MessageFormatHTML
		m.IsEdit = true

		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
					whc.CommandText(trans.COMMAND_TEXT_NEW_BILL, emoji.MEMO_ICON),
					"",
				),
			},
			[]tgbotapi.InlineKeyboardButton{
				shared_group.NewGroupTelegramInlineButton(whc, 0),
			},
		)
		return
	},
	//
	//
	//
	SetMainMenu: setMainMenu,
}

var Router = botsfw.NewWebhookRouter(
	map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText: {
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
		botsfw.WebhookInputCallbackQuery: {
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
			groupSplitCommand,
			joinGroupCommand,
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
			groupSettingsSetCurrencyCommand(botParams),
			groupsCommand,
			settingsCommand,
			groupSettingsChooseCurrencyCommand,
			settleGroupAskForCounterpartyCommand,
			settleGroupCounterpartyChosenCommand,
			settleGroupCounterpartyConfirmedCommand,
			settleBillsCommand,
		},
		botsfw.WebhookInputInlineQuery: {
			inlineQueryCommand,
		},
		botsfw.WebhookInputChosenInlineResult: {
			chosenInlineResultCommand,
		},
		botsfw.WebhookInputNewChatMembers: {
			newChatMembersCommand,
		},
	},
	func() string { return "Please report any errors to @SplitusGroup" },
)

func init() {
	shared_all.AddSharedRoutes(Router, botParams)
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
