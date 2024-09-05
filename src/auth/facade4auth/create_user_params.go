package facade4auth

import (
	"github.com/sneat-co/sneat-go-backend/src/coretodo"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
)

type CreateUserWorkerParams struct {
	*dal4userus.UserWorkerParams
	coretodo.WithRecordChanges
}
