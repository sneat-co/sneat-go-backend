package facade2debtus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"time"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/slices"
)

type ReceiptUsersLinker struct {
	changes *receiptDbChanges
}

func NewReceiptUsersLinker(changes *receiptDbChanges) *ReceiptUsersLinker {
	if changes == nil {
		changes = newReceiptDbChanges()
	}
	return &ReceiptUsersLinker{
		changes: changes,
	}
}

func (linker *ReceiptUsersLinker) LinkReceiptUsers(c context.Context, receiptID, invitedUserID string) (isJustLinked bool, err error) {
	logus.Debugf(c, "ReceiptUsersLinker.LinkReceiptUsers(receiptID=%v, invitedUserID=%v)", receiptID, invitedUserID)
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return false, err
	}
	if invitedUser, err := User.GetUserByID(c, db, invitedUserID); err != nil {
		// TODO: Instead pass user as a parameter? Even better if the user entity was created within following transaction.
		return isJustLinked, err
	} else if invitedUser.Data.DtCreated.After(time.Now().Add(-time.Second / 2)) {
		logus.Debugf(c, "A new user, will wait for half a seconds to cleanup previous transaction")
		time.Sleep(time.Second / 2)
	}
	var invitedContact models.ContactEntry
	attempt := 0
	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
		if attempt += 1; attempt > 1 {
			sleepPeriod := time.Duration(attempt) * time.Second
			logus.Warningf(c, "Transaction retry will sleep for %v, invitedContact.ID: %v", attempt, invitedContact.ID)
			time.Sleep(sleepPeriod)
		}
		changes := linker.changes
		var (
			receipt     models.Receipt
			transfer    models.TransferEntry
			inviterUser models.AppUser
			invitedUser models.AppUser
		)
		if receipt, transfer, inviterUser, invitedUser, err = getReceiptTransferAndUsers(tc, tx, receiptID, invitedUserID); err != nil {
			return
		}
		changes.receipt = &receipt
		changes.transfer = &transfer
		changes.inviterUser = &inviterUser
		changes.invitedUser = &invitedUser
		if invitedContact.ID != "" { // This means we are attempting to retry failed transaction
			if err = workaroundReinsertContact(tc, receipt, invitedContact, changes); err != nil {
				return
			}
		}

		if isJustLinked, err = linker.linkUsersByReceiptWithinTransaction(c, tc, tx); err != nil {
			return
		} else {
			invitedContact = *changes.invitedContact
		}

		// Integrity checks
		{
			invitedContact.Data.MustMatchCounterparty(*changes.inviterContact)
		}

		if entitiesToSave := changes.Changes.Records(); len(entitiesToSave) > 0 {
			if err = tx.SetMulti(c, entitiesToSave); err != nil {
				return
			}
		} else {
			logus.Debugf(c, "Receipt and transfer has not changed")
		}
		return
	}, dal.TxWithCrossGroup())
	if err != nil {
		return
	}
	logus.Debugf(c, "ReceiptUsersLinker.LinkReceiptUsers() => invitedContact: %+v", invitedContact)
	if invitedContact, err = GetContactByID(c, nil, invitedContact.ID); err != nil {
		return
	}
	logus.Debugf(c, "ReceiptUsersLinker.LinkReceiptUsers() => invitedContact from DB: %+v", invitedContact)
	return
}

