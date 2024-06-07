package maintainance

//import (
//	"fmt"
//	"github.com/dal-go/dalgo/dal"
//	"runtime/debug"
//
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
//	"context"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/strongo/log"
//	"time"
//)
//
//type transfersRecreateContacts struct {
//	transfersAsyncJob
//}
//
//func (m *transfersRecreateContacts) Next(c context.Context, counters mapper.Counters, key *dal.Key) (err error) {
//	return m.startTransferWorker(c, counters, key, m.verifyAndFix)
//}
//
//func (m *transfersRecreateContacts) verifyAndFix(c context.Context, tx dal.ReadwriteTransaction, counters *asyncCounters, transfer models.Transfer) (err error) {
//	defer func() {
//		if r := recover(); r != nil {
//			log.Errorf(c, "*transfersRecreateContacts.verifyAndFix() => panic: %v\n\n%v", r, string(debug.Stack()))
//		}
//	}()
//	var fixed bool
//	fixed, err = verifyAndFixMissingTransferContacts(c, tx, transfer)
//	if fixed {
//		counters.Increment("fixed", 1)
//	}
//	return
//}
//
//func verifyAndFixMissingTransferContacts(c context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer) (fixed bool, err error) {
//	isMissingAndCanBeFixed := func(contactID, contactUserID, counterpartyContactID int64) (bool, error) {
//		if contactID != 0 && contactUserID != 0 && counterpartyContactID != 0 {
//			if _, err := facade.GetContactByID(c, tx, contactID); err != nil {
//				if dal.IsNotFound(err) {
//					if user, err := facade.User.GetUserByID(c, tx, contactUserID); err != nil {
//						return false, err
//					} else {
//						for _, c := range user.Data.Contacts() {
//							if c.ID == contactID {
//								return true, nil
//							}
//						}
//						return false, nil
//					}
//				}
//				return false, err
//			}
//		}
//		return false, nil
//	}
//
//	doFix := func(contactInfo *models.TransferCounterpartyInfo, counterpartyInfo *models.TransferCounterpartyInfo) (err error) {
//		var db dal.DB
//		if db, err = facade.GetDatabase(c); err != nil {
//			return
//		}
//		err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
//			log.Debugf(c, "Recreating contact # %v", contactInfo.ContactID)
//			var counterpartyContact models.ContactEntry
//			if counterpartyContact, err = facade.GetContactByID(c, tx, counterpartyInfo.ContactID); err != nil {
//				return
//			}
//			var contactUser, counterpartyUser models.AppUser
//
//			if contactUser, err = facade.User.GetUserByID(c, tx, counterpartyInfo.UserID); err != nil {
//				return
//			}
//
//			if counterpartyUser, err = facade.User.GetUserByID(c, tx, contactInfo.UserID); err != nil {
//				return
//			}
//
//			var contactUserContactJson models.UserContactJson
//
//			for _, c := range contactUser.Data.Contacts() {
//				if c.ID == contactInfo.ContactID {
//					contactUserContactJson = c
//					break
//				}
//			}
//
//			if contactUserContactJson.ID == 0 {
//				log.Errorf(c, "ContactEntry %v info not found in user %v contacts json", contactInfo.ContactID, counterpartyInfo.UserID)
//				return
//			}
//
//			if counterpartyContact.Data.CounterpartyCounterpartyID == 0 {
//				if counterpartyContact.Data.CounterpartyCounterpartyID == 0 {
//					counterpartyContact.Data.CounterpartyCounterpartyID = contactInfo.ContactID
//					counterpartyContact.Data.CounterpartyUserID = counterpartyInfo.UserID
//				} else if counterpartyContact.Data.CounterpartyCounterpartyID != contactInfo.ContactID {
//					log.Errorf(c, "counterpartyContact.CounterpartyCounterpartyID != contact.ID: %v != %v", counterpartyContact.Data.CounterpartyCounterpartyID, contactInfo.ContactID)
//					return
//				}
//				if err = facade.SaveContact(c, counterpartyContact); err != nil {
//					return err
//				}
//			}
//
//			contact := models.NewContact(contactInfo.ContactID, &models.ContactData{
//				UserID:         counterpartyInfo.UserID,
//				DtCreated:      time.Now(),
//				Status:         models.STATUS_ACTIVE,
//				TransfersJson:  counterpartyContact.Data.TransfersJson,
//				ContactDetails: counterpartyUser.Data.ContactDetails,
//				Balanced:       counterpartyContact.Data.Balanced,
//			})
//			if contact.Data.Nickname != contactUserContactJson.Name &&
//				contact.Data.FirstName != contactUserContactJson.Name &&
//				contact.Data.LastName != contactUserContactJson.Name &&
//				contact.Data.ScreenName != contactUserContactJson.Name {
//				contact.Data.Nickname = contactUserContactJson.Name
//			}
//			if err = contact.Data.SetBalance(counterpartyContact.Data.Balance().Reversed()); err != nil {
//				return
//			}
//			if !contact.Data.Balance().Equal(contactUserContactJson.Balance()) {
//				err = fmt.Errorf("contact(%v).Balance != contactUserContactJson.Balance(): %v != %v", contact.ID, contact.Data.Balance(), contactUserContactJson.Balance())
//				return
//			}
//			if err = facade.SaveContact(c, contact); err != nil {
//				return
//			}
//
//			return
//		})
//		if err != nil {
//			return
//		}
//		fixed = true
//		log.Warningf(c, "Counterparty re-created: %v", contactInfo.ContactID)
//		return
//	}
//
//	verifyAndFix := func(contactInfo *models.TransferCounterpartyInfo, counterpartyInfo *models.TransferCounterpartyInfo) error {
//		if toBeFixed, err := isMissingAndCanBeFixed(contactInfo.ContactID, counterpartyInfo.UserID, counterpartyInfo.ContactID); err != nil {
//			return err
//		} else if toBeFixed {
//			return doFix(contactInfo, counterpartyInfo)
//		}
//		return nil
//	}
//
//	from, to := transfer.Data.From(), transfer.Data.To()
//
//	if err = verifyAndFix(from, to); err != nil {
//		return
//	}
//
//	if err = verifyAndFix(to, from); err != nil {
//		return
//	}
//	return
//}
