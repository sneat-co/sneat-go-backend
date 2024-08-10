package transfers

import (
	"context"
	"encoding/json"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus/dto"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	sneatfacade "github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func HandleCreateTransfer(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	var request dto.CreateTransferRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx sneatfacade.User) (interface{}, error) {
			var from, to *models.TransferCounterpartyInfo

			appUser, err := facade2debtus.User.GetUserByID(c, nil, authInfo.UserID)
			if err != nil {
				return nil, err
			}

			newTransfer := dto.NewTransferInput(api.GetEnvironment(r),
				transferSourceSetToAPI{appPlatform: "api", createdOnID: r.Host},
				appUser,
				request,
				from, to,
			)

			output, err := facade2debtus.Transfers.CreateTransfer(c, newTransfer)
			if err != nil {
				return nil, err
			}

			response := dto.CreateTransferResponse{
				Transfer: dto.TransferToDto(authInfo.UserID, output.Transfer),
			}

			var counterparty models.ContactEntry
			switch output.Transfer.Data.CreatorUserID {
			case output.Transfer.Data.From().UserID:
				counterparty = output.To.Contact
			case output.Transfer.Data.To().UserID:
				counterparty = output.From.Contact
			default:
				panic("Unknown direction")
			}
			if counterparty.Data.BalanceJson != "" {
				counterpartyBalance := json.RawMessage(counterparty.Data.BalanceJson)
				response.CounterpartyBalance = &counterpartyBalance
			}
			return response, err
		})
}
