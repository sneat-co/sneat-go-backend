package cmds4listusbot

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/url"
)

func listCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	spaceRef := tghelpers.GetSpaceRef(callbackUrl)
	callbackQuery := callbackUrl.Query()
	listKey := dbo4listus.ListKey(callbackQuery.Get("k"))
	if err = listKey.Validate(); err != nil {
		return
	}
	action := ListAction(callbackQuery.Get("a"))
	tab := ListTab(callbackQuery.Get("t"))

	switch action {
	case ListActionClearConfirmed:
		userCtx := facade.NewUserContext(whc.AppUserID())
		request := dto4listus.ListRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: spaceRef.SpaceID(),
			},
			ListID: listKey,
		}
		var list dal4listus.ListEntry
		if err = dal4listus.RunListWorker(ctx, userCtx, request, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4listus.ListWorkerParams) (err error) {
			params.List.Data.Items = nil
			params.List.Data.Count = 0
			params.ListUpdates = append(params.ListUpdates, dal.Update{
				Field: "items",
				Value: dal.DeleteField,
			})
			params.List.Record.MarkAsChanged()
			list = params.List
			return
		}); err != nil {
			return
		}
		if m, err = getShowListMessage(ctx, whc, spaceRef, listKey, list, action, tab); err != nil {
			return
		}
	case "", ListActionRefresh, ListActionFull, ListActionClear, ListActionClearCancel:
		list := dal4listus.NewListEntry(spaceRef.SpaceID(), listKey)
		var db dal.DB
		if db, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
		if err = db.Get(ctx, list.Record); err != nil && !dal.IsNotFound(err) {
			return
		}
		if m, err = getShowListMessage(ctx, whc, spaceRef, listKey, list, action, tab); err != nil {
			return
		}
	}

	keyboard := m.Keyboard
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard

	m.ResponseChannel = botsfw.BotAPISendMessageOverResponse
	chatData := whc.ChatData()
	awaitingReplyTo := getShowListCallbackData(spaceRef, listKey, "", "")
	chatData.SetAwaitingReplyTo(awaitingReplyTo)
	switch chatData := chatData.(type) {
	case interface {
		SetSpaceRef(core4spaceus.SpaceRef)
	}:
		chatData.SetSpaceRef(spaceRef)
	default:
		err = fmt.Errorf("chatData %T does not support SetSpaceRef(core4spaceus.SpaceRef)", chatData)
	}
	return
}
