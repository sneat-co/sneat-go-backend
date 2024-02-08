package api

import (
	"net/http"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func handleGetUserCurrencies(c context.Context, w http.ResponseWriter, _ *http.Request, _ auth.AuthInfo, user models.AppUser) {
	jsonToResponse(c, w, user.Data.LastCurrencies)
}
