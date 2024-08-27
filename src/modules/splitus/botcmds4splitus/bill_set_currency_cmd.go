package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/dal4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
	"net/url"
)

const setBillCurrencyCommandCode = "set-bill-currency"

var setBillCurrencyCommand = billCallbackCommand(setBillCurrencyCommandCode,
	func(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		logus.Debugf(ctx, "setBillCurrencyCommand.CallbackAction()")
		query := callbackUrl.Query()
		currencyCode := money.CurrencyCode(query.Get("currency"))
		if bill.Data.Currency != currencyCode {
			previousCurrency := bill.Data.Currency
			bill.Data.Currency = currencyCode
			if err = facade4splitus.SaveBill(ctx, tx, bill); err != nil {
				return
			}

			if bill.Data.SpaceID != "" {
				splitusSpace := models4splitus.NewSplitusSpaceEntry(bill.Data.SpaceID)
				if err = dal4splitus.GetSplitusSpace(ctx, tx, splitusSpace); err != nil {
					return
				}
				diff := bill.Data.GetBalance().BillBalanceDifference(make(briefs4splitus.BillBalanceByMember, 0))
				if _, err = splitusSpace.Data.ApplyBillBalanceDifference(bill.Data.Currency, diff); err != nil {
					return
				}
				if previousCurrency != "" {
					if _, err = splitusSpace.Data.ApplyBillBalanceDifference(previousCurrency, diff.Reverse()); err != nil {
						return
					}
				}
				if err = dal4splitus.SaveSplitusSpace(ctx, tx, splitusSpace); err != nil {
					return
				}
			}
		}
		if m.Text, err = getBillCardMessageText(ctx, whc.GetBotCode(), whc, bill, true, whc.Translate(trans.MESSAGE_TEXT_BILL_ASK_WHO_PAID)); err != nil {
			return
		}
		m.Format = botsfw.MessageFormatHTML
		m.Keyboard = getWhoPaidInlineKeyboard(whc, bill.ID)
		m.IsEdit = true

		return
	},
)
