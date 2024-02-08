package transfers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api/dto"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/log"
	"net/http"
	"time"
)

func HandleCreateTransfer(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	var direction models.TransferDirection
	switch r.PostFormValue("direction") {
	case "contact-to-user":
		direction = models.TransferDirectionCounterparty2User
	case "user-to-contact":
		direction = models.TransferDirectionUser2Counterparty
	default:
		w.WriteHeader(http.StatusBadRequest)
		m := "Unknown direction: " + r.PostFormValue("direction")
		log.Debugf(c, m)
		_, _ = w.Write([]byte(m))
		return
	}
	amountValue, err := decimal.ParseDecimal64p2(r.PostFormValue("amount"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if amountValue < 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("amount < 0"))
		return
	}
	currency := r.PostFormValue("currency")
	if len(currency) > 30 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("len(currency) > 30"))
		return
	}

	amountWithCurrency := money.NewAmount(money.CurrencyCode(currency), amountValue)

	contactID := api.GetStrID(c, w, r, "contactID")
	if contactID == "" {
		return
	}

	var (
		returnToTransferID string
		dueOn              time.Time
	)

	isReturn := r.PostFormValue("isReturn") == "true"

	if dueStr := r.PostFormValue("due"); dueStr != "" {
		if dueOn, err = time.Parse("2006-01-02", dueStr); err != nil {
			api.BadRequestError(c, w, err)
			return
		}
	}

	//user, err := facade.User.GetUserByID(c, authInfo.AppUserIntID)
	//if err != nil {
	//	hashedWriter.WriteHeader(http.StatusInternalServerError)
	//	hashedWriter.Write([]byte(errors.Wrap(err, "Failed to get user")))
	//}
	if isReturn {
		if counterparty, err := facade.GetContactByID(c, nil, contactID); err != nil {
			if dal.IsNotFound(err) {
				api.BadRequestError(c, w, err)
			} else {
				api.InternalError(c, w, err)
			}
			return
		} else {
			balance := counterparty.Data.Balance()
			if balanceAmount, ok := balance[amountWithCurrency.Currency]; !ok {
				api.BadRequestMessage(c, w, fmt.Sprintf("No balance for %v", amountWithCurrency.Currency))
			} else {
				switch direction {
				case models.TransferDirectionUser2Counterparty:
					if balanceAmount > 0 {
						api.BadRequestMessage(c, w, fmt.Sprintf("balanceAmount > 0 && direction == %v", direction))
					}
				case models.TransferDirectionCounterparty2User:
					if balanceAmount < 0 {
						api.BadRequestMessage(c, w, fmt.Sprintf("balanceAmount < 0 && direction == %v", direction))
					}
				}
			}
		}
	}

	var from, to *models.TransferCounterpartyInfo

	switch direction {
	case models.TransferDirectionUser2Counterparty:
		from = models.NewFrom(authInfo.UserID, r.PostFormValue("comment"))
		to = models.NewTo(contactID)
	case models.TransferDirectionCounterparty2User:
		from = models.NewTo(contactID)
		to = models.NewFrom(authInfo.UserID, r.PostFormValue("comment"))
	default:
		panic(fmt.Sprintf("Unknown direction: %v", direction))
	}

	platform := r.PostFormValue("platform")
	if len(platform) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("len(platform) > 20"))
	}
	switch platform {
	case "web":
	case "ios":
	case "android":
	default:
		api.BadRequestMessage(c, w, "Unknown platform: "+platform)
		return
	}

	var appUser models.AppUser
	if appUser, err = facade.User.GetUserByID(c, nil, authInfo.UserID); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	newTransfer := facade.NewTransferInput(api.GetEnvironment(r),
		transferSourceSetToAPI{appPlatform: platform, createdOnID: r.Host},
		appUser,
		"",
		isReturn, returnToTransferID,
		from, to,
		amountWithCurrency, dueOn, models.NoInterest())
	output, err := facade.Transfers.CreateTransfer(c, newTransfer)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	//userBalance := json.RawMessage(user.BalanceJson)
	log.Infof(c, "transfer.DtDueOn: %v", output.Transfer.Data.DtDueOn)
	response := dto.CreateTransferResponse{
		Transfer: dto.TransferToDto(authInfo.UserID, output.Transfer),
	}

	var counterparty models.Contact
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
	api.JsonToResponse(c, w, response)
}
