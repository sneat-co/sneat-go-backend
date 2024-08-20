package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_space"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/url"

	"context"
	"errors"
)

func GetBillMembersCallbackData(billID string) string {
	return billCallbackCommandData(billMembersCommandCode, billID)
}

func GetBillID(callbackUrl *url.URL) (billID string, err error) {
	if billID = callbackUrl.Query().Get("bill"); billID == "" {
		err = errors.New("required parameter 'bill' is not passed")
	}
	return
}

func getBill(c context.Context, tx dal.ReadSession, callbackUrl *url.URL) (bill models4splitus.BillEntry, err error) {
	if bill.ID, err = GetBillID(callbackUrl); err != nil {
		return
	}
	if bill, err = facade4splitus.GetBillByID(c, tx, bill.ID); err != nil {
		return
	}
	return
}

type billCallbackActionHandler func(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error)

func billCallbackCommand(code string, f billCallbackActionHandler) (command botsfw.Command) {
	command = botsfw.NewCallbackCommand(code, billCallbackAction(f))
	//if txOptions != nil {
	//	command.CallbackAction = shared_all.TransactionalCallbackAction(txOptions, command.CallbackAction)
	//}
	return
}

func billCallbackAction(f billCallbackActionHandler) func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	return func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
			var bill models4splitus.BillEntry
			if bill, err = getBill(c, tx, callbackUrl); err != nil {
				return
			}
			if bill.Data.GetUserGroupID() == "" {
				if whc.IsInGroup() {
					var splitusSpace models4splitus.SplitusSpaceEntry
					if splitusSpace.ID, err = shared_space.GetUserGroupID(whc); err != nil {
						return
					}
					if bill, splitusSpace, err = facade4splitus.AssignBillToGroup(c, tx, bill, splitusSpace.ID, whc.AppUserID()); err != nil {
						return
					}
				} else {
					logus.Debugf(c, "Not in group")
				}
			}
			m, err = f(whc, tx, callbackUrl, bill)
			return err
		}); err != nil {
			return
		}
		return
	}
}
