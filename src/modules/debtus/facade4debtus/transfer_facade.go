package facade4debtus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sanity-io/litter"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
	"github.com/strongo/logus"
	"github.com/strongo/slice"
	"slices"
	"time"
)

const (
	userBalanceIncreased = 1
	userBalanceDecreased = -1
)

var (
	ErrNotImplemented                      = errors.New("not implemented yet")
	ErrDebtAlreadyReturned                 = errors.New("this debt already has been returned")
	ErrPartialReturnGreaterThenOutstanding = errors.New("an attempt to do partial return for amount greater then outstanding")
	//
	ErrNoOutstandingTransfers                                       = errors.New("no outstanding api4transfers")
	ErrAttemptToCreateDebtWithInterestAffectingOutstandingTransfers = errors.New("an attempt to create a debt with interest that will affect outstanding api4transfers")
)

func TransferCounterparties(direction models4debtus.TransferDirection, creatorInfo models4debtus.TransferCounterpartyInfo) (from, to *models4debtus.TransferCounterpartyInfo) {
	creator := models4debtus.TransferCounterpartyInfo{
		UserID:  creatorInfo.UserID,
		Comment: creatorInfo.Comment,
	}
	counterparty := models4debtus.TransferCounterpartyInfo{
		ContactID:   creatorInfo.ContactID,
		ContactName: creatorInfo.ContactName,
	}
	switch direction {
	case models4debtus.TransferDirectionUser2Counterparty:
		return &creator, &counterparty
	case models4debtus.TransferDirectionCounterparty2User:
		return &counterparty, &creator
	default:
		panic("Unknown direction: " + string(direction))
	}
}

var Transfers = TransfersFacade{}

type TransfersFacade struct {
}

func (TransfersFacade) SaveTransfer(ctx context.Context, tx dal.ReadwriteTransaction, transfer models4debtus.TransferEntry) error {
	return tx.Set(ctx, transfer.Record)
}

func (transferFacade TransfersFacade) CreateTransfer(ctx context.Context, input CreateTransferInput) (
	output CreateTransferOutput, err error,
) {
	now := time.Now()

	logus.Infof(ctx, "CreateTransfer(input=%v)", input)

	var returnToTransferIDs []string

	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return
	}

	//if input.Request.ReturnToTransferID != "" {
	//	if counterparty, err := GetDebtusSpaceContactByID(ctx, nil, contactID); err != nil {
	//		if dal.IsNotFound(err) {
	//			api4debtus.BadRequestError(ctx, w, err)
	//		} else {
	//			api4debtus.InternalError(ctx, w, err)
	//		}
	//		return
	//	} else {
	//		balance := counterparty.Data.Balance()
	//		if balanceAmount, ok := balance[amountWithCurrency.Currency]; !ok {
	//			api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("No balance for %v", amountWithCurrency.Currency))
	//		} else {
	//			switch direction {
	//			case models.TransferDirectionUser2Counterparty:
	//				if balanceAmount > 0 {
	//					api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("balanceAmount > 0 && direction == %v", direction))
	//				}
	//			case models.TransferDirectionCounterparty2User:
	//				if balanceAmount < 0 {
	//					api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("balanceAmount < 0 && direction == %v", direction))
	//				}
	//			}
	//		}
	//	}
	//}

	creatorContactusSpace := dal4contactus.NewContactusSpaceEntry(input.Request.SpaceID)
	if err = db.Get(ctx, creatorContactusSpace.Record); err != nil {
		if !dal.IsNotFound(err) {
			err = fmt.Errorf("failed to get creatorContactusSpace: %w", err)
			return
		}
	}

	creatorDebtusSpace := models4debtus.NewDebtusSpaceEntry(input.Request.SpaceID)
	if err = db.Get(ctx, creatorDebtusSpace.Record); err != nil && !dal.IsNotFound(err) {
		err = fmt.Errorf("faield to get debtusbot space entry: %w", err)
		return
	}

	if input.Request.ReturnToTransferID == "" {
		logus.Debugf(ctx, "input.ReturnToTransferID == 0")
		creatorContactID := input.CreatorContactID()
		if creatorContactID == "" {
			panic(fmt.Errorf("3d party api4transfers are not implemented yet: %w", err))
		}
		var creatorContact dal4contactus.ContactEntry
		var creatorDebtusContact models4debtus.DebtusSpaceContactEntry

		verifyUserDebtusContact := func() (contactBriefFound bool) {
			var debtusContactBrief *models4debtus.DebtusContactBrief
			if debtusContactBrief, contactBriefFound = creatorDebtusSpace.Data.Contacts[creatorContactID]; contactBriefFound {
				contactBalance := debtusContactBrief.Balance
				if v, ok := contactBalance[input.Request.Amount.Currency]; !ok || v == 0 {
					logus.Debugf(ctx, "No need to check for outstanding api4transfers as contacts balance is 0")
				} else {
					if input.Request.Interest.HasInterest() {
						if d := input.Direction(); d == models4debtus.TransferDirectionUser2Counterparty && v < 0 || d == models4debtus.TransferDirectionCounterparty2User && v > 0 {
							err = ErrAttemptToCreateDebtWithInterestAffectingOutstandingTransfers
							return
						}
					}
					if returnToTransferIDs, err = transferFacade.checkOutstandingTransfersForReturns(ctx, now, input); err != nil {
						return
					}
				}
				contactBriefFound = true
				return
			}
			return
		}
		if contactJsonFound := verifyUserDebtusContact(); contactJsonFound {
			if err != nil {
				return
			}
			goto contactFound
		}

		// If Contact not found in user's JSON try to recover from DB record
		if creatorContact, err = dal4contactus.GetContactByID(ctx, nil, input.Request.SpaceID, creatorContactID); err != nil {
			return
		}

		// If Contact not found in user's JSON try to recover from DB record
		if creatorDebtusContact, err = GetDebtusSpaceContactByID(ctx, nil, input.Request.SpaceID, creatorContactID); err != nil {
			return
		}

		if creatorContact.Data.UserID != input.CreatorUser.ID {
			err = fmt.Errorf("creatorContact.UserID != input.CreatorUser.ContactID: %v != %v", creatorContact.Data.UserID, input.CreatorUser.ID)
			return
		}

		_, _ = models4debtus.AddOrUpdateDebtusContact(creatorDebtusSpace, creatorDebtusContact)

		if contactJsonFound := verifyUserDebtusContact(); contactJsonFound {
			if err != nil {
				return
			}
			goto contactFound
		}
		if err == nil {
			err = fmt.Errorf("user Contact not found by creatorContactID=%s", creatorContactID)
		}
		return
	contactFound:
	} else if !input.Request.IsReturn {
		panic("ReturnToTransferID != 0 && !IsReturn")
	}

	if input.Request.ReturnToTransferID != "" {
		var transferToReturn models4debtus.TransferEntry
		if transferToReturn, err = Transfers.GetTransferByID(ctx, db, input.Request.ReturnToTransferID); err != nil {
			err = fmt.Errorf("failed to get returnToTransferID=%s: %w", input.Request.ReturnToTransferID, err)
			return
		}

		if transferToReturn.Data.Currency != input.Request.Amount.Currency {
			panic("transferToReturn.Currency != amount.Currency")
		}

		if transferToReturn.Data.GetOutstandingValue(now) == 0 {
			// When the transfer has been already returned
			err = ErrDebtAlreadyReturned
			return
		}

		if input.Request.Amount.Value > transferToReturn.Data.GetOutstandingValue(now) {
			logus.Debugf(ctx, "amount.Value:%v > transferToReturn.GetOutstandingValue(now):%v", input.Request.Amount.Value, transferToReturn.Data.GetOutstandingValue(now))
			if input.Request.Amount.Value == transferToReturn.Data.AmountInCents {
				// For situations when a transfer was partially returned but user wants to mark it as fully returned.
				logus.Debugf(ctx, "amount.Value (%v) == transferToReturn.AmountInCents (%v)", input.Request.Amount.Value, transferToReturn.Data.AmountInCents)
				input.Request.Amount.Value = transferToReturn.Data.GetOutstandingValue(now)
				logus.Debugf(ctx, "Updated amount.Value: %v", input.Request.Amount.Value)
			} else {
				err = ErrPartialReturnGreaterThenOutstanding
				return
			}
		} else if input.Request.Amount.Value < transferToReturn.Data.GetOutstandingValue(now) {
			logus.Debugf(ctx, "input.Amount.Value < transferToReturn.GetOutstandingValue(now)")
		}

		returnToTransferIDs = append(returnToTransferIDs, input.Request.ReturnToTransferID)
		output.ReturnedTransfers = append(output.ReturnedTransfers, transferToReturn)
	}

	if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		output, err = transferFacade.createTransferWithinTransaction(ctx, tx, now, input, returnToTransferIDs)
		return err
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}

	output.Validate()

	return
}

