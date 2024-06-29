package splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/logus"
	"net/url"
)

const settleBillsCommandCode = "settle"

var settleBillsCommand = botsfw.Command{
	Code:     settleBillsCommandCode,
	Commands: []string{"/settle"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return settleBillsAction(whc)
	},
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return settleBillsAction(whc)
	},
}

func settleBillsAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	logus.Debugf(c, "settleBillsAction()")
	var user models.AppUser
	if user, err = facade.User.GetUserByID(c, nil, whc.AppUserID()); err != nil {
		return
	}

	outstandingBills := user.Data.GetOutstandingBills()

	m.Text = fmt.Sprintf("len(outstandingBills): %v", len(outstandingBills))

	return
}
