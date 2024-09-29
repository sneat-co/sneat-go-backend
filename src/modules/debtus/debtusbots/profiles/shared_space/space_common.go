package shared_space

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"net/url"
)

func SpaceCallbackCommandData(command string, spaceID string) string {
	return command + "?space=" + spaceID
}

type SpaceAction func(whc botsfw.WebhookContext, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error)
type SplitusSpaceAction func(whc botsfw.WebhookContext, splitusSpace models4splitus.SplitusSpaceEntry) (m botsfw.MessageFromBot, err error)

type SpaceCallbackAction func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error)
type SplitusSpaceCallbackAction func(whc botsfw.WebhookContext, callbackUrl *url.URL, splitusSpace models4splitus.SplitusSpaceEntry) (m botsfw.MessageFromBot, err error)

func SpaceCallbackCommand(code string, f SpaceCallbackAction) botsfw.Command {
	return botsfw.NewCallbackCommand(code,
		func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			var space dbo4spaceus.SpaceEntry
			if space, err = GetSpaceEntryByCallbackUrl(whc, callbackUrl); err != nil {
				return
			}
			return f(whc, callbackUrl, space)
		},
	)
}

func NewSpaceAction(f SpaceAction) botsfw.CommandAction {
	return func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		var space dbo4spaceus.SpaceEntry
		if space, err = GetSpaceEntryByCallbackUrl(whc, nil); err != nil {
			return
		}
		return f(whc, space)
	}
}

func NewSplitusSpaceAction(f SplitusSpaceAction) botsfw.CommandAction {
	return func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		var splitusSpace models4splitus.SplitusSpaceEntry
		if splitusSpace, err = GetSplitusSpaceEntryByCallbackUrl(whc, nil); err != nil {
			return
		}
		return f(whc, splitusSpace)
	}
}

func NewSpaceCallbackAction(f SpaceCallbackAction) botsfw.CallbackAction {
	return func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		var space dbo4spaceus.SpaceEntry
		if space, err = GetSpaceEntryByCallbackUrl(whc, nil); err != nil {
			return
		}
		return f(whc, callbackUrl, space)
	}
}

func NewSplitusSpaceCallbackAction(f SplitusSpaceCallbackAction) botsfw.CallbackAction {
	return func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		var splitusSpace models4splitus.SplitusSpaceEntry
		if splitusSpace, err = GetSplitusSpaceEntryByCallbackUrl(whc, nil); err != nil {
			return
		}
		return f(whc, callbackUrl, splitusSpace)
	}
}
