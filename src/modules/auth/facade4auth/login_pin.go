package facade4auth

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"strings"
	"time"

	"context"
	"errors"
)

var _ LoginPinDal = (*LoginPinDalGae)(nil)

type LoginPinDalGae struct {
}

func NewLoginPinDalGae() LoginPinDalGae {
	return LoginPinDalGae{}
}

func (LoginPinDalGae) GetLoginPinByID(ctx context.Context, tx dal.ReadSession, id int) (loginPin models4auth.LoginPin, err error) {
	loginPin = models4auth.NewLoginPin(id, nil)
	return loginPin, tx.Get(ctx, loginPin.Record)
}

func (LoginPinDalGae) SaveLoginPin(ctx context.Context, tx dal.ReadwriteTransaction, loginPin models4auth.LoginPin) (err error) {
	return tx.Set(ctx, loginPin.Record)
}

func (loginPinDalGae LoginPinDalGae) CreateLoginPin(ctx context.Context, tx dal.ReadwriteTransaction, channel, gaClientID string, createdUserID string) (loginPin models4auth.LoginPin, err error) {
	switch strings.ToLower(channel) {
	case "":
		return loginPin, errors.New("parameter 'channel' is not set")
	case "telegram":
	case "viber":
	default:
		return loginPin, fmt.Errorf("Unknown channel: %v", channel)
	}
	if createdUserID != "" {
		user := dbo4userus.NewUserEntry(createdUserID)
		if err = dal4userus.GetUser(ctx, tx, user); err != nil {
			return loginPin, fmt.Errorf("unknown createdUserID=%s: %w", createdUserID, err)
		}
	}

	loginPin = models4auth.NewLoginPin(0, &models4auth.LoginPinData{
		Channel:    channel,
		Created:    time.Now(),
		UserID:     createdUserID,
		GaClientID: gaClientID,
	})
	if err = tx.Insert(ctx, loginPin.Record); err != nil {
		return
	}
	loginPin.ID = loginPin.Record.Key().ID.(int)
	return
}

//func (loginPinDalGae LoginPinDalGae) GetByID(ctx context.Context, loginID int64) (entity *models.LoginPinData, err error) {
//	entity = new(models.LoginPinData)
//	err = gaedb.Get(c, models.NewLoginPinKey(ctx, loginID), entity)
//	return
//}
