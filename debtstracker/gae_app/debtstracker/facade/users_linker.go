package facade

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/logus"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type usersLinker struct {
	// Groups methods for linking 2 users via ContactEntry
	changes *usersLinkingDbChanges
}

func newUsersLinker(changes *usersLinkingDbChanges) usersLinker {
	return usersLinker{
		changes: changes,
	}
}

func (linker usersLinker) linkUsersWithinTransaction(
	tc context.Context, // 'tc' is transactional context, 'c' is not
	tx dal.ReadwriteTransaction,
	linkedBy string,
) (
	err error,
) {
	changes := linker.changes
	inviterUser, invitedUser := changes.inviterUser, changes.invitedUser
	if inviterUser == nil {
		panic("inviterUser == nil")
	}
	if invitedUser == nil {
		panic("invitedUser == nil")
	}
	inviterContact, invitedContact := changes.inviterContact, changes.invitedContact
	if invitedContact == nil {
		invitedContact = new(models.ContactEntry)
		changes.invitedContact = invitedContact
	}

	logus.Debugf(tc, "usersLinker.linkUsersWithinTransaction(inviterUser.ID=%s, invitedUser.ID=%s, inviterContact=%s, inviterContact.UserID=%s)", inviterUser.ID, invitedUser.ID, inviterContact.ID, inviterContact.Data.UserID)

	// First of all lets validate input
	if err = linker.validateInput(inviterUser, invitedUser, inviterContact); err != nil {
		return
	}

	if tx == nil {
		err = errors.New("usersLinker.linkUsersWithinTransaction is called without transaction")
		return
	}

	// Update entities
	{
		if err = linker.getOrCreateInvitedContactByInviterUserAndInviterContact(tc, tx, changes); err != nil {
			return
		}
		invitedContact = changes.invitedContact

		if invitedContact.Data == nil {
			err = fmt.Errorf(
				"getOrCreateInvitedContactByInviterUserAndInviterContact() returned invitedContact.Data == nil, invitedContact.ID: %s",
				invitedContact.ID)
			return
		}
		if invitedContact.Data.UserID != invitedUser.ID {
			return fmt.Errorf("invitedContact.UserID != invitedUser.ID: %v != %v", invitedContact.Data.UserID, invitedUser.ID)
		}

		logus.Debugf(tc, "getOrCreateInvitedContactByInviterUserAndInviterContact() => invitedContact.ID: %v", invitedContact.ID)

		if err = linker.updateInvitedUser(tc, *invitedUser, inviterUser.ID, *inviterContact); err != nil {
			return
		}

		if _, err = linker.updateInviterContact(tc, *inviterUser, *invitedUser, inviterContact, invitedContact, linkedBy); err != nil {
			return
		}
	}

	// verify
	{
		invitedContact.Data.MustMatchCounterparty(*inviterContact)

		addContactJSON := func(user *models.AppUser, contact *models.ContactEntry) (contactJSON *models.UserContactJson) {
			contactJSON = user.Data.ContactByID(invitedContact.ID)
			if contactJSON == nil {
				// err = fmt.Errorf("invitedUserContact == nil, ID=%v", invitedContact.ID)
				userContactJSON, _ := models.AddOrUpdateContact(user, *contact)
				contactJSON = &userContactJSON
				linker.changes.FlagAsChanged(user.Record)
			}
			return contactJSON
		}
		inviterUserContact := addContactJSON(inviterUser, inviterContact)
		invitedUserContact := addContactJSON(invitedUser, invitedContact)

		if !invitedUserContact.Balance().Equal(inviterUserContact.Balance().Reversed()) {
			panic(fmt.Sprintf("users contacts json balances are not equal (invited vs inviter): %v != %v",
				invitedUser.Data.ContactByID(invitedContact.ID).Balance(),
				inviterUser.Data.ContactByID(inviterContact.ID).Balance(),
			))
		}
	}
	return
}

func (linker usersLinker) validateInput(
	inviterUser, invitedUser *models.AppUser,
	inviterContact *models.ContactEntry,
) error {
	if inviterUser.ID == "" {
		panic("inviterUser.ID == 0")
	}
	if invitedUser.ID == "" {
		panic("invitedUser.ID == 0")
	}
	if inviterContact.ID == "" {
		panic("inviterContact.ID == 0")
	}
	if inviterUser.ID == invitedUser.ID {
		panic(fmt.Sprintf("inviterUser.ID == invitedUser.ID: %v", inviterUser.ID))
	}
	if inviterContact.Data.UserID != inviterUser.ID {
		panic(fmt.Sprintf("usersLinker.validateInput(): inviterContact.UserID != inviterUser.ID: %v != %v", inviterContact.Data.UserID, inviterUser.ID))
	}
	return nil
}

