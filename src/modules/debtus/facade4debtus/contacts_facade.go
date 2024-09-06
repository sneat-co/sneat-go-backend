package facade4debtus

import (
	"context"
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sanity-io/litter"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/delays4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"reflect"
	"strconv"
)

func ChangeContactStatus(
	ctx context.Context, userCtx facade.UserContext, spaceID, contactID string, newStatus models4debtus.DebtusContactStatus,
) (
	contact dal4contactus.ContactEntry,
	debtusContact models4debtus.DebtusSpaceContactEntry,
	err error,
) {

	spaceRequest := dto4spaceus.SpaceRequest{
		SpaceID: spaceID,
	}
	err = dal4contactus.RunContactusSpaceWorker(ctx, userCtx, spaceRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {
			if debtusContact, err = GetDebtusSpaceContactByID(ctx, tx, spaceID, contactID); err != nil {
				return err
			}
			if debtusContact.Data.Status != newStatus {
				debtusContact.Data.Status = newStatus
				debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)
				if err = tx.Get(ctx, debtusSpace.Record); err != nil {
					return err
				}
				if _, userChanged := models4debtus.AddOrUpdateDebtusContact(debtusSpace, debtusContact); userChanged {
					err = tx.SetMulti(ctx, []dal.Record{debtusContact.Record, debtusSpace.Record})
				} else {
					err = tx.Set(ctx, debtusContact.Record)
				}
			}
			return err
		})
	return
}

