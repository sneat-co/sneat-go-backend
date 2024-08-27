package facade4debtus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/logus"

	"context"
	"errors"
)

type usersLinker struct {
	// Groups methods for linking 2 users via DebtusSpaceContactEntry
	changes *usersLinkingDbChanges
}

func newUsersLinker(changes *usersLinkingDbChanges) *usersLinker {
	return &usersLinker{
		changes: changes,
	}
}

func (linker *usersLinker) linkUsersWithinTransaction(
	tctx context.Context, // 'tc' is transactional context, 'c' is not
	tx dal.ReadwriteTransaction,
	linkedBy string,
) (
	err error,
) {
	changes := linker.changes
	inviter, invited := changes.inviter, changes.invited
	if inviter.user.ID == "" {
		return errors.New("inviter.user.ID is empty string")
	}
	if invited.user.ID == "" {
		return errors.New("invited.user.ID is empty string")
	}
	//if invited.contact.ID == "" {
	//	//invitedContact = new(models4debtus.DebtusSpaceContactEntry)
	//}

	logus.Debugf(tctx, "usersLinker.linkUsersWithinTransaction(inviterUser.ContactID=%s, invitedUser.ContactID=%s, inviterContact=%s, inviterContact.UserID=%s)",
		inviter.user.ID, invited.user.ID, inviter.contact.ID, inviter.contact.Data.UserID)

	// First lets validate input
	if err = linker.validateInput(changes.inviter, changes.invited); err != nil {
		return
	}

	if tx == nil {
		err = errors.New("usersLinker.linkUsersWithinTransaction is called without transaction")
		return
	}

	// Update entities
	{
		if err = linker.getOrCreateInvitedContactByInviterUserAndInviterContact(tctx, tx, changes); err != nil {
			return
		}

		if invited.contact.Data == nil {
			err = fmt.Errorf(
				"getOrCreateInvitedContactByInviterUserAndInviterContact() returned invitedContact.Data == nil, invitedContact.ContactID: %s",
				invited.contact.ID)
			return
		}
		if invited.contact.Data.UserID != invited.user.ID {
			return fmt.Errorf("invitedContact.UserID != invitedUser.ContactID: %v != %v", invited.contact.Data.UserID, invited.user.ID)
		}

		logus.Debugf(tctx, "getOrCreateInvitedContactByInviterUserAndInviterContact() => invitedContact.ContactID: %v", invited.contact.ID)

		if err = linker.updateInvitedUser(tctx, invited.user, invited.debtusSpace, inviter.user.ID, inviter.debtusContact); err != nil {
			return
		}

		if _, err = linker.updateInviterContact(tctx, inviter, invited, linkedBy); err != nil {
			return
		}
	}

	// verify
	{
		invited.debtusContact.Data.MustMatchCounterparty(inviter.debtusContact)

		addContactJSON := func(debtusSpace models4debtus.DebtusSpaceEntry, debtusContact models4debtus.DebtusSpaceContactEntry) (debtusContactBrief *models4debtus.DebtusContactBrief) {
			//debtusContactBrief = debtusSpace.Data.Contacts[debtusContact.ID]
			debtusContactBrief, _ = models4debtus.AddOrUpdateDebtusContact(debtusSpace, debtusContact)
			linker.changes.FlagAsChanged(debtusSpace.Record)
			return debtusContactBrief
		}
		inviterUserContact := addContactJSON(inviter.debtusSpace, inviter.debtusContact)
		invitedUserContact := addContactJSON(invited.debtusSpace, invited.debtusContact)

		if !invitedUserContact.Balance.Equal(inviterUserContact.Balance.Reversed()) {
			panic(fmt.Sprintf("users contacts json balances are not equal (invited vs inviter): %v != %v",
				invited.debtusSpace.Data.Contacts[invited.contact.ID].Balance,
				inviter.debtusSpace.Data.Contacts[inviter.contact.ID].Balance,
			))
		}
	}
	return
}

