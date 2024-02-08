package transfers

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api/unsorted"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"net/http"
)

func HandleGetTransfer(c context.Context, w http.ResponseWriter, r *http.Request) {
	if transferID := api.GetStrID(c, w, r, "id"); transferID == "" {
		return
	} else {
		transfer, err := facade.Transfers.GetTransferByID(c, nil, transferID)
		if api.HasError(c, w, err, models.TransferKind, transferID, http.StatusBadRequest) {
			return
		}

		var db dal.DB
		if db, err = facade.GetDatabase(c); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
			if err = facade.CheckTransferCreatorNameAndFixIfNeeded(c, tx, transfer); err != nil {
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

		api.JsonToResponse(c, w, &record)
	}
}

type transferSourceSetToAPI struct {
	appPlatform string
	createdOnID string
}

func (s transferSourceSetToAPI) PopulateTransfer(t *models.TransferData) {
	t.CreatedOnPlatform = s.appPlatform
	t.CreatedOnID = s.createdOnID
}
