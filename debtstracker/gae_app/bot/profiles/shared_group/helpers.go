package shared_group

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"net/url"

	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
	"strconv"
)

func GetGroup(whc botsfw.WebhookContext, callbackUrl *url.URL) (group models.Group, err error) {
	if callbackUrl != nil {
		group.ID = callbackUrl.Query().Get("group")
	}
	if group.ID == "" {
		if group.ID, err = GetUserGroupID(whc); err != nil {
			return
		}
	}

	if group.ID != "" {
		return dtdal.Group.GetGroupByID(whc.Context(), nil, group.ID)
	}

	if !whc.IsInGroup() {
		if callbackUrl != nil {
			err = errors.New("An attempt to get group ID outside of group chat without callback parameter 'group'.")
		}
		return
	}

	tgChat := whc.Input().(telegram.TgWebhookInput).TgUpdate().Chat()
	var tgChatEntity *models.DebtusTelegramChatData
	if tgChatEntity, err = getTgChatEntity(whc); err != nil {
		return
	}
	var db dal.DB
	c := whc.Context()
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		group, err = createGroupFromTelegram(c, whc, tx, tgChatEntity, tgChat)
		return err
	})
	return
}

func GetUserGroupID(whc botsfw.WebhookContext) (groupID string, err error) {
	var tgChatEntity *models.DebtusTelegramChatData
	if tgChatEntity, err = getTgChatEntity(whc); err != nil || tgChatEntity == nil {
		return
	}
	if groupID = tgChatEntity.UserGroupID; groupID != "" {
		return
	}
	return
}

func createGroupFromTelegram(c context.Context, whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, chatData *models.DebtusTelegramChatData, tgChat *tgbotapi.Chat) (group models.Group, err error) {
	log.Debugf(c, "createGroupFromTelegram()")
	var user *models.DebutsAppUserDataOBSOLETE
	if user, err = shared_all.GetUser(whc); err != nil {
		return
	}
	var chatInviteLink string

	if tgChat.IsSuperGroup() { // See: https://core.telegram.org/bots/api#exportchatinvitelink
		// TODO: Do this in delayed task - Lets try to get chat  invite link
		msg := botsfw.MessageFromBot{BotMessage: telegram.ExportChatInviteLink{}}
		if tgResponse, err := whc.Responder().SendMessage(c, msg, botsfw.BotAPISendMessageOverHTTPS); err != nil {
			log.Debugf(c, "Not able to export chat invite link: %v", err)
		} else {
			chatInviteLink = string(tgResponse.TelegramMessage.(tgbotapi.APIResponse).Result)
			log.Debugf(c, "exportInviteLink response: %v", chatInviteLink)
		}
	}

	userID := whc.AppUserID()
	groupEntity := models.GroupEntity{
		CreatorUserID: userID,
		Name:          tgChat.Title,
	}
	groupEntity.SetTelegramGroups([]models.GroupTgChatJson{
		{
			ChatID:         tgChat.ID,
			Title:          tgChat.Title,
			ChatInviteLink: chatInviteLink,
		},
	})

	hasTgGroupEntity := false
	beforeGroupInsert := func(c context.Context, groupEntity *models.GroupEntity) (group models.Group, err error) {
		log.Debugf(c, "beforeGroupInsert()")
		var tgGroup models.TgGroup
		if tgGroup, err = dtdal.TgGroup.GetTgGroupByID(c, nil, tgChat.ID); err != nil {
			if dal.IsNotFound(err) {
				err = nil
			} else {
				return
			}
		}
		if tgGroup.Data != nil && tgGroup.Data.UserGroupID != "" {
			hasTgGroupEntity = true
			return dtdal.Group.GetGroupByID(c, tx, tgGroup.Data.UserGroupID)
		}
		_, _, idx, member, members := groupEntity.AddOrGetMember(userID, "", user.FullName())
		member.TgUserID = strconv.FormatInt(int64(whc.Input().GetSender().GetID().(int)), 10)
		members[idx] = member
		groupEntity.SetGroupMembers(members)
		return
	}

	afterGroupInsert := func(c context.Context, group models.Group, user models.AppUser) (err error) {
		log.Debugf(c, "afterGroupInsert()")
		if !hasTgGroupEntity {
			data := &models.TgGroupData{
				UserGroupID: group.ID,
			}
			tgGroup := models.NewTgGroup(tgChat.ID, data)
			if err = dtdal.TgGroup.SaveTgGroup(c, tx, tgGroup); err != nil {
				return
			}
		}

		_ = user.Data.AddGroup(group, whc.GetBotCode())
		chatData.UserGroupID = group.ID // TODO: !!! has to be updated in transaction!!!
		if err = whc.SaveBotChat(c); err != nil {
			return
		}
		return
	}

	if group, _, err = facade.Group.CreateGroup(c, &groupEntity, whc.GetBotCode(), beforeGroupInsert, afterGroupInsert); err != nil {
		return
	}
	return
}

func getTgChatEntity(whc botsfw.WebhookContext) (tgChatEntity *models.DebtusTelegramChatData, err error) {
	chatEntity := whc.ChatData()
	if chatEntity == nil {
		whc.LogRequest()
		log.Debugf(whc.Context(), "can't get group as chatEntity == nil")
		return
	}
	var ok bool
	if tgChatEntity, ok = chatEntity.(*models.DebtusTelegramChatData); !ok {
		log.Debugf(whc.Context(), "whc.ChatData() is not TgChatEntityBase")
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
