package botcmds4splitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/debtusbots/profiles/shared_space"
	"github.com/strongo/logus"
	"net/url"
)

const CURRENCY_PARAM_NAME = "currency"

func currenciesInlineKeyboard(callbackDataPrefix string, more ...[]tgbotapi.InlineKeyboardButton) *tgbotapi.InlineKeyboardMarkup {
	currencyButton := func(code, flag string) tgbotapi.InlineKeyboardButton {
		btn := tgbotapi.InlineKeyboardButton{CallbackData: callbackDataPrefix + "&" + CURRENCY_PARAM_NAME + "=" + code}
		if flag == "" {
			btn.Text = code
		} else {
			btn.Text = flag + " " + code
		}
		return btn
	}

	usdRow := []tgbotapi.InlineKeyboardButton{
		currencyButton("USD", "🇺🇸"),
		currencyButton("AUD", "🇦🇺"),
		currencyButton("CAD", "🇨🇦"),
		currencyButton("GBP", "🇬🇧"),
	}

	eurRow := []tgbotapi.InlineKeyboardButton{
		currencyButton("EUR", "🇪🇺"),
		currencyButton("CHF", "🇨🇭"),
		currencyButton("NOK", "🇳🇴"),
		currencyButton("SEK", "🇸🇪"),
	}

	eurRow2 := []tgbotapi.InlineKeyboardButton{
		currencyButton("BGN", "🇧🇬"),
		currencyButton("HUF", "🇭🇺"),
		currencyButton("PLN", "🇵🇱"),
		currencyButton("RON", "🇷🇴"),
	}

	rubRow := []tgbotapi.InlineKeyboardButton{
		currencyButton("RUB", "🇷🇺"),
		currencyButton("BYN", "🇧🇾"),
		currencyButton("UAH", "🇺🇦"),
		currencyButton("MDL", "🇲🇩"),
	}

	exUSSR := []tgbotapi.InlineKeyboardButton{
		currencyButton("KGS", "🇰🇬"),
		currencyButton("KZT", "🇰🇿"),
		currencyButton("TJS", "🇹🇯"),
		currencyButton("UZS", "🇺🇿"),
	}

	asiaRow := []tgbotapi.InlineKeyboardButton{
		currencyButton("CNY", "🇨🇳"),
		currencyButton("JPY", "🇯🇵"),
		currencyButton("IDR", "🇮🇩"),
		currencyButton("KRW", "🇰🇷"),
		//currencyButton("VND", "🇻🇳"),
	}

	keyboard := append([][]tgbotapi.InlineKeyboardButton{
		usdRow,
		eurRow,
		rubRow,
		exUSSR,
		eurRow2,
		asiaRow,
	}, more...)

	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

const (
	GroupSettingsChooseCurrencyCommandCode = "grp-stngs-chs-ccy"
	GroupSettingsSetCurrencyCommandCode    = "grp-stngs-set-ccy"
)

var groupSettingsChooseCurrencyCommand = shared_space.SpaceCallbackCommand(GroupSettingsChooseCurrencyCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		m.IsEdit = true
		m.Text = whc.Translate(trans.MESSAGE_TEXT_ASK_PRIMARY_CURRENCY)
		m.Keyboard = currenciesInlineKeyboard(
			GroupSettingsSetCurrencyCommandCode+"?space="+space.ID,
			[]tgbotapi.InlineKeyboardButton{
				{
					Text: whc.Translate(trans.BT_OTHER_CURRENCY),
					URL:  fmt.Sprintf("https://t.me/%v?start=", whc.GetBotCode()) + GroupSettingsChooseCurrencyCommandCode,
				},
			},
		)
		return
	},
)

func groupSettingsSetCurrencyCommand() botsfw.Command {
	return botsfw.Command{
		Code: GroupSettingsSetCurrencyCommandCode,
		CallbackAction: shared_space.NewSpaceCallbackAction(func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
			currency := money.CurrencyCode(callbackUrl.Query().Get(CURRENCY_PARAM_NAME))
			if space.Data.PrimaryCurrency != currency {
				ctx := whc.Context()
				user := facade.NewUserContext(whc.AppUserID())
				if err := dal4spaceus.RunSpaceWorkerWithUserContext(ctx, user, space.ID, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceWorkerParams) (err error) {
					params.SpaceUpdates, err = space.Data.SetPrimaryCurrency(currency)
					return
				}); err != nil {
					logus.Errorf(whc.Context(), "failed to change space default currency: %v", err)
				} else {
					logus.Debugf(ctx, "Default currency for space %v updated to: %v", space.ID, currency)
				}
			}
			if callbackUrl.Query().Get("start") == "y" {
				return onStartCallbackInGroup(whc, space)
			} else {
				return SpaceSettingsAction(whc, space, true)
			}
		}),
	}
}

func onStartCallbackInGroup(whc botsfw.WebhookContext, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
	// This links Telegram ChatID and ChatInstance
	err = errors.New("onStartCallbackInGroup is not implemented yet")
	return
	// if twhc, ok := whc.(*telegram.tgWebhookContext); ok {
	// 	if err = twhc.CreateOrUpdateTgChatInstance(); err != nil {
	// 		return
	// 	}
	// }
	// return inGroupWelcomeMessage(whc, group)
}

//func inGroupWelcomeMessage(whc botsfw.WebhookContext, group models.GroupEntry) (m botsfw.MessageFromBot, err error) {
//	m, err = SpaceSettingsAction(whc, group, false)
//	if err != nil {
//		return
//	}
//	if _, err = whc.Responder().SendMessage(whc.Context(), m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
//		return
//	}
//
//	return whc.NewEditMessage(whc.Translate(trans.MESSAGE_TEXT_HI)+
//		"\n\n"+whc.Translate(trans.SPLITUS_TEXT_HI_IN_GROUP)+
//		"\n\n"+whc.Translate(trans.SPLITUS_TEXT_ABOUT_ME_AND_CO),
//		botsfw.MessageFormatHTML)
//}
