package facade

import (
	"context"
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sanity-io/litter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
	"reflect"
	"strconv"
)

func ChangeContactStatus(c context.Context, contactID string, newStatus string) (contact models.Contact, err error) {
	err = RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if contact, err = GetContactByID(c, tx, contactID); err != nil {
			return err
		}
		if contact.Data.Status != newStatus {
			contact.Data.Status = newStatus
			var user models.AppUser
			if user, err = User.GetUserByID(c, tx, contact.Data.UserID); err != nil {
				return err
			}
			if _, userChanged := models.AddOrUpdateContact(&user, contact); userChanged {
				err = tx.SetMulti(c, []dal.Record{contact.Record, user.Record})
			} else {
				err = tx.Set(c, contact.Record)
			}
		}
		return err
	})
	return
}

func createContactWithinTransaction(
	tc context.Context,
	tx dal.ReadwriteTransaction,
	changes *createContactDbChanges,
	counterpartyUserID string,
	contactDetails models.ContactDetails,
) (
	contact models.Contact,
	counterpartyContact models.Contact,
	err error,
) {
	if tx == nil {
		err = errors.New("tx == nil")
		return
	}
	appUser := *changes.user
	if changes.counterpartyContact != nil {
		counterpartyContact = *changes.counterpartyContact
	}

	log.Debugf(tc, "createContactWithinTransaction(appUser.ID=%v, counterpartyDetails=%v)", appUser.ID, contactDetails)
	if appUser.ID == "" {
		err = errors.New("appUser.ID == 0")
		return
	}
	if appUser.Data == nil {
		err = errors.New("appUser.DebutsAppUserDataOBSOLETE == nil")
		return
	}
	if appUser.ID == counterpartyUserID {
		panic(fmt.Sprintf("appUser.ID == counterpartyUserID: %v", counterpartyUserID))
	}
	if counterpartyContact.Data != nil && counterpartyContact.ID == "" {
		panic(fmt.Sprintf("counterpartyContact.DebtusContactData != nil && counterpartyContact.ID == 0: %v", litter.Sdump(counterpartyContact)))
	}

	contact.Data = models.NewDebtusContactData(appUser.ID, contactDetails)
	if counterpartyContact.ID != "" {
		if counterpartyContact.Data == nil {
			if counterpartyContact, err = GetContactByID(tc, tx, counterpartyContact.ID); err != nil {
				return
			}
			changes.counterpartyContact = &counterpartyContact
		}
		if counterpartyContact.Data.UserID != counterpartyUserID {
			if counterpartyUserID == "" {
				counterpartyUserID = counterpartyContact.Data.UserID
			} else {
				panic(fmt.Sprintf("counterpartyContact.UserID != counterpartyUserID: %v != %v", counterpartyContact.Data.UserID, counterpartyUserID))
			}
		}
		contact.Data.CounterpartyUserID = counterpartyUserID
		contact.Data.CounterpartyCounterpartyID = counterpartyContact.ID
		contact.Data.TransfersJson = counterpartyContact.Data.TransfersJson
		contact.Data.Balanced = money.Balanced{
			CountOfTransfers: counterpartyContact.Data.CountOfTransfers,
			LastTransferID:   counterpartyContact.Data.LastTransferID,
			LastTransferAt:   counterpartyContact.Data.LastTransferAt,
		}
		invitedCounterpartyBalance := counterpartyContact.Data.Balance().Reversed()
		log.Debugf(tc, "invitedCounterpartyBalance: %v", invitedCounterpartyBalance)
		if err = contact.Data.SetBalance(invitedCounterpartyBalance); err != nil {
			return
		}
		contact.Data.MustMatchCounterparty(counterpartyContact)
	}

	if contact, err = dtdal.Contact.InsertContact(tc, tx, contact.Data); err != nil {
		return
	}

	if counterpartyContact.ID != "" {
		if counterpartyContact.Data.CounterpartyCounterpartyID == "" {
			counterpartyContact.Data.CounterpartyCounterpartyID = contact.ID
			if counterpartyContact.Data.CounterpartyUserID == "" {
				counterpartyContact.Data.CounterpartyUserID = contact.Data.UserID
			} else {
				err = fmt.Errorf("inviter contact %v already has CounterpartyUserID=%v", counterpartyContact.ID, counterpartyContact.Data.CounterpartyUserID)
				return
			}
			changes.FlagAsChanged(changes.counterpartyContact.Record)
		} else if counterpartyContact.Data.CounterpartyCounterpartyID != contact.ID {
			err = fmt.Errorf("inviter contact %v already has CounterpartyCounterpartyID=%v", counterpartyContact.ID, counterpartyContact.Data.CounterpartyCounterpartyID)
			return
		}
	}

	if _, changed := models.AddOrUpdateContact(&appUser, contact); changed {
		changes.FlagAsChanged(changes.user.Record)
	}

	{ // Verifications for data integrity
		if counterpartyContact.Data != nil {
			contact.Data.MustMatchCounterparty(counterpartyContact)
		}
		if contact.Data.UserID != appUser.ID {
			panic(fmt.Sprintf("contact.UserID != appUser.ID: %v != %v", contact.Data.UserID, appUser.ID))
		}
		if counterpartyContact.Data != nil {
			if counterpartyContact.Data.UserID != counterpartyUserID {
				panic(fmt.Sprintf("counterpartyContact.UserID != counterpartyUserID: %v != %v", counterpartyContact.Data.UserID, counterpartyUserID))
			}
			if contact.ID == counterpartyContact.ID {
				panic(fmt.Sprintf("contact.ID == counterpartyContact.ID: %v", contact.ID))
			}
			if contact.Data.UserID == counterpartyContact.Data.UserID {
				panic(fmt.Sprintf("contact.UserID == counterpartyContact.UserID: %v", contact.Data.UserID))
			}
			if contact.Data.TransfersJson != counterpartyContact.Data.TransfersJson {
				log.Errorf(tc, "contact.TransfersJson != counterpartyContact.TransfersJson\n contact: %v\n counterpartyContact: %v", contact.Data.TransfersJson, counterpartyContact.Data.TransfersJson)
			}
			if contact.Data.BalanceCount != counterpartyContact.Data.BalanceCount {
				panic(fmt.Sprintf("contact.BalanceCount != counterpartyContact.BalanceCount: %v != %v", contact.Data.BalanceCount, counterpartyContact.Data.BalanceCount))
			}
			if cBalance, cpBalance := contact.Data.Balance(), counterpartyContact.Data.Balance(); !cBalance.Equal(cpBalance.Reversed()) {
				panic(fmt.Sprintf("!contact.Balance().Equal(counterpartyContact.Balance())\ncontact.Balance(): %v\n counterpartyContact.Balance(): %v", cBalance, cpBalance))
			}
		}
		appUserContactJson := appUser.Data.ContactByID(contact.ID)
		if ucBalance, cBalance := appUserContactJson.Balance(), contact.Data.Balance(); !ucBalance.Equal(cBalance) {
			panic(fmt.Sprintf("appUserContactJson.Balance().Equal(contact.Balance())\nappUser.ContactByID(contact.ID).Balance(): %v\ncontact.Balance(): %v", ucBalance, cBalance))
		}
	}
	return
}

