package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type UserEmailGaeDal struct {
}

func NewUserEmailGaeDal() UserEmailGaeDal {
	return UserEmailGaeDal{}
}

func (UserEmailGaeDal) GetUserEmailByID(c context.Context, tx dal.ReadSession, email string) (userEmail models.UserEmail, err error) {
	userEmail = models.NewUserEmail(email, nil)
	return userEmail, tx.Get(c, userEmail.Record)
}

func (UserEmailGaeDal) SaveUserEmail(c context.Context, tx dal.ReadwriteTransaction, userEmail models.UserEmail) (err error) {
	return tx.Set(c, userEmail.Record)
}
