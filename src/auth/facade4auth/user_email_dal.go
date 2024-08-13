package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
)

type UserEmailGaeDal struct {
}

func NewUserEmailGaeDal() UserEmailGaeDal {
	return UserEmailGaeDal{}
}

func (UserEmailGaeDal) GetUserEmailByID(c context.Context, tx dal.ReadSession, email string) (userEmail models4auth.UserEmailEntry, err error) {
	userEmail = models4auth.NewUserEmail(email, nil)
	return userEmail, tx.Get(c, userEmail.Record)
}

func (UserEmailGaeDal) SaveUserEmail(c context.Context, tx dal.ReadwriteTransaction, userEmail models4auth.UserEmailEntry) (err error) {
	return tx.Set(c, userEmail.Record)
}
