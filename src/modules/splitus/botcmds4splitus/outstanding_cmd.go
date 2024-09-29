package botcmds4splitus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/const4debtus"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/models4debtus"
	"github.com/strongo/logus"
)

const outstandingBalanceCommandCode = "outstanding-balance"

var outstandingBalanceCommand = botsfw.Command{
	Code:     outstandingBalanceCommandCode,
	Commands: []string{"/outstanding"},
	Action:   outstandingBalanceAction,
}

func outstandingBalanceAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	logus.Debugf(ctx, "outstandingBalanceAction()")
	userID := whc.AppUserID()
	err = dal4userus.RunUserModuleWorker[models4debtus.DebtusUserDbo](
		ctx, userID, const4debtus.ModuleID, new(models4debtus.DebtusUserDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, param *dal4userus.UserModuleWorkerParams[models4debtus.DebtusUserDbo]) (err error) {
			outstandingBalance := param.UserModule.Data.GetOutstandingBalance()
			m.Text = fmt.Sprintf("Outstanding balance: %v", outstandingBalance)
			return err
		})
	if err != nil {
		return m, err
	}
	return
}
