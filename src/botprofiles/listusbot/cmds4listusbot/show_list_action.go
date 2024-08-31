package cmds4listusbot

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func showListAction(
	ctx context.Context, whc botsfw.WebhookContext, tx dal.ReadSession, spaceRef core4spaceus.SpaceRef, listKey dbo4listus.ListKey,
) (
	m botsfw.MessageFromBot, err error,
) {
	spaceID := spaceRef.SpaceID()
	var list dal4listus.ListEntry
	if spaceID == "" {
		list.Data = new(dbo4listus.ListDbo)
	} else {
		list = dal4listus.NewSpaceListEntry(spaceID, listKey)
		if tx == nil {
			if tx, err = facade.GetDatabase(ctx); err != nil {
				return
			}
		}
		if err = tx.Get(ctx, list.Record); dal.IsNotFound(err) {
			err = nil
		} else if err != nil {
			return
		}
	}
	title := list.Data.Title
	emoji := list.Data.Emoji
	if title == "" || emoji == "" {
		switch listKey.ListType() {
		case dbo4listus.ListTypeToBuy:
			if emoji == "" {
				emoji = "🛒"
			}
			if listKey.ListSubID() == "groceries" {
				title = "Groceries to buy"
			}
		case dbo4listus.ListTypeToDo:
			if emoji == "" {
				emoji = "🏗"
			}
			if listKey.ListSubID() == "tasks" {
				title = "Tasks to do"
			}
		case dbo4listus.ListTypeToWatch:
			if emoji == "" {
				emoji = "🎬"
			}
			if title == "" {
				if listKey.ListSubID() == "movies" {
					title = "Movies to watch"
				}
			}
		default:
			if emoji == "" {
				emoji = "📋"
			}
			title = fmt.Sprintf("%s: %s", listKey.ListType(), listKey.ListSubID())
		}
	}
	title = fmt.Sprintf("%s <b>%s</b>", emoji, title)

	m = whc.NewMessage(title)
	m.Format = botsfw.MessageFormatHTML

	if len(list.Data.Items) == 0 {
		m.Text += "\n\n<i>List is empty.</i>"
	} else {
		for _, item := range list.Data.Items {
			emoji := item.Emoji
			if emoji == "" {
				emoji = "•"
			}
			m.Text += fmt.Sprintf("\n\n%s %s", emoji, item.Title)
		}
	}
	var firstRow []tgbotapi.InlineKeyboardButton
	if len(list.Data.Items) > 0 {
		firstRow = append(firstRow, tgbotapi.InlineKeyboardButton{
			Text:         "❌ Clear list",
			CallbackData: fmt.Sprintf("list?k=%s&s=%s&action=clear", listKey, spaceRef),
		})
	}
	firstRow = append(firstRow, tgbotapi.InlineKeyboardButton{
		Text: "💻 Edit list",
		WebApp: &tgbotapi.WebappInfo{
			Url: fmt.Sprintf("https://local-app.sneat.ws/space/family/h4qax/list/%s/%s", listKey.ListType(), listKey.ListSubID()),
		},
	})
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		firstRow,
		//[]tgbotapi.InlineKeyboardButton{
		//	tghelpers.BackToSpaceMenuButton(spaceType, spaceID),
		//},
	)
	m.ResponseChannel = botsfw.BotAPISendMessageOverHTTPS
	return
}