func createContactWithinTransaction(
	tctx context.Context,
	tx dal.ReadwriteTransaction,
	changes *createContactDbChanges,
	spaceID string,
	counterpartyUserID string,
	contactDetails dto4contactus.ContactDetails,
) (
	err error,
) {
	creator := changes.creator
	counterparty := changes.counterparty

	if tx == nil {
		err = errors.New("tx == nil")
		return
	}
	appUser := changes.user

	logus.Debugf(tctx, "createContactWithinTransaction(appUser.ContactID=%v, counterpartyDetails=%v)", appUser.ID, contactDetails)
	if appUser.ID == "" {
		err = errors.New("appUser.ContactID == 0")
		return
	}
	if appUser.Data == nil {
		err = errors.New("appUser.DebutsAppUserDataOBSOLETE == nil")
		return
	}
	if appUser.ID == counterpartyUserID {
		panic(fmt.Sprintf("appUser.ContactID == counterpartyUserID: %v", counterpartyUserID))
	}
	if counterparty.Contact.Data != nil && counterparty.Contact.ID == "" {
		panic(fmt.Sprintf("counterpartyContact.DebtusSpaceContactDbo != nil && counterpartyContact.ContactID == 0: %v", litter.Sdump(counterparty.Contact)))
	}

	creator.DebtusContact.Data = models4debtus.NewDebtusContactDbo(contactDetails)
	creator.DebtusContact.Data.CreatedBy = appUser.ID
	if counterparty.Contact.ID != "" {
		if counterparty.Contact.Data == nil {
			if counterparty.DebtusContact, err = GetDebtusSpaceContactByID(tctx, tx, spaceID, counterparty.Contact.ID); err != nil {
				return
			}
			changes.counterparty.Contact = counterparty.Contact
		}
		if counterparty.Contact.Data.UserID != counterpartyUserID {
			if counterpartyUserID == "" {
				counterpartyUserID = counterparty.Contact.Data.UserID
			} else {
				panic(fmt.Sprintf("counterpartyContact.UserID != counterpartyUserID: %v != %v", counterparty.Contact.Data.UserID, counterpartyUserID))
			}
		}
		creator.DebtusContact.Data.CounterpartyUserID = counterpartyUserID
		creator.DebtusContact.Data.CounterpartyContactID = counterparty.Contact.ID
		creator.DebtusContact.Data.Transfers = counterparty.DebtusContact.Data.Transfers
		creator.DebtusContact.Data.Balanced = money.Balanced{
			CountOfTransfers: counterparty.DebtusContact.Data.CountOfTransfers,
			LastTransferID:   counterparty.DebtusContact.Data.LastTransferID,
			LastTransferAt:   counterparty.DebtusContact.Data.LastTransferAt,
		}
		invitedCounterpartyBalance := counterparty.DebtusContact.Data.Balance.Reversed()
		logus.Debugf(tctx, "invitedCounterpartyBalance: %v", invitedCounterpartyBalance)
		creator.DebtusContact.Data.Balance = invitedCounterpartyBalance
		creator.DebtusContact.Data.MustMatchCounterparty(counterparty.DebtusContact)
	}

	if creator.DebtusContact, err = dtdal.Contact.InsertContact(tctx, tx, creator.DebtusContact.Data); err != nil {
		return
	}

	if counterparty.Contact.ID != "" {
		if counterparty.DebtusContact.Data.CounterpartyContactID == "" {
			counterparty.DebtusContact.Data.CounterpartyContactID = creator.DebtusContact.ID
			if counterparty.Contact.Data.UserID == "" {
				counterparty.Contact.Data.UserID = creator.Contact.Data.UserID
			} else {
				err = fmt.Errorf("inviter DebtusContact %v already has CounterpartyUserID=%v", counterparty.Contact.ID, counterparty.Contact.Data.UserID)
				return
			}
			changes.FlagAsChanged(changes.counterparty.Contact.Record)
		} else if counterparty.DebtusContact.Data.CounterpartyContactID != creator.DebtusContact.ID {
			err = fmt.Errorf("inviter DebtusContact %v already has CounterpartyContactID=%v", counterparty.Contact.ID, counterparty.Contact.Data.UserID)
			return
		}
	}

	if _, changed := models4debtus.AddOrUpdateDebtusContact(changes.debtusSpace, creator.DebtusContact); changed {
		changes.FlagAsChanged(changes.user.Record)
	}

	{ // Verifications for data integrity
		if counterparty.DebtusContact.Data != nil {
			creator.DebtusContact.Data.MustMatchCounterparty(counterparty.DebtusContact)
		}
		if creator.Contact.Data.UserID != appUser.ID {
			panic(fmt.Sprintf("DebtusContact.UserID != appUser.ContactID: %v != %v", creator.Contact.Data.UserID, appUser.ID))
		}
		if counterparty.Contact.Data != nil {
			if counterparty.Contact.Data.UserID != counterpartyUserID {
				panic(fmt.Sprintf("counterpartyContact.UserID != counterpartyUserID: %v != %v", counterparty.Contact.Data.UserID, counterpartyUserID))
			}
			if creator.DebtusContact.ID == counterparty.Contact.ID {
				panic(fmt.Sprintf("DebtusContact.ContactID == counterpartyContact.ContactID: %v", creator.DebtusContact.ID))
			}
			if creator.Contact.Data.UserID == counterparty.Contact.Data.UserID {
				panic(fmt.Sprintf("DebtusContact.UserID == counterpartyContact.UserID: %v", creator.Contact.Data.UserID))
			}
			if creator.DebtusContact.Data.Transfers != counterparty.DebtusContact.Data.Transfers {
				logus.Errorf(tctx, "DebtusContact.TransfersJson != counterpartyContact.TransfersJson\n DebtusContact: %v\n counterpartyContact: %v", creator.DebtusContact.Data.Transfers, counterparty.DebtusContact.Data.Transfers)
			}
			if cBalance, cpBalance := creator.DebtusContact.Data.Balance, counterparty.DebtusContact.Data.Balance; !cBalance.Equal(cpBalance.Reversed()) {
				panic(fmt.Sprintf("!DebtusContact.Balance().Equal(counterpartyContact.Balance())\nDebtusContact.Balance(): %v\n counterpartyContact.Balance(): %v", cBalance, cpBalance))
			}
		}
		appUserContactJson := changes.debtusSpace.Data.Contacts[creator.DebtusContact.ID]
		if ucBalance, cBalance := appUserContactJson.Balance, creator.DebtusContact.Data.Balance; !ucBalance.Equal(cBalance) {
			panic(fmt.Sprintf("appUserContactJson.Balance().Equal(DebtusContact.Balance())\nappUser.ContactByID(DebtusContact.ContactID).Balance(): %v\nDebtusContact.Balance(): %v", ucBalance, cBalance))
		}
	}
	return
}

