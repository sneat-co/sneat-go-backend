package facade

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"time"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

type usersLinkingDbChanges struct {
	// use pointer as we pass it to FlagAsChanged() and IsChanged()
	dal.Changes
	inviterUser, invitedUser       *models.AppUser
	inviterContact, invitedContact *models.Contact
}

func newUsersLinkingDbChanges() *usersLinkingDbChanges {
	return &usersLinkingDbChanges{}
}

type receiptDbChanges struct {
	// use pointer as we pass it to FlagAsChanged() and IsChanged()
	*usersLinkingDbChanges
	receipt  *models.Receipt
	transfer *models.Transfer
}

func newReceiptDbChanges() *receiptDbChanges {
	return &receiptDbChanges{
		usersLinkingDbChanges: newUsersLinkingDbChanges(),
	}
}

func workaroundReinsertContact(c context.Context, receipt models.Receipt, invitedContact models.Contact, changes *receiptDbChanges) (err error) {
	if _, err = GetContactByID(c, nil, invitedContact.ID); err != nil {
		if dal.IsNotFound(err) {
			log.Warningf(c, "workaroundReinsertContact(invitedContact.ID=%v) => %v", invitedContact.ID, err.Error())
			err = nil
			if receipt.Data.Status == models.ReceiptStatusAcknowledged {
				if invitedContactInfo := changes.invitedUser.Data.ContactByID(invitedContact.ID); invitedContactInfo != nil {
					log.Warningf(c, "Transactional retry, contact was not created in DB but invitedUser already has the contact info & receipt is acknowledged")
					changes.invitedContact = &invitedContact
				} else {
					log.Warningf(c, "Transactional retry, contact was not created in DB but receipt is acknowledged & invitedUser has not contact info in JSON")
				}
			}
			changes.FlagAsChanged(changes.invitedContact.Record)
		} else {
			log.Errorf(c, "workaroundReinsertContact(invitedContact.ID=%v) => %v", invitedContact.ID, err.Error())
		}
	} else {
		log.Debugf(c, "workaroundReinsertContact(%v) => contact found by ID!", invitedContact.ID)
	}
	return
}

func AcknowledgeReceipt(c context.Context, receiptID, currentUserID string, operation string) (
	receipt models.Receipt, transfer models.Transfer, isCounterpartiesJustConnected bool, err error,
) {
	log.Debugf(c, "AcknowledgeReceipt(receiptID=%d, currentUserID=%d, operation=%v)", receiptID, currentUserID, operation)
	var transferAckStatus string
	switch operation {
	case dtdal.AckAccept:
		transferAckStatus = models.TransferAccepted
	case dtdal.AckDecline:
		transferAckStatus = models.TransferDeclined
	default:
		err = ErrInvalidAcknowledgeType
		return
	}

	var invitedContact models.Contact

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}

	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
		var inviterUser, invitedUser models.AppUser
		var inviterContact models.Contact

		receipt, transfer, inviterUser, invitedUser, err = getReceiptTransferAndUsers(tc, tx, receiptID, currentUserID)
		if err != nil {
			return
		}

		if transfer.Data.CreatorUserID == currentUserID {
			log.Errorf(tc, "An attempt to claim receipt on self created transfer")
			err = ErrSelfAcknowledgement
			return
		}

		changes := &receiptDbChanges{
			receipt:  &receipt,
			transfer: &transfer,
			usersLinkingDbChanges: &usersLinkingDbChanges{
				inviterUser: &inviterUser,
				invitedUser: &invitedUser,
			},
		}

		if invitedContact.ID != "" { // This means we are attempting to retry failed transaction
			if err = workaroundReinsertContact(tc, receipt, invitedContact, changes); err != nil {
				return
			}
		}

		{ // data integrity checks
			for _, counterpartyTgUserID := range invitedUser.Data.GetTelegramUserIDs() {
				for _, creatorTgUserID := range inviterUser.Data.GetTelegramUserIDs() {
					if counterpartyTgUserID == creatorTgUserID {
						return fmt.Errorf("data integrity issue: counterpartyTgUserID == creatorTgUserID (%v)", counterpartyTgUserID)
					}
				}
			}
		}

		if receipt.Data.Status == models.ReceiptStatusAcknowledged {
			if receipt.Data.AcknowledgedByUserID != currentUserID {
				err = fmt.Errorf("receipt.AcknowledgedByUserID != currentUserID (%s != %s)", receipt.Data.AcknowledgedByUserID, currentUserID)
				return
			}
			log.Debugf(c, "Receipt is already acknowledged")
		} else {
			receipt.Data.DtAcknowledged = time.Now()
			receipt.Data.Status = models.ReceiptStatusAcknowledged
			receipt.Data.AcknowledgedByUserID = currentUserID
			markReceiptAsViewed(receipt.Data, currentUserID)
			changes.FlagAsChanged(changes.receipt.Record)

			transfer.Data.AcknowledgeStatus = transferAckStatus
			transfer.Data.AcknowledgeTime = receipt.Data.DtAcknowledged
			changes.FlagAsChanged(changes.transfer.Record)
		}

		if transfer.Data.Counterparty().UserID == "" {
			if isCounterpartiesJustConnected, err = NewReceiptUsersLinker(changes).linkUsersByReceiptWithinTransaction(c, tc, tx); err != nil {
				return
			}
			invitedContact = *changes.invitedContact
			inviterContact = *changes.inviterContact
			log.Debugf(c, "linkUsersByReceiptWithinTransaction() =>\n\tinvitedContact %v: %+v\n\tinviterContact %v: %v",
				invitedContact.ID, invitedContact.Data, inviterContact.ID, inviterContact.Data)
		} else {
			log.Debugf(c, "No need to link users as already linked")
			inviterContact.ID = transfer.Data.CounterpartyInfoByUserID(inviterUser.ID).ContactID
			invitedContact.ID = transfer.Data.CounterpartyInfoByUserID(invitedUser.ID).ContactID
		}

		inviterUser.Data.CountOfAckTransfersByCounterparties += 1
		invitedUser.Data.CountOfAckTransfersByUser += 1

		if recordsToSave := changes.Records(); len(recordsToSave) > 0 {
			log.Debugf(c, "%v entities to save: %+v", len(recordsToSave), recordsToSave)
			if err = tx.SetMulti(c, recordsToSave); err != nil {
				return
			}
		} else {
			log.Debugf(c, "Nothing to save")
		}

		//if _, err = GetContactByID(c, invitedContact.ID); err != nil {
		//	if dal.IsNotFound(err) {
		//		log.Errorf(c, "Invited contact is not found by ID, let's try to re-insert.")
		//		if err = facade.SaveContact(c, invitedContact); err != nil {
		//			return
		//		}
		//	} else {
		//		return
		//	}
		//}
		return
	}, dal.TxWithCrossGroup())

	if err != nil {
		if err == ErrSelfAcknowledgement {
			err = nil
			return
		}
		err = fmt.Errorf("failed to acknowledge receipt: %w", err)
		return
	}
	log.Infof(c, "Receipt successfully acknowledged")

	{ // verify invitedContact
		if invitedContact, err = GetContactByID(c, nil, invitedContact.ID); err != nil {
			err = fmt.Errorf("failed to load invited contact outside of transaction: %w", err)
			if dal.IsNotFound(err) {
				return
			}
			log.Errorf(c, err.Error())
			err = nil // We are OK to ignore technical issues here
			return
		}
	}
	return
}

