package unsorted

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"net/http"

	"context"
)

func HandleGetUserCurrencies(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ token4auth.AuthInfo, user dbo4userus.UserEntry) {
	api4debtus.JsonToResponse(ctx, w, user.Data.LastCurrencies)
}
