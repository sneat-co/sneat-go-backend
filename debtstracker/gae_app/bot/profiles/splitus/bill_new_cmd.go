package splitus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/url"

	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
)

const (
	NEW_BILL_PARAM_I      = "i"
	NEW_BILL_PARAM_V      = "v"
	NEW_BILL_PARAM_C      = "c"
	NEW_BILL_PARAM_I_OWE  = "owe"
	NEW_BILL_PARAM_I_PAID = "paid"
)

const newBillCommandCode = "new-bill"

var newBillCommand = botsfw.Command{
	Code: newBillCommandCode,
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "newBillCommand.CallbackAction(callbackUrl=%v)", callbackUrl)
		query := callbackUrl.Query()
		paramI := query.Get(NEW_BILL_PARAM_I)
		if paramI != NEW_BILL_PARAM_I_OWE && paramI != NEW_BILL_PARAM_I_PAID {
			err = errors.New("paramI != NEW_BILL_PARAM_I_OWE && paramI != NEW_BILL_PARAM_I_PAID")
			return
		}
		var amountValue, paidAmount decimal.Decimal64p2
		if amountValue, err = decimal.ParseDecimal64p2(query.Get(NEW_BILL_PARAM_V)); err != nil {
			return
		}
		if paramI == NEW_BILL_PARAM_I_PAID {
			paidAmount = amountValue
		}

		strUserID := whc.AppUserID()

		billEntity := models.NewBillEntity(
			models.BillCommon{
				Status:        models.BillStatusDraft,
				SplitMode:     models.SplitModeEqually,
				CreatorUserID: strUserID,
				AmountTotal:   amountValue,
				Currency:      money.CurrencyCode(query.Get("c")),
				UserIDs:       []string{strUserID},
			},
		)
		//tgMessage := whc.Input().(telegram.TelegramWebhookInput).
		//callbackQuery :=
		tgChatMessageID := fmt.Sprintf("%v@%v@%v", whc.Input().(botsfw.WebhookCallbackQuery).GetInlineMessageID(), whc.GetBotCode(), whc.Locale().Code5)
		billEntity.TgChatMessageIDs = []string{tgChatMessageID}

		var appUser botsfwmodels.AppUserData
		if appUser, err = whc.AppUserData(); err != nil {
			return
		}
		user := appUser.(interface{ FullName() string })
		userName := user.FullName()
		if userName == "" {
			err = errors.New("user has no name")
			return
		}

		billMember := models.BillMemberJson{
			Paid: paidAmount,
		}

		//appUserID := whc.AppUserID()

		if err = billEntity.SetBillMembers([]models.BillMemberJson{billMember}); err != nil {
			return
		}

		return m, facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
			var bill models.Bill
			if bill, err = facade2debtus.Bill.CreateBill(c, tx, billEntity); err != nil {
				return
			}
			m, err = ShowBillCard(whc, true, bill, "")
			return err
		})
	},
}
