package facade

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/i18n"
	"github.com/strongo/log"
	"strconv"
)

func GetLocale(c context.Context, botID string, tgChatIntID int64, userID string) (locale i18n.Locale, err error) {
	chatID := botsfwmodels.NewChatID(botID, strconv.FormatInt(tgChatIntID, 10))
	//var tgChatEntity botsfwtgmodels.ChatEntity
	//tgChatBaseData := botsfwtgmodels.NewTelegramChatBaseData()
	//chatID, new(models.DebtusTelegramChatData)
	key := dal.NewKeyWithID(botsfwtgmodels.TgChatCollection, chatID)
	var tgChat = record.NewDataWithID[string, *models.DebtusTelegramChatData](chatID, key, new(models.DebtusTelegramChatData))
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	if err = db.Get(c, tgChat.Record); err != nil {
		log.Debugf(c, "Failed to get TgChat entity by string ID=%v: %v", tgChat.ID, err) // TODO: Replace with error once load by int ID removed
		if dal.IsNotFound(err) {
			panic("TODO: Remove this load by int ID")
			//if err = nds.Get(c, datastore.NewKey(c, botsfwtgmodels.TgChatCollection, "", tgChatIntID, nil), &tgChatEntity); err != nil { // TODO: Remove this load by int ID
			//	log.Errorf(c, "Failed to get TgChat entity by int ID=%v: %v", tgChatIntID, err)
			//	return
			//}
		} else {
			return
		}
	}
	tgChatPreferredLanguage := tgChat.Data.BaseChatData().PreferredLanguage
	if tgChatPreferredLanguage == "" {
		if userID == "" && tgChat.Data.BaseChatData().AppUserID != "" {
			userID = tgChat.Data.BaseChatData().AppUserID
		}
		if userID != "" {
			var db dal.DB
			if db, err = GetDatabase(c); err != nil {
				return
			}
			user, err := User.GetUserByID(c, db, userID)
			if err != nil {
				log.Errorf(c, fmt.Errorf("failed to get user by ID=%v: %w", userID, err).Error())
				return locale, err
			}
			tgChatPreferredLanguage = user.Data.PreferredLanguage
		}
		if tgChatPreferredLanguage == "" {
			tgChatPreferredLanguage = i18n.LocaleCodeEnUS
			log.Warningf(c, "tgChat.PreferredLanguage == '' && user.PreferredLanguage == '', set to %v", i18n.LocaleCodeEnUS)
		}
	}
	locale = i18n.LocalesByCode5[tgChatPreferredLanguage]
	return
}
