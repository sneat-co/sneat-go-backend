package cmds4listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/random"
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

	spaceType, spaceID := tghelpers.GetSpaceParams(callbackUrl)

	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "❌ Clear list",
				CallbackData: fmt.Sprintf("buy?action=clear&spaceType=%s&spaceID=%s", spaceType, spaceID),
			},
			{
				Text: "💻 Edit list",
				WebApp: &tgbotapi.WebappInfo{
					Url: "https://local-app.sneat.ws/space/family/h4qax/budget", // TODO: generate URL
				},
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			tghelpers.BackToSpaceMenuButton(callbackUrl),
		},
	)
	m.ResponseChannel = botsfw.BotAPISendMessageOverHTTPS
	chatData := whc.ChatData()
	chatData.SetAwaitingReplyTo("buy")
	switch chatData := chatData.(type) {
	case interface{ SetSpaceID(string) }:
		chatData.SetSpaceID(spaceID)
	default:
		err = fmt.Errorf("chatData %T does not support SetSpaceID()", chatData)
	}
	return
}

func buyAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	chatData := whc.ChatData()

	sneatAppChatData := chatData.(interface{ GetSpaceID() string })

	input := whc.Input().(botsfw.WebhookTextMessage)
	text := strings.TrimSpace(input.Text())
	firstSpaceIndex := strings.Index(text, " ")
	if firstSpaceIndex > 0 {
		firstWord := text[:]
		if firstWord == "buy" || firstWord == "/buy" {
			text = text[len(firstWord):]
		}
	}
	userCtx := facade.NewUserContext(whc.AppUserID())

	spaceID := sneatAppChatData.GetSpaceID()

	if spaceID == "" {
		userID := userCtx.GetUserID()
		var user dbo4userus.UserEntry
		var db dal.DB
		if db, err = facade.GetDatabase(ctx); err != nil {
			return
		}
		if user, err = dal4userus.GetUserByID(ctx, db, userID); err != nil {
			return
		}
		spaceID, _ = user.Data.GetSpaceBriefByType(core4spaceus.SpaceTypeFamily)
		if spaceID == "" {
			m = whc.NewMessage("You are not a member of any family team")
			return m, nil
		}
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
				ID: random.ID(5), // TODO: should be generated inside transaction?
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
