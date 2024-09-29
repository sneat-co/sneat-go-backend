package facade4anybot

import (
	"context"
)

type refererFacade struct {
}

var Referer = refererFacade{}

//const lastTgReferrers = "lastTgReferrers"

//var errAlreadyReferred = errors.New("already referred")

//func delayedSetUserReferrer(ctx context.Context, userID string, referredBy string) (err error) {
//	userChanged := false
//	if err = dal4userus.RunUserWorker(ctx, facade.NewUserContext(userID), true,
//		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) error {
//			if err != nil {
//				return err
//			}
//			if params.User.Data.ReferredBy != "" {
//				logus.Debugf(ctx, "already referred")
//				return nil
//			}
//			params.User.Data.ReferredBy = referredBy
//			params.User.Record.MarkAsChanged()
//			params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "referredBy", Value: referredBy})
//			userChanged = true
//			return nil
//		}); err != nil {
//		logus.Errorf(ctx, "failed to check & update user: %v", err)
//		return err
//	}
//	if userChanged {
//		logus.Infof(ctx, "User's referrer saved")
//	}
//	return nil
//}

//func delaySetUserReferrer(ctx context.Context, userID string, referredBy string) (err error) {
//	return delayerSetUserReferrer.EnqueueWork(ctx, delaying.With(const4userus.QueueUsers, "set-user-referrer", time.Second/2), userID, referredBy)
//}

//var topReferralsCacheTime = time.Hour

func (f refererFacade) AddTelegramReferrer(ctx context.Context, userID string, tgUsername, botID string) {
	panic("TODO: implement AddTelegramReferrer")
	//tgUsername = strings.ToLower(tgUsername)
	//now := time.Now()
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			logus.Errorf(ctx, "panic in refererFacade.AddTelegramReferrer(): %v", r)
	//		}
	//	}()
	//	if err := facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
	//		user, err := dal4userus.GetUserByID(ctx, tx, userID)
	//		if err != nil {
	//			logus.Errorf(ctx, err.Error())
	//			return nil
	//		}
	//		if user.Data.ReferredBy != "" {
	//			logus.Debugf(ctx, "AddTelegramReferrer() => already referred")
	//			return nil
	//		}
	//
	//		referer := models4debtus.Referer{
	//			Data: &models4debtus.RefererDbo{
	//				Platform:   "tg",
	//				ReferredTo: botID,
	//				DtCreated:  now,
	//				ReferredBy: tgUsername,
	//			},
	//		}
	//
	//		referredBy := "tg:" + tgUsername
	//
	//		if err = delaySetUserReferrer(ctx, userID, referredBy); err != nil {
	//			logus.Errorf(ctx, "Failed to delay set user referrer: %v", err)
	//			if err = delayedSetUserReferrer(ctx, userID, referredBy); err != nil {
	//				logus.Errorf(ctx, "Failed to set user referrer: %v", err)
	//				return nil
	//			}
	//		}
	//
	//		var isLocked bool
	//		item, err := memcache.Get(ctx, lastTgReferrers)
	//		if err != nil {
	//			if err == memcache.ErrCacheMiss {
	//				item = f.lockMemcacheItem(ctx)
	//				isLocked = true
	//				err = nil
	//			} else {
	//				logus.Warningf(ctx, "failed to get last-tg-referrers from memcache: %v", err)
	//			}
	//		}
	//		if err := tx.Insert(ctx, referer.Record); err != nil {
	//			logus.Errorf(ctx, "failed to insert referer to DB: %v", err)
	//		}
	//		if item == nil {
	//			if err = memcache.Delete(ctx, lastTgReferrers); err != nil {
	//				logus.Warningf(ctx, "Failed to clear memcache item: %v", err) // TODO: add a queue task to remove?
	//				return nil
	//			}
	//		} else {
	//			var tgUsernames []string
	//			if isLocked {
	//				tgUsernames = []string{tgUsername}
	//			} else {
	//				tgUsernames = append(strings.Split(string(item.Value), ","), tgUsername)
	//				if len(tgUsernames) > 100 {
	//					tgUsernames = tgUsernames[:100]
	//				}
	//			}
	//			item.Value = []byte(strings.Join(tgUsernames, ","))
	//			item.Expiration = topReferralsCacheTime
	//			if err = memcache.CompareAndSwap(ctx, item); err != nil {
	//				if err = memcache.Delete(ctx, lastTgReferrers); err != nil {
	//					logus.Warningf(ctx, "failed to delete '%v' from memcache", lastTgReferrers)
	//				}
	//			}
	//		}
	//		return nil
	//	}); err != nil {
	//		panic(err)
	//	}
	//
	//}()
}

