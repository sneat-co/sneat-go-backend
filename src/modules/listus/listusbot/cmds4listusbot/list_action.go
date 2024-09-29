package cmds4listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/random"
	"net/url"
	"slices"
	"strings"
)

func listAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	chatData := whc.ChatData()

	var listKey dbo4listus.ListKey

	awaitingReplyTo := chatData.GetAwaitingReplyTo()
	if awaitingReplyTo != "" {
		var awaitingReplyToUrl *url.URL
		if awaitingReplyToUrl, err = url.Parse(awaitingReplyTo); err != nil {
			err = fmt.Errorf("failed to parse awaitingReplyTo as URL: %w", err)
			return
		}
		listKey = dbo4listus.ListKey(awaitingReplyToUrl.Query().Get("k"))
	}
	if listKey == "" {
		listKey = dbo4listus.BuyGroceriesListID
	}

	sneatAppChatData := chatData.(interface{ GetSpaceRef() core4spaceus.SpaceRef })

	input := whc.Input().(botinput.WebhookTextMessage)
	text := strings.TrimSpace(input.Text())
	if slices.Contains(listCommandPrefixes, text) {
		text = ""
	}
	firstSpaceIndex := strings.Index(text, " ")
	if firstSpaceIndex > 0 {
		firstWord := text[:firstSpaceIndex]
		if slices.Contains(listCommandPrefixes, firstWord) {
			text = strings.TrimSpace(text[len(firstWord):])
		}
	}
	userID := whc.AppUserID()
	userCtx := facade.NewUserContext(userID)

	spaceRef := sneatAppChatData.GetSpaceRef()

	spaceID, spaceType := spaceRef.SpaceID(), spaceRef.SpaceType()

	if spaceID == "" {
		var user dbo4userus.UserEntry
		var db dal.DB
		if db, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
		if user, err = dal4userus.GetUserByID(ctx, db, userID); err != nil {
			return
		}
		spaceID, _ = user.Data.GetFirstSpaceBriefBySpaceType(spaceRef.SpaceType())
		if spaceID == "" {
			if spaceType == core4spaceus.SpaceTypeFamily {
				var result facade4spaceus.CreateSpaceParams
				if result, err = facade4spaceus.CreateSpace(ctx, userCtx, dto4spaceus.CreateSpaceRequest{Type: spaceType}); err != nil {
					err = fmt.Errorf("failed to create missing family space: %w", err)
					return
				}
				spaceID = result.Space.ID
				spaceRef = core4spaceus.NewSpaceRef(spaceType, spaceID)
			} else {
				m = whc.NewMessage(fmt.Sprintf("You are not a member of any %s space", spaceType))
				return m, nil
			}
		}
	}

	request := dto4listus.CreateListItemsRequest{
		ListRequest: dto4listus.ListRequest{
			ListID: listKey,
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: spaceRef.SpaceID(),
			},
		},
	}
	for _, itemText := range strings.Split(text, "\n") {
		item := dto4listus.CreateListItemRequest{
			ID: random.ID(5), // TODO: should be generated inside transaction
			ListItemBase: dbo4listus.ListItemBase{
				Title: cleanListItemTitle(itemText),
			},
		}
		if item.Title != "" {
			request.Items = append(request.Items, item)
		}
	}

	var response dto4listus.CreateListItemResponse
	var list dal4listus.ListEntry
	if response, list, err = facade4listus.CreateListItems(ctx, userCtx, request); err != nil {
		return m, fmt.Errorf("failed to create list items: %w", err)
	}

	itemTexts := make([]string, 0, len(response.CreatedItems))
	for _, item := range response.CreatedItems {
		var itemText string
		if item.Emoji == "" {
			itemText = "\t• " + item.Title
		} else {
			itemText = "\t" + item.Emoji + " " + item.Title
		}
		itemTexts = append(itemTexts, itemText)
	}
	m = whc.NewMessage("Added to groceries list:\n" + strings.Join(itemTexts, "\n"))
	m.Format = botsfw.MessageFormatText
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         list.Data.Emoji + " Show full list",
				CallbackData: getShowListCallbackData(spaceRef, request.ListID, ListActionFull, ListTabActive),
			},
		},
	)
	return m, nil
}

func cleanListItemTitle(title string) string {
	return strings.Trim(title, "•- \t")
}
