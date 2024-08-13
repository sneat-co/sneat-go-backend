package botcmds4splitus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/url"

	"errors"
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

		billEntity := models4splitus.NewBillEntity(
			models4splitus.BillCommon{
				Status:        models4splitus.BillStatusDraft,
				SplitMode:     models4splitus.SplitModeEqually,
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
		userData := appUser.(*dbo4userus.UserDbo)
		userName := userData.Names.GetFullName()
		if userName == "" {
			err = errors.New("userData has no name")
			return
		}

		spaceID := userData.GetFamilySpaceID()

		billMember := briefs4splitus.BillMemberBrief{
			Paid: paidAmount,
		}

		//appUserID := whc.AppUserID()

		if err = billEntity.SetBillMembers([]*briefs4splitus.BillMemberBrief{&billMember}); err != nil {
			return
		}

		return m, facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
			var bill models4splitus.BillEntry
			if bill, err = facade4splitus.CreateBill(c, tx, spaceID, billEntity); err != nil {
				return
			}
			m, err = ShowBillCard(whc, true, bill, "")
			return err
		})
	},
}
