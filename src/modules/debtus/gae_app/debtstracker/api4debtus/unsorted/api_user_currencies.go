package unsorted

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/common4all"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"net/http"

	"context"
)

func HandleGetUserCurrencies(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ token4auth.AuthInfo, user dbo4userus.UserEntry) {
	common4all.JsonToResponse(ctx, w, user.Data.LastCurrencies)
}
