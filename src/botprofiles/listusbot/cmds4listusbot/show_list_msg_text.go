package cmds4listusbot

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	botsgocore "github.com/bots-go-framework/bots-go-core"
	"github.com/sneat-co/sneat-go-backend/src/botscore/bothelpers"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/strongo/random"
	"strings"
)

type ListAction string

const (
	ListActionRefresh        ListAction = "refresh"
	ListActionClear          ListAction = "clear"
	ListActionClearCancel    ListAction = "clear0"
	ListActionClearConfirmed ListAction = "clear1"
	ListActionFull           ListAction = "full"
)

type ListTab string

const (
	ListTabActive ListTab = "active"
	ListTabDone   ListTab = "done"
	ListTabAll    ListTab = "all"
)

func getShowListMessage(
	_ context.Context,
	whc botsfw.WebhookContext,
	spaceRef core4spaceus.SpaceRef,
	listKey dbo4listus.ListKey,
	list dal4listus.ListEntry,
	listAction ListAction,
	listTab ListTab,
) (
	m botsfw.MessageFromBot, err error,
) {
	if listTab == "" {
		activeCount, doneCount := getListCounts(list.Data.Items)
		if activeCount == 0 && doneCount > 0 {
			listTab = ListTabDone
		} else {
			listTab = ListTabActive
		}
	}
	title := list.Data.Title
	listEmoji := list.Data.Emoji
	if title == "" || listEmoji == "" {
		switch listKey.ListType() {
		case dbo4listus.ListTypeToBuy:
			if listEmoji == "" {
				listEmoji = "🛒"
			}
			if listKey.ListSubID() == "groceries" {
				title = "Groceries to buy"
			}
		case dbo4listus.ListTypeToDo:
			if listEmoji == "" {
				listEmoji = "🏗"
			}
			if listKey.ListSubID() == "tasks" {
				title = "Tasks to do"
			}
		case dbo4listus.ListTypeToWatch:
			if listEmoji == "" {
				listEmoji = "🎬"
			}
			if title == "" {
				if listKey.ListSubID() == "movies" {
					title = "Movies to watch"
				}
			}
		default:
			if listEmoji == "" {
				listEmoji = "📋"
			}
			title = fmt.Sprintf("%s: %s", listKey.ListType(), listKey.ListSubID())
		}
	}
	title = fmt.Sprintf("%s <b>%s</b>", listEmoji, title)
	switch listTab {
	case ListTabActive:
		title += " (active)"
	case ListTabDone:
		title += " (done)"
	case ListTabAll:
		title += " (all)"
	}

	m = whc.NewMessage(title)
	m.Format = botsfw.MessageFormatHTML

	if len(list.Data.Items) == 0 {
		m.Text += "\n\n<i>List is empty.</i>"
	} else {
		itemTexts := make([]string, 0, len(list.Data.Items))
		for _, item := range list.Data.Items {
			if listTab == ListTabActive && item.IsDone || listTab == ListTabDone && !item.IsDone {
				continue
			}
			emoji := item.Emoji
			if emoji == "" {
				emoji = "•"
			}
			itemText := fmt.Sprintf("%s %s", emoji, item.Title)
			if item.IsDone && listTab == ListTabAll {
				itemText += " ✔️"
			}
			itemTexts = append(itemTexts, itemText)
		}
		if len(itemTexts) == 0 {
			m.Text += fmt.Sprintf("\n\n<i>No \"%s\" items on the list.</i>", listTab)
		} else {
			m.Text += "\n\n" + strings.Join(itemTexts, "\n")
		}
	}

	switch listAction {
	case ListActionClear:
		if len(list.Data.Items) == 0 {
			m.Keyboard = getShowListStandardKeyboard(spaceRef, listKey, list.Data.Items, listTab)
		} else {
			m.Keyboard = getShowListClearKeyboard(spaceRef, listKey, listTab)
		}
	default:
		m.Keyboard = getShowListStandardKeyboard(spaceRef, listKey, list.Data.Items, listTab)
	}
	return
}

func getShowListClearKeyboard(spaceRef core4spaceus.SpaceRef, listKey dbo4listus.ListKey, listTab ListTab) *tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "❌ Yes, clear list",
				CallbackData: getShowListCallbackData(spaceRef, listKey, ListActionClearConfirmed, listTab),
			},
			{
				Text:         "✖️ Cancel",
				CallbackData: getShowListCallbackData(spaceRef, listKey, ListActionClearCancel, listTab),
			},
		},
	)
}

func getListCounts(listItems []*dbo4listus.ListItemBrief) (activeCount, completedCount int) {
	for _, item := range listItems {
		if item.IsDone {
			completedCount++
		} else {
			activeCount++
		}
	}
	return
}

func getShowListStandardKeyboard(spaceRef core4spaceus.SpaceRef, listKey dbo4listus.ListKey, items []*dbo4listus.ListItemBrief, listTab ListTab) botsgocore.Keyboard {

	activeCount, completedCount := getListCounts(items)

	row0 := []tgbotapi.InlineKeyboardButton{
		{
			Text:         fmt.Sprintf("Active: %d", activeCount),
			CallbackData: getShowListCallbackData(spaceRef, listKey, "", ListTabActive),
		},
		{
			Text:         fmt.Sprintf("Done: %d", completedCount),
			CallbackData: getShowListCallbackData(spaceRef, listKey, "", ListTabDone),
		},
		{
			Text:         fmt.Sprintf("Total: %d", len(items)),
			CallbackData: getShowListCallbackData(spaceRef, listKey, "", ListTabAll),
		},
	}
	const checkMark = "✅ "
	switch listTab {
	case ListTabActive:
		row0[0].Text = checkMark + row0[0].Text
	case ListTabDone:
		row0[1].Text = checkMark + row0[1].Text
	default:
		row0[2].Text = checkMark + row0[2].Text
	}
	row := make([]tgbotapi.InlineKeyboardButton, 1, 2)
	row[0] = tgbotapi.InlineKeyboardButton{
		Text:         "🔃 Refresh",
		CallbackData: getShowListCallbackData(spaceRef, listKey, ListActionRefresh, listTab) + "&!=" + random.ID(5),
	}

	if len(items) > 0 {
		row = append(row, tgbotapi.InlineKeyboardButton{
			Text:         "❌ Clear list",
			CallbackData: getShowListCallbackData(spaceRef, listKey, ListActionClear, listTab),
		})
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if len(items) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row0)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		[]tgbotapi.InlineKeyboardButton{{
			Text: "💻 Edit list",
			WebApp: &tgbotapi.WebappInfo{
				Url: bothelpers.GetBotWebAppUrl() + fmt.Sprintf("space/%s/%s/list/%s/%s", spaceRef.SpaceType(), spaceRef.SpaceID(), listKey.ListType(), listKey.ListSubID()),
			}},
		},
		row,
		[]tgbotapi.InlineKeyboardButton{
			tghelpers.BackToSpaceMenuButton(spaceRef),
		},
	)
	return keyboard
}

func getShowListCallbackData(spaceRef core4spaceus.SpaceRef, listKey dbo4listus.ListKey, action ListAction, tab ListTab) (callbackData string) {
	callbackData = fmt.Sprintf("list?k=%s&s=%s", listKey, spaceRef)
	if action != "" {
		callbackData += "&a=" + string(action)
	}
	if tab != "" {
		callbackData += "&t=" + string(tab)
	}
	return
}