// Purpose of the function is an attempt to link existing counterparties
func (linker usersLinker) getOrCreateInvitedContactByInviterUserAndInviterContact(
	tc context.Context,
	tx dal.ReadwriteTransaction,
	changes *usersLinkingDbChanges,
) (err error) {
	inviterUser, invitedUser := *changes.inviterUser, *changes.invitedUser
	inviterContact := *changes.inviterContact
	logus.Debugf(tc, "getOrCreateInvitedContactByInviterUserAndInviterContact()\n\tinviterContact.ID: %v", inviterContact.ID)
	if inviterUser.ID == invitedUser.ID {
		panic(fmt.Sprintf("inviterUser.ID == invitedUser.ID: %v", inviterUser.ID))
	}

	var invitedContact models.ContactEntry
	if changes.invitedContact != nil && changes.invitedContact.ID != "" {
		invitedContact = *changes.invitedContact
	} else {
		changes.invitedContact = &invitedContact
	}

	if invitedUser.Data.ContactsCount > 0 {
		var invitedUserContacts []models.ContactEntry
		// Use non transaction context
		invitedUserContacts, err = GetContactsByIDs(tc, tx, invitedUser.Data.ContactIDs())
		if err != nil {
			err = fmt.Errorf("failed to call facade.GetContactsByIDs(): %w", err)
			return
		}
		for _, invitedUserContact := range invitedUserContacts {
			if invitedUserContact.Data.CounterpartyUserID == inviterUser.ID {
				// We re-get the entity of the found invitedContact using transactional context
				// and store it to output var
				if invitedContact, err = GetContactByID(tc, tx, invitedUserContact.ID); err != nil {
					err = fmt.Errorf("failed to call GetContactByID(%s): %w", invitedUserContact.ID, err)
					return
				}
				if invitedContact.Data.FirstName == "" {
					invitedContact.Data.FirstName = inviterUser.Data.FirstName
				}
				if invitedContact.Data.LastName == "" {
					invitedContact.Data.LastName = inviterUser.Data.LastName
				}
				break
			}
		}
	}

	if invitedContact.ID == "" {
		logus.Debugf(tc, "getOrCreateInvitedContactByInviterUserAndInviterContact(): creating new contact for invited user")
		invitedContactDetails := models.ContactDetails{
			FirstName:  inviterUser.Data.FirstName,
			LastName:   inviterUser.Data.LastName,
			Nickname:   inviterUser.Data.Nickname,
			ScreenName: inviterUser.Data.ScreenName,
			Username:   inviterUser.Data.Username,
		}
		createContactDbChanges := &createContactDbChanges{
			user:                changes.invitedUser,
			counterpartyContact: changes.inviterContact,
		}
		if invitedContact, inviterContact, err = createContactWithinTransaction(tc, tx, createContactDbChanges, inviterUser.ID, invitedContactDetails); err != nil {
			return
		}
		if changes.inviterContact == nil {
			changes.inviterContact = &inviterContact
			linker.changes.FlagAsChanged(inviterContact.Record)
		}
		if invitedUser.Data.LastTransferAt.Before(inviterContact.Data.LastTransferAt) {
			invitedUser.Data.LastTransferID = inviterContact.Data.LastTransferID
			invitedUser.Data.LastTransferAt = inviterContact.Data.LastTransferAt
			linker.changes.FlagAsChanged(invitedUser.Record)
		}
	} else {
		logus.Debugf(tc, "getOrCreateInvitedContactByInviterUserAndInviterContact(): linking existing contact: %v", invitedContact)
		// TODO: How do we merge existing contacts?
		invitedContact.Data.CountOfTransfers = inviterContact.Data.CountOfTransfers
		invitedContact.Data.LastTransferID = inviterContact.Data.LastTransferID
		invitedContact.Data.LastTransferAt = inviterContact.Data.LastTransferAt
		if err = invitedContact.Data.SetBalance(inviterContact.Data.Balance().Reversed()); err != nil {
			return
		}
		linker.changes.FlagAsChanged(invitedContact.Record)
	}
	invitedContact.Data.MustMatchCounterparty(inviterContact)
	return
}

func (linker usersLinker) updateInvitedUser(c context.Context,
	invitedUser models.AppUser,
	inviterUserID string, inviterContact models.ContactEntry,
) (err error) {
	logus.Debugf(c, "usersLinker.updateInvitedUser()")
	var invitedUserChanged bool

	if invitedUser.Data.InvitedByUserID == "" {
		invitedUser.Data.InvitedByUserID = inviterUserID
		invitedUserChanged = true
	}

	if inviterContact.Data.LastTransferAt.After(invitedUser.Data.LastTransferAt) {
		invitedUser.Data.LastTransferID = inviterContact.Data.LastTransferID
		invitedUser.Data.LastTransferAt = inviterContact.Data.LastTransferAt
		invitedUserChanged = true
	}

	if invitedUserChanged {
		linker.changes.FlagAsChanged(invitedUser.Record)
	}
	return
}

