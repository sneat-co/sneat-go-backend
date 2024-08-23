package api4transfers

import (
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus/dto4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"net/http"

	"context"
)

func HandleAdminLatestTransfers(c context.Context, w http.ResponseWriter, r *http.Request, _ token4auth.AuthInfo) {
	transfers, err := dtdal.Transfer.LoadLatestTransfers(c, 0, 20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(([]byte)(err.Error()))
	}
	transfersToResponse(c, w, "", transfers, true)
}

func HandleUserTransfers(c context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo, user dbo4userus.UserEntry) {
	transfers, hasMore, err := dtdal.Transfer.LoadTransfersByUserID(c, user.ID, 0, 100)
	if api4debtus.HasError(c, w, err, "", "", http.StatusInternalServerError) {
		return
	}
	transfersToResponse(c, w, user.ID, transfers, hasMore)
}

func transfersToResponse(c context.Context, w http.ResponseWriter, userID string, transfers []models4debtus.TransferEntry, hasMore bool) {
	api4debtus.JsonToResponse(c, w, dto4debtus.TransfersResultDto{
		Transfers:        dto4debtus.TransfersToDto(userID, transfers),
		HasMoreTransfers: hasMore,
	})
}
