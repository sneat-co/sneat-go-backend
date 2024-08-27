package maintainance

//import (
//	"fmt"
//	"github.com/dal-go/dalgo/dal"
//	"runtime/debug"
//
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/facade4debtus"
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
//	"context"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/strongo/logus"
//	"time"
//)
//
//type transfersRecreateContacts struct {
//	transfersAsyncJob
//}
//
//func (m *transfersRecreateContacts) Next(ctx context.Context, counters mapper.Counters, key *dal.Key) (err error) {
//	return m.startTransferWorker(ctx, counters, key, m.verifyAndFix)
//}
//
//func (m *transfersRecreateContacts) verifyAndFix(ctx context.Context, tx dal.ReadwriteTransaction, counters *asyncCounters, transfer models.Transfer) (err error) {
//	defer func() {
//		if r := recover(); r != nil {
//			logus.Errorf(ctx, "*transfersRecreateContacts.verifyAndFix() => panic: %v\n\n%v", r, string(debug.Stack()))
//		}
//	}()
//	var fixed bool
//	fixed, err = verifyAndFixMissingTransferContacts(ctx, tx, transfer)
//	if fixed {
//		counters.Increment("fixed", 1)
//	}
//	return
//}
//
//func verifyAndFixMissingTransferContacts(ctx context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer) (fixed bool, err error) {
//	isMissingAndCanBeFixed := func(contactID, contactUserID, counterpartyContactID int64) (bool, error) {
//		if contactID != 0 && contactUserID != 0 && counterpartyContactID != 0 {
//			if _, err := facade4debtus.GetContactByID(ctx, tx, contactID); err != nil {
//				if dal.IsNotFound(err) {
//					if user, err := dal4userus.GetUserByID(c, tx, contactUserID); err != nil {
//						return false, err
//					} else {
//						for _, c := range user.Data.Contacts() {
//							if c.ContactID == contactID {
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
//		if db, err = facade4debtus.GetDatabase(ctx); err != nil {
//			return
//		}
//		err = db.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {
//			logus.Debugf(c, "Recreating contact # %v", contactInfo.ContactID)
//			var counterpartyContact models.DebtusSpaceContactEntry
//			if counterpartyContact, err = facade4debtus.GetContactByID(tctx, tx, counterpartyInfo.ContactID); err != nil {
//				return
//			}
//			var contactUser, counterpartyUser models.AppUser
//
//			if contactUser, err = dal4userus.GetUserByID(tctx, tx, counterpartyInfo.UserID); err != nil {
//				return
//			}
//
//			if counterpartyUser, err = dal4userus.GetUserByID(tctx, tx, contactInfo.UserID); err != nil {
//				return
//			}
//
//			var contactUserContactJson models.UserContactJson
//
//			for _, c := range contactUser.Data.Contacts() {
//				if c.ContactID == contactInfo.ContactID {
//					contactUserContactJson = c
//					break
//				}
//			}
//
//			if contactUserContactJson.ContactID == 0 {
//				logus.Errorf(c, "DebtusSpaceContactEntry %v info not found in user %v contacts json", contactInfo.ContactID, counterpartyInfo.UserID)
//				return
//			}
//
//			if counterpartyContact.Data.CounterpartyContactID == 0 {
//				if counterpartyContact.Data.CounterpartyContactID == 0 {
//					counterpartyContact.Data.CounterpartyContactID = contactInfo.ContactID
//					counterpartyContact.Data.CounterpartyUserID = counterpartyInfo.UserID
//				} else if counterpartyContact.Data.CounterpartyContactID != contactInfo.ContactID {
//					logus.Errorf(c, "counterpartyContact.CounterpartyContactID != contact.ContactID: %v != %v", counterpartyContact.Data.CounterpartyContactID, contactInfo.ContactID)
//					return
//				}
//				if err = facade4debtus.SaveContact(tctx, counterpartyContact); err != nil {
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
//				err = fmt.Errorf("contact(%v).Balance != contactUserContactJson.Balance(): %v != %v", contact.ContactID, contact.Data.Balance(), contactUserContactJson.Balance())
//				return
//			}
//			if err = facade4debtus.SaveContact(c, contact); err != nil {
//				return
//			}
//
//			return
//		})
//		if err != nil {
//			return
//		}
//		fixed = true
//		logus.Warningf(ctx, "Counterparty re-created: %v", contactInfo.ContactID)
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
