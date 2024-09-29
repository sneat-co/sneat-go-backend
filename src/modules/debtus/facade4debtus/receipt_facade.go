package facade4debtus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/anybot/facade4anybot"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"time"
)

type userLinkingParty struct {
	spaceID        string
	contact        dal4contactus.ContactEntry
	contactusSpace dal4contactus.ContactusSpaceEntry
	debtusSpace    models4debtus.DebtusSpaceEntry
	debtusContact  models4debtus.DebtusSpaceContactEntry
	debtusUser     models4debtus.DebtusUserEntry // TODO: DO we need this?
	user           dbo4userus.UserEntry          // TODO: DO we need this? Would debtusUser be enough?
}

type usersLinkingDbChanges struct { // use as a pointer as we pass it to FlagAsChanged() and IsChanged()
	dal.Changes
	inviter *userLinkingParty
	invited *userLinkingParty
}

func newUsersLinkingDbChanges() *usersLinkingDbChanges {
	return &usersLinkingDbChanges{}
}

type receiptDbChanges struct { // use as a pointer as we pass it to FlagAsChanged() and IsChanged()
	*usersLinkingDbChanges
	receipt  models4debtus.ReceiptEntry
	transfer models4debtus.TransferEntry
}

func newReceiptDbChanges() *receiptDbChanges {
	return &receiptDbChanges{
		usersLinkingDbChanges: newUsersLinkingDbChanges(),
	}
}

func workaroundReinsertContact(ctx context.Context, receipt models4debtus.ReceiptEntry, invitedContact models4debtus.DebtusSpaceContactEntry, changes *receiptDbChanges) (err error) {
	if _, err = GetDebtusSpaceContactByID(ctx, nil, receipt.Data.SpaceID, invitedContact.ID); err != nil {
		if dal.IsNotFound(err) {
			logus.Warningf(ctx, "workaroundReinsertContact(invitedContact.ContactID=%s) => %v", invitedContact.ID, err)
			err = nil
			if receipt.Data.Status == models4debtus.ReceiptStatusAcknowledged {
				if invitedContactInfo := changes.inviter.contactusSpace.Data.GetContactBriefByContactID(invitedContact.ID); invitedContactInfo != nil {
					logus.Warningf(ctx, "Transactional retry, Contact was not created in DB but invitedUser already has the Contact info & receipt is acknowledged")
					changes.invited.debtusContact = invitedContact
				} else {
					logus.Warningf(ctx, "Transactional retry, Contact was not created in DB but receipt is acknowledged & invitedUser has not Contact info in JSON")
				}
			}
			changes.FlagAsChanged(changes.invited.contact.Record)
		} else {
			logus.Errorf(ctx, "workaroundReinsertContact(invitedContact.ContactID=%s) => %v", invitedContact.ID, err)
		}
	} else {
		logus.Debugf(ctx, "workaroundReinsertContact(%s) => Contact found by ContactID!", invitedContact.ID)
	}
	return
}