// Updates counterparty entity that belongs to inviter user (inviterContact.UserID == inviterUser.ID)
func (linker usersLinker) updateInviterContact(
	tc context.Context,
	inviterUser, invitedUser models.AppUser,
	inviterContact, invitedContact *models.ContactEntry,
	linkedBy string,
) (
	isJustConnected bool, err error,
) {
	logus.Debugf(tc, "usersLinker.updateInviterContact(), inviterContact.CounterpartyUserID: %s, inviterContact.CountOfTransfers: %d", inviterContact.Data.CounterpartyUserID, inviterContact.Data.CountOfTransfers)
	// validate input
	{
		if inviterUser.ID == "" {
			panic("inviterUser.ID == 0")
		}
		if invitedUser.ID == "" {
			panic("invitedUser.ID == 0")
		}
		if inviterContact.Data.UserID != inviterUser.ID {
			panic(fmt.Sprintf("usersLinker.updateInviterContact(): inviterContact.UserID != inviterUser.ID: %v != %v\ninvitedContact.UserID: %v, invitedUser.ID: %v",
				inviterContact.Data.UserID, inviterUser.ID, invitedContact.Data.UserID, invitedUser.ID))
		}
		if invitedContact.Data.UserID != invitedUser.ID {
			panic(fmt.Errorf("invitedContact.UserID != invitedUser.ID: %v != %v\ninviterContact.UserID: %v, inviterUser.ID: %v",
				invitedContact.Data.UserID, invitedContact.ID, inviterContact.Data.UserID, inviterUser.ID))
		}
		if invitedContact.ID == inviterContact.ID {
			panic(fmt.Sprintf("invitedContact.ID == inviterContact.ID: %v", invitedContact.ID))
		}
		if invitedUser.ID == inviterUser.ID {
			panic(fmt.Sprintf("invitedUser.ID == inviterUser.ID: %v", invitedUser.ID))
		}
	}
	var inviterContactChanged bool
	if inviterContact.Data.FirstName == "" {
		inviterContact.Data.FirstName = invitedUser.Data.FirstName
		inviterContactChanged = true
	}
	if inviterContact.Data.LastName == "" {
		inviterContact.Data.LastName = invitedUser.Data.LastName
		inviterContactChanged = true
	}
	//if inviterContactChanged {
	//	inviterContact.UpdateSearchName()
	//}
	if inviterContactChanged {
		linker.changes.FlagAsChanged(inviterContact.Record)
	} else {
		defer func() {
			if inviterContactChanged {
				linker.changes.FlagAsChanged(inviterContact.Record)
			}
		}()
	}
	switch inviterContact.Data.CounterpartyUserID {
	case "":
		logus.Debugf(tc, "Updating inviterUser.ContactEntry* fields...")
		isJustConnected = true
		inviterContactChanged = true
		inviterContact.Data.CounterpartyUserID = invitedUser.ID
		inviterContact.Data.CounterpartyCounterpartyID = invitedContact.ID
		inviterContact.Data.LinkedBy = linkedBy
		inviterUserContacts := inviterUser.Data.Contacts()
		for i, inviterUserContact := range inviterUserContacts {
			if inviterUserContact.ID == inviterContact.ID {
				if inviterUserContact.UserID == "" {
					inviterUserContact.UserID = inviterContact.Data.CounterpartyUserID
					inviterUserContacts[i] = inviterUserContact
					inviterUser.Data.SetContacts(inviterUserContacts)
					linker.changes.FlagAsChanged(inviterUser.Record)
				} else if inviterUserContact.UserID == inviterContact.Data.CounterpartyUserID {
					// do nothing
				} else {
					err = fmt.Errorf(
						"data integrity issue for contact %v: inviterUserContact.UserID != inviterContact.CounterpartyUserID: %v != %v",
						inviterContact.ID, inviterUserContact.UserID, inviterContact.Data.CounterpartyUserID)
					return
				}
				goto inviterUserContactFound
			}
		}
		if _, changed := models.AddOrUpdateContact(&inviterUser, *inviterContact); changed {
			linker.changes.FlagAsChanged(inviterUser.Record)
		}
	inviterUserContactFound:
		// Queue task to update all existing transfers
		if inviterContact.Data.CountOfTransfers > 0 {
			if err = dtdal.Transfer.DelayUpdateTransfersWithCounterparty(
				tc,
				invitedContact.ID,
				inviterContact.ID,
			); err != nil {
				err = fmt.Errorf("failed to enqueue delayUpdateTransfersWithCounterparty(): %w", err)
				return
			}
		} else {
			logus.Debugf(tc, "No need to update transfers of inviter as inviterContact.CountOfTransfers == 0")
		}
	case invitedUser.ID:
		logus.Infof(tc, "inviterContact.CounterpartyUserID is already set, updateInviterContact() did nothing")
	default:
		err = fmt.Errorf("inviterContact.CounterpartyUserID is different from current user. inviterContact.CounterpartyUserID: %v, currentUserID: %v", inviterContact.Data.CounterpartyUserID, invitedUser.ID)
		return
	}
	return
}
