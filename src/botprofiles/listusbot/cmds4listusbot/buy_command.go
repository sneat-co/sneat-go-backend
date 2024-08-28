package cmds4listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/listusbot/dal4listusbot"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/url"
	"strings"
)

var buyCommand = botsfw.Command{
	Code:     "buy",
	Commands: []string{"/buy"},
	Icon:     "🛒",
	InputTypes: []botsfw.WebhookInputType{
		botsfw.WebhookInputText,
		botsfw.WebhookInputCallbackQuery,
	},
	Matcher: func(_ botsfw.Command, context botsfw.WebhookContext) bool {
		input := context.Input()
		if input.InputType() == botsfw.WebhookInputText {
			text := strings.ToLower(strings.TrimSpace(input.(botsfw.WebhookTextMessage).Text()))
			return strings.HasPrefix(text, "buy ") || strings.HasPrefix(text, "купить ")
		}
		return false
	},
	Action:         buyAction,
	CallbackAction: buyCallbackAction,
}

func buyCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	m.Format = botsfw.MessageFormatHTML
	m.Text = "🛒 <b>Groceries to buy</b>"
	if callbackUrl.Query().Get("action") == "clear" {
		m.Text += "\n\n<i>List is empty.</i>"
	} else {
		m.Text += "\n\n🥛 Milk"
		m.Text += "\n\n🍞 Bread"
	}
	m.Text += "\n\nSent text to add it to the \"To-Buy\" list."
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "Clear list",
				CallbackData: "buy?action=clear",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			tghelpers.BackToSpaceMenuButton(callbackUrl),
		},
	)
	m.ResponseChannel = botsfw.BotAPISendMessageOverHTTPS
	return
}

func buyAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	chatData := whc.ChatData()

	listusChatData := chatData.(*dal4listus.ListusChatData)

	input := whc.Input().(botsfw.WebhookTextMessage)
	text := strings.TrimSpace(input.Text())
	text = text[strings.Index(text, " ")+1:]
	userCtx := facade.NewUserContext(whc.AppUserID())

	user := dbo4userus.NewUserEntry(userCtx.GetUserID())

	if err = dal4userus.GetUser(ctx, whc.DB(), user); err != nil {
		return m, err
	}

	spaceID := listusChatData.SpaceID

	if spaceID == "" {
		familySpaceID, familySpaceBrief := user.Data.GetSpaceBriefByType(core4spaceus.SpaceTypeFamily)
		if familySpaceBrief == nil {
			m = whc.NewMessage("You are not a member of any family team")
			return m, nil
		}
		spaceID = familySpaceID
	}

	request := facade4listus.CreateListItemsRequest{
		ListRequest: facade4listus.ListRequest{
			ListID: dbo4listus.GetFullListID(dbo4listus.ListTypeToBuy, "groceries"),
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: spaceID,
			},
		},
		Items: []facade4listus.CreateListItemRequest{
			{
				ListItemBase: dbo4listus.ListItemBase{
					Title: text,
				},
			},
		},
	}

	if _, err = facade4listus.CreateListItems(ctx, userCtx, request); err != nil {
		return m, err
	}
	responseText := fmt.Sprintf("Added to groceries list: %s", text)
	m = whc.NewMessage(responseText)
	return m, nil
}