func (linker *usersLinker) validateInput(
	inviter *userLinkingParty,
	invited *userLinkingParty,
) error {
	if inviter.user.ID == "" {
		return errors.New("inviter.user.ID is empty string")
	}
	if invited.user.ID == "" {
		return errors.New("invitedUser.ID is empty string")
	}
	if inviter.debtusContact.ID == "" {
		return errors.New("inviter.debtusContact.ID is empty string")
	}
	if inviter.user.ID == invited.user.ID {
		return fmt.Errorf("inviter.user.ID == invited.user.ID: %s", inviter.user.ID)
	}
	if inviter.contact.Data.UserID != inviter.user.ID {
		return fmt.Errorf("usersLinker.validateInput(): inviterDebtusContact.UserID != inviterUser.ContactID: %s != %s", inviter.contact.Data.UserID, inviter.user.ID)
	}
	return nil
}

// Purpose of the function is an attempt to link existing counterparties
func (linker *usersLinker) getOrCreateInvitedContactByInviterUserAndInviterContact(
	tctx context.Context,
	tx dal.ReadwriteTransaction,
	changes *usersLinkingDbChanges,
) (err error) {
	inviter, invited := changes.inviter, changes.invited
	logus.Debugf(tctx, "getOrCreateInvitedContactByInviterUserAndInviterContact()\n\tinviterContact.ContactID: %v", inviter.contact.ID)
	if inviter.user.ID == invited.user.ID {
		panic(fmt.Sprintf("inviterUser.ContactID == invitedUser.ContactID: %v", inviter.user.ID))
	}

	if len(invited.contactusSpace.Data.Contacts) > 0 {
		var invitedUserContacts []models4debtus.DebtusSpaceContactEntry
		// Use non transaction context
		invitedUserContacts, err = GetDebtusSpaceContactsByIDs(tctx, tx, invited.spaceID, invited.contactusSpace.Data.ContactIDs())
		if err != nil {
			err = fmt.Errorf("failed to call facade4debtus.GetDebtusSpaceContactsByIDs(): %w", err)
			return
		}
		for _, invitedUserContact := range invitedUserContacts {
			if invitedUserContact.Data.CounterpartyUserID == inviter.user.ID {
				// We re-get the entity of the found invitedContact using transactional context
				// and store it to output var
				if invited.debtusContact, err = GetDebtusSpaceContactByID(tctx, tx, invited.spaceID, invitedUserContact.ID); err != nil {
					err = fmt.Errorf("failed to call GetDebtusSpaceContactByID(%s): %w", invitedUserContact.ID, err)
					return
				}
				if invited.contact.Data.Names.FirstName == "" {
					invited.contact.Data.Names.FirstName = inviter.user.Data.Names.FirstName
				}
				if invited.contact.Data.Names.LastName == "" {
					invited.contact.Data.Names.LastName = inviter.user.Data.Names.LastName
				}
				break
			}
		}
	}

	if invited.debtusContact.ID == "" {
		logus.Debugf(tctx, "getOrCreateInvitedContactByInviterUserAndInviterContact(): creating new Contact for invited user")
		invitedContactDetails := dto4contactus.ContactDetails{
			NameFields: *inviter.user.Data.Names,
		}
		createContactDbChanges := &createContactDbChanges{
			user: changes.invited.user,
			counterparty: ParticipantEntries{
				User:          changes.invited.user,
				Contact:       changes.inviter.contact,
				DebtusSpace:   changes.inviter.debtusSpace,
				DebtusContact: changes.inviter.debtusContact,
			},
		}
		if err = createContactWithinTransaction(tctx, tx, createContactDbChanges, inviter.spaceID, inviter.user.ID, invitedContactDetails); err != nil {
			return
		}
		//if changes.inviter.contact.ID == "" {
		//	//linker.changes.FlagAsChanged(inviterContact.Record)
		//}
		if invited.debtusSpace.Data.LastTransferAt.Before(inviter.debtusContact.Data.LastTransferAt) {
			invited.debtusSpace.Data.LastTransferID = inviter.debtusContact.Data.LastTransferID
			invited.debtusSpace.Data.LastTransferAt = inviter.debtusContact.Data.LastTransferAt
			linker.changes.FlagAsChanged(invited.debtusSpace.Record)
		}
	} else {
		logus.Debugf(tctx, "getOrCreateInvitedContactByInviterUserAndInviterContact(): linking existing Contact: %v", invited.contact)
		// TODO: How do we merge existing contacts?
		invited.debtusContact.Data.CountOfTransfers = inviter.debtusContact.Data.CountOfTransfers
		invited.debtusContact.Data.LastTransferID = inviter.debtusContact.Data.LastTransferID
		invited.debtusContact.Data.LastTransferAt = inviter.debtusContact.Data.LastTransferAt
		invited.debtusContact.Data.Balance = inviter.debtusContact.Data.Balance.Reversed()
		linker.changes.FlagAsChanged(invited.debtusContact.Record)
	}
	invited.debtusContact.Data.MustMatchCounterparty(inviter.debtusContact)
	return
}