func (transferFacade TransfersFacade) checkOutstandingTransfersForReturns(ctx context.Context, now time.Time, input CreateTransferInput) (returnToTransferIDs []string, err error) {
	logus.Debugf(ctx, "TransfersFacade.checkOutstandingTransfersForReturns()")
	var (
		outstandingTransfers []models4debtus.TransferEntry
	)

	creatorUserID := input.CreatorUser.ID
	creatorContactID := input.CreatorContactID()

	reversedDirection := input.Direction().Reverse()

	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return
	}
	outstandingTransfers, err = dtdal.Transfer.LoadOutstandingTransfers(ctx, db, now, creatorUserID, creatorContactID, input.Request.Amount.Currency, reversedDirection)
	if err != nil {
		err = fmt.Errorf("failed to load outstanding api4transfers: %w", err)
		return
	}
	if input.Request.IsReturn && len(outstandingTransfers) == 0 {
		err = ErrNoOutstandingTransfers
		return
	}

	logus.Debugf(ctx, "facade4debtus.checkOutstandingTransfersForReturns() => dtdal.TransferEntry.LoadOutstandingTransfers(userID=%v, currency=%v) => %d api4transfers", input.CreatorUser.ID, input.Request.Amount.Currency, len(outstandingTransfers))

	if outstandingTransfersCount := len(outstandingTransfers); outstandingTransfersCount > 0 { // Assign the return to specific api4transfers
		var (
			assignedValue             decimal.Decimal64p2
			outstandingRightDirection int
		)
		buf := new(bytes.Buffer)
		_, _ = fmt.Fprintf(buf, "%v outstanding api4transfers\n", outstandingTransfersCount)
		for i, outstandingTransfer := range outstandingTransfers {
			_, _ = fmt.Fprintf(buf, "\t[%v]: %v", i, litter.Sdump(outstandingTransfer))
			outstandingTransferID := outstandingTransfers[i].ID
			outstandingValue := outstandingTransfer.Data.GetOutstandingValue(now)
			if outstandingValue == input.Request.Amount.Value { // A check for exact match that has higher priority then earlie api4transfers
				logus.Infof(ctx, " - found outstanding transfer %v with exact amount match: %v", outstandingTransfer.ID, outstandingValue)
				assignedValue = input.Request.Amount.Value
				returnToTransferIDs = []string{outstandingTransferID}
				break
			}
			if assignedValue < input.Request.Amount.Value { // Do not break so we check all outstanding api4transfers for exact match
				returnToTransferIDs = append(returnToTransferIDs, outstandingTransferID)
				assignedValue += outstandingValue
			}
			outstandingRightDirection += 1
			buf.WriteString("\n")
		}
		logus.Debugf(ctx, buf.String())
		if input.Request.IsReturn && assignedValue < input.Request.Amount.Value {
			logus.Warningf(ctx,
				"There are not enough outstanding api4transfers to return %v. All outstanding count: %v, Right direction: %v, Assigned amount: %v. Could be data integrity issue.",
				input.Request.Amount, len(outstandingTransfers), outstandingRightDirection, assignedValue,
			)
		}
	}
	return
}

