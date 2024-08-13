package botcmds4splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"net/url"
)

const deleteBillCommandCode = "delete_bill"

var deleteBillCommand = billCallbackCommand(deleteBillCommandCode,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		if _, err = facade4splitus.DeleteBill(c, bill.ID, whc.AppUserID()); err != nil {
			if err == facade4splitus.ErrSettledBillsCanNotBeDeleted {
				m.Text = whc.Translate(err.Error())
				err = nil
			}
			return
		}
		m.Text = fmt.Sprintf("BillEntry #%v has been deleted", bill.ID)
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
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		if _, err = facade4splitus.RestoreBill(c, bill.ID, whc.AppUserID()); err != nil {
			if err == facade4splitus.ErrSettledBillsCanNotBeDeleted {
				m.Text = whc.Translate(err.Error())
				err = nil
			}
			return
		}
		if m.Text, err = getBillCardMessageText(c, whc.GetBotCode(), whc, bill, false, "BillEntry has been restored"); err != nil {
			return
		}
		m.Format = botsfw.MessageFormatHTML
		m.IsEdit = true
		return
	},
)
