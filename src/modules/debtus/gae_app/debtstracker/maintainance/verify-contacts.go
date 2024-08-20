package maintainance

//import (
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/facade4debtus"
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
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
//	//logus.Debugf(c, "*verifyContacts.Next(id: %v)", key.IntID())
//	return m.startContactWorker(c, counters, key, m.processContact)
//}
//
//func (m *verifyContacts) processContact(c context.Context, counters *asyncCounters, contact models.DebtusSpaceContactEntry) (err error) {
//	if _, err = dal4userus.GetUserByID(c, nil, contact.Data.UserID); dal.IsNotFound(err) {
//		counters.Increment("wrong_UserID", 1)
//		logus.Warningf(c, "DebtusSpaceContactEntry %d reference unknown user %d", contact.ContactID, contact.Data.UserID)
//	} else if err != nil {
//		logus.Errorf(c, err.Error())
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
//func (m *verifyContacts) verifyLinking(c context.Context, counters *asyncCounters, contact models.DebtusSpaceContactEntry) (err error) {
//	if contact.Data.CounterpartyContactID != 0 {
//		var counterpartyContact models.DebtusSpaceContactEntry
//		if counterpartyContact, err = facade4debtus.GetContactByID(c, nil, contact.Data.CounterpartyContactID); err != nil {
//			logus.Errorf(c, err.Error())
//			return
//		}
//		if counterpartyContact.Data.CounterpartyContactID == 0 || counterpartyContact.Data.CounterpartyUserID == 0 {
//			if err = m.linkContacts(c, counters, contact); err != nil {
//				return
//			}
//		} else if counterpartyContact.Data.CounterpartyContactID == contact.ContactID && counterpartyContact.Data.CounterpartyUserID == contact.Data.UserID {
//			// Pass, we are OK
//		} else {
//			logus.Warningf(c, "Wrongly linked contacts: %v=>%v != %v=>%v",
//				contact.ContactID, contact.Data.CounterpartyContactID,
//				counterpartyContact.ContactID, counterpartyContact.Data.CounterpartyContactID)
//		}
//	}
//	return
//}
//
//func (m *verifyContacts) linkContacts(c context.Context, counters *asyncCounters, contact models.DebtusSpaceContactEntry) (err error) {
//	var counterpartyContact models.DebtusSpaceContactEntry
//	var db dal.DB
//	if db, err = facade4debtus.GetDatabase(c); err != nil {
//		return
//	}
//	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if counterpartyContact, err = facade4debtus.GetContactByID(c, tx, contact.Data.CounterpartyContactID); err != nil {
//			logus.Errorf(c, err.Error())
//			return
//		}
//		if counterpartyContact.Data.CounterpartyContactID == 0 {
//			counterpartyContact.Data.CounterpartyContactID = contact.ContactID
//			if counterpartyContact.Data.CounterpartyUserID == 0 {
//				counterpartyContact.Data.CounterpartyUserID = contact.Data.UserID
//			} else if counterpartyContact.Data.CounterpartyUserID != contact.Data.UserID {
//				err = fmt.Errorf("counterpartyContact(id=%v).CounterpartyUserID != contact(id=%v).UserID: %v != %v",
//					counterpartyContact.ContactID, contact.ContactID, counterpartyContact.Data.CounterpartyUserID, contact.Data.UserID)
//				return
//			}
//			if err = facade4debtus.SaveContact(c, counterpartyContact); err != nil {
//				return
//			}
//		} else if counterpartyContact.Data.CounterpartyContactID != contact.ContactID {
//			logus.Warningf(c, "in tx: wrongly linked contacts: %v=>%v != %v=>%v",
//				contact.ContactID, contact.Data.CounterpartyContactID,
//				counterpartyContact.ContactID, counterpartyContact.Data.CounterpartyContactID)
//		}
//		return
//	}); err != nil {
//		logus.Errorf(c, err.Error())
//		return
//	}
//	counters.Increment("linked_contacts", 1)
//	logus.Infof(c, "Successfully linked contact %v to %v", counterpartyContact.ContactID, contact.ContactID)
//	return
//}
//
//func (m *verifyContacts) verifyBalance(c context.Context, counters *asyncCounters, contact models.DebtusSpaceContactEntry) (err error) {
//	balance := contact.Data.Balance()
//	if FixBalanceCurrencies(balance) {
//		if err = nds.RunInTransaction(c, func(c context.Context) (err error) {
//			if contact, err = facade4debtus.GetContactByID(c, nil, contact.ContactID); err != nil {
//				return err
//			}
//			if balance := contact.Data.Balance(); FixBalanceCurrencies(balance) {
//				if err = contact.Data.SetBalance(balance); err != nil {
//					return err
//				}
//				if err = facade4debtus.SaveContact(c, contact); err != nil {
//					return err
//				}
//				logus.Infof(c, "Fixed contact balance currencies: %d", contact.ContactID)
//			}
//			return nil
//		}, nil); err != nil {
//			return
//		}
//	}
//	return
//}