func (transferFacade TransfersFacade) createTransferWithinTransaction(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	dtCreated time.Time,
	input CreateTransferInput,
	returnToTransferIDs []string,
) (
	output CreateTransferOutput, err error,
) {
	logus.Debugf(ctx, "createTransferWithinTransaction(input=%v, returnToTransferIDs=%v)", input, returnToTransferIDs)

	if err = input.Validate(); err != nil {
		return
	}
	// if len(returnToTransferIDs) > 0 && !input.IsReturn { // TODO: It's OK to have api4transfers without isReturn=true
	// 	panic("len(returnToTransferIDs) > 0 && !isReturn")
	// }

	output.From = new(ParticipantEntries)
	output.To = new(ParticipantEntries)
	from, to := input.From, input.To

	records := make([]dal.Record, 0, 4+len(returnToTransferIDs))
	if from.UserID != "" {
		output.From.User.ID = from.UserID
		records = append(records, output.From.User.Record)
	}
	if to.UserID != "" {
		output.To.User.ID = to.UserID
		records = append(records, output.To.User.Record)
	}
	if from.ContactID != "" {
		output.From.Contact.ID = from.ContactID
		records = append(records, output.From.Contact.Record)
	}
	if to.ContactID != "" {
		output.To.Contact.ID = to.ContactID
		records = append(records, output.To.Contact.Record)
	}

	if err = tx.GetMulti(ctx, records); err != nil {
		err = fmt.Errorf("failed to get user & counterparty records from datastore by keys: %w", err)
		return
	}
	fromContact, toContact := output.From.Contact, output.To.Contact
	fromUser, toUser := output.From.User, output.To.User

	if from.ContactID != "" && output.From.Contact.Data.UserID == "" {
		err = fmt.Errorf("got bad counterparty entity from DB by id=%s, fromCounterparty.UserID == 0", from.ContactID)
		return
	}

	if to.ContactID != "" && output.To.Contact.Data.UserID == "" {
		err = fmt.Errorf("got bad counterparty entity from DB by id=%s, toCounterparty.UserID == 0", to.ContactID)
		return
	}

	if to.ContactID != "" && from.ContactID != "" {
		if fromContact.Data.UserID != toContact.Data.UserID {
			err = fmt.Errorf("fromCounterparty.CounterpartyUserID != toCounterparty.UserID (%s != %s)",
				fromContact.Data.UserID, toContact.Data.UserID)
		}
		if toContact.Data.UserID != fromContact.Data.UserID {
			err = fmt.Errorf("toCounterparty.CounterpartyUserID != fromCounterparty.UserID (%s != %s)",
				toContact.Data.UserID, fromContact.Data.UserID)
		}
		return
	}

	// Check if counterparties are linked and if yes load the missing DebtusSpaceContactEntry
	{

		link := func(sideName, countersideName string, side, counterparty *models4debtus.TransferCounterpartyInfo, sideContact dal4contactus.ContactEntry, sideDebtusContact models4debtus.DebtusSpaceContactEntry,
		) (counterpartyContact dal4contactus.ContactEntry, counterpartyDebtusContact models4debtus.DebtusSpaceContactEntry, err error) {
			logus.Debugf(ctx, "link(%v=%v, %v=%v, %vContact=%v)", sideName, side, countersideName, counterparty, sideName, sideContact)
			if side.ContactID != "" && sideDebtusContact.Data.CounterpartyContactID != "" && counterparty.ContactID == "" {
				var spaceID string
				if counterpartyContact, err = dal4contactus.GetContactByID(ctx, tx, spaceID, sideDebtusContact.Data.CounterpartyContactID); err != nil {
					err = fmt.Errorf("failed to get counterparty by 'fromCounterparty.CounterpartyContactID': %w", err)
					return
				}
				if counterpartyDebtusContact, err = GetDebtusSpaceContactByID(ctx, tx, spaceID, sideDebtusContact.Data.CounterpartyContactID); err != nil {
					err = fmt.Errorf("failed to get counterparty by 'fromCounterparty.CounterpartyContactID': %w", err)
					return
				}
				counterparty.ContactID = counterpartyContact.ID
				counterparty.ContactName = counterpartyContact.Data.Names.GetFullName()
				side.UserID = counterpartyContact.Data.UserID
				records = append(records, counterpartyContact.Record)
			}
			return
		}

		var linkedContact dal4contactus.ContactEntry // TODO: This smells. TODO: explain why smells or remove comment
		if linkedContact, _, err = link("from", "to", from, to, fromContact, output.From.DebtusContact); err != nil {
			return
		} else if linkedContact.Data != nil {
			toContact = linkedContact
			output.To.Contact = linkedContact
		}

		logus.Debugf(ctx, "toContact: %v", toContact.Data == nil)
		if linkedContact, _, err = link("to", "from", to, from, toContact, output.To.DebtusContact); err != nil {
			return
		} else if linkedContact.Data != nil {
			fromContact = linkedContact
			output.From.Contact = fromContact
		}

		// //// When: toCounterparty == nil, fromUser == nil,
		// if from.ContactID != 0 && fromContact.CounterpartyContactID != 0 && to.ContactID == 0 {
		// 	// Get toCounterparty and fill to.DebtusSpaceContactEntry* fields
		// 	if toContact, err = GetDebtusSpaceContactByID(ctx, fromContact.CounterpartyContactID); err != nil {
		// 		err = fmt.Errorf("%w: Failed to get 'To' counterparty by 'fromCounterparty.CounterpartyContactID'", err)
		// 		return
		// 	}
		// 	output.To.DebtusSpaceContactEntry = toContact
		// 	logus.Debugf(ctx, "Got toContact id=%d: %v", toContact.ContactID, toContact.DebtusSpaceContactDbo)
		// 	to.ContactID = toContact.ContactID
		// 	to.ContactName = toContact.GetFullName()
		// 	from.UserID = toContact.UserID
		// 	records = append(records, &toContact)
		// }
		// if to.ContactID != 0 && toCounterparty.CounterpartyContactID != 0 && from.ContactID == 0 {
		// 	if fromCounterparty, err = GetDebtusSpaceContactByID(ctx, toCounterparty.CounterpartyContactID); err != nil {
		// 		err = fmt.Errorf("failed to get 'From' counterparty by 'toCounterparty.CounterpartyContactID' == %d: %w", fromCounterparty.CounterpartyContactID, err)
		// 		return
		// 	}
		// 	output.From.DebtusSpaceContactEntry = fromCounterparty
		// 	logus.Debugf(ctx, "Got fromCounterparty id=%d: %v", fromCounterparty.ContactID, fromCounterparty.DebtusSpaceContactDbo)
		// 	from.ContactID = fromCounterparty.ContactID
		// 	from.ContactName = fromCounterparty.GetFullName()
		// 	to.UserID = fromCounterparty.UserID
		// 	records = append(records, &fromCounterparty)
		// }
	}

	// In case if we just loaded above missing counterparty we need to check for missing user
	{
		loadUserIfNeeded := func(who string, appUser dbo4userus.UserEntry) (dbo4userus.UserEntry, error) {
			logus.Debugf(ctx, "%v.UserID: %s, %sUser.DebutsAppUserDataOBSOLETE: %+v", who, appUser.ID, who, appUser.Data)
			if appUser.ID != "" {
				if appUser.Data == nil {
					appUser = dbo4userus.NewUserEntry(appUser.ID)
					if err = dal4userus.GetUser(ctx, tx, appUser); err != nil {
						err = fmt.Errorf("failed to get %vUser for linked counterparty: %w", who, err)
						return appUser, err
					}
					records = append(records, appUser.Record)
				}
			}
			return appUser, err
		}

		if output.From.User, err = loadUserIfNeeded("from", fromUser); err != nil {
			return
		}
		if output.To.User, err = loadUserIfNeeded("to", toUser); err != nil {
			return
		}
	}

	transferData := models4debtus.NewTransferData(input.CreatorUser.ID, input.Request.IsReturn, input.Request.Amount, input.From, input.To)
	transferData.DtCreated = dtCreated
	output.Transfer.Data = transferData
	input.Source.PopulateTransfer(transferData)
	transferData.TransferInterest = *input.Request.Interest

	type TransferReturnInfo struct {
		Transfer       models4debtus.TransferEntry
		ReturnedAmount decimal.Decimal64p2
	}

	var (
		transferReturnInfos             = make([]TransferReturnInfo, 0, len(returnToTransferIDs))
		returnedValue, returnedInterest decimal.Decimal64p2
		closedTransferIDs               []string
	)

	// For api4transfers to specific api4transfers
	if len(returnToTransferIDs) > 0 {
		transferData.ReturnToTransferIDs = returnToTransferIDs

		returnToTransfers := models4debtus.NewTransfers(returnToTransferIDs)

		if err = tx.GetMulti(ctx, models4debtus.TransferRecords(returnToTransfers)); err != nil { // TODO: This can exceed limit on TX entity groups
			err = fmt.Errorf("failed to load returnToTransfers by keys (%v): %w", returnToTransferIDs, err)
			return
		}
		logus.Debugf(ctx, "Loaded %d returnToTransfers by keys", len(returnToTransfers))
		amountToAssign := input.Request.Amount.Value
		assignedToExistingTransfers := false
		for _, returnToTransfer := range returnToTransfers {
			//returnToTransfer := returnToTransfer.Data().(*models.TransferData)
			returnToTransferOutstandingValue := returnToTransfer.Data.GetOutstandingValue(dtCreated)
			if !returnToTransfer.Data.IsOutstanding {
				logus.Warningf(ctx, "TransferEntry(%v).IsOutstanding: false, returnToTransferOutstandingValue: %v", returnToTransfer.ID, returnToTransferOutstandingValue)
				continue
			} else if returnToTransferOutstandingValue == 0 {
				logus.Warningf(ctx, "TransferEntry(%s) => returnToTransferOutstandingValue == %d", returnToTransfer.ID, returnToTransferOutstandingValue)
				continue
			} else if returnToTransferOutstandingValue < 0 {
				panic(fmt.Sprintf("TransferEntry(%v) => returnToTransferOutstandingValue:%d <= 0", returnToTransfer.ID, returnToTransferOutstandingValue))
			}
			var amountReturnedToTransfer decimal.Decimal64p2
			if amountToAssign < returnToTransferOutstandingValue {
				amountReturnedToTransfer = amountToAssign
			} else {
				amountReturnedToTransfer = returnToTransferOutstandingValue
			}
			interestReturnedToTransfer := returnToTransfer.Data.GetInterestValue(dtCreated)
			if interestReturnedToTransfer > 0 {
				if interestReturnedToTransfer > amountReturnedToTransfer {
					interestReturnedToTransfer = amountReturnedToTransfer
				}
				returnedInterest += interestReturnedToTransfer
			}
			transferReturnInfos = append(transferReturnInfos, TransferReturnInfo{Transfer: returnToTransfer, ReturnedAmount: amountReturnedToTransfer})
			amountToAssign -= amountReturnedToTransfer
			returnedValue += amountReturnedToTransfer

			if err = transferData.AddReturn(models4debtus.TransferReturnJson{
				TransferID: returnToTransfer.ID,
				Amount:     amountReturnedToTransfer,
				Time:       returnToTransfer.Data.DtCreated,
			}); err != nil {
				return
			}

			assignedToExistingTransfers = true
			records = append(records, returnToTransfer.Record) // TODO: Potentially can exceed max number of records in GAE transaction

			if transferData.CreatorUserID == returnToTransfer.Data.CreatorUserID && transferData.Direction() == returnToTransfer.Data.Direction() {
				panic(fmt.Sprintf(
					"transfer.CreatorUserID == returnToTransfer.CreatorUserID && transfer.Direction == returnToTransfer.Direction, userID=%v, direction=%v, returnToTransfer=%v",
					transferData.CreatorUserID, transferData.Direction(), returnToTransfer.ID))
			}

			if transferData.CreatorUserID == returnToTransfer.Data.Counterparty().UserID && transferData.Direction() != returnToTransfer.Data.Direction() {
				panic(fmt.Sprintf(
					"transfer.CreatorUserID == returnToTransfer.CounterpartyUserID && transfer.Direction=%v != returnToTransfer.Direction=%v, userID=%v",
					transferData.Direction(), returnToTransfer.Data.Direction(), transferData.CreatorUserID))
			}

			if amountToAssign == 0 {
				break
			}
		}
		if assignedToExistingTransfers {
			if returnedValue > 0 {
				if returnedValue > input.Request.Amount.Value {
					panic("returnedAmount > input.Amount.Value")
				}
				if returnedValue == input.Request.Amount.Value && !transferData.IsReturn {
					transferData.IsReturn = true
					// transferData.AmountInCentsOutstanding = 0
					// transferData.AmountInCentsReturned = 0
					logus.Debugf(ctx, "TransferEntry marked IsReturn=true as it's amount less or equal to outstanding debt(s)")
				}
				// if returnedValue != input.Amount.Value {
				// 	// transferData.AmountInCentsOutstanding = input.Amount.Value - returnedAmount
				// 	transferData.AmountInCentsReturned = returnedValue
				// }
			}
			if output.From.User.ID != "" {
				if err = DelayUpdateHasDueTransfers(ctx, output.From.User.ID, output.From.SpaceID); err != nil {
					return
				}
			}
			if output.To.User.ID != "" {
				if err = DelayUpdateHasDueTransfers(ctx, output.To.User.ID, output.To.SpaceID); err != nil {
					return
				}
			}
		}
	}

	if input.Request.DueOn != nil && !input.Request.DueOn.IsZero() {
		transferData.DtDueOn = *input.Request.DueOn
		if from.UserID != "" {
			output.From.DebtusSpace.Data.HasDueTransfers = true
		}
		if to.UserID != "" {
			output.To.DebtusSpace.Data.HasDueTransfers = true
		}
	}

	// Set from & to names if needed
	{
		fixUserName := func(counterparty *models4debtus.TransferCounterpartyInfo, user dbo4userus.UserEntry) {
			if counterparty.UserID != "" && counterparty.UserName == "" {
				counterparty.UserName = user.Data.GetFullName()
			}
		}
		fixUserName(input.From, output.From.User)
		fixUserName(input.To, output.To.User)

		fixContactName := func(counterparty *models4debtus.TransferCounterpartyInfo, contact dal4contactus.ContactEntry) {
			if counterparty.ContactID != "" && counterparty.ContactName == "" {
				counterparty.ContactName = contact.Data.Names.GetFullName()
			}
		}
		fixContactName(input.From, output.From.Contact)
		fixContactName(input.To, output.To.Contact)
	}

	logus.Debugf(ctx, "from: %v", input.From)
	logus.Debugf(ctx, "to: %v", input.To)
	transferData.AmountInCentsInterest = returnedInterest

	// logus.Debugf(ctx, "transferData before insert: %v", litter.Sdump(transferData))
	if output.Transfer, err = InsertTransfer(ctx, tx, transferData); err != nil {
		err = fmt.Errorf("failed to save transfer entity: %w", err)
		return
	}

	createdTransfer := output.Transfer

	if output.Transfer.ID == "" {
		panic(fmt.Sprintf("Can't proceed creating transfer as InsertTransfer() returned transfer.ContactID == 0, err: %v", err))
	}

	logus.Infof(ctx, "TransferEntry inserted to DB with ContactID=%s, %+v", output.Transfer.ID, createdTransfer.Data)

	if len(transferReturnInfos) > 2 {
		transferReturnUpdates := make([]dtdal.TransferReturnUpdate, len(transferReturnInfos))
		for i, tri := range transferReturnInfos {
			transferReturnUpdates[i] = dtdal.TransferReturnUpdate{TransferID: tri.Transfer.ID, ReturnedAmount: tri.ReturnedAmount}
		}
		if err = dtdal.Transfer.DelayUpdateTransfersOnReturn(ctx, createdTransfer.ID, transferReturnUpdates); err != nil {
			return
		}
	} else {
		for _, transferReturnInfo := range transferReturnInfos {
			if err = Transfers.UpdateTransferOnReturn(ctx, tx, createdTransfer, transferReturnInfo.Transfer, transferReturnInfo.ReturnedAmount); err != nil {
				return
			}
			if !transferReturnInfo.Transfer.Data.IsOutstanding {
				closedTransferIDs = append(closedTransferIDs, transferReturnInfo.Transfer.ID)
			}
		}
	}

	// Update user and counterparty records with transfer info
	{
		var amountWithoutInterest money.Amount
		if returnedValue > 0 {
			amountWithoutInterest = money.Amount{Currency: input.Request.Amount.Currency, Value: input.Request.Amount.Value - returnedInterest}
		} else if returnedValue < 0 {
			panic(fmt.Sprintf("returnedValue < 0: %v", returnedValue))
		} else {
			amountWithoutInterest = input.Request.Amount
		}

		logus.Debugf(ctx, "closedTransferIDs: %v", closedTransferIDs)

		if output.From.User.ID == output.To.User.ID {
			panic(fmt.Sprintf("output.From.UserEntry.ContactID == output.To.UserEntry.ContactID: %v", output.From.User.ID))
		}
		if output.From.Contact.ID == output.To.Contact.ID {
			panic(fmt.Sprintf("output.From.DebtusSpaceContactEntry.ContactID == output.To.DebtusSpaceContactEntry.ContactID: %v", output.From.Contact.ID))
		}

		if output.From.User.ID != "" {
			if err = transferFacade.updateDebtusSpaceAndCounterpartyWithTransferInfo(ctx, amountWithoutInterest, output.Transfer, output.From.DebtusSpace, output.To.DebtusContact, closedTransferIDs); err != nil {
				return
			}
		}
		if output.To.User.ID != "" {
			if err = transferFacade.updateDebtusSpaceAndCounterpartyWithTransferInfo(ctx, amountWithoutInterest, output.Transfer, output.To.DebtusSpace, output.From.DebtusContact, closedTransferIDs); err != nil {
				return
			}
		}
	}

	{ // Integrity checks
		checkContacts := func(c1, c2 string, contact models4debtus.DebtusSpaceContactEntry, space models4debtus.DebtusSpaceEntry) {
			contactBalance := contact.Data.Balance
			contactBrief := space.Data.Contacts[contact.ID]
			if contactBrief == nil {
				panic(fmt.Sprintf("DebtusSpaceContactEntry.ContactID not found in counterparty Contacts(): %v", contact.ID))
			}
			cBalance := contactBrief.Balance
			for currency, val := range contactBalance {
				if cVal := cBalance[currency]; cVal != val {
					m := "balance inconsistency"
					panic(m)
				}
			}
		}

		if output.From.User.Data != nil {
			checkContacts("to", "from", output.To.DebtusContact, output.From.DebtusSpace)
		}
		if output.To.User.Data != nil {
			checkContacts("from", "to", output.From.DebtusContact, output.To.DebtusSpace)
		}
		if output.From.User.Data != nil && output.To.User.Data != nil {
			currency := output.Transfer.Data.Currency
			fromBalance := output.From.DebtusContact.Data.Balance[currency]
			toBalance := output.To.DebtusContact.Data.Balance[currency]
			if fromBalance != -toBalance {
				panic(fmt.Sprintf("from.DebtusSpaceContactEntry.Balance != -1*to.DebtusSpaceContactEntry.Balance => %v != -1*%v", fromBalance, -toBalance))
			}
		}
	}

	if err = tx.SetMulti(ctx, records); err != nil {
		err = fmt.Errorf("failed to update records: %w", err)
		return
	}

	if output.Transfer.Data.Counterparty().UserID != "" {
		if err = dtdal.Receipt.DelayCreateAndSendReceiptToCounterpartyByTelegram(ctx, input.Env, createdTransfer.ID, createdTransfer.Data.Counterparty().UserID); err != nil {
			// TODO: Send by any available channel
			err = fmt.Errorf("failed to delay sending receipt to counterpartyEntity by Telegram: %w", err)
			return
		}
	} else {
		logus.Debugf(ctx, "No receipt to counterpartyEntity: [%v]", createdTransfer.Data.Counterparty().ContactName)
	}

	if createdTransfer.Data.IsOutstanding && dtdal.Reminder != nil { // TODO: check for nil is temporary workaround for unittest
		if err = dtdal.Reminder.DelayCreateReminderForTransferUser(ctx, createdTransfer.ID, createdTransfer.Data.CreatorUserID); err != nil {
			err = fmt.Errorf("failed to delay reminder creation for creator: %w", err)
			return
		}
	}

	logus.Debugf(ctx, "createTransferWithinTransaction(): transferID=%v", createdTransfer.ID)
	return
}

