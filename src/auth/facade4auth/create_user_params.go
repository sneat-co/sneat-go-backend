package facade4auth

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
)

type CreateUserWorkerParams struct {
	*dal4userus.UserWorkerParams
	record.WithRecordChanges
}
