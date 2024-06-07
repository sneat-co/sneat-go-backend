package shared_group

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"net/url"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func GroupCallbackCommandData(command string, groupID string) string {
	return command + "?group=" + groupID
}

type GroupAction func(whc botsfw.WebhookContext, group models.GroupEntry) (m botsfw.MessageFromBot, err error)
type GroupCallbackAction func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.GroupEntry) (m botsfw.MessageFromBot, err error)

func GroupCallbackCommand(code string, f GroupCallbackAction) botsfw.Command {
	return botsfw.NewCallbackCommand(code,
		func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			var group models.GroupEntry
			if group, err = GetGroup(whc, callbackUrl); err != nil {
				return
			}
			return f(whc, callbackUrl, group)
		},
	)
}

func NewGroupAction(f GroupAction) botsfw.CommandAction {
	return func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		var group models.GroupEntry
		if group, err = GetGroup(whc, nil); err != nil {
			return
		}
		return f(whc, group)
	}
}

func NewGroupCallbackAction(f GroupCallbackAction) botsfw.CallbackAction {
	return func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		var group models.GroupEntry
		if group, err = GetGroup(whc, nil); err != nil {
			return
		}
		return f(whc, callbackUrl, group)
	}
}
