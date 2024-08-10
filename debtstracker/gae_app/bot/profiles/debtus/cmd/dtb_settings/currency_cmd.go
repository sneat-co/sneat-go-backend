package dtb_settings

import (
	"context"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
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
		c := whc.Context()
		logus.Debugf(c, "SetPrimaryCurrency.Action()")
		whc.ChatData().SetAwaitingReplyTo("")
		primaryCurrency := whc.Input().(botsfw.WebhookTextMessage).Text()
		if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
			var user models.AppUser

			//goland:noinspection GoDeprecation
			if user, err = facade2debtus.User.GetUserByID(c, tx, whc.AppUserID()); err != nil {
				return
			}
			user.Data.PrimaryCurrency = primaryCurrency
			return facade2debtus.User.SaveUser(c, tx, user)
		}, nil); err != nil {
			return
		}
		return whc.NewMessageByCode(trans.MESSAGE_TEXT_PRIMARY_CURRENCY_IS_SET_TO, whc.Input().(botsfw.WebhookTextMessage).Text()), nil
	},
}
