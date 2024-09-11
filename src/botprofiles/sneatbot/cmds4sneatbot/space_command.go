package cmds4sneatbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/bothelpers"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/url"
	"strings"
)

var spaceCommand = botsfw.Command{
	Code:           "space",
	Commands:       []string{"/space"},
	InputTypes:     []botinput.WebhookInputType{botinput.WebhookInputCallbackQuery},
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
	if appUserID == "" {
		return
	}
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
	appUserID := whc.AppUserID()
	if spaceID == "" && appUserID != "" {
		if spaceID, _, err = getSpaceID(whc, spaceType); err != nil {
			return
		}
		if spaceID == "" {
			ctx := whc.Context()

			userCtx := facade.NewUserContext(appUserID)

			var createSpaceParams facade4spaceus.CreateSpaceParams
			if createSpaceParams, err = facade4spaceus.CreateSpace(ctx, userCtx, dto4spaceus.CreateSpaceRequest{Type: spaceType}); err != nil {
				return
			}
			spaceID = createSpaceParams.Space.ID
		}
		spaceRef = core4spaceus.NewSpaceRef(spaceType, spaceID)
	}
	var spaceIcon string

	var switchSpaceCallbackData string
	var switchSpaceTitle string
	var switchSpaceIcon string

	switch spaceType {
	case core4spaceus.SpaceTypeFamily:
		spaceIcon = "👪"
		switchSpaceIcon = "🔒"
		switchSpaceTitle = "Private"
		switchSpaceCallbackData = "space?s=" + string(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypePrivate, ""))
	case core4spaceus.SpaceTypePrivate:
		spaceIcon = "🔒"
		switchSpaceIcon = "👪"
		switchSpaceTitle = "Family"
		switchSpaceCallbackData = "space?s=" + string(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, ""))
	}

	spaceTitle := strings.ToUpper(string(spaceType)[:1]) + string(spaceType)[1:]
	m.Text += fmt.Sprintf("Current space: %s <b>%s</b>", spaceIcon, spaceTitle)
	m.Format = botsfw.MessageFormatHTML

	spaceCallbackParams := "s=" + string(core4spaceus.NewSpaceRef(spaceType, spaceID))

	if spaceRef.SpaceID() == "" {
		err = fmt.Errorf("spaceRef.SpaceID() is empty string")
		return
	}

	var spaceUrlPath = spaceRef.UrlPath()

	botWebAppUrl := bothelpers.GetBotWebAppUrl()

	spacePageUrl := func(page string) string {
		return fmt.Sprintf("%s/space/%s/%s", botWebAppUrl, spaceUrlPath, page)
	}

	firstRow := []tgbotapi.InlineKeyboardButton{
		{
			Text: "📇 Contacts",
			WebApp: &tgbotapi.WebappInfo{
				Url: spacePageUrl("contacts"),
			},
		},
	}
	if spaceType != core4spaceus.SpaceTypePrivate {
		firstRow = append(firstRow, tgbotapi.InlineKeyboardButton{
			Text: "👪 Members",
			WebApp: &tgbotapi.WebappInfo{
				Url: spacePageUrl("members"),
			},
		})
	}

	listCallbackData := func(id string) string {
		return fmt.Sprintf("list?k=%s&%s", id, spaceCallbackParams)
	}

	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		firstRow,
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "🚗 Assets",
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("assets"),
				},
			},
			{
				Text: "💰 Budget",
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("budget"),
				},
			},
			{
				Text: "💸 Debts",
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("debts"),
				},
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🛒 Buy",
				CallbackData: listCallbackData(dbo4listus.BuyGroceriesListID),
			},
			{
				Text:         "🏗️ ToDo",
				CallbackData: listCallbackData(dbo4listus.DoTasksListID),
			},
			{
				Text:         "📽️ Watch",
				CallbackData: listCallbackData(dbo4listus.WatchMoviesListID),
			},
			{
				Text:         "📘 Read",
				CallbackData: listCallbackData(dbo4listus.ReadBooksListID),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "🗓️ Calendar",
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("calendar"),
				},
			},
			{
				Text:         "⚙️ Settings",
				CallbackData: "settings?" + spaceCallbackParams,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("🔀 Switch to \"%s\" space %s", switchSpaceTitle, switchSpaceIcon),
				CallbackData: switchSpaceCallbackData,
			},
		},
	)
	return
}