func (TransfersFacade) GetTransferByID(ctx context.Context, tx dal.ReadSession, id string) (transfer models4debtus.TransferEntry, err error) {
	if tx == nil {
		if tx, err = facade.GetDatabase(ctx); err != nil {
			return
		}
	}
	transfer = models4debtus.NewTransfer(id, nil)
	err = tx.Get(ctx, transfer.Record)
	return
}

func (TransfersFacade) updateDebtusSpaceAndCounterpartyWithTransferInfo(
	ctx context.Context,
	amount money.Amount,
	transfer models4debtus.TransferEntry,
	debtusSpace models4debtus.DebtusSpaceEntry,
	contact models4debtus.DebtusSpaceContactEntry,
	closedTransferIDs []string,
) (err error) {
	logus.Debugf(ctx, "updateDebtusSpaceAndCounterpartyWithTransferInfo(debtusSpace=%v, Contact=%v)", debtusSpace, contact)
	var val decimal.Decimal64p2
	switch debtusSpace.ID {
	case transfer.Data.From().UserID:
		val = amount.Value * userBalanceIncreased
	case transfer.Data.To().UserID:
		val = amount.Value * userBalanceDecreased
	default:
		panic(fmt.Sprintf("debtusSpace is not related to transfer: %v", debtusSpace.ID))
	}
	logus.Debugf(ctx, "Updating balance with [%d %v] for debtusSpace #%s, Contact #%s", val, amount.Currency, debtusSpace.ID, contact.ID)

	if err = updateContactWithTransferInfo(ctx, val, transfer, contact, closedTransferIDs); err != nil {
		return
	}
	if err = updateDebtusSpaceWithTransferInfo(ctx, val, transfer, debtusSpace, contact, closedTransferIDs); err != nil {
		return
	}
	return
}