func (linker *usersLinker) updateInvitedUser(ctx context.Context,
	invitedUser dbo4userus.UserEntry,
	invitedDebtusSpace models4debtus.DebtusSpaceEntry,
	inviterUserID string,
	inviterDebtusContact models4debtus.DebtusSpaceContactEntry,
) (err error) {
	logus.Debugf(ctx, "usersLinker.updateInvitedUser()")
	var invitedUserChanged bool

	if invitedUser.Data.InvitedByUserID == "" {
		invitedUser.Data.InvitedByUserID = inviterUserID
		invitedUser.Record.MarkAsChanged()
	}

	if inviterDebtusContact.Data.LastTransferAt.After(invitedDebtusSpace.Data.LastTransferAt) {
		invitedDebtusSpace.Data.LastTransferID = inviterDebtusContact.Data.LastTransferID
		invitedDebtusSpace.Data.LastTransferAt = inviterDebtusContact.Data.LastTransferAt
		invitedUserChanged = true
	}

	if invitedUserChanged {
		linker.changes.FlagAsChanged(invitedDebtusSpace.Record)
	}
	return
}

// Updates counterparty entity that belongs to inviter user (inviterContact.UserID == inviterUser.ContactID)
func (linker *usersLinker) updateInviterContact(
	tctx context.Context,
	inviter *userLinkingParty,
	invited *userLinkingParty,
	linkedBy string,
) (
	isJustConnected bool, err error,
) {
	logus.Debugf(tctx, "usersLinker.updateInviterContact(), inviterContact.CounterpartyUserID: %s, inviterContact.CountOfTransfers: %d", inviter.debtusContact.Data.CounterpartyUserID, inviter.debtusContact.Data.CountOfTransfers)
	// validate input
	{
		if inviter.user.ID == "" {
			err = errors.New("inviter.user.ID == 0")
			return
		}
		if inviter.contact.Data.UserID != inviter.user.ID {
			panic(fmt.Sprintf("usersLinker.updateInviterContact(): inviterContact.UserID != inviterUser.ContactID: %s != %s\ninvitedContact.UserID: %s, invitedUser.ContactID: %s",
				inviter.contact.Data.UserID, inviter.user.ID, invited.contact.Data.UserID, invited.user.ID))
		}
		if invited.contact.Data.UserID != invited.user.ID {
			panic(fmt.Errorf("invitedContact.UserID != invitedUser.ContactID: %s != %s\ninviterContact.UserID: %s, inviterUser.ContactID: %s",
				invited.contact.Data.UserID, invited.contact.ID, inviter.contact.Data.UserID, inviter.user.ID))
		}
		if invited.contact.ID == inviter.contact.ID {
			panic(fmt.Sprintf("invitedContact.ContactID == inviterContact.ContactID: %v", invited.contact.ID))
		}
		if invited.user.ID == inviter.user.ID {
			panic(fmt.Sprintf("invitedUser.ContactID == inviterUser.ContactID: %v", invited.user.ID))
		}
	}
	var inviterContactChanged bool
	if inviter.contact.Data.Names.FirstName == "" && invited.user.Data.Names.FirstName != "" {
		inviter.contact.Data.Names.FirstName = invited.user.Data.Names.FirstName
		inviterContactChanged = true
	}
	if inviter.contact.Data.Names.LastName == "" && invited.user.Data.Names.LastName != "" {
		inviter.contact.Data.Names.LastName = invited.user.Data.Names.LastName
		inviterContactChanged = true
	}
	//if inviterContactChanged {
	//	inviterContact.UpdateSearchName()
	//}
	if inviterContactChanged {
		linker.changes.FlagAsChanged(inviter.contact.Record)
	} else {
		defer func() {
			if inviterContactChanged {
				linker.changes.FlagAsChanged(inviter.contact.Record)
			}
		}()
	}
	switch inviter.contact.Data.UserID {
	case "":
		logus.Debugf(tctx, "Updating inviterUser.DebtusSpaceContactEntry* fields...")
		isJustConnected = true
		inviterContactChanged = true
		inviter.contact.Data.UserID = invited.user.ID
		inviter.debtusContact.Data.CounterpartySpaceID = invited.spaceID
		inviter.debtusContact.Data.CounterpartyContactID = invited.contact.ID
		inviter.debtusContact.Data.LinkedBy = linkedBy
		inviterContacts := inviter.contactusSpace.Data.Contacts
		for inviterSpaceContactID, inviterSpaceContact := range inviterContacts {
			if inviterSpaceContactID == inviter.contact.ID {
				if inviterSpaceContact.UserID == "" {
					inviterSpaceContact.UserID = inviter.debtusContact.Data.CounterpartyUserID
					inviterContacts[inviterSpaceContactID] = inviterSpaceContact
					linker.changes.FlagAsChanged(inviter.contactusSpace.Record)
				} else if inviterSpaceContact.UserID == inviter.debtusContact.Data.CounterpartyUserID {
					// do nothing
				} else {
					err = fmt.Errorf(
						"data integrity issue for Contact %s: inviterSpaceContact.UserID != inviterContact.CounterpartyUserID: %s != %v",
						inviter.contact.ID, inviterSpaceContact.UserID, inviter.debtusContact.Data.CounterpartyUserID)
					return
				}
				goto inviterUserContactFound
			}
		}
		if _, changed := models4debtus.AddOrUpdateDebtusContact(inviter.debtusSpace, inviter.debtusContact); changed {
			linker.changes.FlagAsChanged(inviter.debtusSpace.Record)
		}
	inviterUserContactFound:
		// Queue task to update all existing api4transfers
		if inviter.debtusContact.Data.CountOfTransfers > 0 {
			if err = dtdal.Transfer.DelayUpdateTransfersWithCounterparty(
				tctx,
				invited.contact.ID,
				inviter.contact.ID,
			); err != nil {
				err = fmt.Errorf("failed to enqueue delayUpdateTransfersWithCounterparty(): %w", err)
				return
			}
		} else {
			logus.Debugf(tctx, "No need to update api4transfers of inviter as inviterContact.CountOfTransfers == 0")
		}
	case invited.user.ID:
		logus.Infof(tctx, "inviterContact.CounterpartyUserID is already set, updateInviterContact() did nothing")
	default:
		err = fmt.Errorf("inviterContact.CounterpartyUserID is different from current user. inviterContact.CounterpartyUserID: %s, currentUserID: %s", inviter.debtusContact.Data.CounterpartyUserID, invited.user.ID)
		return
	}
	return
}
