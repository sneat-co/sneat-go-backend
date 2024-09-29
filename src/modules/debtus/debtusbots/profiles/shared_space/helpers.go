package shared_space

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/anybot"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/url"
	"strconv"
)

func GetSpaceIdFromCallbackUrl(callbackUrl *url.URL) string {
	if callbackUrl == nil {
		return ""
	}
	return callbackUrl.Query().Get("space")
}

func GetSplitusSpaceEntryByCallbackUrl(whc botsfw.WebhookContext, callbackUrl *url.URL) (splitusSpace models4splitus.SplitusSpaceEntry, err error) {
	err = errors.New("func GetSplitusSpaceEntryByCallbackUrl() is not implemented yet")
	return
}

func GetSpaceEntryByCallbackUrl(whc botsfw.WebhookContext, callbackUrl *url.URL) (space dbo4spaceus.SpaceEntry, err error) {
	space.ID = GetSpaceIdFromCallbackUrl(callbackUrl)
	if space.ID == "" {
		if space.ID, err = GetUserGroupID(whc); err != nil {
			return
		}
	}

	var isInGroup bool
	if isInGroup, err = whc.IsInGroup(); err != nil {
		return
	} else if isInGroup {
		if callbackUrl != nil {
			err = errors.New("an attempt to get space ContactID outside of space chat without callback parameter 'space'")
		}
		return
	}

	if space.ID != "" {
		ctx := whc.Context()
		var db dal.DB
		if db, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
		err = db.Get(ctx, space.Record)
		return
	}

	// TODO: document who we can get space callback without space ContactID

	tgChat := whc.Input().(telegram.TgWebhookInput).TgUpdate().Chat()
	var tgChatEntity *anybot.SneatAppTgChatDbo
	if tgChatEntity, err = getTgChatEntity(whc); err != nil {
		return
	}
	ctx := whc.Context()
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		space, err = createSpaceForTelegramGroup(ctx, whc, tx, tgChatEntity, tgChat)
		return err
	})
	return
}

func GetUserGroupID(whc botsfw.WebhookContext) (groupID string, err error) {
	var tgChatEntity *anybot.SneatAppTgChatDbo
	if tgChatEntity, err = getTgChatEntity(whc); err != nil || tgChatEntity == nil {
		return
	}
	if groupID = tgChatEntity.UserGroupID; groupID != "" {
		return
	}
	return
}

