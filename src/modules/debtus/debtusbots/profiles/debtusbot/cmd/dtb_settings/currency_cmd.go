package dtb_settings

import (
	"context"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
)

const ASK_CURRENCY_SETTING_COMMAND = "ask-currency-settings"

var AskCurrencySettingsCommand = botsfw.Command{ // TODO: make used
	Code:     ASK_CURRENCY_SETTING_COMMAND,
	Replies:  []botsfw.Command{SetPrimaryCurrency},
	Commands: []string{"\xF0\x9F\x92\xB1"},
	Icon:     emoji.CURRENCY_EXCAHNGE_ICON,
	Title:    trans.COMMAND_TEXT_SETTINGS_PRIMARY_CURRENCY,
	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
		m := whc.NewMessageByCode(trans.MESSAGE_TEXT_ASK_PRIMARY_CURRENCY)
		m.Keyboard = tgbotapi.NewReplyKeyboardUsingStrings([][]string{
			{
				"€ - Euro ",
				"$ - USD",
				"₽ - RUB",
			},
			{
				"Other",
			},
		})
		whc.ChatData().SetAwaitingReplyTo(ASK_CURRENCY_SETTING_COMMAND)
		return m, nil
	},
}

const SET_PRIMARY_CURRENCY_COMMAND = "set-primary-currency"

var SetPrimaryCurrency = botsfw.Command{
	Code: SET_PRIMARY_CURRENCY_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		logus.Debugf(ctx, "SetPrimaryCurrency.Action()")
		whc.ChatData().SetAwaitingReplyTo("")
		primaryCurrency := whc.Input().(botinput.WebhookTextMessage).Text()
		userID := whc.AppUserID()
		userContext := facade.NewUserContext(userID)
		err = dal4userus.RunUserWorker(ctx, userContext, true, func(ctx context.Context, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) error {
			userWorkerParams.UserUpdates, err = userWorkerParams.User.Data.SetPrimaryCurrency(money.CurrencyCode(primaryCurrency))
			return err
		})
		if err != nil {
			return m, err
		}
		return whc.NewMessageByCode(trans.MESSAGE_TEXT_PRIMARY_CURRENCY_IS_SET_TO, whc.Input().(botinput.WebhookTextMessage).Text()), err
	},
}
