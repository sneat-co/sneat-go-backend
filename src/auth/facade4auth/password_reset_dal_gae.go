package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-core/facade"
)

func NewPasswordResetDalGae() PasswordResetDalGae {
	return PasswordResetDalGae{}
}

type PasswordResetDalGae struct {
}

var _ PasswordResetDal = (*PasswordResetDalGae)(nil)

func (PasswordResetDalGae) GetPasswordResetByID(c context.Context, tx dal.ReadSession, id int) (passwordReset models4auth.PasswordReset, err error) {
	passwordReset = models4auth.NewPasswordReset(id, nil)
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	return passwordReset, tx.Get(c, passwordReset.Record)
}

func (PasswordResetDalGae) CreatePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, entity *models4auth.PasswordResetData) (passwordReset models4auth.PasswordReset, err error) {
	passwordReset = models4auth.NewPasswordReset(0, entity)
	if err = tx.Insert(c, passwordReset.Record); err != nil {
		return
	}
	passwordReset.ID = passwordReset.Key.ID.(int)
	return
}

func (PasswordResetDalGae) SavePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, passwordReset models4auth.PasswordReset) (err error) {
	return tx.Set(c, passwordReset.Record)
}