func AcknowledgeReceipt(ctx context.Context, userCtx facade.UserContext, receiptID string, operation string) (
	receipt models4debtus.ReceiptEntry, transfer models4debtus.TransferEntry, isCounterpartiesJustConnected bool, err error,
) {
	currentUserID := userCtx.GetUserID()
	logus.Debugf(ctx, "AcknowledgeReceipt(receiptID=%s, currentUserID=%s, operation=%s)", receiptID, currentUserID, operation)

	var transferAckStatus string
	switch operation {
	case dtdal.AckAccept:
		transferAckStatus = models4debtus.TransferAccepted
	case dtdal.AckDecline:
		transferAckStatus = models4debtus.TransferDeclined
	default:
		err = facade4anybot.ErrInvalidAcknowledgeType
		return
	}

	var invitedContact dal4contactus.ContactEntry

	err = facade.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {

		var inviterUser, invitedUser dbo4userus.UserEntry
		var inviterDebtusUser, invitedDebtusUser models4debtus.DebtusUserEntry
		var inviterContact dal4contactus.ContactEntry

		receipt, transfer, inviterUser, inviterDebtusUser, invitedUser, invitedDebtusUser, err = getReceiptTransferAndUsers(tctx, tx, receiptID, currentUserID)
		if err != nil {
			return
		}

		spaceID := invitedUser.Data.GetFamilySpaceID()
		var invitedDebtusSpace models4debtus.DebtusSpaceEntry
		if spaceID != "" {
			invitedDebtusSpace = models4debtus.NewDebtusSpaceEntry(spaceID)
			if err = tx.Get(tctx, invitedDebtusSpace.Record); err != nil {
				return
			}
		}

		if transfer.Data.CreatorUserID == currentUserID {
			logus.Errorf(tctx, "An attempt to claim receipt on self created transfer")
			err = facade4anybot.ErrSelfAcknowledgement
			return
		}

		inviterSpaceID := inviterUser.Data.GetFamilySpaceID()
		invitedSpaceID := invitedUser.Data.GetFamilySpaceID()

		changes := &receiptDbChanges{
			receipt:  receipt,
			transfer: transfer,
			usersLinkingDbChanges: &usersLinkingDbChanges{
				inviter: &userLinkingParty{
					contact:        inviterContact,
					contactusSpace: dal4contactus.NewContactusSpaceEntry(inviterSpaceID),
					debtusSpace:    models4debtus.NewDebtusSpaceEntry(inviterSpaceID),
					debtusContact:  models4debtus.NewDebtusSpaceContactEntry(inviterSpaceID, inviterUser.ID, nil),
					debtusUser:     inviterDebtusUser,
				},
				invited: &userLinkingParty{
					contact:        invitedContact,
					contactusSpace: dal4contactus.NewContactusSpaceEntry(invitedSpaceID),
					debtusSpace:    invitedDebtusSpace,
					debtusContact:  models4debtus.NewDebtusSpaceContactEntry(invitedSpaceID, invitedUser.ID, nil),
					debtusUser:     invitedDebtusUser,
				},
			},
		}

		if invitedContact.ID != "" { // This means we are attempting to retry failed transaction
			if err = workaroundReinsertContact(tctx, receipt, changes.invited.debtusContact, changes); err != nil {
				return
			}
		}

		{ // data integrity checks
			for _, counterpartyTgUserID := range invitedUser.Data.GetTelegramUserIDs() {
				for _, creatorTgUserID := range inviterUser.Data.GetTelegramUserIDs() {
					if counterpartyTgUserID == creatorTgUserID {
						return fmt.Errorf("data integrity issue: counterpartyTgUserID == creatorTgUserID (%d)", counterpartyTgUserID)
					}
				}
			}
		}

		if receipt.Data.Status == models4debtus.ReceiptStatusAcknowledged {
			if receipt.Data.AcknowledgedByUserID != currentUserID {
				err = fmt.Errorf("receipt.AcknowledgedByUserID != currentUserID (%s != %s)", receipt.Data.AcknowledgedByUserID, currentUserID)
				return
			}
			logus.Debugf(ctx, "ReceiptEntry is already acknowledged")
		} else {
			receipt.Data.DtAcknowledged = time.Now()
			receipt.Data.Status = models4debtus.ReceiptStatusAcknowledged
			receipt.Data.AcknowledgedByUserID = currentUserID
			markReceiptAsViewed(receipt.Data, currentUserID)
			changes.FlagAsChanged(changes.receipt.Record)

			transfer.Data.AcknowledgeStatus = transferAckStatus
			transfer.Data.AcknowledgeTime = receipt.Data.DtAcknowledged
			changes.FlagAsChanged(changes.transfer.Record)
		}

		if transfer.Data.Counterparty().UserID == "" {
			if isCounterpartiesJustConnected, err = NewReceiptUsersLinker(changes).linkUsersByReceiptWithinTransaction(ctx, tctx, tx); err != nil {
				return
			}
			invitedContact = changes.invited.contact
			inviterContact = changes.inviter.contact
			//logus.Debugf(ctx, "linkUsersByReceiptWithinTransaction() =>\n\tinvitedContact %s: %+v\n\tinviterContact %s: %v",
			//	invitedContact.ContactID, invitedContact.Data, inviterContact.ContactID, inviterContact.Data)
		} else {
			logus.Debugf(ctx, "No need to link users as already linked")
			inviterContact.ID = transfer.Data.CounterpartyInfoByUserID(inviterDebtusUser.ID).ContactID
			invitedContact.ID = transfer.Data.CounterpartyInfoByUserID(invitedDebtusUser.ID).ContactID
		}

		inviterDebtusUser.Data.CountOfAckTransfersByCounterparties += 1
		invitedDebtusUser.Data.CountOfAckTransfersByUser += 1

		if recordsToSave := changes.Records(); len(recordsToSave) > 0 {
			//logus.Debugf(ctx, "%d entities to save: %+v", len(recordsToSave), recordsToSave)
			if err = tx.SetMulti(ctx, recordsToSave); err != nil {
				return
			}
		} else {
			logus.Debugf(ctx, "Nothing to save")
		}

		//if _, err = GetDebtusSpaceContactByID(ctx, invitedContact.ContactID); err != nil {
		//	if dal.IsNotFound(err) {
		//		logus.Errorf(ctx, "Invited Contact is not found by ContactID, let's try to re-insert.")
		//		if err = facade4debtus.SaveContact(ctx, invitedContact); err != nil {
		//			return
		//		}
		//	} else {
		//		return
		//	}
		//}
		return
	}, dal.TxWithCrossGroup())

	if err != nil {
		if errors.Is(err, facade4anybot.ErrSelfAcknowledgement) {
			err = nil
			return
		}
		err = fmt.Errorf("failed to acknowledge receipt: %w", err)
		return
	}
	logus.Infof(ctx, "ReceiptEntry successfully acknowledged")

	{ // verify invitedContact
		if _, err = GetDebtusSpaceContactByID(ctx, nil, receipt.Data.SpaceID, invitedContact.ID); err != nil {
			err = fmt.Errorf("failed to load invited Contact outside of transaction: %w", err)
			if dal.IsNotFound(err) {
				return
			}
			logus.Errorf(ctx, err.Error())
			err = nil // We are OK to ignore technical issues here
			return
		}
	}
	return
}