type createContactDbChanges struct {
	dal.Changes
	user                *models.AppUser
	counterpartyContact *models.Contact
}

func CreateContact(c context.Context, tx dal.ReadwriteTransaction, userID string, contactDetails models.ContactDetails) (contact models.Contact, user models.AppUser, err error) {
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	var contactIDs []string
	if contactIDs, err = dtdal.Contact.GetContactIDsByTitle(c, tx, userID, contactDetails.Username, false); err != nil {
		return
	}
	switch len(contactIDs) {
	case 0:
		err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
			if user, err = User.GetUserByID(tc, tx, userID); err != nil {
				return
			}
			changes := &createContactDbChanges{
				user:                &user,
				counterpartyContact: new(models.Contact),
			}
			if contact, _, err = createContactWithinTransaction(tc, tx, changes, "", contactDetails); err != nil {
				err = fmt.Errorf("failed to create contact within transaction: %w", err)
				return
			}

			if changes.HasChanges() {
				//db, err := dtdal.DB.GetDB(tc)
				if err = tx.SetMulti(tc, changes.Records()); err != nil {
					err = fmt.Errorf("failed to save entity related to new contact: %w", err)
					return
				}
				// TODO: move calls of delays to createContactWithinTransaction() ?
				if err = dtdal.User.DelayUpdateUserWithContact(tc, userID, contact.ID); err != nil { // Just in case
					return
				}
				if changes.counterpartyContact != nil && changes.counterpartyContact.ID > "" {
					counterpartyContact := *changes.counterpartyContact
					if err = dtdal.User.DelayUpdateUserWithContact(tc, counterpartyContact.Data.UserID, counterpartyContact.ID); err != nil { // Just in case
						return
					}
				}
			}
			return
		}, dal.TxWithCrossGroup())
		if err != nil {
			if err = dtdal.User.DelayUpdateUserWithContact(c, contact.Data.UserID, contact.ID); err != nil {
				return
			}
			return
		}
		return
	case 1:
		if contact, err = GetContactByID(c, tx, contactIDs[0]); err != nil {
			return
		}
		user.ID = userID
		return
	default:
		err = fmt.Errorf("too many counterparties (%d), IDs: %v", len(contactIDs), contactIDs)
		return
	}
}

