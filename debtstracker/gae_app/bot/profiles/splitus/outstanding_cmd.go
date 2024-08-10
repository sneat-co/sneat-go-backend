package splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/logus"
)

const outstandingBalanceCommandCode = "outstanding-balance"

var outstandingBalanceCommand = botsfw.Command{
	Code:     outstandingBalanceCommandCode,
	Commands: []string{"/outstanding"},
	Action:   outstandingBalanceAction,
}

func outstandingBalanceAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	logus.Debugf(c, "outstandingBalanceAction()")
	var user models.AppUser
	if user, err = facade2debtus.User.GetUserByID(c, nil, whc.AppUserID()); err != nil {
		return
	}

	outstandingBalance := user.Data.GetOutstandingBalance()
	m.Text = fmt.Sprintf("Outstanding balance: %v", outstandingBalance)
	return
}
