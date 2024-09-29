package cmds4anybot

import "github.com/bots-go-framework/bots-fw/botsfw"

type SetMainMenuFunc = func(whc botsfw.WebhookContext, messageText string, showHint bool) (m botsfw.MessageFromBot, err error)
type StartInBotActionFunc = func(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error)

type WelcomeMessageProvider = func(whc botsfw.WebhookContext) (text string, err error)

// BotParams defines parameters to be defined by a bot to be able to use shared_all package
// This is supposed to be passed only to AddSharedCommands and not to be passed down to other functions
type BotParams struct {
	GetWelcomeMessageText WelcomeMessageProvider
	StartInBotAction      StartInBotActionFunc
	StartInGroupAction    botsfw.CommandAction
	HelpCommandAction     botsfw.CommandAction
	HelpCallbackAction    botsfw.CallbackAction
	SetMainMenu           SetMainMenuFunc

	//GetGroupBillCardInlineKeyboard   func(translator i18n.SingleLocaleTranslator, bill models.Bill) *tgbotapi.InlineKeyboardMarkup
	//GetPrivateBillCardInlineKeyboard func(translator i18n.SingleLocaleTranslator, botCode string, bill models.Bill) *tgbotapi.InlineKeyboardMarkup
	//OnAfterBillCurrencySelected      func(translator i18n.SingleLocaleTranslator, billID string) *tgbotapi.InlineKeyboardMarkup
	//DelayUpdateBillCardOnUserJoin    func(ctx context.Context, billID string, message string) error
	//ShowGroupMembers                 func(whc botsfw.WebhookContext, group models.Group, isEdit bool) (m botsfw.MessageFromBot, err error)
	//InGroupWelcomeMessage func(whc botsfw.WebhookContext, group models.Group) (m botsfw.MessageFromBot, err error)
}

func (v *BotParams) Validate() {
	if v.StartInBotAction == nil {
		panic("StartInBotAction is not set")
	}
	if v.StartInGroupAction == nil {
		panic("StartInGroupAction is not set")
	}
	if v.GetWelcomeMessageText == nil {
		panic("GetWelcomeMessageText is not set")
	}
	if v.HelpCommandAction == nil {
		panic("HelpCommandAction is not set")
	}
	if v.SetMainMenu == nil {
		panic("SetMainMenu is not set")
	}
}