func UpdateContact(c context.Context, contactID string, values map[string]string) (contactEntity *models.DebtusContactData, err error) {
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if contact, err := GetContactByID(c, tx, contactID); err != nil {
			return err
		} else {
			var changed bool
			for name, value := range values {
				switch name {
				case "Username":
					if contact.Data.Username != value {
						contact.Data.Username = value
						changed = true
					}
				case "FirstName":
					if contact.Data.FirstName != value {
						contact.Data.FirstName = value
						changed = true
					}
				case "LastName":
					if contact.Data.LastName != value {
						contact.Data.LastName = value
						changed = true
					}
				case "ScreenName":
					if contact.Data.ScreenName != value {
						contact.Data.ScreenName = value
						changed = true
					}
				case "EmailAddress":
					if contact.Data.EmailAddressOriginal != value {
						contact.Data.EmailAddressOriginal = value
						changed = true
					}
				case "PhoneNumber":
					if phoneNumber, err := strconv.ParseInt(value, 10, 64); err != nil {
						return err
					} else if contact.Data.PhoneNumber != phoneNumber {
						contact.Data.PhoneNumber = phoneNumber
						contact.Data.PhoneNumberConfirmed = false
						changed = true
					}
				default:
					log.Debugf(c, "Unknown field: %v", name)
				}
			}
			if changed {
				if user, err := User.GetUserByID(c, tx, contact.Data.UserID); err != nil {
					return fmt.Errorf("failed to get user by ID=%v: %w", contact.Data.UserID, err)
				} else {
					models.AddOrUpdateContact(&user, contact)
					return tx.SetMulti(c, []dal.Record{contact.Record, user.Record})
				}
			}
		}
		return nil
	}, dal.TxWithCrossGroup())
	return
}

var ErrContactIsNotDeletable = errors.New("contact is not deletable")

func DeleteContact(c context.Context, contactID string) (user models.AppUser, err error) {
	log.Warningf(c, "ContactDalGae.DeleteContact(%d)", contactID)
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	var contact models.Contact
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if contact, err = GetContactByID(c, tx, contactID); err != nil {
			if dal.IsNotFound(err) {
				log.Warningf(c, "Contact not found by ID: %v", contactID)
				err = nil
			}
			return
		}
		if contact.Data != nil && contact.Data.CounterpartyUserID != "" {
			return ErrContactIsNotDeletable
		}
		if user, err = User.GetUserByID(c, tx, contact.Data.UserID); err != nil {
			return
		}
		if userContact := user.Data.ContactByID(contactID); userContact != nil {
			userContactBalance := userContact.Balance()
			contactBalance := contact.Data.Balance()
			if !reflect.DeepEqual(userContactBalance, contactBalance) {
				return fmt.Errorf("Data integrity issue: userContactBalance != contactBalance\n\tuserContactBalance: %v\n\tcontactBalance: %v", userContactBalance, contactBalance)
			}
			if !user.Data.RemoveContact(contactID) {
				return errors.New("implementation error: user not changed on removing contact")
			}
			if contact.Data.BalanceCount > 0 {
				userBalance := user.Data.Balance()
				for k, v := range contactBalance {
					userBalance[k] -= v
				}
				if err = user.Data.SetBalance(userBalance); err != nil {
					return err
				}
			}
			if err = User.SaveUser(c, tx, user); err != nil {
				return err
			}
		}
		key := models.NewDebtusContactKey(contactID)
		if err = tx.Delete(c, key); err != nil {
			return err
		}
		return nil
	}, dal.TxWithCrossGroup())
	return
}

func SaveContact(c context.Context, contact models.Contact) error {
	db, err := GetDatabase(c)
	if err != nil {
		return err
	}
	return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(c, contact.Record)
	})
}

func GetContactsByIDs(c context.Context, tx dal.ReadSession, contactsIDs []string) (contacts []models.Contact, err error) {
	if tx == nil {
		if tx, err = GetDatabase(c); err != nil {
			return
		}
	}
	contacts = models.NewDebtusContacts(contactsIDs...)
	records := models.DebtusContactRecords(contacts)
	return contacts, tx.GetMulti(c, records)
}

func GetContactByID(c context.Context, tx dal.ReadSession, contactID string) (contact models.Contact, err error) {
	contact = models.NewDebtusContact(contactID, nil)
	if tx == nil {
		tx, err = GetDatabase(c)
	}
	if err != nil {
		return contact, err
	}
	return contact, tx.Get(c, contact.Record)
}