type createContactDbChanges struct {
	dal.Changes
	user           dbo4userus.UserEntry
	contactusSpace dal4contactus.ContactusSpaceEntry
	debtusSpace    models4debtus.DebtusSpaceEntry
	creator        ParticipantEntries
	counterparty   ParticipantEntries
}

func CreateContact(
	ctx context.Context, tx dal.ReadwriteTransaction, userID, spaceID string, contactDetails dto4contactus.ContactDetails,
) (
	contact dal4contactus.ContactEntry,
	contactusSpace dal4contactus.ContactusSpaceEntry,
	debtusContact models4debtus.DebtusSpaceContactEntry,
	err error,
) {
	var contactIDs []string
	if contactIDs, err = dtdal.Contact.GetContactIDsByTitle(ctx, tx, spaceID, userID, contactDetails.UserName, false); err != nil {
		return
	}
	userCtx := facade.NewUserContext(userID)
	spaceRequest := dto4spaceus.SpaceRequest{
		SpaceID: spaceID,
	}

	switch len(contactIDs) {
	case 0:
		err = dal4contactus.RunContactusSpaceWorker(ctx, userCtx, spaceRequest, func(tctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) (err error) {
			changes := &createContactDbChanges{
				//user:                user,
				debtusSpace:    models4debtus.NewDebtusSpaceEntry(spaceID),
				contactusSpace: dal4contactus.NewContactusSpaceEntry(spaceID),
				counterparty:   ParticipantEntries{},
			}
			if err = createContactWithinTransaction(tctx, tx, changes, spaceRequest.SpaceID, "", contactDetails); err != nil {
				err = fmt.Errorf("failed to create Contact within transaction: %w", err)
				return
			}

			if changes.HasChanges() {
				//db, err := dtdal.DB.GetDB(tctx)
				if err = tx.SetMulti(tctx, changes.Records()); err != nil {
					err = fmt.Errorf("failed to save entity related to new Contact: %w", err)
					return
				}
				// TODO: move calls of delays to createContactWithinTransaction() ?
				if err = delays4contactus.DelayUpdateContactusSpaceDboWithContact(tctx, 0, userID, contact.ID); err != nil { // Just in case
					return
				}
				if changes.counterparty.Contact.ID != "" {
					if err = delays4contactus.DelayUpdateContactusSpaceDboWithContact(tctx, 0, changes.counterparty.Contact.Data.UserID, changes.counterparty.Contact.ID); err != nil { // Just in case
						return
					}
				}
			}
			return
		})
		if err != nil {
			if err = delays4contactus.DelayUpdateContactusSpaceDboWithContact(ctx, 0, contact.Data.UserID, contact.ID); err != nil {
				return
			}
			return
		}
		return
	case 1:
		if debtusContact, err = GetDebtusSpaceContactByID(ctx, tx, spaceID, contactIDs[0]); err != nil {
			return
		}
		if err = tx.Get(ctx, contact.Record); err != nil {
			return
		}
		return
	default:
		err = fmt.Errorf("too many counterparties (%d), IDs: %v", len(contactIDs), contactIDs)
		return
	}
}

func UpdateContact(ctx context.Context, spaceID, contactID string, values map[string]string) (debtusSpaceContact models4debtus.DebtusSpaceContactEntry, err error) {
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if debtusSpaceContact, err = GetDebtusSpaceContactByID(ctx, tx, spaceID, contactID); err != nil {
			return err
		} else {
			var changed bool
			for name, value := range values {
				switch name {
				case "Username":
					if debtusSpaceContact.Data.UserName != value {
						debtusSpaceContact.Data.UserName = value
						changed = true
					}
				case "FirstName":
					if debtusSpaceContact.Data.FirstName != value {
						debtusSpaceContact.Data.FirstName = value
						changed = true
					}
				case "LastName":
					if debtusSpaceContact.Data.LastName != value {
						debtusSpaceContact.Data.LastName = value
						changed = true
					}
				case "ScreenName":
					if debtusSpaceContact.Data.ScreenName != value {
						debtusSpaceContact.Data.ScreenName = value
						changed = true
					}
				case "EmailAddress":
					if debtusSpaceContact.Data.EmailAddressOriginal != value {
						debtusSpaceContact.Data.EmailAddressOriginal = value
						changed = true
					}
				case "PhoneNumber":
					if phoneNumber, err := strconv.ParseInt(value, 10, 64); err != nil {
						return err
					} else if debtusSpaceContact.Data.PhoneNumber != phoneNumber {
						debtusSpaceContact.Data.PhoneNumber = phoneNumber
						debtusSpaceContact.Data.PhoneNumberConfirmed = false
						changed = true
					}
				default:
					logus.Debugf(ctx, "Unknown field: %v", name)
				}
			}
			if changed {
				debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)
				if err = tx.Get(ctx, debtusSpace.Record); err != nil {
					return err
				} else {
					models4debtus.AddOrUpdateDebtusContact(debtusSpace, debtusSpaceContact)
					return tx.SetMulti(ctx, []dal.Record{debtusSpaceContact.Record, debtusSpace.Record})
				}
			}
		}
		return nil
	}, dal.TxWithCrossGroup())
	return
}