func MarkReceiptAsViewed(ctx context.Context, receiptID, userID string) (receipt models4debtus.ReceiptEntry, err error) {
	err = facade.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) error {
		receipt, err = dtdal.Receipt.GetReceiptByID(tctx, tx, receiptID)
		if err != nil {
			return err
		}
		changed := markReceiptAsViewed(receipt.Data, userID)

		if receipt.Data.DtViewed.IsZero() {
			receipt.Data.DtViewed = time.Now()
			changed = true
		}
		if changed {
			if err = dtdal.Receipt.UpdateReceipt(ctx, tx, receipt); err != nil {
				return err
			}
		}
		return err
	}, dal.TxWithCrossGroup())
	return
}

func markReceiptAsViewed(receipt *models4debtus.ReceiptDbo, userID string) (changed bool) {
	alreadyViewedByUser := false
	for _, uid := range receipt.ViewedByUserIDs {
		if uid == userID {
			alreadyViewedByUser = true
			break
		}
	}
	if !alreadyViewedByUser {
		receipt.ViewedByUserIDs = append(receipt.ViewedByUserIDs, userID)
		changed = true
	}
	return
}

func getReceiptTransferAndUsers(ctx context.Context, tx dal.ReadSession, receiptID, userID string) (
	receipt models4debtus.ReceiptEntry,
	transfer models4debtus.TransferEntry,
	creatorUser dbo4userus.UserEntry, creatorDebtusUser models4debtus.DebtusUserEntry,
	counterpartyUser dbo4userus.UserEntry, counterpartyDebtusUser models4debtus.DebtusUserEntry,
	err error,
) {
	logus.Debugf(ctx, "getReceiptTransferAndUsers(receiptID=%s, userID=%s)", receiptID, userID)

	if receipt, err = dtdal.Receipt.GetReceiptByID(ctx, tx, receiptID); err != nil {
		return
	}

	if transfer, err = Transfers.GetTransferByID(ctx, tx, receipt.Data.TransferID); err != nil {
		return
	}

	if receipt.Data.CreatorUserID != transfer.Data.CreatorUserID {
		err = errors.New("data integrity issue: receipt.CreatorUserID != transfer.CreatorUserID")
		return
	}

	creatorUser = dbo4userus.NewUserEntry(receipt.Data.CreatorUserID)
	creatorDebtusUser = models4debtus.NewDebtusUserEntry(receipt.Data.CreatorUserID)

	var recordsToGet []dal.Record = []dal.Record{creatorUser.Record, creatorDebtusUser.Record}

	if counterpartyUserID := transfer.Data.Counterparty().UserID; counterpartyUserID != "" {
		counterpartyUser = dbo4userus.NewUserEntry(counterpartyUserID)
		counterpartyDebtusUser = models4debtus.NewDebtusUserEntry(counterpartyUserID)
		recordsToGet = append(recordsToGet, counterpartyUser.Record, counterpartyDebtusUser.Record)
	}

	if err = tx.GetMulti(ctx, recordsToGet); err != nil {
		return
	}

	logus.Debugf(ctx, "getReceiptTransferAndUsers(receiptID=%s, userID=%s) =>\n\tcreatorDebtusUser(id=%s): %+v\n\tcounterpartyDebtusUser(id=%s): %+v",
		receiptID, userID,
		creatorDebtusUser.ID, creatorDebtusUser.Data,
		counterpartyDebtusUser.ID, counterpartyDebtusUser.Data,
	)

	if creatorDebtusUser.Data == nil {
		err = fmt.Errorf("creatorDebtusUser(id=%s) == nil - data integrity or app logic issue", transfer.Data.CreatorUserID)
		return
	}
	return
}
