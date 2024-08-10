package splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/url"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
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

func getBill(c context.Context, tx dal.ReadSession, callbackUrl *url.URL) (bill models.Bill, err error) {
	if bill.ID, err = GetBillID(callbackUrl); err != nil {
		return
	}
	if bill, err = facade2debtus.GetBillByID(c, tx, bill.ID); err != nil {
		return
	}
	return
}

type billCallbackActionHandler func(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error)

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
			var bill models.Bill
			if bill, err = getBill(c, tx, callbackUrl); err != nil {
				return
			}
			if bill.Data.GetUserGroupID() == "" {
				if whc.IsInGroup() {
					var group models.GroupEntry
					if group.ID, err = shared_group.GetUserGroupID(whc); err != nil {
						return
					}
					if bill, group, err = facade2debtus.Bill.AssignBillToGroup(c, tx, bill, group.ID, whc.AppUserID()); err != nil {
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
