package maintainance

//import (
//	"bytes"
//	"fmt"
//	"github.com/crediterra/money"
//	"github.com/dal-go/dalgo/dal"
//	"strings"
//
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
//	"context"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/strongo/logus"
//	"github.com/strongo/nds"
//)
//
//type verifyTransfers struct {
//	transfersAsyncJob
//}
//
//func (m *verifyTransfers) Next(c context.Context, counters mapper.Counters, key *dal.Key) (err error) {
//	return m.startTransferWorker(c, counters, key, m.verifyTransfer)
//}
//
//func (m *verifyTransfers) verifyTransfer(c context.Context, tx dal.ReadwriteTransaction, counters *asyncCounters, transfer models.Transfer) (err error) {
//	buf := new(bytes.Buffer)
//	if err = m.verifyTransferUsers(c, tx, transfer, buf, counters); err != nil {
//		logus.Errorf(c, fmt.Errorf("verifyTransferUsers:transfer=%d: %w", transfer.ID, err).Error())
//		return
//	}
//	if err = m.verifyTransferContacts(c, tx, transfer, buf, counters); err != nil {
//		logus.Errorf(c, fmt.Errorf("verifyTransferContacts:transfer=%d: %w", transfer.ID, err).Error())
//		return
//	}
//	if err = m.verifyTransferCurrency(c, tx, transfer, buf, counters); err != nil {
//		logus.Errorf(c, fmt.Errorf("verifyTransferCurrency:transfer=%d: %w", transfer.ID, err).Error())
//		return
//	}
//	if err = m.verifyReturnsToTransferIDs(c, tx, transfer, buf, counters); err != nil {
//		logus.Errorf(c, fmt.Errorf("verifyReturnsToTransferIDs:transfer=%d: %w", transfer.ID, err).Error())
//		return
//	}
//	if buf.Len() > 0 {
//		logus.Warningf(c, fmt.Errorf("transfer.ID: %v, Created: %v\n", transfer.ID, transfer.Data.DtCreated).Error()+buf.String())
//	}
//	return
//}
//
//func (m *verifyTransfers) verifyTransferUsers(c context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer, buf *bytes.Buffer, counters *asyncCounters) (err error) {
//	for _, userID := range transfer.Data.BothUserIDs {
//		if userID != 0 {
//			if _, err2 := facade2debtus.User.GetUserByID(c, tx, userID); dal.IsNotFound(err2) {
//				counters.Increment(fmt.Sprintf("User:%d", userID), 1)
//				fmt.Fprintf(buf, "Unknown user %d\n", userID)
//			} else if err2 != nil {
//				err = fmt.Errorf("failed to get user by ID=%v: %w", userID, err2)
//				return
//			}
//		}
//	}
//	return
//}
//
//func (m *verifyTransfers) verifyTransferContacts(c context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer, buf *bytes.Buffer, counters *asyncCounters) (err error) {
//	for _, contactID := range transfer.Data.BothCounterpartyIDs {
//		if contactID != 0 {
//			if _, err2 := facade2debtus.GetContactByID(c, tx, contactID); dal.IsNotFound(err2) {
//				counters.Increment(fmt.Sprintf("ContactEntry:%d", contactID), 1)
//				_, _ = fmt.Fprintf(buf, "Unknown contact %d\n", contactID)
//			} else if err2 != nil {
//				err = fmt.Errorf("failed to get contact by ID=%v: %w", contactID, err2)
//				return
//			}
//		}
//	}
//	from := transfer.Data.From()
//	to := transfer.Data.To()
//
//	if from.UserID != 0 && to.UserID != 0 {
//		fixContactID := func(toFix, toUse *models.TransferCounterpartyInfo) (changed bool, err error) {
//			if toFix.ContactID != 0 {
//				panic("toFix.ContactID != 0")
//			}
//			var user models.AppUser
//			if user, err = facade2debtus.User.GetUserByID(c, tx, toUse.UserID); err != nil {
//				return changed, fmt.Errorf("failed to get user by ID: %w", err)
//			}
//			contactIDs := make([]int64, 0, user.Data.ContactsCount)
//			for _, c := range user.Data.Contacts() {
//				contactIDs = append(contactIDs, c.ID)
//			}
//			contacts, err := facade2debtus.GetContactsByIDs(c, nil, contactIDs)
//			if err != nil {
//				return false, fmt.Errorf("failed to get contacts by IDs=%+v: %w", contactIDs, err)
//			}
//			for _, contact := range contacts {
//				if contact.Data.CounterpartyUserID == toFix.UserID {
//					toFix.ContactID = contact.ID
//					changed = true
//					_, _ = fmt.Fprintf(buf, "will assign ContactID=%v, ContactName=%v for UserID=%v, UserName=%v", contact.ID, contact.Data.FullName(), from.UserID, from.UserName)
//					break
//				}
//			}
//			return changed, nil
//		}
//		//var transferChanged, changed bool
//
//		if from.ContactID == 0 {
//			if /*changed*/ _, err = fixContactID(from, to); err != nil {
//				return
//				//} else if changed {
//				//	transferChanged = transferChanged || changed
//			}
//		}
//		if to.ContactID == 0 {
//			if /*changed*/ _, err = fixContactID(to, from); err != nil {
//				return
//				//} else if changed {
//				//	transferChanged = transferChanged || changed
//			}
//		}
//		//changed = changed || transferChanged
//	}
//	return nil
//}
//
//func (*verifyTransfers) verifyTransferCurrency(c context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer, buf *bytes.Buffer, counters *asyncCounters) (err error) {
//	var currency money.CurrencyCode
//	if transfer.Data.Currency == money.CurrencyCode("euro") {
//		currency = money.CurrencyCode("EUR")
//	} else if len(transfer.Data.Currency) == 3 {
//		if v2 := money.CurrencyCode(strings.ToUpper(string(transfer.Data.Currency))); v2 != transfer.Data.Currency && v2.IsMoney() {
//			currency = v2
//		}
//	}
//	if currency != "" {
//		if err = nds.RunInTransaction(c, func(c context.Context) error {
//			if transfer, err = facade2debtus.Transfers.GetTransferByID(c, tx, transfer.ID); err != nil {
//				return fmt.Errorf("failed to get transfer by ID=%d: %w", transfer.ID, err)
//			}
//			if transfer.Data.Currency != currency {
//				transfer.Data.Currency = currency
//				if err = facade2debtus.Transfers.SaveTransfer(c, tx, transfer); err != nil {
//					return fmt.Errorf("failed to save transfer: %w", err)
//				}
//				_, _ = fmt.Fprintf(buf, "Currency fixed: %d\n", transfer.ID)
//			}
//			return nil
//		}, nil); err != nil {
//			return err
//		}
//	}
//	return
//}
//
//func (*verifyTransfers) verifyReturnsToTransferIDs(c context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer, buf *bytes.Buffer, counters *asyncCounters) (err error) {
//	if len(transfer.Data.ReturnToTransferIDs) == 0 {
//		return
//	}
//	var returnToTransfers []models.Transfer
//	if returnToTransfers, err = dtdal.Transfer.GetTransfersByID(c, tx, transfer.Data.ReturnToTransferIDs); err != nil {
//		return fmt.Errorf("failed to get transfers by IDs=%+v: %w", transfer.Data.ReturnToTransferIDs, err)
//	}
//	for _, returnToTransfer := range returnToTransfers {
//		if transfer.Data.From().ContactID != returnToTransfer.Data.To().ContactID {
//			_, _ = fmt.Fprintf(buf, "returnToTransfer(id=%v).To().ContactID != From().ContactID: %v != %v\n", returnToTransfer.ID, returnToTransfer.Data.To().ContactID, transfer.Data.From().ContactID)
//		}
//		if transfer.Data.To().ContactID != returnToTransfer.Data.From().ContactID {
//			_, _ = fmt.Fprintf(buf, "returnToTransfer(id=%v).From().ContactID != To().ContactID: %v != %v\n", returnToTransfer.ID, returnToTransfer.Data.From().ContactID, transfer.Data.To().ContactID)
//		}
//	}
//	return
//}
