package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
)

type EmailDalGae struct {
}

func NewEmailDalGae() EmailDalGae {
	return EmailDalGae{}
}

var _ dtdal.EmailDal = (*EmailDalGae)(nil)

func (EmailDalGae) InsertEmail(c context.Context, tx dal.ReadwriteTransaction, data *models4auth.EmailData) (email models4auth.Email, err error) {
	email = models4auth.NewEmail(0, data)
	if err = tx.Insert(c, email.Record); err != nil {
		return
	}
	email.ID = email.Record.Key().ID.(int64)
	email.Data = data
	return
}

func (EmailDalGae) UpdateEmail(c context.Context, tx dal.ReadwriteTransaction, email models4auth.Email) (err error) {
	return tx.Set(c, email.Record)
}

func (EmailDalGae) GetEmailByID(c context.Context, tx dal.ReadSession, id int64) (email models4auth.Email, err error) {
	email = models4auth.NewEmail(id, nil)
	return email, tx.Get(c, email.Record)
}