func (linker *ReceiptUsersLinker) linkUsersByReceiptWithinTransaction(
	c context.Context, // non-transactional context
	tc context.Context, // transactional context,
	tx dal.ReadwriteTransaction,
) (
	isCounterpartiesJustConnected bool,
	err error,
) {
	changes := linker.changes
	receipt := changes.receipt
	transfer := changes.transfer
	inviterUser, invitedUser := *changes.inviterUser, *changes.invitedUser
	var invitedContact, inviterContact models.ContactEntry
	if changes.inviterContact != nil {
		inviterContact = *changes.inviterContact
	}
	if changes.invitedContact != nil {
		invitedContact = *changes.invitedContact
	}

	logus.Debugf(c,
		"ReceiptUsersLinker.linkUsersByReceiptWithinTransaction(receipt.ID=%s, transfer.ID=%s, inviterUser.ID=%s, invitedUser.ID=%s, inviterContact.ID=%s, invitedContact.ID=%s)",
		receipt.ID, transfer.ID, inviterUser.ID, invitedUser.ID, inviterContact.ID, invitedContact.ID)

	{ // validate inputs
		if err = linker.validateInput(changes); err != nil {
			return
		}
		if receipt.Data.TransferID != transfer.ID {
			panic(fmt.Sprintf("receipt.TransferID != transfer.ID: %v != %v", receipt.Data.TransferID, transfer.ID))
		}
		if transferCreatorUserID := transfer.Data.Creator().UserID; transferCreatorUserID == "" {
			panic("transfer.Creator().UserID is zero")
		} else if transferCreatorUserID != inviterUser.ID {
			panic(fmt.Sprintf("transfer.Creator().UserID != inviterUser.ID:  %v != %v", transferCreatorUserID, inviterUser.ID))
		} else if transferCreatorUserID == invitedUser.ID {
			panic(fmt.Sprintf("transfer.Creator().UserID == invitedUser.ID:  %v != %v", transferCreatorUserID, invitedUser.ID))
		}
	}

	fromOriginal := *transfer.Data.From()
	toOriginal := *transfer.Data.To()
	//logus.Debugf(c, "transferEntity: %v", transfer.Data)
	//logus.Debugf(c, "transfer.From(): %v", fromOriginal)
	//logus.Debugf(c, "transfer.To(): %v",toOriginal)

	transferCreatorCounterparty := transfer.Data.Counterparty()

	if inviterContact, err = GetContactByID(tc, tx, transferCreatorCounterparty.ContactID); err != nil {
		return
	} else if inviterContact.Data.UserID != inviterUser.ID {
		panic(fmt.Errorf("inviterContact.UserID !=  inviterUser.ID: %v != %v", inviterContact.Data.UserID, inviterUser.ID))
	} else {
		changes.inviterContact = &inviterContact
	}

	if err = newUsersLinker(changes.usersLinkingDbChanges).linkUsersWithinTransaction(tc, tx, receipt.Record.Key().String()); err != nil {
		err = fmt.Errorf("failed to link users: %w", err)
		return
	} else {
		invitedContact = *changes.invitedContact // as was updated
	}
	{ // Update invited user's last currency
		var invitedUserChanged bool
		if invitedUser.Data.LastCurrencies, invitedUserChanged = slices.MergeStrings(invitedUser.Data.LastCurrencies, []string{string(transfer.Data.Currency)}); invitedUserChanged {
			changes.FlagAsChanged(changes.invitedUser.Record)
		}
	}

	logus.Debugf(c, "linkUsersWithinTransaction() => invitedContact.ID: %v, inviterContact.ID: %v", invitedContact.ID, inviterContact.ID)

	// Update entities
	{
		if err = linker.updateReceipt(); err != nil {
			return
		} else if err = linker.updateTransfer(); err != nil {
			return
		} else if linker.changes.IsChanged(linker.changes.transfer.Record) {
			logus.Debugf(c, "transfer changed:\n\tFrom(): %v\n\tTo(): %v", transfer.Data.From(), transfer.Data.To())
			// Just double check we did not screw up
			{
				if fromOriginal.UserID != "" && fromOriginal.UserID != transfer.Data.From().UserID {
					err = errors.New("fromOriginal.UserID != 0 && fromOriginal.UserID != transfer.From().UserID")
					return
				}
				if fromOriginal.ContactID != "" && fromOriginal.ContactID != transfer.Data.From().ContactID {
					err = errors.New("fromOriginal.ContactID != 0 && fromOriginal.ContactID != transfer.From().ContactID")
					return
				}
				if toOriginal.UserID != "" && toOriginal.UserID != transfer.Data.To().UserID {
					err = errors.New("toOriginal.UserID != 0 && toOriginal.UserID != transfer.To().UserID")
					return
				}
				if toOriginal.ContactID != "" && toOriginal.ContactID != transfer.Data.To().ContactID {
					err = errors.New("toOriginal.ContactID != 0 && toOriginal.ContactID != transfer.To().ContactID")
					return
				}
			}
		}
	}

	if transfer.Data.DtDueOn.After(time.Now()) {
		if err = dtdal.Reminder.DelayCreateReminderForTransferUser(tc, receipt.Data.TransferID, transfer.Data.Counterparty().UserID); err != nil {
			err = fmt.Errorf("failed to delay creation of reminder for transfer coutnerparty: %w", err)
			return
		}
	} else {
		if transfer.Data.DtDueOn.IsZero() {
			logus.Debugf(tc, "No need to create reminder for counterparty as no due date")
		} else {
			logus.Debugf(tc, "No need to create reminder for counterparty as due date in past")
		}
	}
	return
}

