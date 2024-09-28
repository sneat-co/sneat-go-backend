package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_space"
	"net/url"
)

const LeaveGroupCommandCode = "leave-group"

var leaveGroupCommand = shared_space.SpaceCallbackCommand(LeaveGroupCommandCode,
	func(whc botsfw.WebhookContext, _ *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		return
	},
)
