package facade2debtus

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2/memcache"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

type refererFacade struct {
}

var Referer = refererFacade{}

const lastTgReferrers = "lastTgReferrers"

//var errAlreadyReferred = errors.New("already referred")

func setUserReferrer(c context.Context, userID string, referredBy string) (err error) {
	userChanged := false
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		user, err := User.GetUserByID(c, tx, userID)
		if err != nil {
			return err
		}
		if user.Data.ReferredBy != "" {
			logus.Debugf(c, "already referred")
			return nil
		}
		user.Data.ReferredBy = referredBy
		userChanged = true
		return User.SaveUser(c, tx, user)
	}); err != nil {
		logus.Errorf(c, "failed to check & update user: %v", err)
		return err
	}
	if userChanged {
		logus.Infof(c, "User's referrer saved")
	}
	return nil
}

func delaySetUserReferrer(c context.Context, userID string, referredBy string) (err error) {
	return delayedSetUserReferrer.EnqueueWork(c, delaying.With(common.QUEUE_USERS, "set-user-referrer", time.Second/2), userID, referredBy)
}

var topReferralsCacheTime = time.Hour

func (f refererFacade) AddTelegramReferrer(c context.Context, userID string, tgUsername, botID string) {
	tgUsername = strings.ToLower(tgUsername)
	now := time.Now()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logus.Errorf(c, "panic in refererFacade.AddTelegramReferrer(): %v", r)
			}
		}()
		if err := facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
			user, err := User.GetUserByID(c, tx, userID)
			if err != nil {
				logus.Errorf(c, err.Error())
				return nil
			}
			if user.Data.ReferredBy != "" {
				logus.Debugf(c, "AddTelegramReferrer() => already referred")
				return nil
			}

			referer := models.Referer{
				Data: &models.RefererDbo{
					Platform:   "tg",
					ReferredTo: botID,
					DtCreated:  now,
					ReferredBy: tgUsername,
				},
			}

			referredBy := "tg:" + tgUsername

			if err = delaySetUserReferrer(c, userID, referredBy); err != nil {
				logus.Errorf(c, "Failed to delay set user referrer: %v", err)
				if err = setUserReferrer(c, userID, referredBy); err != nil {
					logus.Errorf(c, "Failed to set user referrer: %v", err)
					return nil
				}
			}

			var isLocked bool
			item, err := memcache.Get(c, lastTgReferrers)
			if err != nil {
				if err == memcache.ErrCacheMiss {
					item = f.lockMemcacheItem(c)
					isLocked = true
					err = nil
				} else {
					logus.Warningf(c, "failed to get last-tg-referrers from memcache: %v", err)
				}
			}
			if err := tx.Insert(c, referer.Record); err != nil {
				logus.Errorf(c, "failed to insert referer to DB: %v", err)
			}
			if item == nil {
				if err = memcache.Delete(c, lastTgReferrers); err != nil {
					logus.Warningf(c, "Failed to clear memcache item: %v", err) // TODO: add a queue task to remove?
					return nil
				}
			} else {
				var tgUsernames []string
				if isLocked {
					tgUsernames = []string{tgUsername}
				} else {
					tgUsernames = append(strings.Split(string(item.Value), ","), tgUsername)
					if len(tgUsernames) > 100 {
						tgUsernames = tgUsernames[:100]
					}
				}
				item.Value = []byte(strings.Join(tgUsernames, ","))
				item.Expiration = topReferralsCacheTime
				if err = memcache.CompareAndSwap(c, item); err != nil {
					if err = memcache.Delete(c, lastTgReferrers); err != nil {
						logus.Warningf(c, "failed to delete '%v' from memcache", lastTgReferrers)
					}
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}

	}()
}

func (refererFacade) lockMemcacheItem(c context.Context) (item *memcache.Item) {
	lock := make([]byte, 9)
	lock[0] = []byte("_")[0]
	binary.LittleEndian.PutUint64(lock[1:], rand.Uint64())
	item = &memcache.Item{
		Key:        lastTgReferrers,
		Value:      lock,
		Expiration: time.Second * 10,
	}

	if err := memcache.Set(c, item); err == nil {
		if item, err = memcache.Get(c, item.Key); err != nil {
			logus.Warningf(c, "memcache error: %v", err)
			item = nil
		} else if !bytes.Equal(lock, item.Value) {
			item = nil
		}
	}
	return
}

func (f refererFacade) TopTelegramReferrers(c context.Context, botID string, limit int) (topTelegramReferrers []string, err error) {
	var item *memcache.Item
	var tgUsernames []string

	isLockItem := func() bool {
		return item != nil && len(item.Value) == 9 && item.Value[0] == []byte("_")[0]
	}
	if item, err = memcache.Get(c, lastTgReferrers); err == nil && !isLockItem() {
		tgUsernames = strings.Split(string(item.Value), ",")
		item = nil
	} else {
		q := dal.From(models.RefererKind).
			WhereField("p", "=", "tg").
			WhereField("to", "=", botID).
			OrderBy(dal.DescendingField("t")).
			Limit(100).
			SelectInto(func() dal.Record {
				return dal.NewRecordWithIncompleteKey(models.RefererKind, reflect.String, new(models.RefererDbo))
			})
		var reader dal.Reader
		if reader, err = dtdal.DB.QueryReader(c, q); err != nil {
			return nil, err
		}
		for {
			var record dal.Record
			if record, err = reader.Next(); err != nil {
				if errors.Is(err, dal.ErrNoMoreRecords) {
					err = nil
					break
				}
				return
			}
			tgUsernames = append(tgUsernames, record.Data().(*models.RefererDbo).ReferredBy)
		}
		if !isLockItem() {
			if item, err = memcache.Get(c, lastTgReferrers); err == nil && isLockItem() {
				item.Value = []byte(strings.Join(tgUsernames, ","))
				item.Expiration = topReferralsCacheTime
				if err = memcache.CompareAndSwap(c, item); err != nil {
					logus.Warningf(c, "Failed to set top referrals to memcache: %v", err)
					err = nil
				}
			} else { // We don't care about error here
				err = nil
			}

		}
	}
	counts := make(map[string]int, len(tgUsernames))
	for _, tgUsername := range tgUsernames {
		counts[tgUsername] += 1
	}

	//count := len(counts)
	//if count > limit {
	//	count = limit
	//}

	topTelegramReferrers = rankByCount(counts, limit)

	return
}

func rankByCount(countsByName map[string]int, limit int) (names []string) {
	pl := make(PairList, len(countsByName))
	i := 0
	for k, v := range countsByName {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	if len(pl) <= limit {
		names = make([]string, len(pl))
	} else {
		names = make([]string, limit)
	}
	for i := range pl {
		names[i] = pl[i].Key
	}
	return
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
