package maintainance

import (
	"context"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/datastore"
	"net/http"
	"strings"
	"sync"
)

func mergeContactsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	targetContactID := q.Get("target")
	if targetContactID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("target contact ID is empty"))
		return
	}
	sourceContactIDs := strings.Split(q.Get("source"), ",")
	var db dal.DB
	var err error
	db, err = facade.GetDatabase(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	err = db.RunReadwriteTransaction(r.Context(), func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return mergeContacts(appengine.NewContext(r), tx, targetContactID, sourceContactIDs...)
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte("done"))
}

func mergeContacts(c context.Context, tx dal.ReadwriteTransaction, targetContactID string, sourceContactIDs ...string) (err error) {
	if len(sourceContactIDs) == 0 {
		panic("len(sourceContactIDs) == 0")
	}

	var (
		targetContact models.ContactEntry
		user          models.AppUser
	)

	if targetContact, err = facade2debtus.GetContactByID(c, tx, targetContactID); err != nil {
		if dal.IsNotFound(err) && len(sourceContactIDs) == 1 {
			if targetContact, err = facade2debtus.GetContactByID(c, tx, sourceContactIDs[0]); err != nil {
				return
			}
			targetContact.ID = targetContactID
			if err = facade2debtus.SaveContact(c, targetContact); err != nil {
				return
			}
		} else {
			return
		}
	}

	if user, err = facade2debtus.User.GetUserByID(c, nil, targetContact.Data.UserID); err != nil {
		return
	}

	for _, sourceContactID := range sourceContactIDs {
		if sourceContactID == targetContactID {
			err = fmt.Errorf("sourceContactID == targetContactID: %v", sourceContactID)
			return
		}
		var sourceContact models.ContactEntry
		if sourceContact, err = facade2debtus.GetContactByID(c, tx, sourceContactID); err != nil {
			if dal.IsNotFound(err) {
				continue
			}
			return
		}
		if sourceContact.Data.UserID != targetContact.Data.UserID {
			err = fmt.Errorf("sourceContact.UserID != targetContact.UserID: %v != %v",
				sourceContact.Data.UserID, targetContact.Data.UserID)
			return
		}
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(sourceContactIDs))

	for _, sourceContactID := range sourceContactIDs {
		go func(sourceContactID string) {
			if err2 := mergeContactTransfers(c, tx, wg, targetContactID, sourceContactID); err2 != nil {
				logus.Errorf(c, "failed to merge transfers for contact %v: %v", sourceContactID, err2)
				if err == nil {
					err = err2
				}
			}
		}(sourceContactID)
	}
	wg.Wait()

	if err != nil {
		return
	}

	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if user, err = facade2debtus.User.GetUserByID(c, tx, user.ID); err != nil {
			return
		}
		var contacts []models.UserContactJson
		userContacts := user.Data.Contacts()
		targetContactBalance := targetContact.Data.Balance()
		for _, contact := range userContacts {
			for _, sourceContactID := range sourceContactIDs {
				if contact.ID == sourceContactID {
					for currency, value := range contact.Balance() {
						targetContactBalance.Add(money.NewAmount(currency, value))
					}
					var sourceContact models.ContactEntry
					if sourceContact, err = facade2debtus.GetContactByID(c, tx, sourceContactID); err != nil {
						if !dal.IsNotFound(err) {
							return
						}
					} else {
						targetContact.Data.CountOfTransfers += sourceContact.Data.CountOfTransfers
						if targetContact.Data.LastTransferAt.Before(sourceContact.Data.LastTransferAt) {
							targetContact.Data.LastTransferAt = sourceContact.Data.LastTransferAt
							targetContact.Data.LastTransferID = sourceContact.Data.LastTransferID
						}
						if sourceContact.Data.CounterpartyCounterpartyID != "" {
							var counterpartyContact models.ContactEntry
							if counterpartyContact, err = facade2debtus.GetContactByID(c, tx, sourceContact.Data.CounterpartyCounterpartyID); err != nil {
								if !dal.IsNotFound(err) {
									return
								}
							} else if counterpartyContact.Data.CounterpartyCounterpartyID == sourceContactID {
								counterpartyContact.Data.CounterpartyCounterpartyID = targetContactID
								if err = facade2debtus.SaveContact(c, counterpartyContact); err != nil {
									return
								}
							} else if counterpartyContact.Data.CounterpartyCounterpartyID != "" && counterpartyContact.Data.CounterpartyCounterpartyID != targetContactID {
								err = fmt.Errorf(
									"data integrity issue : counterpartyContact(id=%v).CounterpartyCounterpartyID != sourceContactID: %v != %v",
									counterpartyContact.ID, counterpartyContact.Data.CounterpartyCounterpartyID, sourceContactID)
								return
							}
						}
					}
					if _, err = facade2debtus.DeleteContact(c, sourceContactID); err != nil {
						return
					}
				} else {
					contacts = append(contacts, contact)
				}
			}
		}
		for i := range contacts {
			if contacts[i].ID == targetContactID {
				if err = contacts[i].SetBalance(targetContactBalance); err != nil {
					return
				}
				user.Data.SetContacts(contacts)
				break
			}
		}

		if err = facade2debtus.User.SaveUser(c, tx, user); err != nil {
			return
		}
		return
	}); err != nil {
		return fmt.Errorf("%w: failed to update user entity", err)
	}

	return
}

func mergeContactTransfers(c context.Context, tx dal.ReadwriteTransaction, wg *sync.WaitGroup, targetContactID string, sourceContactID string) (err error) {
	defer func() {
		wg.Done()
	}()
	transfersQ := dal.From(models.TransfersCollection).
		Where(dal.Field("BothCounterpartyIDs").EqualTo(sourceContactID)).
		SelectInto(func() dal.Record {
			return models.NewTransfer("", nil).Record
		})
	transfers, err := tx.QueryReader(c, transfersQ)
	if err != nil {
		return fmt.Errorf("failed to select transfers: %w", err)
	}
	var (
		record   dal.Record
		transfer models.TransferEntry
	)
	for {
		if record, err = transfers.Next(); err != nil {
			if err == datastore.Done {
				err = nil
				break
			}
			logus.Errorf(c, "Failed to get next transfer: %v", err)
		}
		transfer.ID = record.Key().ID.(string)
		switch sourceContactID {
		case transfer.Data.From().ContactID:
			transfer.Data.From().ContactID = targetContactID
		case transfer.Data.To().ContactID:
			transfer.Data.To().ContactID = targetContactID
		}
		switch sourceContactID {
		case transfer.Data.BothCounterpartyIDs[0]:
			transfer.Data.BothCounterpartyIDs[0] = targetContactID
		case transfer.Data.BothCounterpartyIDs[1]:
			transfer.Data.BothCounterpartyIDs[1] = targetContactID
		}
		if err = facade2debtus.Transfers.SaveTransfer(c, tx, transfer); err != nil {
			logus.Errorf(c, "Failed to save transfer #%v: %v", transfer.ID, err)
		}
	}
	return
}
