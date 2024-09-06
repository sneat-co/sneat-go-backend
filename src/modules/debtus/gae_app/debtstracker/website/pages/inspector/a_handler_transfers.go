package inspector

import (
	"context"
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/datastore"
	"net/http"
	"sync"
	"time"
)

type transfersPage struct {
}

func (h transfersPage) transfersPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)

	now := time.Now()

	urlQuery := r.URL.Query()

	currency := money.CurrencyCode(urlQuery.Get("currency"))

	contactID := urlQuery.Get("debtusContact")
	if contactID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, errors.New("debtusContact ContactID is empty"))
	}

	spaceID := urlQuery.Get("space")

	debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)

	var transfers []models4debtus.TransferEntry

	var transfersTotalWithoutInterest decimal.Decimal64p2

	debtusContact := models4debtus.NewDebtusSpaceContactEntry(spaceID, contactID, nil)

	wg := new(sync.WaitGroup)

	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var db dal.ReadSession
		if db, err = facade.GetSneatDB(c); err != nil {
			return
		}
		if err = db.GetMulti(c, []dal.Record{debtusContact.Record, debtusSpace.Record}); err != nil {
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var db dal.ReadSession
		if db, err = facade.GetSneatDB(c); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}
		if transfers, transfersTotalWithoutInterest, err = h.processTransfers(c, db, contactID, currency); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}
	}()

	wg.Wait()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
	}

	balancesWithoutInterest := balanceRow{
		user:      debtusSpace.Data.Contacts[contactID].Balance[currency],
		contacts:  debtusContact.Data.Balance[currency],
		transfers: transfersTotalWithoutInterest,
	}

	balancesWithInterest := balanceRow{}
	if balance, err := debtusSpace.Data.Contacts[contactID].BalanceWithInterest(c, now); err == nil {
		balancesWithInterest.user = balance[currency]
	} else {
		balancesWithInterest.userContactBalanceErr = err
	}
	if balance, err := debtusContact.Data.BalanceWithInterest(c, now); err == nil {
		balancesWithInterest.contacts = balance[currency]
	} else {
		balancesWithInterest.contactBalanceErr = err
	}

	renderTransfersPage(debtusContact, currency, balancesWithoutInterest, balancesWithInterest, transfers, w)
}

func (h transfersPage) processTransfers(ctx context.Context, tx dal.ReadSession, contactID string, currency money.CurrencyCode) (
	transfers []models4debtus.TransferEntry,
	balanceWithoutInterest decimal.Decimal64p2,
	err error,
) {
	query := dal.From(models4debtus.TransfersCollection).
		Where(
			dal.Field("BothCounterpartyIDs").EqualTo(contactID),
			dal.Field("Currency").EqualTo(currency),
		).
		OrderBy(dal.DescendingField("DtCreated")).
		SelectInto(models4debtus.NewTransferRecord)

	var reader dal.Reader
	if reader, err = tx.QueryReader(ctx, query); err != nil {
		return
	}
	for {
		var record dal.Record
		if record, err = reader.Next(); err != nil {
			if err == datastore.Done {
				err = nil
				break
			}
			panic(err)
		}
		transfer := models4debtus.NewTransfer(record.Key().ID.(string), record.Data().(*models4debtus.TransferData))
		transfers = append(transfers, transfer)
		switch contactID {
		case transfer.Data.From().ContactID:
			balanceWithoutInterest -= transfer.Data.AmountInCents
		case transfer.Data.To().ContactID:
			balanceWithoutInterest += transfer.Data.AmountInCents
		default:
			panic(fmt.Sprintf("contactID != from && contactID != to: contactID=%v, from=%v, to=%v",
				contactID, transfer.Data.From().ContactID, transfer.Data.To().ContactID))
		}
	}

	return
}
