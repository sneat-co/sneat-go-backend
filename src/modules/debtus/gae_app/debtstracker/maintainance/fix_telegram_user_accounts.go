package maintainance

//import (
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
//	"context"
//	"errors"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/dal-go/dalgo/dal"
//	"net/http"
//)
//
//type verifyTelegramUserAccounts struct {
//	asyncMapper
//	entity *models.DebtusTelegramChatData
//}
//
//func (m *verifyTelegramUserAccounts) Make() interface{} {
//	m.entity = new(models.DebtusTelegramChatData)
//	return m.entity
//}
//
//func (m *verifyTelegramUserAccounts) Query(r *http.Request) (query *mapper.Query, err error) {
//	var filtered bool
//	if query, filtered, err = filterByStrID(r, "TgChat", "tgchat"); err != nil {
//		return
//	} else {
//		paramsCount := len(r.URL.Query())
//		if filtered {
//			paramsCount -= 1
//		}
//		if paramsCount != 1 {
//			err = errors.New("unexpected params: " + r.URL.RawQuery)
//		}
//	}
//	return
//}
//
//func (m *verifyTelegramUserAccounts) Next(ctx context.Context, counters mapper.Counters, key *datastore.Key) (err error) {
//	entity := *m.entity
//	if key.StringID() == "" {
//		if key.IntID() != 0 {
//			counters.Increment("integer-keys", 1)
//			return m.startWorker(c, counters, func() Worker {
//				return func(counters *asyncCounters) error {
//					return m.dealWithIntKey(c, counters, key, &entity)
//				}
//			})
//		}
//	} else {
//		tgChat := models.DebtusTelegramChat{Data: &entity}
//		tgChat.Key = dal.NewKeyWithID("TgChat", key.StringID())
//		return m.startWorker(c, counters, func() Worker {
//			return func(counters *asyncCounters) error {
//				return m.processTelegramChat(c, tgChat, counters)
//			}
//		})
//	}
//	return
//}
//
//func (m *verifyTelegramUserAccounts) dealWithIntKey(ctx context.Context, counters *asyncCounters, key *datastore.Key, tgChatEntity *models.DebtusTelegramChatData) (err error) {
//	panic("TODO: implement me")
//	//if tgChatEntity.BotID == "" {
//	//	counters.Increment("empty_BotID_count", 1)
//	//	if err = datastore.Delete(c, key); err != nil {
//	//		logus.Errorf(c, "failed to delete %v: %v", key.IntID(), err)
//	//		return nil
//	//	}
//	//	counters.Increment("empty_BotID_deleted", 1)
//	//}
//	//var tgChat models.DebtusTelegramChat
//	//if tgChat, err = dtdal.TgChat.GetTgChatByID(ctx, tgChatEntity.BotID, tgChatEntity.TelegramUserID); err != nil {
//	//	if dal.IsNotFound(err) {
//	//		//tgChat.SetID(tgChatEntity.BotID, tgChatEntity.TelegramUserID)
//	//		//tgChat.SetEntity(tgChatEntity)
//	//		if err = dtdal.DB.Update(ctx, &tgChat); err != nil {
//	//			logus.Errorf(c, "failed to created entity with fixed key %v: %v", tgChat.ContactID, err)
//	//			return nil
//	//		}
//	//		if err = datastore.Delete(ctx, key); err != nil {
//	//			logus.Errorf(c, "failed to delete migrated %v: %v", key.IntID(), err)
//	//			return nil
//	//		}
//	//		counters.Increment("migrated", 1)
//	//	}
//	//} else if tgChat.BotID == tgChatEntity.BotID && tgChat.TelegramUserID == tgChatEntity.TelegramUserID {
//	//	if err = datastore.Delete(c, key); err != nil {
//	//		logus.Errorf(c, "failed to delete already migrated %v: %v", key.IntID(), err)
//	//		return nil
//	//	}
//	//	counters.Increment("already_migrated_so_deleted", 1)
//	//} else {
//	//	counters.Increment("mismatches", 1)
//	//	if tgChat.BotID != tgChatEntity.BotID {
//	//		logus.Warningf(c, "%v: tgChat.BotID != tgChatEntity.BotID: %v != %v", key.IntID(), tgChat.BotID, tgChatEntity.BotID)
//	//	} else if tgChat.TelegramUserID != tgChatEntity.TelegramUserID {
//	//		logus.Warningf(c, "%v: tgChat.TelegramUserID != tgChatEntity.TelegramUserID: %v != %v", key.IntID(), tgChat.TelegramUserID, tgChatEntity.TelegramUserID)
//	//	}
//	//}
//	//return
//}
//
//func (m *verifyTelegramUserAccounts) processTelegramChat(ctx context.Context, tgChat models.DebtusTelegramChat, counters *asyncCounters) (err error) {
//	panic("TODO: implement")
//	//var (
//	//	user        models.AppUser
//	//	userChanged bool
//	//)
//	//if tgChat.BotID == "" || tgChat.TelegramUserID == 0 {
//	//	logus.Warningf(c, "TgChat(%v) => BotID=%v, TelegramUserID=%v", tgChat.ContactID, tgChat.TelegramUserID)
//	//	if strings.Contains(tgChat.ContactID, ":") {
//	//		botID := strings.Split(tgChat.ContactID, ":")[0]
//	//		tgUserID := tgChat.TelegramUserID
//	//		if tgUserID == 0 {
//	//			if tgUserID, err = strconv.ParseInt(strings.Split(tgChat.ContactID, ":")[1], 10, 64); err != nil {
//	//				return
//	//			}
//	//		}
//	//		if err = dtdal.DB.RunInTransaction(c, func(ctx context.Context) (err error) {
//	//			if tgChat, err = dtdal.TgChat.GetTgChatByID(ctx, botID, tgUserID); err != nil {
//	//				return
//	//			}
//	//			tgChat.TelegramUserID = tgUserID
//	//			tgChat.BotID = botID
//	//			return dtdal.DB.Update(c, &tgChat)
//	//		}, db.CrossGroupTransaction); err != nil {
//	//			logus.Errorf(c, "Failed to fix TgChat(%v): %v", tgChat.ContactID, err)
//	//			err = nil
//	//			return
//	//		}
//	//		logus.Infof(c, "Fixed TgChat(%v)", tgChat.ContactID)
//	//	} else {
//	//		return
//	//	}
//	//}
//	//if tgChat.AppUserIntID == 0 {
//	//	logus.Warningf(c, "TgChat(%v).AppUserIntID == 0", tgChat.ContactID)
//	//	return
//	//}
//	//if err = dtdal.DB.RunInTransaction(ctx, func(ctx context.Context) (err error) {
//	//	if user, err = dal4userus.GetUserByID(ctx, tgChat.AppUserIntID); err != nil {
//	//		if dal.IsNotFound(err) {
//	//			logus.Errorf(c, "Failed to process %v: %v", tgChat.ContactID, err)
//	//			err = nil
//	//		}
//	//		return
//	//	}
//	//	telegramAccounts := user.Data.GetTelegramAccounts()
//	//	tgChatStrID := strconv.FormatInt(tgChat.TelegramUserID, 10)
//	//	for _, ua := range telegramAccounts {
//	//		if ua.ContactID == tgChatStrID {
//	//			if ua.App == tgChat.BotID {
//	//				goto userAccountFound
//	//			} else if ua.App == "" {
//	//				//logus.Debugf(c, "will be fixed")
//	//				user.Data.RemoveAccount(ua)
//	//				ua.App = tgChat.BotID
//	//				userChanged = user.Data.AddAccount(ua) || userChanged
//	//				goto userAccountFound
//	//			}
//	//		}
//	//	}
//	//	userChanged = user.Data.AddAccount(users.Account{
//	//		ContactID:       strconv.FormatInt(tgChat.TelegramUserID, 10),
//	//		App:      tgChat.BotID,
//	//		Provider: telegram.PlatformID,
//	//	}) || userChanged
//	//userAccountFound:
//	//	if userChanged {
//	//		//logus.Debugf(c, "user changed %v", user.ContactID)
//	//		defer func() {
//	//			if r := recover(); r != nil {
//	//				logus.Errorf(c, "panic on saving user %v: %v", user.ContactID, r)
//	//				err = fmt.Errorf("panic on saving user %v: %v", user.ContactID, r)
//	//			}
//	//		}()
//	//		if err = facade4debtus.User.SaveUserOBSOLETE(c, tx, user); err != nil {
//	//			return
//	//		}
//	//		//} else {
//	//		//	logus.Debugf(c, "user NOT changed %v", user.ContactID)
//	//	}
//	//	return
//	//}, db.CrossGroupTransaction); err != nil {
//	//	counters.Increment("failed", 1)
//	//	return
//	//} else if userChanged {
//	//	logus.Infof(c, "User %v fixed", user.ContactID)
//	//	counters.Increment("users-changed", 1)
//	//}
//	//return
//}
