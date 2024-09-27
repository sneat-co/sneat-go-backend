package facade4auth

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/auth/dto4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
)

type AuthClient interface {
	CreateUser(ctx context.Context, userToCreate dto4auth.DataToCreateUser) (user dbo4userus.UserEntry, err error)
}