func createSpaceForTelegramGroup(ctx context.Context, whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, chatData *anybot.SneatAppTgChatDbo, tgChat *tgbotapi.Chat) (space dbo4spaceus.SpaceEntry, err error) {
	logus.Debugf(ctx, "createSpaceForTelegramGroup()")
	err = errors.New("creation of space from not implemented yet")
	//var user *models4debtus.DebutsAppUserDataOBSOLETE
	//if user, err = shared_all.GetUser(whc); err != nil {
	//	return
	//}
	//var chatInviteLink string
	//
	//if tgChat.IsSuperGroup() { // See: https://core.telegram.org/bots/api#exportchatinvitelink
	//	// TODO: Do this in delayed task - Lets try to get chat  invite link
	//	msg := botsfw.MessageFromBot{BotMessage: telegram.ExportChatInviteLink{}}
	//	if tgResponse, err := whc.Responder().SendMessage(ctx, msg, botsfw.BotAPISendMessageOverHTTPS); err != nil {
	//		logus.Debugf(ctx, "Not able to export chat invite link: %v", err)
	//	} else {
	//		chatInviteLink = string(tgResponse.TelegramMessage.(tgbotapi.APIResponse).Result)
	//		logus.Debugf(ctx, "exportInviteLink response: %v", chatInviteLink)
	//	}
	//}
	//
	//userID := whc.AppUserID()
	//groupEntity := models4debtus.GroupDbo{
	//	CreatorUserID: userID,
	//	Name:          tgChat.Title,
	//}
	//groupEntity.SetTelegramGroups([]models4debtus.GroupTgChatJson{
	//	{
	//		ChatID:         tgChat.ContactID,
	//		Title:          tgChat.Title,
	//		ChatInviteLink: chatInviteLink,
	//	},
	//})
	//
	//hasTgGroupEntity := false
	//beforeGroupInsert := func(ctx context.Context, groupEntity *models4debtus.GroupDbo) (group models4debtus.GroupEntry, err error) {
	//	logus.Debugf(ctx, "beforeGroupInsert()")
	//	var tgGroup models4auth.TgGroup
	//	if tgGroup, err = dtdal.TgGroup.GetTgGroupByID(ctx, nil, tgChat.ContactID); err != nil {
	//		if dal.IsNotFound(err) {
	//			err = nil
	//		} else {
	//			return
	//		}
	//	}
	//	if tgGroup.Data != nil && tgGroup.Data.SpaceRef != "" {
	//		hasTgGroupEntity = true
	//		return dtdal.Group.GetGroupByID(ctx, tx, tgGroup.Data.SpaceRef)
	//	}
	//	_, _, idx, member, members := groupEntity.AddOrGetMember(userID, "", user.FullName())
	//	member.TgUserID = strconv.FormatInt(int64(whc.Input().GetSender().GetID().(int)), 10)
	//	members[idx] = member
	//	groupEntity.SetGroupMembers(members)
	//	return
	//}
	//
	//afterGroupInsert := func(ctx context.Context, group models4debtus.GroupEntry, user models4debtus.AppUserOBSOLETE) (err error) {
	//	logus.Debugf(ctx, "afterGroupInsert()")
	//	if !hasTgGroupEntity {
	//		data := &models4auth.TgGroupData{
	//			SpaceRef: group.ContactID,
	//		}
	//		tgGroup := models4auth.NewTgGroup(tgChat.ContactID, data)
	//		if err = dtdal.TgGroup.SaveTgGroup(ctx, tx, tgGroup); err != nil {
	//			return
	//		}
	//	}
	//
	//	_ = user.Data.AddGroup(group, whc.GetBotCode())
	//	chatData.SpaceRef = group.ContactID // TODO: !!! has to be updated in transaction!!!
	//	if err = whc.SaveBotChat(ctx); err != nil {
	//		return
	//	}
	//	return
	//}
	//
	//if space, _, err = facade4debtus.Group.CreateGroup(ctx, &groupEntity, whc.GetBotCode(), beforeGroupInsert, afterGroupInsert); err != nil {
	//	return
	//}
	return
}

func getTgChatEntity(whc botsfw.WebhookContext) (tgChatEntity *anybot.SneatAppTgChatDbo, err error) {
	chatEntity := whc.ChatData()
	if chatEntity == nil {
		whc.Input().LogRequest()
		logus.Debugf(whc.Context(), "can't get group as chatEntity == nil")
		return
	}
	var ok bool
	if tgChatEntity, ok = chatEntity.(*anybot.SneatAppTgChatDbo); !ok {
		logus.Debugf(whc.Context(), "whc.ChatData() is not TgChatEntityBase")
		return
	}
	return tgChatEntity, nil
}

func NewGroupTelegramInlineButton(whc botsfw.WebhookContext, groupsMessageID int) tgbotapi.InlineKeyboardButton {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "https://t.me/%v?startgroup=utm_s=%v__utm_m=%v__l=%v", whc.GetBotCode(), whc.GetBotCode(), "tgbot", whc.Locale().Code5)
	if groupsMessageID != 0 {
		buf.WriteString("__grpsMsgID=")
		buf.WriteString(strconv.Itoa(groupsMessageID))
	}
	return tgbotapi.InlineKeyboardButton{
		Text: whc.CommandText(trans.COMMAND_TEXT_ADD_GROUP, ""),
		URL:  buf.String(),
	}
}
