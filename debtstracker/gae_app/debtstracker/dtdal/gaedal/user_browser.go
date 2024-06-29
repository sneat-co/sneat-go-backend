package gaedal

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/strongo/logus"
	"strings"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type UserBrowserDalGae struct {
}

func NewUserBrowserDalGae() UserBrowserDalGae {
	return UserBrowserDalGae{}
}

func (UserBrowserDalGae) insertUserBrowser(c context.Context, data *models.UserBrowserData) (userBrowser models.UserBrowser, err error) {

	userBrowser = models.NewUserBrowserWithIncompleteKey(data)
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	return userBrowser, db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		return tx.Insert(c, userBrowser.Record)
	})
}

func (userBrowserDalGae UserBrowserDalGae) SaveUserBrowser(c context.Context, userID string, userAgent string) (userBrowser models.UserBrowser, err error) {
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" {
		panic("Missign required parameter userAgent")
	}
	const limit = 1
	q := dal.From(models.UserBrowserKind).
		WhereField("AppUserIntID", dal.Equal, userID).
		WhereField("UserAgent", dal.Equal, userAgent)
	query := q.Limit(limit).SelectInto(models.NewUserBrowserRecord)

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}

	var records []dal.Record
	if records, err = db.QueryAllRecords(c, query); err != nil {
		return
	}

	switch len(records) {
	case 0:
		ub := models.UserBrowserData{
			UserID:      userID,
			UserAgent:   userAgent,
			LastUpdated: time.Now(),
		}
		userBrowser, err = userBrowserDalGae.insertUserBrowser(c, &ub)
		return
	case 1:
		userBrowser := records[0].Data().(*models.UserBrowserData)
		if userBrowser.LastUpdated.Before(time.Now().Add(-24 * time.Hour)) {
			err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
				if err := tx.Get(c, records[0]); err != nil {
					return err
				}
				if userBrowser.LastUpdated.Before(time.Now().Add(-time.Hour)) {
					userBrowser.LastUpdated = time.Now()
					err = tx.Set(c, records[0])
				}
				return err
			})
		}
	default:
		logus.Errorf(c, "Loaded too many records: %v", len(records))
	}
	return
}
