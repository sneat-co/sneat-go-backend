package splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"net/url"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

const LEAVE_GROUP_COMMAND = "leave-group"

var leaveGroupCommand = shared_group.GroupCallbackCommand(LEAVE_GROUP_COMMAND,
	func(whc botsfw.WebhookContext, _ *url.URL, group models.Group) (m botsfw.MessageFromBot, err error) {
		return
	},
)