var ErrContactIsNotDeletable = errors.New("Contact is not deletable")

func DeleteContact(ctx context.Context, userCtx facade.UserContext, spaceID, contactID string) (err error) {
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return DeleteContactTx(ctx, userCtx, tx, spaceID, contactID)
	})
}

func DeleteContactTx(ctx context.Context, userCtx facade.UserContext, tx dal.ReadwriteTransaction, spaceID, contactID string) (err error) {
	logus.Warningf(ctx, "ContactDalGae.DeleteContact(%s)", contactID)
	var contact models4debtus.DebtusSpaceContactEntry
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if contact, err = GetDebtusSpaceContactByID(ctx, tx, spaceID, contactID); err != nil {
			if dal.IsNotFound(err) {
				logus.Warningf(ctx, "DebtusSpaceContactEntry not found by ContactID: %v", contactID)
				err = nil
			}
			return
		}
		if contact.Data != nil && contact.Data.CounterpartyUserID != "" {
			return ErrContactIsNotDeletable
		}

		debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)

		if err = tx.Get(ctx, debtusSpace.Record); err != nil {
			return
		}
		if userContact := debtusSpace.Data.Contacts[contactID]; userContact != nil {
			userContactBalance := userContact.Balance
			contactBalance := contact.Data.Balance
			if !reflect.DeepEqual(userContactBalance, contactBalance) {
				return fmt.Errorf("Data integrity issue: userContactBalance != contactBalance\n\tuserContactBalance: %v\n\tcontactBalance: %v", userContactBalance, contactBalance)
			}
			delete(debtusSpace.Data.Contacts, contactID)
			if len(contact.Data.Balance) > 0 {
				err = errors.New("removing Contact with non-zero balance is not implemented yet")
				return
				//userBalance := user.Data.Balance()
				//for k, v := range contactBalance {
				//	userBalance[k] -= v
				//}
				//if err = user.Data.SetBalance(userBalance); err != nil {
				//	return err
				//}
			}
		}
		key := models4debtus.NewDebtusContactKey(spaceID, contactID)
		if err = tx.Delete(ctx, key); err != nil {
			return err
		}
		return nil
	}, dal.TxWithCrossGroup())
	return
}

func SaveContact(ctx context.Context, contact models4debtus.DebtusSpaceContactEntry) error {
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(ctx, contact.Record)
	})
}

func GetDebtusSpaceContactsByIDs(ctx context.Context, tx dal.ReadSession, spaceID string, contactsIDs []string) (debtusContacts []models4debtus.DebtusSpaceContactEntry, err error) {
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
	}
	debtusContacts = models4debtus.NewDebtusSpaceContacts(spaceID, contactsIDs...)
	records := models4debtus.DebtusContactRecords(debtusContacts)
	return debtusContacts, tx.GetMulti(ctx, records)
}

func GetDebtusSpaceContactByID(ctx context.Context, tx dal.ReadSession, spaceID, contactID string) (contact models4debtus.DebtusSpaceContactEntry, err error) {
	contact = models4debtus.NewDebtusSpaceContactEntry(spaceID, contactID, nil)
	return contact, GetDebtusSpaceContact(ctx, tx, contact)
}

func GetDebtusSpaceContact(ctx context.Context, tx dal.ReadSession, contact models4debtus.DebtusSpaceContactEntry) (err error) {
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
	}
	return tx.Get(ctx, contact.Record)
}
