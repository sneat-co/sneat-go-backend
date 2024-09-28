package facade4auth

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/dto4auth"
)

type AuthClient interface {
	CreateUser(ctx context.Context, userToCreate dto4auth.DataToCreateUser) (user dbo4userus.UserEntry, err error)
}
