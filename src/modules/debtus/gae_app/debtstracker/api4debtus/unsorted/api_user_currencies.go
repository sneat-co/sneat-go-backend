package unsorted

import (
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"net/http"

	"context"
)

func HandleGetUserCurrencies(c context.Context, w http.ResponseWriter, _ *http.Request, _ auth.AuthInfo, user dbo4userus.UserEntry) {
	api4debtus.JsonToResponse(c, w, user.Data.LastCurrencies)
}