func MarkReceiptAsViewed(c context.Context, receiptID, userID string) (receipt models.Receipt, err error) {
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		receipt, err = dtdal.Receipt.GetReceiptByID(tc, tx, receiptID)
		if err != nil {
			return err
		}
		changed := markReceiptAsViewed(receipt.Data, userID)

		if receipt.Data.DtViewed.IsZero() {
			receipt.Data.DtViewed = time.Now()
			changed = true
		}
		if changed {
			if err = dtdal.Receipt.UpdateReceipt(c, tx, receipt); err != nil {
				return err
			}
		}
		return err
	}, dal.TxWithCrossGroup())
	return
}

func markReceiptAsViewed(receipt *models.ReceiptData, userID string) (changed bool) {
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

func getReceiptTransferAndUsers(c context.Context, tx dal.ReadSession, receiptID, userID string) (
	receipt models.Receipt,
	transfer models.Transfer,
	creatorUser models.AppUser,
	counterpartyUser models.AppUser,
	err error,
) {
	log.Debugf(c, "getReceiptTransferAndUsers(receiptID=%v, userID=%v)", receiptID, userID)

	if receipt, err = dtdal.Receipt.GetReceiptByID(c, tx, receiptID); err != nil {
		return
	}

	if transfer, err = Transfers.GetTransferByID(c, tx, receipt.Data.TransferID); err != nil {
		return
	}

	if receipt.Data.CreatorUserID != transfer.Data.CreatorUserID {
		err = errors.New("data integrity issue: receipt.CreatorUserID != transfer.CreatorUserID")
		return
	}

	if creatorUser, err = User.GetUserByID(c, tx, transfer.Data.CreatorUserID); err != nil {
		return
	}

	if counterpartyUser.ID = transfer.Data.Counterparty().UserID; counterpartyUser.ID == "" && userID != creatorUser.ID {
		counterpartyUser.ID = userID
	}

	if counterpartyUser.ID != "" {
		if counterpartyUser, err = User.GetUserByID(c, tx, counterpartyUser.ID); err != nil {
			return
		}
	}

	log.Debugf(c, "getReceiptTransferAndUsers(receiptID=%v, userID=%v) =>\n\tcreatorUser(id=%v): %+v\n\tcounterpartyUser(id=%v): %+v",
		receiptID, userID,
		creatorUser.ID, creatorUser.Data,
		counterpartyUser.ID, counterpartyUser.Data,
	)

	if creatorUser.Data == nil {
		err = fmt.Errorf("creatorUser(id=%v) == nil - data integrity or app logic issue", transfer.Data.CreatorUserID)
		return
	}
	return
}
