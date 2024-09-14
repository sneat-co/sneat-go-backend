package cmds4sneatbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-go-backend/src/botscore/bothelpers"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/i18n"
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

	var switchSpace switchSpaceArgs

	switch spaceType {
	case core4spaceus.SpaceTypeFamily:
		spaceIcon = "👪"
		switchSpace.icon = "🔒"
		switchSpace.title = "Private"
		switchSpace.callbackData = "space?s=" + string(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypePrivate, ""))
	case core4spaceus.SpaceTypePrivate:
		spaceIcon = "🔒"
		switchSpace.icon = "👪"
		switchSpace.title = "Family"
		switchSpace.callbackData = "space?s=" + string(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, ""))
	}

	spaceTitle := strings.ToUpper(string(spaceType)[:1]) + string(spaceType)[1:]
	m.Text += whc.Translate(trans.SPACE_CMD_TEXT, spaceIcon, spaceTitle)
	m.Text += "\n" + strings.Repeat(".", 100)
	m.Format = botsfw.MessageFormatHTML

	if spaceRef.SpaceID() == "" {
		err = fmt.Errorf("spaceRef.SpaceID() is empty string")
		return
	}

	m.Keyboard = spaceInlineKeyboard(whc, spaceRef, switchSpace)
	return
}

type switchSpaceArgs struct {
	callbackData string
	title        string
	icon         string
}

func spaceInlineKeyboard(translator i18n.SingleLocaleTranslator, spaceRef core4spaceus.SpaceRef, switchSpace switchSpaceArgs) *tgbotapi.InlineKeyboardMarkup {
	spaceType, spaceID := spaceRef.SpaceType(), spaceRef.SpaceID()

	spaceCallbackParams := "s=" + string(core4spaceus.NewSpaceRef(spaceType, spaceID))

	var spaceUrlPath = spaceRef.UrlPath()

	botWebAppUrl := bothelpers.GetBotWebAppUrl()

	spacePageUrl := func(page string) string {
		return fmt.Sprintf("%s/space/%s/%s", botWebAppUrl, spaceUrlPath, page)
	}

	firstRow := []tgbotapi.InlineKeyboardButton{
		{
			Text: "📇 " + translator.Translate(trans.SPACE_CMD_BTN_CONTACTS),
			WebApp: &tgbotapi.WebappInfo{
				Url: spacePageUrl("contacts"),
			},
		},
	}
	if spaceType != core4spaceus.SpaceTypePrivate {
		firstRow = append(firstRow, tgbotapi.InlineKeyboardButton{
			Text: "👪 " + translator.Translate(trans.SPACE_CMD_BTN_MEMBER),
			WebApp: &tgbotapi.WebappInfo{
				Url: spacePageUrl("members"),
			},
		})
	}

	listCallbackData := func(id string) string {
		return fmt.Sprintf("list?k=%s&%s", id, spaceCallbackParams)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		firstRow,
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "🚗 " + translator.Translate(trans.SPACE_CMD_BTN_ASSETS),
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("assets"),
				},
			},
			{
				Text: "💰 " + translator.Translate(trans.SPACE_CMD_BTN_BUDGET),
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("budget"),
				},
			},
			{
				Text: "💸 " + translator.Translate(trans.SPACE_CMD_BTN_DEBTS),
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("debts"),
				},
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🛒 " + translator.Translate(trans.LIST_CMD_BUY),
				CallbackData: listCallbackData(dbo4listus.BuyGroceriesListID),
			},
			{
				Text:         "🏗️ " + translator.Translate(trans.LIST_CMD_TODO),
				CallbackData: listCallbackData(dbo4listus.DoTasksListID),
			},
			{
				Text:         "📽️ " + translator.Translate(trans.LIST_CMD_WATCH),
				CallbackData: listCallbackData(dbo4listus.WatchMoviesListID),
			},
			{
				Text:         "📘 " + translator.Translate(trans.LIST_CMD_READ),
				CallbackData: listCallbackData(dbo4listus.ReadBooksListID),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "🗓️ " + translator.Translate(trans.SPACE_CMD_BTN_CALENDAR),
				WebApp: &tgbotapi.WebappInfo{
					Url: spacePageUrl("calendar"),
				},
			},
			{
				Text:         "⚙️ " + translator.Translate(trans.SPACE_CMD_BTN_SETTINGS),
				CallbackData: cmds4anybot.SpaceSettingsCommandCode + "?" + spaceCallbackParams,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s Go %s space", switchSpace.icon, switchSpace.title),
				CallbackData: switchSpace.callbackData,
			},
			{
				Text:         "🌌 " + translator.Translate(trans.BTN_SPACES),
				CallbackData: "spaces?s=" + string(spaceRef),
			},
		},
	)
}
