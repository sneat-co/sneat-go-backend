package unsorted

import (
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"net/http"

	"context"
)

func HandleGetUserCurrencies(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ token4auth.AuthInfo, user dbo4userus.UserEntry) {
	common4all.JsonToResponse(ctx, w, user.Data.LastCurrencies)
}
