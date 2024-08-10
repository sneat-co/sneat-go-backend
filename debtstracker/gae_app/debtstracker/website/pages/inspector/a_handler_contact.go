package inspector

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
	"google.golang.org/appengine/v2"
	"net/http"
	//"sync"

	"sync"

	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"google.golang.org/appengine/v2/datastore"
)

type contactPage struct {
}

func (h contactPage) contactPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)

	contactID := r.URL.Query().Get("id")
	if contactID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, errors.New("contact ID is empty"))
		return
	}

	var contact models.ContactEntry
	var err error

	if contact, err = facade2debtus.GetContactByID(c, nil, contactID); err != nil {
		_, _ = fmt.Fprint(w, err)
		return
	}

	//var user, counterpartyUser models.AppUser
	//var counterpartyContact models.ContactEntry
	//

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		//var transfers []models.TransferEntry
		if _, err = h.verifyTransfers(c, contactID); err != nil {
			panic(err)
		}
	}()

	//
	//wg.Add(1)
	//go func() {
	//	if user, err = facade2debtus.User.GetUserByID(c, contact.UserID); err != nil {
	//		return
	//	}
	//
	//}()
	//
	//if contact.CounterpartyUserID != 0 {
	//	wg.Add(1)
	//	if user, err = facade2debtus.User.GetUserByID(c, contact.CounterpartyUserID); err != nil {
	//		return
	//	}
	//}
	//
	//if contact.CounterpartyCounterpartyID != 0 {
	//	wg.Add(1)
	//	if counterpartyContact, err = facade2debtus.GetContactByID(c, tx, contact.CounterpartyCounterpartyID); err != nil {
	//		return
	//	}
	//}

	RenderContactPage(contact, w)

	//renderContactUsers(w, user, counterpartyUser)

}

func (contactPage) verifyTransfers(c context.Context, contactID string) (
	transfers []models.TransferEntry, err error,
) {

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	//select := dal.Select{
	//	From: &dal.CollectionRef{Name: models.TransfersCollection},
	//}
	query := dal.From(models.TransfersCollection).
		Where(dal.Field("BothCounterpartyIDs").EqualTo(contactID)).
		SelectInto(models.NewTransferRecord)

	var reader dal.Reader
	if reader, err = db.QueryReader(c, query); err != nil {
		return
	}

	for {
		//transferEntity := new(models.TransferData)
		//var key *datastore.Key
		var record dal.Record
		if record, err = reader.Next(); err != nil {
			if err == datastore.Done {
				break
			}
			panic(err)
		}
		transfers = append(transfers, models.NewTransfer(
			record.Key().ID.(string),
			record.Data().(*models.TransferData),
		))
	}

	return
}

//func renderContactUsers(w http.ResponseWriter, user, counterpartyUser models.AppUser) {
//
//}
//
//func renderCounterparty(w http.ResponseWriter, counterpartyUser models.AppUser, counterpartyContact models.ContactEntry) {
//
//}
