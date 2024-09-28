package facade4debtus

import (
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/slice"
	"time"

	"context"
	"errors"
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

func (linker *ReceiptUsersLinker) LinkReceiptUsers(ctx context.Context, receiptID, invitedUserID string) (isJustLinked bool, err error) {
	logus.Debugf(ctx, "ReceiptUsersLinker.LinkReceiptUsers(receiptID=%v, invitedUserID=%v)", receiptID, invitedUserID)
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return false, err
	}
	invitedUser := dbo4userus.NewUserEntry(invitedUserID)

	spaceID := invitedUser.Data.GetFamilySpaceID()

	if err = dal4userus.GetUser(ctx, db, invitedUser); err != nil {
		// TODO: Instead pass user as a parameter? Even better if the user entity was created within following transaction.
		return isJustLinked, err
	} else if invitedUser.Data.CreatedAt.After(time.Now().Add(-time.Second / 2)) {
		logus.Debugf(ctx, "A new user, will wait for half a seconds to cleanup previous transaction")
		time.Sleep(time.Second / 2)
	}
	var invitedContact dal4contactus.ContactEntry
	var invitedDebtusContact models4debtus.DebtusSpaceContactEntry
	attempt := 0
	err = db.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if attempt += 1; attempt > 1 {
			sleepPeriod := time.Duration(attempt) * time.Second
			logus.Warningf(ctx, "Transaction retry will sleep for %v, invitedContact.ContactID: %v", attempt, invitedContact.ID)
			time.Sleep(sleepPeriod)
		}
		changes := linker.changes
		if changes.receipt, changes.transfer, changes.inviter.user, _, changes.invited.user, _, err =
			getReceiptTransferAndUsers(tctx, tx, receiptID, invitedUserID); err != nil {
			return
		}
		if invitedContact.ID != "" { // This means we are attempting to retry failed transaction
			if err = workaroundReinsertContact(tctx, changes.receipt, invitedDebtusContact, changes); err != nil {
				return
			}
		}

		if isJustLinked, err = linker.linkUsersByReceiptWithinTransaction(ctx, tctx, tx); err != nil {
			return
		} else {
			invitedContact = changes.invited.contact
		}

		// Integrity checks
		{
			invitedDebtusContact.Data.MustMatchCounterparty(changes.inviter.debtusContact)
		}

		if entitiesToSave := changes.Changes.Records(); len(entitiesToSave) > 0 {
			if err = tx.SetMulti(ctx, entitiesToSave); err != nil {
				return
			}
		} else {
			logus.Debugf(ctx, "ReceiptEntry and transfer has not changed")
		}
		return
	}, dal.TxWithCrossGroup())
	if err != nil {
		return
	}
	logus.Debugf(ctx, "ReceiptUsersLinker.LinkReceiptUsers() => invitedContact: %+v", invitedContact)
	if invitedDebtusContact, err = GetDebtusSpaceContactByID(ctx, nil, spaceID, invitedContact.ID); err != nil {
		return
	}
	logus.Debugf(ctx, "ReceiptUsersLinker.LinkReceiptUsers() => invitedContact from DB: %+v", invitedContact)
	return
}