func (linker *ReceiptUsersLinker) validateInput(changes *receiptDbChanges) error {

	if changes.receipt.Data.CounterpartyUserID != "" {
		if changes.receipt.Data.CounterpartyUserID != changes.invitedUser.ID { // Already linked
			return errors.New("an attempt to link 3d user to a receipt")
		}

		transferCounterparty := changes.transfer.Data.Counterparty()

		if transferCounterparty.UserID != "" && transferCounterparty.UserID != changes.invitedUser.ID {
			return fmt.Errorf(
				"transferCounterparty.UserID != invitedUser.ID : %s != %s",
				transferCounterparty.UserID, changes.invitedUser.ID,
			)
		}
	}
	return nil
}

func (linker *ReceiptUsersLinker) updateReceipt() (err error) {
	receipt := *linker.changes.receipt
	counterpartyUser := *linker.changes.invitedUser
	if receipt.Data.CounterpartyUserID != counterpartyUser.ID {
		receipt.Data.CounterpartyUserID = counterpartyUser.ID
		linker.changes.FlagAsChanged(linker.changes.receipt.Record)
	}
	return
}

func (linker *ReceiptUsersLinker) updateTransfer() (err error) {
	changes := linker.changes

	transfer := changes.transfer
	inviterUser, invitedUser := *changes.inviterUser, *changes.invitedUser
	inviterContact, invitedContact := *changes.inviterContact, *changes.invitedContact
	{ // Validate input parameters
		if transfer.ID == "" || transfer.Data == nil {
			panic(fmt.Sprintf("Invalid parameter: transfer: %v", transfer))
		}
		validateSide := func(side string, user models.AppUser, contact models.ContactEntry) {
			if user.ID == "" || user.Data == nil {
				panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => %vUser: %v", side, user))
			}
			if contact.ID == "" || contact.Data == nil {
				panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => %vContact: %v", side, contact))
			} else if contact.Data.UserID != user.ID {
				panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => %vContact.UserID != %vUser.ID: %v != %v", side, side, contact.Data.UserID, invitedUser.ID))
			}
		}
		validateSide("inviter", inviterUser, inviterContact)
		validateSide("invited", invitedUser, invitedContact)
		if transfer.Data.CreatorUserID != inviterUser.ID {
			panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => transfer.CreatorUserID != inviterUser.ID: %v != %v", transfer.Data.CreatorUserID, invitedUser.ID))
		}
	}

	transferCounterparty := transfer.Data.Counterparty()

	if transferCounterparty.UserID != invitedUser.ID {
		if transferCounterparty.UserID != "" {
			err = fmt.Errorf("transfer.ContactEntry().UserID != counterpartyUserID : %s != %s",
				transfer.Data.Counterparty().UserID, invitedUser.ID)
			return
		}
		transfer.Data.Counterparty().UserID = invitedUser.ID
		linker.changes.FlagAsChanged(linker.changes.transfer.Record)
	}

	updateTransferCounterpartyInfo := func(
		side string,
		counterparty *models.TransferCounterpartyInfo,
		user models.AppUser,
		contact models.ContactEntry,
	) {
		if contact.Data.UserID == user.ID {
			panic(fmt.Sprintf(
				"updateTransferCounterpartyInfo() => %sContact.UserID == %sUser.ID : %s, counterparty.UserID: %s",
				side, side, contact.Data.UserID, counterparty.UserID))
		}
		if counterparty.UserID == "" {
			counterparty.UserID = user.ID
		} else if counterparty.UserID != user.ID {
			panic(fmt.Sprintf("updateTransferCounterpartyInfo() => counterparty.UserID != %sUser.ID : %s != %s, %sContact.UserID: %s", side, counterparty.UserID, user.ID, side, contact.Data.UserID))
		}
		counterparty.UserName = user.Data.FullName()

		if counterparty.ContactID == "" {
			counterparty.ContactID = contact.ID
		} else if counterparty.ContactID != contact.ID {
			panic(fmt.Sprintf(
				"ReceiptUsersLinker.updateTransfer() => counterparty.ContactID != %sContact.ID : %s != %s",
				side, counterparty.ContactID, contact.ID))
		}
		counterparty.ContactName = contact.Data.FullName()
	}

	updateTransferCounterpartyInfo("inviter", transfer.Data.Creator(), inviterUser, invitedContact)
	updateTransferCounterpartyInfo("invited", transfer.Data.Counterparty(), invitedUser, inviterContact)

	//if inlineMessageID != "" {
	//	transfer.CounterpartyTgReceiptInlineMessageID = inlineMessageID
	//}
	//transferAmount := transfer.Data.GetAmount()
	if transfer.Data.Direction() == models.TransferDirectionUser2Counterparty {
		transfer.Data.AmountInCents *= -1
	}

	return
}
