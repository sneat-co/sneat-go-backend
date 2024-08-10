package gaedal

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"math/rand"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type LoginCodeDalGae struct {
}

func NewLoginCodeDalGae() LoginCodeDalGae {
	return LoginCodeDalGae{}
}

func (LoginCodeDalGae) NewLoginCode(c context.Context, userID string) (code int, err error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	for i := 1; i < 20; i++ {
		code = random.Intn(99999) + 1
		loginCode := models.NewLoginCode(code, nil)
		if err = db.Get(c, loginCode.Record); dal.IsNotFound(err) {
			var created bool
			err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
				if err := tx.Get(c, loginCode.Record); dal.IsNotFound(err) || err == nil && loginCode.Data.Created.Add(time.Hour).Before(time.Now()) {
					loginCode.Data.Created = time.Now()
					loginCode.Data.UserID = userID
					loginCode.Data.Claimed = time.Time{}
					if err = tx.Set(c, loginCode.Record); err != nil {
						logus.Errorf(c, err.Error())
						return err
					}
					created = true
					return nil
				} else if err != nil {
					return fmt.Errorf("failed to get entity within transaction: %w", err)
				} else {
					logus.Warningf(c, "This logic code already creted outside of the current transaction")
					return nil
				}
			}, nil)
			if err != nil {
				logus.Errorf(c, fmt.Errorf("%w: transaction failed", err).Error())
			} else if created {
				return code, nil
			}
		} else if err != nil {
			return
		}
	}
	return 0, fmt.Errorf("failed to create new login code: %w", err)
}

func (LoginCodeDalGae) ClaimLoginCode(c context.Context, code int) (userID string, err error) {
	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		loginCode := models.NewLoginCode(code, nil)
		if err = tx.Get(c, loginCode.Record); err != nil {
			if dal.IsNotFound(err) {
				return err
			} else {
				return err
			}
		}
		if loginCode.Data.Created.Add(time.Minute).Before(time.Now()) {
			return models.ErrLoginCodeExpired
		}
		var emptyTime time.Time
		if loginCode.Data.Claimed == emptyTime {
			return models.ErrLoginCodeAlreadyClaimed
		}
		loginCode.Data.Claimed = time.Now()
		if err = tx.Set(c, loginCode.Record); err != nil {
			return fmt.Errorf("failed to save login code record: %w", err)
		}
		userID = loginCode.Data.UserID
		return nil
	}, nil)
	return
}
