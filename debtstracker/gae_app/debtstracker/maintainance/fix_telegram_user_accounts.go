package maintainance

//import (
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
//	"context"
//	"errors"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/dal-go/dalgo/dal"
//	"google.golang.org/appengine/v2/datastore"
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
//func (m *verifyTelegramUserAccounts) Next(c context.Context, counters mapper.Counters, key *datastore.Key) (err error) {
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
//func (m *verifyTelegramUserAccounts) dealWithIntKey(c context.Context, counters *asyncCounters, key *datastore.Key, tgChatEntity *models.DebtusTelegramChatData) (err error) {
//	panic("TODO: implement me")
//	//if tgChatEntity.BotID == "" {
//	//	counters.Increment("empty_BotID_count", 1)
//	//	if err = datastore.Delete(c, key); err != nil {
//	//		log.Errorf(c, "failed to delete %v: %v", key.IntID(), err)
//	//		return nil
//	//	}
//	//	counters.Increment("empty_BotID_deleted", 1)
//	//}
//	//var tgChat models.DebtusTelegramChat
//	//if tgChat, err = dtdal.TgChat.GetTgChatByID(c, tgChatEntity.BotID, tgChatEntity.TelegramUserID); err != nil {
//	//	if dal.IsNotFound(err) {
//	//		//tgChat.SetID(tgChatEntity.BotID, tgChatEntity.TelegramUserID)
//	//		//tgChat.SetEntity(tgChatEntity)
//	//		if err = dtdal.DB.Update(c, &tgChat); err != nil {
//	//			log.Errorf(c, "failed to created entity with fixed key %v: %v", tgChat.ID, err)
//	//			return nil
//	//		}
//	//		if err = datastore.Delete(c, key); err != nil {
//	//			log.Errorf(c, "failed to delete migrated %v: %v", key.IntID(), err)
//	//			return nil
//	//		}
//	//		counters.Increment("migrated", 1)
//	//	}
//	//} else if tgChat.BotID == tgChatEntity.BotID && tgChat.TelegramUserID == tgChatEntity.TelegramUserID {
//	//	if err = datastore.Delete(c, key); err != nil {
//	//		log.Errorf(c, "failed to delete already migrated %v: %v", key.IntID(), err)
//	//		return nil
//	//	}
//	//	counters.Increment("already_migrated_so_deleted", 1)
//	//} else {
//	//	counters.Increment("mismatches", 1)
//	//	if tgChat.BotID != tgChatEntity.BotID {
//	//		log.Warningf(c, "%v: tgChat.BotID != tgChatEntity.BotID: %v != %v", key.IntID(), tgChat.BotID, tgChatEntity.BotID)
//	//	} else if tgChat.TelegramUserID != tgChatEntity.TelegramUserID {
//	//		log.Warningf(c, "%v: tgChat.TelegramUserID != tgChatEntity.TelegramUserID: %v != %v", key.IntID(), tgChat.TelegramUserID, tgChatEntity.TelegramUserID)
//	//	}
//	//}
//	//return
//}
//
//func (m *verifyTelegramUserAccounts) processTelegramChat(c context.Context, tgChat models.DebtusTelegramChat, counters *asyncCounters) (err error) {
//	panic("TODO: implement")
//	//var (
//	//	user        models.AppUser
//	//	userChanged bool
//	//)
//	//if tgChat.BotID == "" || tgChat.TelegramUserID == 0 {
//	//	log.Warningf(c, "TgChat(%v) => BotID=%v, TelegramUserID=%v", tgChat.ID, tgChat.TelegramUserID)
//	//	if strings.Contains(tgChat.ID, ":") {
//	//		botID := strings.Split(tgChat.ID, ":")[0]
//	//		tgUserID := tgChat.TelegramUserID
//	//		if tgUserID == 0 {
//	//			if tgUserID, err = strconv.ParseInt(strings.Split(tgChat.ID, ":")[1], 10, 64); err != nil {
//	//				return
//	//			}
//	//		}
//	//		if err = dtdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
//	//			if tgChat, err = dtdal.TgChat.GetTgChatByID(c, botID, tgUserID); err != nil {
//	//				return
//	//			}
//	//			tgChat.TelegramUserID = tgUserID
//	//			tgChat.BotID = botID
//	//			return dtdal.DB.Update(c, &tgChat)
//	//		}, db.CrossGroupTransaction); err != nil {
//	//			log.Errorf(c, "Failed to fix TgChat(%v): %v", tgChat.ID, err)
//	//			err = nil
//	//			return
//	//		}
//	//		log.Infof(c, "Fixed TgChat(%v)", tgChat.ID)
//	//	} else {
//	//		return
//	//	}
//	//}
//	//if tgChat.AppUserIntID == 0 {
//	//	log.Warningf(c, "TgChat(%v).AppUserIntID == 0", tgChat.ID)
//	//	return
//	//}
//	//if err = dtdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
//	//	if user, err = facade.User.GetUserByID(c, tgChat.AppUserIntID); err != nil {
//	//		if dal.IsNotFound(err) {
//	//			log.Errorf(c, "Failed to process %v: %v", tgChat.ID, err)
//	//			err = nil
//	//		}
//	//		return
//	//	}
//	//	telegramAccounts := user.Data.GetTelegramAccounts()
//	//	tgChatStrID := strconv.FormatInt(tgChat.TelegramUserID, 10)
//	//	for _, ua := range telegramAccounts {
//	//		if ua.ID == tgChatStrID {
//	//			if ua.App == tgChat.BotID {
//	//				goto userAccountFound
//	//			} else if ua.App == "" {
//	//				//log.Debugf(c, "will be fixed")
//	//				user.Data.RemoveAccount(ua)
//	//				ua.App = tgChat.BotID
//	//				userChanged = user.Data.AddAccount(ua) || userChanged
//	//				goto userAccountFound
//	//			}
//	//		}
//	//	}
//	//	userChanged = user.Data.AddAccount(users.Account{
//	//		ID:       strconv.FormatInt(tgChat.TelegramUserID, 10),
//	//		App:      tgChat.BotID,
//	//		Provider: telegram.PlatformID,
//	//	}) || userChanged
//	//userAccountFound:
//	//	if userChanged {
//	//		//log.Debugf(c, "user changed %v", user.ID)
//	//		defer func() {
//	//			if r := recover(); r != nil {
//	//				log.Errorf(c, "panic on saving user %v: %v", user.ID, r)
//	//				err = fmt.Errorf("panic on saving user %v: %v", user.ID, r)
//	//			}
//	//		}()
//	//		if err = facade.User.SaveUser(c, tx, user); err != nil {
//	//			return
//	//		}
//	//		//} else {
//	//		//	log.Debugf(c, "user NOT changed %v", user.ID)
//	//	}
//	//	return
//	//}, db.CrossGroupTransaction); err != nil {
//	//	counters.Increment("failed", 1)
//	//	return
//	//} else if userChanged {
//	//	log.Infof(c, "User %v fixed", user.ID)
//	//	counters.Increment("users-changed", 1)
//	//}
//	//return
//}
