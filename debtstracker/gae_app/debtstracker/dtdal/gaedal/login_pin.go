package gaedal

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"strings"
	"time"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

var _ dtdal.LoginPinDal = (*LoginPinDalGae)(nil)

type LoginPinDalGae struct {
}

func NewLoginPinDalGae() LoginPinDalGae {
	return LoginPinDalGae{}
}

func (LoginPinDalGae) GetLoginPinByID(c context.Context, tx dal.ReadSession, id int) (loginPin models.LoginPin, err error) {
	loginPin = models.NewLoginPin(id, nil)
	return loginPin, tx.Get(c, loginPin.Record)
}

func (LoginPinDalGae) SaveLoginPin(c context.Context, tx dal.ReadwriteTransaction, loginPin models.LoginPin) (err error) {
	return tx.Set(c, loginPin.Record)
}

func (loginPinDalGae LoginPinDalGae) CreateLoginPin(c context.Context, tx dal.ReadwriteTransaction, channel, gaClientID string, createdUserID string) (loginPin models.LoginPin, err error) {
	switch strings.ToLower(channel) {
	case "":
		return loginPin, errors.New("parameter 'channel' is not set")
	case "telegram":
	case "viber":
	default:
		return loginPin, fmt.Errorf("Unknown channel: %v", channel)
	}
	if createdUserID != "" {
		if _, err := facade.User.GetUserByID(c, nil, createdUserID); err != nil {
			return loginPin, fmt.Errorf("unknown createdUserID=%s: %w", createdUserID, err)
		}
	}

	loginPin = models.NewLoginPin(0, &models.LoginPinData{
		Channel:    channel,
		Created:    time.Now(),
		UserID:     createdUserID,
		GaClientID: gaClientID,
	})
	if err = tx.Insert(c, loginPin.Record); err != nil {
		return
	}
	loginPin.ID = loginPin.Record.Key().ID.(int)
	return
}

//func (loginPinDalGae LoginPinDalGae) GetByID(c context.Context, loginID int64) (entity *models.LoginPinData, err error) {
//	entity = new(models.LoginPinData)
//	err = gaedb.Get(c, models.NewLoginPinKey(c, loginID), entity)
//	return
//}
