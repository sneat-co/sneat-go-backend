package splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"net/url"
)

const deleteBillCommandCode = "delete_bill"

var deleteBillCommand = billCallbackCommand(deleteBillCommandCode,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		if _, err = facade2debtus.Bill.DeleteBill(c, bill.ID, whc.AppUserID()); err != nil {
			if err == facade2debtus.ErrSettledBillsCanNotBeDeleted {
				m.Text = whc.Translate(err.Error())
				err = nil
			}
			return
		}
		m.Text = fmt.Sprintf("Bill #%v has been deleted", bill.ID)
		m.IsEdit = true
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         "Restore",
					CallbackData: billCallbackCommandData(restoreBillCommandCode, bill.ID),
				},
			},
		)
		return
	},
)

const restoreBillCommandCode = "restore_bill"

var restoreBillCommand = billCallbackCommand(restoreBillCommandCode,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		if _, err = facade2debtus.Bill.RestoreBill(c, bill.ID, whc.AppUserID()); err != nil {
			if err == facade2debtus.ErrSettledBillsCanNotBeDeleted {
				m.Text = whc.Translate(err.Error())
				err = nil
			}
			return
		}
		if m.Text, err = getBillCardMessageText(c, whc.GetBotCode(), whc, bill, false, "Bill has been restored"); err != nil {
			return
		}
		m.Format = botsfw.MessageFormatHTML
		m.IsEdit = true
		return
	},
)
