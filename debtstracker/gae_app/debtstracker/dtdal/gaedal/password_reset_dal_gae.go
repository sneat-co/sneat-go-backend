package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func NewPasswordResetDalGae() PasswordResetDalGae {
	return PasswordResetDalGae{}
}

type PasswordResetDalGae struct {
}

var _ dtdal.PasswordResetDal = (*PasswordResetDalGae)(nil)

func (PasswordResetDalGae) GetPasswordResetByID(c context.Context, tx dal.ReadSession, id int) (passwordReset models.PasswordReset, err error) {
	passwordReset = models.NewPasswordReset(id, nil)
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	return passwordReset, tx.Get(c, passwordReset.Record)
}

func (PasswordResetDalGae) CreatePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, entity *models.PasswordResetData) (passwordReset models.PasswordReset, err error) {
	passwordReset = models.NewPasswordReset(0, entity)
	if err = tx.Insert(c, passwordReset.Record); err != nil {
		return
	}
	passwordReset.ID = passwordReset.Key.ID.(int)
	return
}

func (PasswordResetDalGae) SavePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, passwordReset models.PasswordReset) (err error) {
	return tx.Set(c, passwordReset.Record)
}
