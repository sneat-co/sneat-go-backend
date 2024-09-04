package facade4auth

import (
	"github.com/sneat-co/sneat-go-backend/src/coretodo"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"time"
)

type CreateUserWorkerParams struct {
	started time.Time
	*dal4userus.UserWorkerParams
	coretodo.WithRecordChanges
}
