package api4transfers

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus/unsorted"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func HandleGetTransfer(c context.Context, w http.ResponseWriter, r *http.Request) {
	if transferID := api4debtus.GetStrID(c, w, r, "id"); transferID == "" {
		return
	} else {
		transfer, err := facade4debtus.Transfers.GetTransferByID(c, nil, transferID)
		if api4debtus.HasError(c, w, err, models4debtus.TransfersCollection, transferID, http.StatusBadRequest) {
			return
		}

		err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
			if err = facade4debtus.CheckTransferCreatorNameAndFixIfNeeded(c, tx, transfer); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			return err
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		record := unsorted.NewReceiptTransferDto(c, transfer)

		api4debtus.JsonToResponse(c, w, &record)
	}
}

type transferSourceSetToAPI struct {
	appPlatform string
	createdOnID string
}

func (s transferSourceSetToAPI) PopulateTransfer(t *models4debtus.TransferData) {
	t.CreatedOnPlatform = s.appPlatform
	t.CreatedOnID = s.createdOnID
}
