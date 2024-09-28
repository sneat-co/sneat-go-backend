package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
)

type UserEmailGaeDal struct {
}

func NewUserEmailGaeDal() UserEmailGaeDal {
	return UserEmailGaeDal{}
}

func (UserEmailGaeDal) GetUserEmailByID(ctx context.Context, tx dal.ReadSession, email string) (userEmail models4auth.UserEmailEntry, err error) {
	userEmail = models4auth.NewUserEmail(email, nil)
	return userEmail, tx.Get(ctx, userEmail.Record)
}

func (UserEmailGaeDal) SaveUserEmail(ctx context.Context, tx dal.ReadwriteTransaction, userEmail models4auth.UserEmailEntry) (err error) {
	return tx.Set(ctx, userEmail.Record)
}
