package facade4debtus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"strconv"
)

func GetLocale(c context.Context, botID string, tgChatIntID int64, userID string) (locale i18n.Locale, err error) {
	chatID := botsfwmodels.NewChatID(botID, strconv.FormatInt(tgChatIntID, 10))
	//var tgChatEntity botsfwtgmodels.ChatEntity
	//tgChatBaseData := botsfwtgmodels.NewTelegramChatBaseData()
	//chatID, new(models.DebtusTelegramChatData)
	key := dal.NewKeyWithID(botsfwtgmodels.TgChatCollection, chatID)
	var tgChat = record.NewDataWithID[string, *models4debtus.DebtusTelegramChatData](chatID, key, new(models4debtus.DebtusTelegramChatData))
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	if err = db.Get(c, tgChat.Record); err != nil {
		logus.Debugf(c, "Failed to get TgChat entity by string ContactID=%v: %v", tgChat.ID, err) // TODO: Replace with error once load by int ContactID removed
		if dal.IsNotFound(err) {
			panic("TODO: Remove this load by int ContactID")
			//if err = nds.Get(c, datastore.NewKey(c, botsfwtgmodels.TgChatCollection, "", tgChatIntID, nil), &tgChatEntity); err != nil { // TODO: Remove this load by int ContactID
			//	logus.Errorf(c, "Failed to get TgChat entity by int ContactID=%v: %v", tgChatIntID, err)
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
			user, err := dal4userus.GetUserByID(c, db, userID)
			if err != nil {
				logus.Errorf(c, fmt.Errorf("failed to get user by ContactID=%v: %w", userID, err).Error())
				return locale, err
			}
			tgChatPreferredLanguage = user.Data.PreferredLocale
		}
		if tgChatPreferredLanguage == "" {
			tgChatPreferredLanguage = i18n.LocaleCodeEnUS
			logus.Warningf(c, "tgChat.PreferredLanguage == '' && user.PreferredLanguage == '', set to %v", i18n.LocaleCodeEnUS)
		}
	}
	locale = i18n.LocalesByCode5[tgChatPreferredLanguage]
	return
}
