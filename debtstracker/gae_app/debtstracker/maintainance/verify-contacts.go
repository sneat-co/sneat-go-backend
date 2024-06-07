package maintainance

//import (
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
//	"context"
//	"fmt"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/dal-go/dalgo/dal"
//	"github.com/strongo/nds"
//	"google.golang.org/appengine/v2/datastore"
//	"google.golang.org/appengine/v2/log"
//)
//
//type verifyContacts struct {
//	contactsAsyncJob
//}
//
//func (m *verifyContacts) Next(c context.Context, counters mapper.Counters, key *datastore.Key) error {
//	//log.Debugf(c, "*verifyContacts.Next(id: %v)", key.IntID())
//	return m.startContactWorker(c, counters, key, m.processContact)
//}
//
//func (m *verifyContacts) processContact(c context.Context, counters *asyncCounters, contact models.ContactEntry) (err error) {
//	if _, err = facade.User.GetUserByID(c, nil, contact.Data.UserID); dal.IsNotFound(err) {
//		counters.Increment("wrong_UserID", 1)
//		log.Warningf(c, "ContactEntry %d reference unknown user %d", contact.ID, contact.Data.UserID)
//	} else if err != nil {
//		log.Errorf(c, err.Error())
//		return
//	}
//
//	if err = m.verifyLinking(c, counters, contact); err != nil {
//		return
//	}
//
//	if err = m.verifyBalance(c, counters, contact); err != nil {
//		return
//	}
//	return
//}
//
//func (m *verifyContacts) verifyLinking(c context.Context, counters *asyncCounters, contact models.ContactEntry) (err error) {
//	if contact.Data.CounterpartyCounterpartyID != 0 {
//		var counterpartyContact models.ContactEntry
//		if counterpartyContact, err = facade.GetContactByID(c, nil, contact.Data.CounterpartyCounterpartyID); err != nil {
//			log.Errorf(c, err.Error())
//			return
//		}
//		if counterpartyContact.Data.CounterpartyCounterpartyID == 0 || counterpartyContact.Data.CounterpartyUserID == 0 {
//			if err = m.linkContacts(c, counters, contact); err != nil {
//				return
//			}
//		} else if counterpartyContact.Data.CounterpartyCounterpartyID == contact.ID && counterpartyContact.Data.CounterpartyUserID == contact.Data.UserID {
//			// Pass, we are OK
//		} else {
//			log.Warningf(c, "Wrongly linked contacts: %v=>%v != %v=>%v",
//				contact.ID, contact.Data.CounterpartyCounterpartyID,
//				counterpartyContact.ID, counterpartyContact.Data.CounterpartyCounterpartyID)
//		}
//	}
//	return
//}
//
//func (m *verifyContacts) linkContacts(c context.Context, counters *asyncCounters, contact models.ContactEntry) (err error) {
//	var counterpartyContact models.ContactEntry
//	var db dal.DB
//	if db, err = facade.GetDatabase(c); err != nil {
//		return
//	}
//	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if counterpartyContact, err = facade.GetContactByID(c, tx, contact.Data.CounterpartyCounterpartyID); err != nil {
//			log.Errorf(c, err.Error())
//			return
//		}
//		if counterpartyContact.Data.CounterpartyCounterpartyID == 0 {
//			counterpartyContact.Data.CounterpartyCounterpartyID = contact.ID
//			if counterpartyContact.Data.CounterpartyUserID == 0 {
//				counterpartyContact.Data.CounterpartyUserID = contact.Data.UserID
//			} else if counterpartyContact.Data.CounterpartyUserID != contact.Data.UserID {
//				err = fmt.Errorf("counterpartyContact(id=%v).CounterpartyUserID != contact(id=%v).UserID: %v != %v",
//					counterpartyContact.ID, contact.ID, counterpartyContact.Data.CounterpartyUserID, contact.Data.UserID)
//				return
//			}
//			if err = facade.SaveContact(c, counterpartyContact); err != nil {
//				return
//			}
//		} else if counterpartyContact.Data.CounterpartyCounterpartyID != contact.ID {
//			log.Warningf(c, "in tx: wrongly linked contacts: %v=>%v != %v=>%v",
//				contact.ID, contact.Data.CounterpartyCounterpartyID,
//				counterpartyContact.ID, counterpartyContact.Data.CounterpartyCounterpartyID)
//		}
//		return
//	}); err != nil {
//		log.Errorf(c, err.Error())
//		return
//	}
//	counters.Increment("linked_contacts", 1)
//	log.Infof(c, "Successfully linked contact %v to %v", counterpartyContact.ID, contact.ID)
//	return
//}
//
//func (m *verifyContacts) verifyBalance(c context.Context, counters *asyncCounters, contact models.ContactEntry) (err error) {
//	balance := contact.Data.Balance()
//	if FixBalanceCurrencies(balance) {
//		if err = nds.RunInTransaction(c, func(c context.Context) (err error) {
//			if contact, err = facade.GetContactByID(c, nil, contact.ID); err != nil {
//				return err
//			}
//			if balance := contact.Data.Balance(); FixBalanceCurrencies(balance) {
//				if err = contact.Data.SetBalance(balance); err != nil {
//					return err
//				}
//				if err = facade.SaveContact(c, contact); err != nil {
//					return err
//				}
//				log.Infof(c, "Fixed contact balance currencies: %d", contact.ID)
//			}
//			return nil
//		}, nil); err != nil {
//			return
//		}
//	}
//	return
//}
