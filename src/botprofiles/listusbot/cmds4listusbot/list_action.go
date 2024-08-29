package cmds4listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/random"
	"slices"
	"strings"
)

func listAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	chatData := whc.ChatData()

	sneatAppChatData := chatData.(interface{ GetSpaceRef() core4spaceus.SpaceRef })

	input := whc.Input().(botsfw.WebhookTextMessage)
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
	userCtx := facade.NewUserContext(whc.AppUserID())

	spaceRef := sneatAppChatData.GetSpaceRef()

	if spaceRef == "" {
		userID := userCtx.GetUserID()
		var user dbo4userus.UserEntry
		var db dal.DB
		if db, err = facade.GetDatabase(ctx); err != nil {
			return
		}
		if user, err = dal4userus.GetUserByID(ctx, db, userID); err != nil {
			return
		}
		var spaceID string
		spaceID, _ = user.Data.GetSpaceBriefByType(core4spaceus.SpaceTypeFamily)
		if spaceID == "" {
			m = whc.NewMessage("You are not a member of any family team")
			return m, nil
		}
		spaceRef = core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, spaceID)
	}

	request := dto4listus.CreateListItemsRequest{
		ListRequest: dto4listus.ListRequest{
			ListID: dbo4listus.NewListKey(dbo4listus.ListTypeToBuy, "groceries"),
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: spaceRef.SpaceID(),
			},
		},
		Items: []dto4listus.CreateListItemRequest{
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
