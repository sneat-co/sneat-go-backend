package facade4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"strings"
	"time"

	"context"
)

func insertUserBrowser(ctx context.Context, data *models4auth.UserBrowserData) (userBrowser models4auth.UserBrowser, err error) {

	userBrowser = models4auth.NewUserBrowserWithIncompleteKey(data)
	return userBrowser, facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return tx.Insert(ctx, userBrowser.Record)
	})
}

func SaveUserBrowser(ctx context.Context, userID string, userAgent string) (userBrowser models4auth.UserBrowser, err error) {
	userAgent = strings.TrimSpace(userAgent)
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" {
		return
	}
	const limit = 1
	q := dal.From(models4auth.UserBrowserKind).
		WhereField("AppUserIntID", dal.Equal, userID).
		WhereField("UserAgent", dal.Equal, userAgent)
	query := q.Limit(limit).SelectInto(models4auth.NewUserBrowserRecord)

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}

	var records []dal.Record
	if records, err = db.QueryAllRecords(ctx, query); err != nil {
		return
	}

	switch len(records) {
	case 0:
		ub := models4auth.UserBrowserData{
			UserID:      userID,
			UserAgent:   userAgent,
			LastUpdated: time.Now(),
		}
		userBrowser, err = insertUserBrowser(ctx, &ub)
		return
	case 1:
		userBrowser := records[0].Data().(*models4auth.UserBrowserData)
		if userBrowser.LastUpdated.Before(time.Now().Add(-24 * time.Hour)) {
			err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
				if err := tx.Get(ctx, records[0]); err != nil {
					return err
				}
				if userBrowser.LastUpdated.Before(time.Now().Add(-time.Hour)) {
					userBrowser.LastUpdated = time.Now()
					err = tx.Set(ctx, records[0])
				}
				return err
			})
		}
	default:
		logus.Errorf(ctx, "Loaded too many records: %v", len(records))
	}
	return
}
