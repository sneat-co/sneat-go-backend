package inspector

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
	"sync"
)

type contactPage struct {
}

func (h contactPage) contactPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := r.Context()

	contactID := r.URL.Query().Get("id")
	if contactID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, errors.New("contact ContactID is empty"))
		return
	}

	spaceID := r.URL.Query().Get("space")

	var contact models4debtus.DebtusSpaceContactEntry
	var err error

	if contact, err = facade4debtus.GetDebtusSpaceContactByID(c, nil, spaceID, contactID); err != nil {
		_, _ = fmt.Fprint(w, err)
		return
	}

	//var user, counterpartyUser models.AppUserOBSOLETE
	//var counterpartyContact models.DebtusSpaceContactEntry
	//

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		//var api4transfers []models.TransferEntry
		if _, err = h.verifyTransfers(c, contactID); err != nil {
			panic(err)
		}
	}()

	//
	//wg.Add(1)
	//go func() {
	//	if user, err = dal4userus.GetUserByID(c, contact.UserID); err != nil {
	//		return
	//	}
	//
	//}()
	//
	//if contact.CounterpartyUserID != 0 {
	//	wg.Add(1)
	//	if user, err = dal4userus.GetUserByID(c, contact.CounterpartyUserID); err != nil {
	//		return
	//	}
	//}
	//
	//if contact.CounterpartyContactID != 0 {
	//	wg.Add(1)
	//	if counterpartyContact, err = facade4debtus.GetDebtusSpaceContactByID(c, tx, contact.CounterpartyContactID); err != nil {
	//		return
	//	}
	//}

	RenderContactPage(contact, w)

	//renderContactUsers(w, user, counterpartyUser)

}

func (contactPage) verifyTransfers(ctx context.Context, contactID string) (
	transfers []models4debtus.TransferEntry, err error,
) {

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	//select := dal.Select{
	//	From: &dal.CollectionRef{Name: models.TransfersCollection},
	//}
	query := dal.From(models4debtus.TransfersCollection).
		Where(dal.Field("BothCounterpartyIDs").EqualTo(contactID)).
		SelectInto(models4debtus.NewTransferRecord)

	var reader dal.Reader
	if reader, err = db.QueryReader(ctx, query); err != nil {
		return
	}

	for {
		//transferEntity := new(models.TransferData)
		//var key *datastore.Key
		var record dal.Record
		if record, err = reader.Next(); err != nil {

			if errors.Is(err, dal.ErrNoMoreRecords) {
				break
			}
			panic(err)
		}
		transfers = append(transfers, models4debtus.NewTransfer(
			record.Key().ID.(string),
			record.Data().(*models4debtus.TransferData),
		))
	}

	return
}

//func renderContactUsers(w http.ResponseWriter, user, counterpartyUser models.AppUserOBSOLETE) {
//
//}
//
//func renderCounterparty(w http.ResponseWriter, counterpartyUser models.AppUserOBSOLETE, counterpartyContact models.DebtusSpaceContactEntry) {
//
//}
