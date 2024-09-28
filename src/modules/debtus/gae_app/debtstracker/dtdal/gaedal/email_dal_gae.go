package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
)

type EmailDalGae struct {
}

func NewEmailDalGae() EmailDalGae {
	return EmailDalGae{}
}

var _ dtdal.EmailDal = (*EmailDalGae)(nil)

func (EmailDalGae) InsertEmail(ctx context.Context, tx dal.ReadwriteTransaction, data *models4auth.EmailData) (email models4auth.Email, err error) {
	email = models4auth.NewEmail(0, data)
	if err = tx.Insert(ctx, email.Record); err != nil {
		return
	}
	email.ID = email.Record.Key().ID.(int64)
	email.Data = data
	return
}

func (EmailDalGae) UpdateEmail(ctx context.Context, tx dal.ReadwriteTransaction, email models4auth.Email) (err error) {
	return tx.Set(ctx, email.Record)
}

func (EmailDalGae) GetEmailByID(ctx context.Context, tx dal.ReadSession, id int64) (email models4auth.Email, err error) {
	email = models4auth.NewEmail(id, nil)
	return email, tx.Get(ctx, email.Record)
}
