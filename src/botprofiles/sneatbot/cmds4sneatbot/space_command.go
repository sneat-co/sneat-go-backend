package cmds4sneatbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"net/url"
	"strings"
	"time"
)

var spaceCommand = botsfw.Command{
	Code:           "space",
	Commands:       []string{"/space"},
	InputTypes:     []botsfw.WebhookInputType{botsfw.WebhookInputCallbackQuery},
	CallbackAction: spaceCallbackAction,
}

func spaceCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	spaceRef := tghelpers.GetSpaceRef(callbackUrl)
	if m, err = spaceAction(whc, spaceRef); err != nil {
		return
	}
	keyboard := m.Keyboard
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard
	if m.EditMessageUID, err = tghelpers.GetEditMessageUID(whc); err != nil {
		return
	}
	whc.ChatData().SetAwaitingReplyTo("")
	return
}

func getSpaceID(whc botsfw.WebhookContext, spaceType core4spaceus.SpaceType) (spaceID string, user dbo4userus.UserEntry, err error) {
	appUserID := whc.AppUserID()
	user = dbo4userus.NewUserEntry(appUserID)
	ctx := whc.Context()
	tx := whc.Tx()
	if err = tx.Get(ctx, user.Record); err != nil {
		return
	}
	spaceID, _ = user.Data.GetFirstSpaceBriefBySpaceType(spaceType)
	return
}

func spaceAction(whc botsfw.WebhookContext, spaceRef core4spaceus.SpaceRef) (m botsfw.MessageFromBot, err error) {
	spaceID, spaceType := spaceRef.SpaceID(), spaceRef.SpaceType()
	if spaceID == "" {
		var user dbo4userus.UserEntry
		if spaceID, user, err = getSpaceID(whc, spaceType); err != nil {
			return
		}
		if spaceID == "" {
			var spaceEntry dbo4spaceus.SpaceEntry
			ctx := whc.Context()
			tx := whc.Tx()
			request := dto4spaceus.CreateSpaceRequest{Type: spaceType}
			params := dal4userus.UserWorkerParams{
				Started: time.Now(), // TODO: get from tx
				User:    user,
			}
			if spaceEntry, _, err = facade4spaceus.CreateSpaceTxWorker(ctx, tx, request, &params); err != nil {
				return
			}
			if err = params.User.Data.Validate(); err != nil {
				err = fmt.Errorf("user record is not valid after CreateSpaceTxWorker: %w", err)
				return
			}
			if err = tx.Update(ctx, params.User.Key, params.UserUpdates); err != nil {
				return
			}
			spaceID = spaceEntry.ID
		}
	}
	var spaceIcon string

	var switchSpaceCallbackData string
	var switchSpaceTitle string
	switch spaceType {
	case core4spaceus.SpaceTypeFamily:
		spaceIcon = "👪"
		switchSpaceTitle = "Private"
		switchSpaceCallbackData = "space?s=" + string(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypePrivate, ""))
	case core4spaceus.SpaceTypePrivate:
		spaceIcon = "🔒"
		switchSpaceTitle = "Family"
		switchSpaceCallbackData = "space?s=" + string(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, ""))
	}

	spaceTitle := strings.ToUpper(string(spaceType)[:1]) + string(spaceType)[1:]
	m.Text += fmt.Sprintf("Current space: %s <b>%s</b>", spaceIcon, spaceTitle)
	m.Format = botsfw.MessageFormatHTML

	spaceCallbackParams := "s=" + string(core4spaceus.NewSpaceRef(spaceType, spaceID))

	firstRow := []tgbotapi.InlineKeyboardButton{
		{
			Text:         "📇 Contacts",
			CallbackData: "contacts?" + spaceCallbackParams,
		},
	}
	if spaceID != "private" {
		firstRow = append(firstRow, tgbotapi.InlineKeyboardButton{
			Text:         "👪 Members",
			CallbackData: "members?" + spaceCallbackParams,
		})
	}

	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		firstRow,
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🚗 Assets",
				CallbackData: "assets?" + spaceCallbackParams,
			},
			{
				Text:         "💰 Budget",
				CallbackData: "budget?" + spaceCallbackParams,
			},
			{
				Text:         "💸 Debts",
				CallbackData: "debts?" + spaceCallbackParams,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🛒 Buy",
				CallbackData: "list?k=buy!groceries&" + spaceCallbackParams,
			},
			{
				Text:         "🏗️ ToDo",
				CallbackData: "list?k=do!tasks&" + spaceCallbackParams,
			},
			{
				Text:         "📽️ Watch",
				CallbackData: "list?k=watch!movies&" + spaceCallbackParams,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🗓️ Calendar",
				CallbackData: "calendar?" + spaceCallbackParams,
			},
			{
				Text:         "⚙️ Settings",
				CallbackData: "settings?" + spaceCallbackParams,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("🔀 Switch to \"%s\" space", switchSpaceTitle),
				CallbackData: switchSpaceCallbackData,
			},
		},
	)
	return
}