func updateDebtusSpaceWithTransferInfo(
	_ context.Context,
	val decimal.Decimal64p2,
	// curr money.CurrencyCode,
	transfer models4debtus.TransferEntry,
	debtusSpace models4debtus.DebtusSpaceEntry,
	contact models4debtus.DebtusSpaceContactEntry,
	// Contact models.DebtusSpaceContactEntry,
	closedTransferIDs []string,
) (err error) {
	debtusSpace.Data.LastTransferID = transfer.ID
	debtusSpace.Data.LastTransferAt = transfer.Data.DtCreated
	_, err = debtusSpace.Data.SetLastCurrency(transfer.Data.Currency)
	if err != nil {
		return err
	}

	// var updateBalanceAndContactTransfersInfo = func(curr money.CurrencyCode, val decimal.Decimal64p2, debtusSpace models.AppUserOBSOLETE, Contact models.DebtusSpaceContactEntry) (err error) {

	debtusSpace.Data.AddToBalance(transfer.Data.Currency, val)
	debtusSpace.Data.CountOfTransfers += 1
	_, _ = models4debtus.AddOrUpdateDebtusContact(debtusSpace, contact)
	return
}

func updateContactWithTransferInfo(
	ctx context.Context,
	val decimal.Decimal64p2,
	transfer models4debtus.TransferEntry,
	contact models4debtus.DebtusSpaceContactEntry,
	closedTransferIDs []string,
) (err error) {
	contact.Data.LastTransferID = transfer.ID
	contact.Data.LastTransferAt = transfer.Data.DtCreated

	contact.Data.AddToBalance(transfer.Data.Currency, val)
	contact.Data.CountOfTransfers += 1

	if contactTransfersInfo := contact.Data.GetTransfersInfo(); contactTransfersInfo.Last.ID != transfer.ID {
		contactTransfersInfo.Count += 1
		contactTransfersInfo.Last.ID = transfer.ID
		contactTransfersInfo.Last.At = transfer.Data.DtCreated
		if transfer.Data.HasInterest() {
			contactTransfersInfo.OutstandingWithInterest = append(contactTransfersInfo.OutstandingWithInterest, models4debtus.TransferWithInterestJson{
				TransferID:       transfer.ID,
				Amount:           transfer.Data.AmountInCents,
				Currency:         transfer.Data.Currency,
				Starts:           transfer.Data.DtCreated,
				TransferInterest: transfer.Data.TransferInterest,
			})
		}
		logus.Debugf(ctx, "len(contactTransfersInfo.OutstandingWithInterest): %v", len(contactTransfersInfo.OutstandingWithInterest))
		if len(contactTransfersInfo.OutstandingWithInterest) > 0 {
			if len(closedTransferIDs) > 0 {
				logus.Debugf(ctx, "removeClosedTransfersFromOutstandingWithInterest(closedTransferIDs: %v)", closedTransferIDs)
				contactTransfersInfo.OutstandingWithInterest = removeClosedTransfersFromOutstandingWithInterest(contactTransfersInfo.OutstandingWithInterest, closedTransferIDs)
			}
			logus.Debugf(ctx, "transfer.ReturnToTransferIDs: %v", transfer.Data.ReturnToTransferIDs)

			isClosed := func(transferID string) bool {
				return slice.Index(closedTransferIDs, transferID) >= 0
			}

		OuterLoop:
			for _, returnToTransferID := range transfer.Data.ReturnToTransferIDs {
				if isClosed(returnToTransferID) {
					logus.Debugf(ctx, "transfer %v is closed", returnToTransferID)
					continue
				}
				for i, outstanding := range contactTransfersInfo.OutstandingWithInterest {
					if outstanding.TransferID == returnToTransferID {
						if len(transfer.Data.ReturnToTransferIDs) == 1 {
							outstanding.Returns = append(outstanding.Returns, models4debtus.TransferReturnJson{
								TransferID: transfer.ID,
								Amount:     transfer.Data.AmountInCents,
								Time:       transfer.Data.DtCreated,
							})
							contactTransfersInfo.OutstandingWithInterest[i] = outstanding
						} else {
							err = fmt.Errorf("%w: return to multiple debts if at least one of them have interest is not implemented yet, please return debts with interest one by one", ErrNotImplemented)
							return
						}
						continue OuterLoop
					}
				}
				logus.Debugf(ctx, "transfer %v is not listed in contactTransfersInfo.OutstandingWithInterest", returnToTransferID)
			}
		}

		logus.Debugf(ctx, "transfer.HasInterest(): %v, contactTransfersInfo: %v", transfer.Data.HasInterest(), litter.Sdump(*contactTransfersInfo))
		if err = contact.Data.SetTransfersInfo(*contactTransfersInfo); err != nil {
			err = fmt.Errorf("failed to call SetTransfersInfo(): %w", err)
			return
		}
	}
	return
}

