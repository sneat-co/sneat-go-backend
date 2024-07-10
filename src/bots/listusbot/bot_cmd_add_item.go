package listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/bots/listusbot/dal4listusbot"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"strings"
)

var addBuyItemCommand = botsfw.Command{
	Code:     "buy",
	Commands: []string{"/buy"},
	Matcher: func(_ botsfw.Command, context botsfw.WebhookContext) bool {
		input := context.Input()
		if input.InputType() == botsfw.WebhookInputText {
			text := strings.ToLower(strings.TrimSpace(input.(botsfw.WebhookTextMessage).Text()))
			return strings.HasPrefix(text, "buy ") || strings.HasPrefix(text, "купить ")
		}
		return false
	},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputInlineQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()

		chatData := whc.ChatData()

		listusChatData := chatData.(*dal4listus.ListusChatData)

		input := whc.Input().(botsfw.WebhookTextMessage)
		text := strings.TrimSpace(input.Text())
		text = text[strings.Index(text, " ")+1:]
		userCtx := facade.NewUser(whc.AppUserID())

		user := dbo4userus.NewUser(userCtx.GetID())
		if err = facade4userus.GetUserByID(ctx, facade.GetDatabase(ctx), user.Record); err != nil {
			return m, err
		}

		spaceID := listusChatData.SpaceID

		if spaceID == "" {
			familySpaceID, familySpaceBrief := user.Data.GetSpaceBriefByType(core4teamus.SpaceTypeFamily)
			if familySpaceBrief == nil {
				m = whc.NewMessage("You are not a member of any family team")
				return m, nil
			}
			spaceID = familySpaceID
		}

		request := facade4listus.CreateListItemsRequest{
			ListRequest: facade4listus.ListRequest{
				ListID: dbo4listus.GetFullListID(dbo4listus.ListTypeToBuy, "groceries"),
				SpaceRequest: dto4teamus.SpaceRequest{
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
	},
}
