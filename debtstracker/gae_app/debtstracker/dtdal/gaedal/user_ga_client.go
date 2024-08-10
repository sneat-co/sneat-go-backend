package gaedal

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"time"
)

type UserGaClientDalGae struct {
}

func NewUserGaClientDalGae() UserGaClientDalGae {
	return UserGaClientDalGae{}
}

func (UserGaClientDalGae) SaveGaClient(c context.Context, gaClientId, userAgent, ipAddress string) (gaClient models.GaClient, err error) {
	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		data := new(models.GaClientEntity)
		key := dal.NewKeyWithID(models.GaClientKind, gaClientId)
		r := dal.NewRecordWithData(key, data)

		if err = tx.Get(c, r); err != nil && !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get UserGaClient by ID: %w", err)
		}
		if data.UserAgent != userAgent || data.IpAddress != ipAddress {
			data.UserAgent = userAgent
			data.IpAddress = ipAddress
			data.Created = time.Now()
			if err = tx.Set(c, r); err != nil {
				err = fmt.Errorf("failed to save UserGaClient: %w", err)
				return err
			}
		}
		return nil
	}, nil)
	return
}