//func (refererFacade) lockMemcacheItem(ctx context.Context) (item *memcache.Item) {
//	lock := make([]byte, 9)
//	lock[0] = []byte("_")[0]
//	binary.LittleEndian.PutUint64(lock[1:], rand.Uint64())
//	item = &memcache.Item{
//		Key:        lastTgReferrers,
//		Value:      lock,
//		Expiration: time.Second * 10,
//	}
//
//	if err := memcache.Set(ctx, item); err == nil {
//		if item, err = memcache.Get(ctx, item.Key); err != nil {
//			logus.Warningf(ctx, "memcache error: %v", err)
//			item = nil
//		} else if !bytes.Equal(lock, item.Value) {
//			item = nil
//		}
//	}
//	return
//}

func (f refererFacade) TopTelegramReferrers(ctx context.Context, botID string, limit int) (topTelegramReferrers []string, err error) {
	panic("TODO: implement TopTelegramReferrers")
	//var item *memcache.Item
	//var tgUsernames []string
	//
	//isLockItem := func() bool {
	//	return item != nil && len(item.Value) == 9 && item.Value[0] == []byte("_")[0]
	//}
	//if item, err = memcache.Get(ctx, lastTgReferrers); err == nil && !isLockItem() {
	//	tgUsernames = strings.Split(string(item.Value), ",")
	//	item = nil
	//} else {
	//	q := dal.From(models4debtus.RefererKind).
	//		WhereField("p", "=", "tg").
	//		WhereField("to", "=", botID).
	//		OrderBy(dal.DescendingField("t")).
	//		Limit(100).
	//		SelectInto(func() dal.Record {
	//			return dal.NewRecordWithIncompleteKey(models4debtus.RefererKind, reflect.String, new(models4debtus.RefererDbo))
	//		})
	//	var reader dal.Reader
	//	if reader, err = dtdal.DB.QueryReader(ctx, q); err != nil {
	//		return nil, err
	//	}
	//	for {
	//		var record dal.Record
	//		if record, err = reader.Next(); err != nil {
	//			if errors.Is(err, dal.ErrNoMoreRecords) {
	//				err = nil
	//				break
	//			}
	//			return
	//		}
	//		tgUsernames = append(tgUsernames, record.Data().(*models4debtus.RefererDbo).ReferredBy)
	//	}
	//	if !isLockItem() {
	//		if item, err = memcache.Get(ctx, lastTgReferrers); err == nil && isLockItem() {
	//			item.Value = []byte(strings.Join(tgUsernames, ","))
	//			item.Expiration = topReferralsCacheTime
	//			if err = memcache.CompareAndSwap(ctx, item); err != nil {
	//				logus.Warningf(ctx, "Failed to set top referrals to memcache: %v", err)
	//				err = nil
	//			}
	//		} else { // We don't care about error here
	//			err = nil
	//		}
	//
	//	}
	//}
	//counts := make(map[string]int, len(tgUsernames))
	//for _, tgUsername := range tgUsernames {
	//	counts[tgUsername] += 1
	//}
	//
	////count := len(counts)
	////if count > limit {
	////	count = limit
	////}
	//
	//topTelegramReferrers = rankByCount(counts, limit)
	//
	//return
}

//func rankByCount(countsByName map[string]int, limit int) (names []string) {
//	pl := make(PairList, len(countsByName))
//	i := 0
//	for k, v := range countsByName {
//		pl[i] = Pair{k, v}
//		i++
//	}
//	sort.Sort(sort.Reverse(pl))
//	if len(pl) <= limit {
//		names = make([]string, len(pl))
//	} else {
//		names = make([]string, limit)
//	}
//	for i := range pl {
//		names[i] = pl[i].Key
//	}
//	return
//}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