func (linker *ReceiptUsersLinker) linkUsersByReceiptWithinTransaction(
	ctx context.Context, // non-transactional context
	tctx context.Context, // transactional context,
	tx dal.ReadwriteTransaction,
) (
	isCounterpartiesJustConnected bool,
	err error,
) {
	changes := linker.changes
	receipt := changes.receipt
	transfer := changes.transfer
	inviterUser, invitedUser := changes.inviter.user, changes.invited.user
	inviterContact, invitedContact := changes.inviter.contact, changes.invited.contact
	//var inviterDebtusContact models4debtus.DebtusSpaceContactEntry
	//if changes.inviterDebtusContact != nil {
	//	inviterDebtusContact = *changes.inviterDebtusContact
	//}
	//if changes.invitedDebtusContact != nil {
	//	invitedDebtusContact = *changes.invitedDebtusContact
	//}

	logus.Debugf(ctx,
		"ReceiptUsersLinker.linkUsersByReceiptWithinTransaction(receipt.ContactID=%s, transfer.ContactID=%s, inviterUser.ContactID=%s, invitedUser.ContactID=%s, inviterContact.ContactID=%s, invitedContact.ContactID=%s)",
		receipt.ID, transfer.ID, inviterUser.ID, invitedUser.ID, inviterContact.ID, invitedContact.ID)

	{ // validate inputs
		if err = linker.validateInput(changes); err != nil {
			return
		}
		if receipt.Data.TransferID != transfer.ID {
			panic(fmt.Sprintf("receipt.TransferID != transfer.ContactID: %v != %v", receipt.Data.TransferID, transfer.ID))
		}
		if transferCreatorUserID := transfer.Data.Creator().UserID; transferCreatorUserID == "" {
			panic("transfer.Creator().UserID is zero")
		} else if transferCreatorUserID != inviterUser.ID {
			panic(fmt.Sprintf("transfer.Creator().UserID != inviterUser.ContactID:  %v != %v", transferCreatorUserID, inviterUser.ID))
		} else if transferCreatorUserID == invitedUser.ID {
			panic(fmt.Sprintf("transfer.Creator().UserID == invitedUser.ContactID:  %v != %v", transferCreatorUserID, invitedUser.ID))
		}
	}

	fromOriginal := *transfer.Data.From()
	toOriginal := *transfer.Data.To()
	//logus.Debugf(ctx, "transferEntity: %v", transfer.Data)
	//logus.Debugf(ctx, "transfer.From(): %v", fromOriginal)
	//logus.Debugf(ctx, "transfer.To(): %v",toOriginal)

	transferCreatorCounterparty := transfer.Data.Counterparty()

	spaceID := invitedUser.Data.GetFamilySpaceID()
	if _, err = GetDebtusSpaceContactByID(tctx, tx, spaceID, transferCreatorCounterparty.ContactID); err != nil {
		return
	} else if inviterContact.Data.UserID != inviterUser.ID {
		panic(fmt.Errorf("inviterContact.UserID !=  inviterUser.ContactID: %v != %v", inviterContact.Data.UserID, inviterUser.ID))
	} else {
		changes.inviter.contact = inviterContact
	}

	if err = newUsersLinker(changes.usersLinkingDbChanges).linkUsersWithinTransaction(tctx, tx, receipt.Record.Key().String()); err != nil {
		err = fmt.Errorf("failed to link users: %w", err)
		return
	} else {
		invitedContact = changes.invited.contact // as was updated
	}
	{ // Update invited user's last currency
		var invitedUserChanged bool
		if invitedUser.Data.LastCurrencies, invitedUserChanged = slice.Merge(invitedUser.Data.LastCurrencies, []money.CurrencyCode{transfer.Data.Currency}); invitedUserChanged {
			changes.FlagAsChanged(changes.invited.user.Record)
		}
	}

	logus.Debugf(ctx, "linkUsersWithinTransaction() => invitedContact.ContactID: %v, inviterContact.ContactID: %v", invitedContact.ID, inviterContact.ID)

	// Update entities
	{
		if err = linker.updateReceipt(); err != nil {
			return
		} else if err = linker.updateTransfer(); err != nil {
			return
		} else if linker.changes.IsChanged(linker.changes.transfer.Record) {
			logus.Debugf(ctx, "transfer changed:\n\tFrom(): %v\n\tTo(): %v", transfer.Data.From(), transfer.Data.To())
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
		if err = dtdal.Reminder.DelayCreateReminderForTransferUser(tctx, receipt.Data.TransferID, transfer.Data.Counterparty().UserID); err != nil {
			err = fmt.Errorf("failed to delay creation of reminder for transfer coutnerparty: %w", err)
			return
		}
	} else {
		if transfer.Data.DtDueOn.IsZero() {
			logus.Debugf(tctx, "No need to create reminder for counterparty as no due date")
		} else {
			logus.Debugf(tctx, "No need to create reminder for counterparty as due date in past")
		}
	}
	return
}

func (linker *ReceiptUsersLinker) validateInput(changes *receiptDbChanges) error {

	if changes.receipt.Data.CounterpartyUserID != "" {
		if changes.receipt.Data.CounterpartyUserID != changes.invited.user.ID { // Already linked
			return errors.New("an attempt to link 3d user to a receipt")
		}

		transferCounterparty := changes.transfer.Data.Counterparty()

		if transferCounterparty.UserID != "" && transferCounterparty.UserID != changes.invited.user.ID {
			return fmt.Errorf(
				"transferCounterparty.UserID != invitedUser.ContactID : %s != %s",
				transferCounterparty.UserID, changes.invited.user.ID,
			)
		}
	}
	return nil
}

func (linker *ReceiptUsersLinker) updateReceipt() (err error) {
	receipt := linker.changes.receipt
	counterpartyUser := linker.changes.invited.user
	if receipt.Data.CounterpartyUserID != counterpartyUser.ID {
		receipt.Data.CounterpartyUserID = counterpartyUser.ID
		linker.changes.FlagAsChanged(linker.changes.receipt.Record)
	}
	return
}

func (linker *ReceiptUsersLinker) updateTransfer() (err error) {
	changes := linker.changes
	inviter, invited := changes.inviter, changes.invited
	transfer := changes.transfer
	{ // Validate input parameters
		if transfer.ID == "" || transfer.Data == nil {
			panic(fmt.Sprintf("Invalid parameter: transfer: %v", transfer))
		}
		validateSide := func(side string, user dbo4userus.UserEntry, contact dal4contactus.ContactEntry, debtusContact models4debtus.DebtusSpaceContactEntry) {
			if user.ID == "" || user.Data == nil {
				panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => %vUser: %v", side, user))
			}
			if contact.ID == "" || contact.Data == nil {
				panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => %vContact: %v", side, contact))
			} else if contact.Data.UserID != user.ID {
				panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => %vContact.UserID != %vUser.ContactID: %v != %v", side, side, contact.Data.UserID, invited.user.ID))
			}
		}
		validateSide("inviter", inviter.user, inviter.contact, inviter.debtusContact)
		validateSide("invited", invited.user, invited.contact, invited.debtusContact)
		if transfer.Data.CreatorUserID != inviter.user.ID {
			panic(fmt.Sprintf("ReceiptUsersLinker.updateTransfer() => transfer.CreatorUserID != inviterUser.ContactID: %v != %v", transfer.Data.CreatorUserID, invited.user.ID))
		}
	}

	transferCounterparty := transfer.Data.Counterparty()

	if transferCounterparty.UserID != invited.user.ID {
		if transferCounterparty.UserID != "" {
			err = fmt.Errorf("transfer.DebtusSpaceContactEntry().UserID != counterpartyUserID : %s != %s",
				transfer.Data.Counterparty().UserID, invited.user.ID)
			return
		}
		transfer.Data.Counterparty().UserID = invited.user.ID
		linker.changes.FlagAsChanged(linker.changes.transfer.Record)
	}

	updateTransferCounterpartyInfo := func(
		side string,
		counterparty *models4debtus.TransferCounterpartyInfo,
		user dbo4userus.UserEntry,
		contact dal4contactus.ContactEntry,
	) {
		if contact.Data.UserID == user.ID {
			panic(fmt.Sprintf(
				"updateTransferCounterpartyInfo() => %sContact.UserID == %sUser.ContactID : %s, counterparty.UserID: %s",
				side, side, contact.Data.UserID, counterparty.UserID))
		}
		if counterparty.UserID == "" {
			counterparty.UserID = user.ID
		} else if counterparty.UserID != user.ID {
			panic(fmt.Sprintf("updateTransferCounterpartyInfo() => counterparty.UserID != %sUser.ContactID : %s != %s, %sContact.UserID: %s", side, counterparty.UserID, user.ID, side, contact.Data.UserID))
		}
		counterparty.UserName = user.Data.GetFullName()

		if counterparty.ContactID == "" {
			counterparty.ContactID = contact.ID
		} else if counterparty.ContactID != contact.ID {
			panic(fmt.Sprintf(
				"ReceiptUsersLinker.updateTransfer() => counterparty.ContactID != %sContact.ContactID : %s != %s",
				side, counterparty.ContactID, contact.ID))
		}
		counterparty.ContactName = contact.Data.Names.GetFullName()
	}

	updateTransferCounterpartyInfo("inviter", transfer.Data.Creator(), inviter.user, invited.contact)
	updateTransferCounterpartyInfo("invited", transfer.Data.Counterparty(), invited.user, inviter.contact)

	//if inlineMessageID != "" {
	//	transfer.CounterpartyTgReceiptInlineMessageID = inlineMessageID
	//}
	//transferAmount := transfer.Data.GetAmount()
	if transfer.Data.Direction() == models4debtus.TransferDirectionUser2Counterparty {
		transfer.Data.AmountInCents *= -1
	}

	return
}
