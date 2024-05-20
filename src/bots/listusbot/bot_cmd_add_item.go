package listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/bots/listusbot/dal4listusbot"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/models4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
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

		user := models4userus.NewUser(userCtx.GetID())
		if err = facade4userus.GetUserByID(ctx, facade.GetDatabase(ctx), user.Record); err != nil {
			return m, err
		}

		teamID := listusChatData.TeamID

		if teamID == "" {
			familyTeamID, familyTeamBrief := user.Data.GetTeamBriefByType(core4teamus.TeamTypeFamily)
			if familyTeamBrief == nil {
				m = whc.NewMessage("You are not a member of any family team")
				return m, nil
			}
			teamID = familyTeamID
		}

		request := facade4listus.CreateListItemsRequest{
			ListRequest: facade4listus.ListRequest{
				ListID: models4listus.GetFullListID(models4listus.ListTypeToBuy, "groceries"),
				TeamRequest: dto4teamus.TeamRequest{
					TeamID: teamID,
				},
			},
			Items: []facade4listus.CreateListItemRequest{
				{
					ListItemBase: models4listus.ListItemBase{
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
