package common4all

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
)

type EmailDal interface {
	InsertEmail(ctx context.Context, tx dal.ReadwriteTransaction, entity *models4auth.EmailData) (models4auth.Email, error)
	UpdateEmail(ctx context.Context, tx dal.ReadwriteTransaction, email models4auth.Email) error
	GetEmailByID(ctx context.Context, tx dal.ReadSession, id int64) (models4auth.Email, error)
}
