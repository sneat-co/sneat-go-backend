package gaedal

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
)

type AdminDalGae struct {
}

func NewAdminDalGae() AdminDalGae {
	return AdminDalGae{}
}

func (AdminDalGae) LatestUsers(c context.Context) (users []dbo4userus.UserEntry, err error) {
	return nil, errors.New("not implemented")
	//var (
	//	userKeys     []*datastore.Key
	//	userEntities []*models.DebutsAppUserDataOBSOLETE
	//)
	//query := datastore.NewQuery(models.AppUserKind).Order("-DtCreated").Limit(20)
	//if userKeys, err = query.GetAll(c, &userEntities); err != nil {
	//	return
	//}
	//users = make([]models.AppUserOBSOLETE, len(userKeys))
	//for i, userEntity := range userEntities {
	//	users[i] = models.NewAppUserOBSOLETE(userKeys[i].IntID(), userEntity)
	//}
	//return
}

func (AdminDalGae) DeleteAll(c context.Context, botCode, botChatID string) error {
	panic("not implemented")
	//tasksCount := 7
	//await := make(chan string, tasksCount)
	//allErrors := make(chan error, tasksCount)
	//
	//deleteAllEntitiesByKind := func(kind string, completion chan string) {
	//	logus.Debugf(c, "Deleting: %v...", kind)
	//	if keys, err := datastore.NewQuery(kind).KeysOnly().GetAll(c, nil); err != nil {
	//		allErrors <- err
	//		logus.Errorf(c, "Failed to load %v entities: %v", kind, err)
	//	} else if len(keys) > 0 {
	//		logus.Debugf(c, "Loaded %v key(s) of %v kind: %v", len(keys), kind, keys)
	//		if err := gaedb.DeleteMulti(c, keys); err != nil {
	//			logus.Errorf(c, "Failed to delete %v entities of %v kind: %v", len(keys), kind, err)
	//			allErrors <- err
	//		}
	//	} else {
	//		logus.Debugf(c, "Noting to delete for: %v", kind)
	//	}
	//	completion <- kind
	//}
	//
	//kindsToDelete := []string{
	//	telegram.TgUserKind,
	//	telegram.ChatKind,
	//	telegram.ChatInstanceKind,
	//	models.TgGroupKind,
	//	models.InviteKind,
	//	models.InviteClaimKind,
	//	models.FeedbackKind,
	//	models.AppUserKind,
	//	models.TransfersCollection,
	//	models.DebtusContactsCollection,
	//	models.ReminderKind,
	//	models.ReceiptKind,
	//	models.UserBrowserKind,
	//	models.TwilioSmsKind,
	//	fbm.ChatKind,
	//	fbm.BotUserKind,
	//	models.UserFacebookCollection,
	//	models.UserGoogleCollection,
	//	models.UserOneSignalKind,
	//	models.LoginCodeKind,
	//	models.LoginPinKind,
	//	models.GroupKind,
	//	models.BillKind,
	//	models.BillScheduleKind,
	//	models.BillsHistoryKind,
	//	//viber.ViberChatKind,
	//	//viber.ViberUserKind,
	//	viber.UserChatKind,
	//	models.UserVkKind,
	//}
	//
	//for _, kind := range kindsToDelete {
	//	go deleteAllEntitiesByKind(kind, await)
	//}
	//
	//for i := 0; i < len(kindsToDelete); i++ {
	//	logus.Debugf(c, "%v - deleted: %v", i, <-await)
	//}
	//
	//close(allErrors)
	//
	//errs := make([]string, 0)
	//for err := range allErrors {
	//	errs = append(errs, err.Error())
	//}
	//
	//if err := memcache.Flush(c); err != nil {
	//	logus.Errorf(c, "Failed to flush memcache: %v", err)
	//	// Do not return
	//}
	//
	//if len(errs) > 0 {
	//	return fmt.Errorf("There were %v errors: %v", len(errs), strings.Join(errs, "\n"))
	//}
	//
	//// We need to delay deletion of chat entity as it will be put by bot framework on reply.
	//chatKey := gaehost.NewGaeTelegramChatStore(shared.TheAppContext.GetBotChatEntityFactory("telegram")).NewBotChatKey(c, botCode, botChatID)
	//if t, err := delayTgChatDeletion.Task(chatKey.StringID()); err != nil {
	//	err = fmt.Errorf("failed to create delay task for Telegram chat deletion: %w", err)
	//	return err
	//} else {
	//	t.Delay = time.Second
	//	if _, err = taskqueue.Add(c, t, shared.QUEUE_SUPPORT); err != nil {
	//		err = fmt.Errorf("failed to delay TgChat deletion: %w", err)
	//		return err
	//	}
	//}
	//
	//return nil
}

//var delayTgChatDeletion = delaying.MustRegisterFunc("delete-%v", func(c context.Context, id string) error {
//	logus.Debugf(c, "delayTgChatDeletion(id=%v)", id)
//	panic("not implemented")
//	key := gaedb.NewKey(c, telegram.ChatKind, id, 0, nil)
//	if err := gaedb.Delete(c, key); err != nil {
//		logus.Errorf(c, "Failed to delete %v: %v", key, err)
//		return err
//	}
//	if err := memcache.Flush(c); err != nil {
//		logus.Errorf(c, "Failed to clear memcache: %v", err)
//	}
//	logus.Infof(c, "%v deleted", key)
//	return nil
//})
