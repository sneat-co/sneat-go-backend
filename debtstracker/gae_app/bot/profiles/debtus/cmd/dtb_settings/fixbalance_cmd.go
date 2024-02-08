package dtb_settings

import (
	"context"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
)

var FixBalanceCommand = botsfw.Command{
	Code:     "fixbalance",
	Commands: []string{"/fixbalance"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		var db dal.DB
		if db, err = facade.GetDatabase(whc.Context()); err != nil {
			return
		}
		if err = db.RunReadwriteTransaction(whc.Context(), func(c context.Context, tx dal.ReadwriteTransaction) error {
			//goland:noinspection GoDeprecation
			user, err := facade.User.GetUserByID(c, tx, whc.AppUserID())
			if err != nil {
				return err
			}
			contacts := user.Data.Contacts()
			balance := make(money.Balance, user.Data.BalanceCount)
			for _, contact := range contacts {
				b := contact.Balance()
				for k, v := range b {
					balance[k] += v
				}
			}
			if err = user.Data.SetBalance(balance); err != nil {
				return err
			}
			return facade.User.SaveUser(c, tx, user)
		}); err != nil {
			return
		}
		m = whc.NewMessage("Balance fixed")
		return
	},
}
