package splitus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
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

var groupSettingsChooseCurrencyCommand = shared_group.GroupCallbackCommand(GroupSettingsChooseCurrencyCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.GroupEntry) (m botsfw.MessageFromBot, err error) {
		m.IsEdit = true
		m.Text = whc.Translate(trans.MESSAGE_TEXT_ASK_PRIMARY_CURRENCY)
		m.Keyboard = currenciesInlineKeyboard(
			GroupSettingsSetCurrencyCommandCode+"?group="+group.ID,
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

func groupSettingsSetCurrencyCommand(params shared_all.BotParams) botsfw.Command {
	return botsfw.Command{
		Code: GroupSettingsSetCurrencyCommandCode,
		CallbackAction: shared_group.NewGroupCallbackAction(func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.GroupEntry) (m botsfw.MessageFromBot, err error) {
			currency := money.CurrencyCode(callbackUrl.Query().Get(CURRENCY_PARAM_NAME))
			if group.Data.DefaultCurrency != currency {
				c := whc.Context()
				if err := facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
					if group, err = dtdal.Group.GetGroupByID(c, tx, group.ID); err != nil {
						return
					}
					if group.Data.DefaultCurrency != currency {
						group.Data.DefaultCurrency = currency
						if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
							return
						}
					}
					return
				}); err != nil {
					logus.Errorf(whc.Context(), "failed to change group default currency: %v", err)
				} else {
					logus.Debugf(c, "Default currency for group %v updated to: %v", group.ID, currency)
				}
			}
			if callbackUrl.Query().Get("start") == "y" {
				return onStartCallbackInGroup(whc, group)
			} else {
				return GroupSettingsAction(whc, group, true)
			}
		}),
	}
}

func onStartCallbackInGroup(whc botsfw.WebhookContext, group models.GroupEntry) (m botsfw.MessageFromBot, err error) {
	// This links Telegram ChatID and ChatInstance
	panic("not implemeted")
	// if twhc, ok := whc.(*telegram.tgWebhookContext); ok {
	// 	if err = twhc.CreateOrUpdateTgChatInstance(); err != nil {
	// 		return
	// 	}
	// }
	// return inGroupWelcomeMessage(whc, group)
}

//func inGroupWelcomeMessage(whc botsfw.WebhookContext, group models.GroupEntry) (m botsfw.MessageFromBot, err error) {
//	m, err = GroupSettingsAction(whc, group, false)
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