func removeClosedTransfersFromOutstandingWithInterest(
	transfersWithInterest []models4debtus.TransferWithInterestJson,
	closedTransferIDs []string,
) []models4debtus.TransferWithInterestJson {
	var i int
	for _, outstanding := range transfersWithInterest {
		if !slices.Contains(closedTransferIDs, outstanding.TransferID) {
			transfersWithInterest[i] = outstanding
			i += 1
		}
	}
	return transfersWithInterest[:i]
}

func InsertTransfer(ctx context.Context, tx dal.ReadwriteTransaction, transferEntity *models4debtus.TransferData) (transfer models4debtus.TransferEntry, err error) {
	transfer = models4debtus.NewTransfer("", transferEntity)
	err = tx.Insert(ctx, transfer.Record)
	return
}

func (TransfersFacade) UpdateTransferOnReturn(ctx context.Context, tx dal.ReadwriteTransaction, returnTransfer, transfer models4debtus.TransferEntry, returnedAmount decimal.Decimal64p2) (err error) {
	logus.Debugf(ctx, "UpdateTransferOnReturn(\n\treturnTransfer=%v,\n\ttransfer=%v,\n\treturnedAmount=%v)", litter.Sdump(returnTransfer), litter.Sdump(transfer), returnedAmount)

	if returnTransfer.Data.Currency != transfer.Data.Currency {
		panic(fmt.Sprintf("returnTransfer(id=%v).Currency != transfer.Currency => %v != %v", returnTransfer.ID, returnTransfer.Data.Currency, transfer.Data.Currency))
	} else if cID := returnTransfer.Data.From().ContactID; cID != "" && cID != transfer.Data.To().ContactID {
		if transfer.Data.To().ContactID == "" && returnTransfer.Data.From().UserID == transfer.Data.To().UserID {
			transfer.Data.To().ContactID = cID
			logus.Warningf(ctx, "Fixed TransferEntry(%v).To().ContactID: 0 => %v", transfer.ID, cID)
		} else {
			panic(fmt.Sprintf("returnTransfer(id=%v).From().ContactID != transfer.To().ContactID => %v != %v", returnTransfer.ID, cID, transfer.Data.To().ContactID))
		}
	} else if cID := returnTransfer.Data.To().ContactID; cID != "" && cID != transfer.Data.From().ContactID {
		if transfer.Data.From().ContactID == "" && returnTransfer.Data.To().UserID == transfer.Data.From().UserID {
			transfer.Data.From().ContactID = cID
			logus.Warningf(ctx, "Fixed TransferEntry(%v).From().ContactID: 0 => %v", transfer.ID, cID)
		} else {
			panic(fmt.Errorf("returnTransfer(id=%v).To().ContactID != transfer.From().ContactID => %v != %v", returnTransfer.ID, cID, transfer.Data.From().ContactID))
		}
	}

	for _, previousReturn := range transfer.Data.GetReturns() {
		if previousReturn.TransferID == returnTransfer.ID {
			logus.Infof(ctx, "TransferEntry already has information about return transfer")
			return
		}
	}

	if outstandingValue := transfer.Data.GetOutstandingValue(returnTransfer.Data.DtCreated); outstandingValue < returnedAmount {
		logus.Errorf(ctx, "transfer.GetOutstandingValue() < returnedAmount: %v <  %v", outstandingValue, returnedAmount)
		if outstandingValue <= 0 {
			return
		}
		returnedAmount = outstandingValue
	}

	if err = transfer.Data.AddReturn(models4debtus.TransferReturnJson{
		TransferID: returnTransfer.ID,
		Time:       returnTransfer.Data.DtCreated, // TODO: Replace with DtActual?
		Amount:     returnedAmount,
	}); err != nil {
		return
	}

	transfer.Data.IsOutstanding = transfer.Data.GetOutstandingValue(time.Now()) > 0

	if err = Transfers.SaveTransfer(ctx, tx, transfer); err != nil {
		return
	}

	if dtdal.Reminder != nil {
		if err = dtdal.Reminder.DelayDiscardReminders(ctx, []string{transfer.ID}, returnTransfer.ID); err != nil {
			err = fmt.Errorf("failed to delay task to discard reminders: %w", err)
			return
		}
	}

	return
}
