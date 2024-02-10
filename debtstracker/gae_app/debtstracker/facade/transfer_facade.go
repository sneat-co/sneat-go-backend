package facade

import (
	"bytes"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade/dto"
	"github.com/strongo/slice"
	"time"

	"context"
	"errors"
	"github.com/crediterra/money"
	"github.com/sanity-io/litter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/log"
)

const (
	userBalanceIncreased = 1
	userBalanceDecreased = -1
)

type TransfersFacade interface {
	GetTransferByID(c context.Context, tx dal.ReadSession, id string) (transfer models.Transfer, err error)
	SaveTransfer(c context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer) error
	CreateTransfer(c context.Context, input dto.CreateTransferInput) (output dto.CreateTransferOutput, err error)
	UpdateTransferOnReturn(c context.Context, tx dal.ReadwriteTransaction, returnTransfer, transfer models.Transfer, returnedAmount decimal.Decimal64p2) (err error)
}

var (
	ErrNotImplemented                      = errors.New("not implemented yet")
	ErrDebtAlreadyReturned                 = errors.New("This debt already has been returned")
	ErrPartialReturnGreaterThenOutstanding = errors.New("An attempt to do partial return for amount greater then outstanding")
	//
	ErrNoOutstandingTransfers                                       = errors.New("no outstanding transfers")
	ErrAttemptToCreateDebtWithInterestAffectingOutstandingTransfers = errors.New("You are trying to create a debt with interest that will affect outstanding transfers. Please close them first.")
)

func TransferCounterparties(direction models.TransferDirection, creatorInfo models.TransferCounterpartyInfo) (from, to *models.TransferCounterpartyInfo) {
	creator := models.TransferCounterpartyInfo{
		UserID:  creatorInfo.UserID,
		Comment: creatorInfo.Comment,
	}
	counterparty := models.TransferCounterpartyInfo{
		ContactID:   creatorInfo.ContactID,
		ContactName: creatorInfo.ContactName,
	}
	switch direction {
	case models.TransferDirectionUser2Counterparty:
		return &creator, &counterparty
	case models.TransferDirectionCounterparty2User:
		return &counterparty, &creator
	default:
		panic("Unknown direction: " + string(direction))
	}
}

type transferFacade struct {
}

var Transfers TransfersFacade = transferFacade{}

func (transferFacade) SaveTransfer(c context.Context, tx dal.ReadwriteTransaction, transfer models.Transfer) error {
	return tx.Set(c, transfer.Record)
}

func (transferFacade transferFacade) CreateTransfer(c context.Context, input dto.CreateTransferInput) (
	output dto.CreateTransferOutput, err error,
) {
	now := time.Now()

	log.Infof(c, "CreateTransfer(input=%v)", input)

	var returnToTransferIDs []string

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}

	//if input.Request.ReturnToTransferID != "" {
	//	if counterparty, err := GetContactByID(c, nil, contactID); err != nil {
	//		if dal.IsNotFound(err) {
	//			api.BadRequestError(c, w, err)
	//		} else {
	//			api.InternalError(c, w, err)
	//		}
	//		return
	//	} else {
	//		balance := counterparty.Data.Balance()
	//		if balanceAmount, ok := balance[amountWithCurrency.Currency]; !ok {
	//			api.BadRequestMessage(c, w, fmt.Sprintf("No balance for %v", amountWithCurrency.Currency))
	//		} else {
	//			switch direction {
	//			case models.TransferDirectionUser2Counterparty:
	//				if balanceAmount > 0 {
	//					api.BadRequestMessage(c, w, fmt.Sprintf("balanceAmount > 0 && direction == %v", direction))
	//				}
	//			case models.TransferDirectionCounterparty2User:
	//				if balanceAmount < 0 {
	//					api.BadRequestMessage(c, w, fmt.Sprintf("balanceAmount < 0 && direction == %v", direction))
	//				}
	//			}
	//		}
	//	}
	//}

	if input.Request.ReturnToTransferID == "" {
		log.Debugf(c, "input.ReturnToTransferID == 0")
		contacts := input.CreatorUser.Data.Contacts()
		creatorContactID := input.CreatorContactID()
		if creatorContactID == "" {
			panic(fmt.Errorf("3d party transfers are not implemented yet: %w", err))
		}
		log.Debugf(c, "creatorContactID=%v, contacts: %+v", creatorContactID, contacts)
		var creatorContact models.Contact
		verifyUserContactJson := func() (contactJsonFound bool) {
			for _, contact := range contacts {
				if contact.ID == creatorContactID {
					contactBalance := contact.Balance()
					if v, ok := contactBalance[input.Request.Amount.Currency]; !ok || v == 0 {
						log.Debugf(c, "No need to check for outstanding transfers as contacts balance is 0")
					} else {
						if input.Request.Interest.HasInterest() {
							if d := input.Direction(); d == models.TransferDirectionUser2Counterparty && v < 0 || d == models.TransferDirectionCounterparty2User && v > 0 {
								err = ErrAttemptToCreateDebtWithInterestAffectingOutstandingTransfers
								return
							}
						}
						if returnToTransferIDs, err = transferFacade.checkOutstandingTransfersForReturns(c, now, input); err != nil {
							return
						}
					}
					contactJsonFound = true
					return
				}
			}
			return
		}
		if contactJsonFound := verifyUserContactJson(); contactJsonFound {
			goto contactFound
		}
		// If contact not found in user's JSON try to recover from DB record
		if creatorContact, err = GetContactByID(c, nil, creatorContactID); err != nil {
			return
		}

		log.Warningf(c, "data integrity issue: contact found by ID in database but is missing in user's JSON: creatorContactID=%v, creatorContact.UserID=%v, user.ID=%v, user.ContactsJsonActive: %v",
			creatorContactID, creatorContact.Data.UserID, input.CreatorUser.ID, input.CreatorUser.Data.ContactsJsonActive)

		if creatorContact.Data.UserID != input.CreatorUser.ID {
			err = fmt.Errorf("creatorContact.UserID != input.CreatorUser.ID: %v != %v", creatorContact.Data.UserID, input.CreatorUser.ID)
			return
		}

		if _, changed := models.AddOrUpdateContact(&input.CreatorUser, creatorContact); changed {
			contacts = input.CreatorUser.Data.Contacts()
		}
		if contactJsonFound := verifyUserContactJson(); contactJsonFound {
			goto contactFound
		}
		if err == nil {
			err = fmt.Errorf("user contact not found by ID=%v, contacts: %v", creatorContactID, litter.Sdump(contacts))
		}
		return
	contactFound:
	} else if !input.Request.IsReturn {
		panic("ReturnToTransferID != 0 && !IsReturn")
	}

	if input.Request.ReturnToTransferID != "" {
		var transferToReturn models.Transfer
		if transferToReturn, err = Transfers.GetTransferByID(c, db, input.Request.ReturnToTransferID); err != nil {
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
			log.Debugf(c, "amount.Value:%v > transferToReturn.GetOutstandingValue(now):%v", input.Request.Amount.Value, transferToReturn.Data.GetOutstandingValue(now))
			if input.Request.Amount.Value == transferToReturn.Data.AmountInCents {
				// For situations when a transfer was partially returned but user wants to mark it as fully returned.
				log.Debugf(c, "amount.Value (%v) == transferToReturn.AmountInCents (%v)", input.Request.Amount.Value, transferToReturn.Data.AmountInCents)
				input.Request.Amount.Value = transferToReturn.Data.GetOutstandingValue(now)
				log.Debugf(c, "Updated amount.Value: %v", input.Request.Amount.Value)
			} else {
				err = ErrPartialReturnGreaterThenOutstanding
				return
			}
		} else if input.Request.Amount.Value < transferToReturn.Data.GetOutstandingValue(now) {
			log.Debugf(c, "input.Amount.Value < transferToReturn.GetOutstandingValue(now)")
		}

		returnToTransferIDs = append(returnToTransferIDs, input.Request.ReturnToTransferID)
		output.ReturnedTransfers = append(output.ReturnedTransfers, transferToReturn)
	}

	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		output, err = transferFacade.createTransferWithinTransaction(c, tx, now, input, returnToTransferIDs)
		return err
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}

	output.Validate()

	return
}

func (transferFacade transferFacade) checkOutstandingTransfersForReturns(c context.Context, now time.Time, input dto.CreateTransferInput) (returnToTransferIDs []string, err error) {
	log.Debugf(c, "transferFacade.checkOutstandingTransfersForReturns()")
	var (
		outstandingTransfers []models.Transfer
	)

	creatorUserID := input.CreatorUser.ID
	creatorContactID := input.CreatorContactID()

	reversedDirection := input.Direction().Reverse()

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	outstandingTransfers, err = dtdal.Transfer.LoadOutstandingTransfers(c, db, now, creatorUserID, creatorContactID, input.Request.Amount.Currency, reversedDirection)
	if err != nil {
		err = fmt.Errorf("failed to load outstanding transfers: %w", err)
		return
	}
	if input.Request.IsReturn && len(outstandingTransfers) == 0 {
		err = ErrNoOutstandingTransfers
		return
	}

	log.Debugf(c, "facade.checkOutstandingTransfersForReturns() => dtdal.Transfer.LoadOutstandingTransfers(userID=%v, currency=%v) => %d transfers", input.CreatorUser.ID, input.Request.Amount.Currency, len(outstandingTransfers))

	if outstandingTransfersCount := len(outstandingTransfers); outstandingTransfersCount > 0 { // Assign the return to specific transfers
		var (
			assignedValue             decimal.Decimal64p2
			outstandingRightDirection int
		)
		buf := new(bytes.Buffer)
		_, _ = fmt.Fprintf(buf, "%v outstanding transfers\n", outstandingTransfersCount)
		for i, outstandingTransfer := range outstandingTransfers {
			_, _ = fmt.Fprintf(buf, "\t[%v]: %v", i, litter.Sdump(outstandingTransfer))
			outstandingTransferID := outstandingTransfers[i].ID
			outstandingValue := outstandingTransfer.Data.GetOutstandingValue(now)
			if outstandingValue == input.Request.Amount.Value { // A check for exact match that has higher priority then earlie transfers
				log.Infof(c, " - found outstanding transfer %v with exact amount match: %v", outstandingTransfer.ID, outstandingValue)
				assignedValue = input.Request.Amount.Value
				returnToTransferIDs = []string{outstandingTransferID}
				break
			}
			if assignedValue < input.Request.Amount.Value { // Do not break so we check all outstanding transfers for exact match
				returnToTransferIDs = append(returnToTransferIDs, outstandingTransferID)
				assignedValue += outstandingValue
			}
			outstandingRightDirection += 1
			buf.WriteString("\n")
		}
		log.Debugf(c, buf.String())
		if input.Request.IsReturn && assignedValue < input.Request.Amount.Value {
			log.Warningf(c,
				"There are not enough outstanding transfers to return %v. All outstanding count: %v, Right direction: %v, Assigned amount: %v. Could be data integrity issue.",
				input.Request.Amount, len(outstandingTransfers), outstandingRightDirection, assignedValue,
			)
		}
	}
	return
}

func (transferFacade transferFacade) createTransferWithinTransaction(
	c context.Context,
	tx dal.ReadwriteTransaction,
	dtCreated time.Time,
	input dto.CreateTransferInput,
	returnToTransferIDs []string,
) (
	output dto.CreateTransferOutput, err error,
) {
	log.Debugf(c, "createTransferWithinTransaction(input=%v, returnToTransferIDs=%v)", input, returnToTransferIDs)

	if err = input.Validate(); err != nil {
		return
	}
	// if len(returnToTransferIDs) > 0 && !input.IsReturn { // TODO: It's OK to have transfers without isReturn=true
	// 	panic("len(returnToTransferIDs) > 0 && !isReturn")
	// }

	output.From = new(dto.CreateTransferOutputCounterparty)
	output.To = new(dto.CreateTransferOutputCounterparty)
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

	if err = tx.GetMulti(c, records); err != nil {
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
		if fromContact.Data.CounterpartyUserID != toContact.Data.UserID {
			err = fmt.Errorf("fromCounterparty.CounterpartyUserID != toCounterparty.UserID (%s != %s)",
				fromContact.Data.CounterpartyUserID, toContact.Data.UserID)
		}
		if toContact.Data.CounterpartyUserID != fromContact.Data.UserID {
			err = fmt.Errorf("toCounterparty.CounterpartyUserID != fromCounterparty.UserID (%s != %s)",
				toContact.Data.CounterpartyUserID, fromContact.Data.UserID)
		}
		return
	}

	// Check if counterparties are linked and if yes load the missing Contact
	{
		link := func(sideName, countersideName string, side, counterside *models.TransferCounterpartyInfo, sideContact models.Contact) (countersideContact models.Contact, err error) {
			log.Debugf(c, "link(%v=%v, %v=%v, %vContact=%v)", sideName, side, countersideName, counterside, sideName, sideContact)
			if side.ContactID != "" && sideContact.Data.CounterpartyCounterpartyID != "" && counterside.ContactID == "" {
				if countersideContact, err = GetContactByID(c, tx, sideContact.Data.CounterpartyCounterpartyID); err != nil {
					err = fmt.Errorf("failed to get counterparty by 'fromCounterparty.CounterpartyCounterpartyID': %w", err)
					return
				}
				counterside.ContactID = countersideContact.ID
				counterside.ContactName = countersideContact.Data.FullName()
				side.UserID = countersideContact.Data.UserID
				records = append(records, countersideContact.Record)
			}
			return
		}

		var linkedContact models.Contact // TODO: This smells
		if linkedContact, err = link("from", "to", from, to, fromContact); err != nil {
			return
		} else if linkedContact.Data != nil {
			toContact = linkedContact
			output.To.Contact = linkedContact
		}

		log.Debugf(c, "toContact: %v", toContact.Data == nil)
		if linkedContact, err = link("to", "from", to, from, toContact); err != nil {
			return
		} else if linkedContact.Data != nil {
			fromContact = linkedContact
			output.From.Contact = fromContact
		}

		// //// When: toCounterparty == nil, fromUser == nil,
		// if from.ContactID != 0 && fromContact.CounterpartyCounterpartyID != 0 && to.ContactID == 0 {
		// 	// Get toCounterparty and fill to.Contact* fields
		// 	if toContact, err = GetContactByID(c, fromContact.CounterpartyCounterpartyID); err != nil {
		// 		err = fmt.Errorf("%w: Failed to get 'To' counterparty by 'fromCounterparty.CounterpartyCounterpartyID'", err)
		// 		return
		// 	}
		// 	output.To.Contact = toContact
		// 	log.Debugf(c, "Got toContact id=%d: %v", toContact.ID, toContact.DebtusContactData)
		// 	to.ContactID = toContact.ID
		// 	to.ContactName = toContact.GetFullName()
		// 	from.UserID = toContact.UserID
		// 	records = append(records, &toContact)
		// }
		// if to.ContactID != 0 && toCounterparty.CounterpartyCounterpartyID != 0 && from.ContactID == 0 {
		// 	if fromCounterparty, err = GetContactByID(c, toCounterparty.CounterpartyCounterpartyID); err != nil {
		// 		err = fmt.Errorf("failed to get 'From' counterparty by 'toCounterparty.CounterpartyCounterpartyID' == %d: %w", fromCounterparty.CounterpartyCounterpartyID, err)
		// 		return
		// 	}
		// 	output.From.Contact = fromCounterparty
		// 	log.Debugf(c, "Got fromCounterparty id=%d: %v", fromCounterparty.ID, fromCounterparty.DebtusContactData)
		// 	from.ContactID = fromCounterparty.ID
		// 	from.ContactName = fromCounterparty.GetFullName()
		// 	to.UserID = fromCounterparty.UserID
		// 	records = append(records, &fromCounterparty)
		// }
	}

	// In case if we just loaded above missing counterparty we need to check for missing user
	{
		loadUserIfNeeded := func(who string, userID string, appUser models.AppUser) (models.AppUser, error) {
			log.Debugf(c, "%v.UserID: %d, %vUser.DebutsAppUserDataOBSOLETE: %v", who, userID, who, appUser.Data)
			if userID != "" {
				if appUser.Data == nil {
					if appUser, err = User.GetUserByID(c, tx, userID); err != nil {
						err = fmt.Errorf("failed to get %vUser for linked counterparty: %w", who, err)
						return appUser, err
					}
					records = append(records, appUser.Record)
				} else if userID != appUser.ID {
					panic("userID != appUser.ID")
				}
			}
			return appUser, err
		}

		if output.From.User, err = loadUserIfNeeded("from", from.UserID, fromUser); err != nil {
			return
		}
		if output.To.User, err = loadUserIfNeeded("to", to.UserID, toUser); err != nil {
			return
		}
	}

	transferData := models.NewTransferData(input.CreatorUser.ID, input.Request.IsReturn, input.Request.Amount, input.From, input.To)
	transferData.DtCreated = dtCreated
	output.Transfer.Data = transferData
	input.Source.PopulateTransfer(transferData)
	transferData.TransferInterest = *input.Request.Interest

	type TransferReturnInfo struct {
		Transfer       models.Transfer
		ReturnedAmount decimal.Decimal64p2
	}

	var (
		transferReturnInfos             = make([]TransferReturnInfo, 0, len(returnToTransferIDs))
		returnedValue, returnedInterest decimal.Decimal64p2
		closedTransferIDs               []string
	)

	// For transfers to specific transfers
	if len(returnToTransferIDs) > 0 {
		transferData.ReturnToTransferIDs = returnToTransferIDs

		returnToTransfers := models.NewTransfers(returnToTransferIDs)

		if err = tx.GetMulti(c, models.TransferRecords(returnToTransfers)); err != nil { // TODO: This can exceed limit on TX entity groups
			err = fmt.Errorf("failed to load returnToTransfers by keys (%v): %w", returnToTransferIDs, err)
			return
		}
		log.Debugf(c, "Loaded %d returnToTransfers by keys", len(returnToTransfers))
		amountToAssign := input.Request.Amount.Value
		assignedToExistingTransfers := false
		for _, returnToTransfer := range returnToTransfers {
			//returnToTransfer := returnToTransfer.Data().(*models.TransferData)
			returnToTransferOutstandingValue := returnToTransfer.Data.GetOutstandingValue(dtCreated)
			if !returnToTransfer.Data.IsOutstanding {
				log.Warningf(c, "Transfer(%v).IsOutstanding: false, returnToTransferOutstandingValue: %v", returnToTransfer.ID, returnToTransferOutstandingValue)
				continue
			} else if returnToTransferOutstandingValue == 0 {
				log.Warningf(c, "Transfer(%v) => returnToTransferOutstandingValue == 0", returnToTransfer.ID, returnToTransferOutstandingValue)
				continue
			} else if returnToTransferOutstandingValue < 0 {
				panic(fmt.Sprintf("Transfer(%v) => returnToTransferOutstandingValue:%d <= 0", returnToTransfer.ID, returnToTransferOutstandingValue))
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

			if err = transferData.AddReturn(models.TransferReturnJson{
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
					log.Debugf(c, "Transfer marked IsReturn=true as it's amount less or equal to outstanding debt(s)")
				}
				// if returnedValue != input.Amount.Value {
				// 	// transferData.AmountInCentsOutstanding = input.Amount.Value - returnedAmount
				// 	transferData.AmountInCentsReturned = returnedValue
				// }
			}
			if output.From.User.ID != "" {
				if err = dtdal.User.DelayUpdateUserHasDueTransfers(c, output.From.User.ID); err != nil {
					return
				}
			}
			if output.To.User.ID != "" {
				if err = dtdal.User.DelayUpdateUserHasDueTransfers(c, output.To.User.ID); err != nil {
					return
				}
			}
		}
	}

	if input.Request.DueOn != nil && !input.Request.DueOn.IsZero() {
		transferData.DtDueOn = *input.Request.DueOn
		if from.UserID != "" {
			output.From.User.Data.HasDueTransfers = true
		}
		if to.UserID != "" {
			output.To.User.Data.HasDueTransfers = true
		}
	}

	// Set from & to names if needed
	{
		fixUserName := func(counterparty *models.TransferCounterpartyInfo, user models.AppUser) {
			if counterparty.UserID != "" && counterparty.UserName == "" {
				counterparty.UserName = user.Data.FullName()
			}
		}
		fixUserName(input.From, output.From.User)
		fixUserName(input.To, output.To.User)

		fixContactName := func(counterparty *models.TransferCounterpartyInfo, contact models.Contact) {
			if counterparty.ContactID != "" && counterparty.ContactName == "" {
				counterparty.ContactName = contact.Data.FullName()
			}
		}
		fixContactName(input.From, output.From.Contact)
		fixContactName(input.To, output.To.Contact)
	}

	log.Debugf(c, "from: %v", input.From)
	log.Debugf(c, "to: %v", input.To)
	transferData.AmountInCentsInterest = returnedInterest

	// log.Debugf(c, "transferData before insert: %v", litter.Sdump(transferData))
	if output.Transfer, err = InsertTransfer(c, tx, transferData); err != nil {
		err = fmt.Errorf("failed to save transfer entity: %w", err)
		return
	}

	createdTransfer := output.Transfer

	if output.Transfer.ID == "" {
		panic(fmt.Sprintf("Can't proceed creating transfer as InsertTransfer() returned transfer.ID == 0, err: %v", err))
	}

	log.Infof(c, "Transfer inserted to DB with ID=%d, %+v", output.Transfer.ID, createdTransfer.Data)

	if len(transferReturnInfos) > 2 {
		transferReturnUpdates := make([]dtdal.TransferReturnUpdate, len(transferReturnInfos))
		for i, tri := range transferReturnInfos {
			transferReturnUpdates[i] = dtdal.TransferReturnUpdate{TransferID: tri.Transfer.ID, ReturnedAmount: tri.ReturnedAmount}
		}
		if err = dtdal.Transfer.DelayUpdateTransfersOnReturn(c, createdTransfer.ID, transferReturnUpdates); err != nil {
			return
		}
	} else {
		for _, transferReturnInfo := range transferReturnInfos {
			if err = Transfers.UpdateTransferOnReturn(c, tx, createdTransfer, transferReturnInfo.Transfer, transferReturnInfo.ReturnedAmount); err != nil {
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

		log.Debugf(c, "closedTransferIDs: %v", closedTransferIDs)

		if output.From.User.ID == output.To.User.ID {
			panic(fmt.Sprintf("output.From.User.ID == output.To.User.ID: %v", output.From.User.ID))
		}
		if output.From.Contact.ID == output.To.Contact.ID {
			panic(fmt.Sprintf("output.From.Contact.ID == output.To.Contact.ID: %v", output.From.Contact.ID))
		}

		if output.From.User.ID != "" {
			if err = transferFacade.updateUserAndCounterpartyWithTransferInfo(c, amountWithoutInterest, output.Transfer, output.From.User, output.To.Contact, closedTransferIDs); err != nil {
				return
			}
		}
		if output.To.User.ID != "" {
			if err = transferFacade.updateUserAndCounterpartyWithTransferInfo(c, amountWithoutInterest, output.Transfer, output.To.User, output.From.Contact, closedTransferIDs); err != nil {
				return
			}
		}
	}

	{ // Integrity checks
		checkContacts := func(c1, c2 string, contact models.Contact, user models.AppUser) {
			contacts := user.Data.Contacts()
			contactBalance := contact.Data.Balance()
			for _, c := range contacts {
				if c.ID == contact.ID {
					cBalance := c.Balance()
					for currency, val := range contactBalance {
						if cVal := cBalance[currency]; cVal != val {
							panic(fmt.Sprintf(
								"balance inconsistency for (user=%v&contact=%v VS user=%v&contact=%v) => "+
									"%v: %v != %v\n%v.Balance: %v\n\n%v.Balance: %v",
								contact.Data.UserID, contact.ID, user.ID, c.ID, currency, cVal, val, c1, contactBalance, c2, cBalance))
						}
					}
					return
				}
			}
			panic(fmt.Sprintf("Contact.ID not found in counterparty Contacts(): %v", contact.ID))
		}

		if output.From.User.Data != nil {
			checkContacts("to", "from", output.To.Contact, output.From.User)
		}
		if output.To.User.Data != nil {
			checkContacts("from", "to", output.From.Contact, output.To.User)
		}
		if output.From.User.Data != nil && output.To.User.Data != nil {
			currency := output.Transfer.Data.Currency
			fromBalance := output.From.Contact.Data.Balance()[currency]
			toBalance := output.To.Contact.Data.Balance()[currency]
			if fromBalance != -toBalance {
				panic(fmt.Sprintf("from.Contact.Balance != -1*to.Contact.Balance => %v != -1*%v", fromBalance, -toBalance))
			}
		}
	}

	if err = tx.SetMulti(c, records); err != nil {
		err = fmt.Errorf("failed to update records: %w", err)
		return
	}

	if output.Transfer.Data.Counterparty().UserID != "" {
		if err = dtdal.Receipt.DelayCreateAndSendReceiptToCounterpartyByTelegram(c, input.Env, createdTransfer.ID, createdTransfer.Data.Counterparty().UserID); err != nil {
			// TODO: Send by any available channel
			err = fmt.Errorf("failed to delay sending receipt to counterpartyEntity by Telegram: %w", err)
			return
		}
	} else {
		log.Debugf(c, "No receipt to counterpartyEntity: [%v]", createdTransfer.Data.Counterparty().ContactName)
	}

	if createdTransfer.Data.IsOutstanding && dtdal.Reminder != nil { // TODO: check for nil is temporary workaround for unittest
		if err = dtdal.Reminder.DelayCreateReminderForTransferUser(c, createdTransfer.ID, createdTransfer.Data.CreatorUserID); err != nil {
			err = fmt.Errorf("failed to delay reminder creation for creator: %w", err)
			return
		}
	}

	log.Debugf(c, "createTransferWithinTransaction(): transferID=%v", createdTransfer.ID)
	return
}

func (transferFacade) GetTransferByID(c context.Context, tx dal.ReadSession, id string) (transfer models.Transfer, err error) {
	if tx == nil {
		if tx, err = GetDatabase(c); err != nil {
			return
		}
	}
	transfer = models.NewTransfer(id, nil)
	err = tx.Get(c, transfer.Record)
	return
}

func (transferFacade) updateUserAndCounterpartyWithTransferInfo(
	c context.Context,
	amount money.Amount,
	transfer models.Transfer,
	user models.AppUser,
	contact models.Contact,
	closedTransferIDs []string,
) (err error) {
	log.Debugf(c, "updateUserAndCounterpartyWithTransferInfo(user=%v, contact=%v)", user, contact)
	if user.ID != contact.Data.UserID {
		panic(fmt.Errorf("user.ID != contact.UserID (%s != %s)", user.ID, contact.Data.UserID))
	}
	var val decimal.Decimal64p2
	switch user.ID {
	case transfer.Data.From().UserID:
		val = amount.Value * userBalanceIncreased
	case transfer.Data.To().UserID:
		val = amount.Value * userBalanceDecreased
	default:
		panic(fmt.Sprintf("user is not related to transfer: %v", user.ID))
	}
	log.Debugf(c, "Updating balance with [%v %v] for user #%d, contact #%d", val, amount.Currency, user.ID, contact.ID)

	if err = updateContactWithTransferInfo(c, val, transfer, contact, closedTransferIDs); err != nil {
		return
	}
	if err = updateUserWithTransferInfo(c, val, transfer, user, contact, closedTransferIDs); err != nil {
		return
	}
	return
}

func updateUserWithTransferInfo(
	c context.Context,
	val decimal.Decimal64p2,
	// curr money.CurrencyCode,
	transfer models.Transfer,
	user models.AppUser,
	contact models.Contact,
	// contact models.Contact,
	closedTransferIDs []string,
) (err error) {
	user.Data.LastTransferID = transfer.ID
	user.Data.LastTransferAt = transfer.Data.DtCreated
	user.Data.SetLastCurrency(string(transfer.Data.Currency))

	// var updateBalanceAndContactTransfersInfo = func(curr money.CurrencyCode, val decimal.Decimal64p2, user models.AppUser, contact models.Contact) (err error) {

	var balance money.Balance
	if balance, err = user.Data.AddToBalance(transfer.Data.Currency, val); err != nil {
		err = fmt.Errorf("failed to add %s=%d to balance for user %v: %w", transfer.Data.Currency, val, user.ID, err)
		return
	} else {
		user.Data.CountOfTransfers += 1
		userBalance := user.Data.Balance()
		log.Debugf(c, "Updated balance to %v | %v for user #%d", balance, userBalance, user.ID)
	}
	log.Debugf(c, "user.ContactsJsonActive (before): %v\ncontact: %v", user.Data.ContactsJsonActive, litter.Sdump(contact))
	_, userContactsChanged := models.AddOrUpdateContact(&user, contact)
	log.Debugf(c, "user.ContactsJson (changed=%v): %v, closedTransferIDs: %+v", userContactsChanged, user.Data.ContactsJsonActive, closedTransferIDs)
	return
}

func updateContactWithTransferInfo(
	c context.Context,
	val decimal.Decimal64p2,
	transfer models.Transfer,
	contact models.Contact,
	closedTransferIDs []string,
) (err error) {
	contact.Data.LastTransferID = transfer.ID
	contact.Data.LastTransferAt = transfer.Data.DtCreated

	var balance money.Balance
	if balance, err = contact.Data.AddToBalance(transfer.Data.Currency, val); err != nil {
		err = fmt.Errorf("failed to add (%s %v) to balance for contact #%s: %w", transfer.Data.Currency, val, contact.ID, err)
		return
	} else {
		contact.Data.CountOfTransfers += 1
		cpBalance := contact.Data.Balance()
		log.Debugf(c, "Updated balance to %v | %v for contact #%d", balance, cpBalance, contact.ID)
	}

	if contactTransfersInfo := contact.Data.GetTransfersInfo(); contactTransfersInfo.Last.ID != transfer.ID {
		contactTransfersInfo.Count += 1
		contactTransfersInfo.Last.ID = transfer.ID
		contactTransfersInfo.Last.At = transfer.Data.DtCreated
		if transfer.Data.HasInterest() {
			contactTransfersInfo.OutstandingWithInterest = append(contactTransfersInfo.OutstandingWithInterest, models.TransferWithInterestJson{
				TransferID:       transfer.ID,
				Amount:           transfer.Data.AmountInCents,
				Currency:         transfer.Data.Currency,
				Starts:           transfer.Data.DtCreated,
				TransferInterest: transfer.Data.TransferInterest,
			})
		}
		log.Debugf(c, "len(contactTransfersInfo.OutstandingWithInterest): %v", len(contactTransfersInfo.OutstandingWithInterest))
		if len(contactTransfersInfo.OutstandingWithInterest) > 0 {
			if len(closedTransferIDs) > 0 {
				log.Debugf(c, "removeClosedTransfersFromOutstandingWithInterest(closedTransferIDs: %v)", closedTransferIDs)
				contactTransfersInfo.OutstandingWithInterest = removeClosedTransfersFromOutstandingWithInterest(contactTransfersInfo.OutstandingWithInterest, closedTransferIDs)
			}
			log.Debugf(c, "transfer.ReturnToTransferIDs: %v", transfer.Data.ReturnToTransferIDs)

			isClosed := func(transferID string) bool {
				return slice.Index(closedTransferIDs, transferID) >= 0
			}

		OuterLoop:
			for _, returnToTransferID := range transfer.Data.ReturnToTransferIDs {
				if isClosed(returnToTransferID) {
					log.Debugf(c, "transfer %v is closed", returnToTransferID)
					continue
				}
				for i, outstanding := range contactTransfersInfo.OutstandingWithInterest {
					if outstanding.TransferID == returnToTransferID {
						if len(transfer.Data.ReturnToTransferIDs) == 1 {
							outstanding.Returns = append(outstanding.Returns, models.TransferReturnJson{
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
				log.Debugf(c, "transfer %v is not listed in contactTransfersInfo.OutstandingWithInterest", returnToTransferID)
			}
		}

		log.Debugf(c, "transfer.HasInterest(): %v, contactTransfersInfo: %v", transfer.Data.HasInterest(), litter.Sdump(*contactTransfersInfo))
		if err = contact.Data.SetTransfersInfo(*contactTransfersInfo); err != nil {
			err = fmt.Errorf("failed to call SetTransfersInfo(): %w", err)
			return
		}
	}
	return
}

func removeClosedTransfersFromOutstandingWithInterest(
	transfersWithInterest []models.TransferWithInterestJson,
	closedTransferIDs []string,
) []models.TransferWithInterestJson {
	var i int
	for _, outstanding := range transfersWithInterest {
		if !slice.Contains(closedTransferIDs, outstanding.TransferID) {
			transfersWithInterest[i] = outstanding
			i += 1
		}
	}
	return transfersWithInterest[:i]
}

func InsertTransfer(c context.Context, tx dal.ReadwriteTransaction, transferEntity *models.TransferData) (transfer models.Transfer, err error) {
	transfer = models.NewTransfer("", transferEntity)
	err = tx.Insert(c, transfer.Record)
	return
}

func (transferFacade) UpdateTransferOnReturn(c context.Context, tx dal.ReadwriteTransaction, returnTransfer, transfer models.Transfer, returnedAmount decimal.Decimal64p2) (err error) {
	log.Debugf(c, "UpdateTransferOnReturn(\n\treturnTransfer=%v,\n\ttransfer=%v,\n\treturnedAmount=%v)", litter.Sdump(returnTransfer), litter.Sdump(transfer), returnedAmount)

	if returnTransfer.Data.Currency != transfer.Data.Currency {
		panic(fmt.Sprintf("returnTransfer(id=%v).Currency != transfer.Currency => %v != %v", returnTransfer.ID, returnTransfer.Data.Currency, transfer.Data.Currency))
	} else if cID := returnTransfer.Data.From().ContactID; cID != "" && cID != transfer.Data.To().ContactID {
		if transfer.Data.To().ContactID == "" && returnTransfer.Data.From().UserID == transfer.Data.To().UserID {
			transfer.Data.To().ContactID = cID
			log.Warningf(c, "Fixed Transfer(%v).To().ContactID: 0 => %v", transfer.ID, cID)
		} else {
			panic(fmt.Sprintf("returnTransfer(id=%v).From().ContactID != transfer.To().ContactID => %v != %v", returnTransfer.ID, cID, transfer.Data.To().ContactID))
		}
	} else if cID := returnTransfer.Data.To().ContactID; cID != "" && cID != transfer.Data.From().ContactID {
		if transfer.Data.From().ContactID == "" && returnTransfer.Data.To().UserID == transfer.Data.From().UserID {
			transfer.Data.From().ContactID = cID
			log.Warningf(c, "Fixed Transfer(%v).From().ContactID: 0 => %v", transfer.ID, cID)
		} else {
			panic(fmt.Errorf("returnTransfer(id=%v).To().ContactID != transfer.From().ContactID => %v != %v", returnTransfer.ID, cID, transfer.Data.From().ContactID))
		}
	}

	for _, previousReturn := range transfer.Data.GetReturns() {
		if previousReturn.TransferID == returnTransfer.ID {
			log.Infof(c, "Transfer already has information about return transfer")
			return
		}
	}

	if outstandingValue := transfer.Data.GetOutstandingValue(returnTransfer.Data.DtCreated); outstandingValue < returnedAmount {
		log.Errorf(c, "transfer.GetOutstandingValue() < returnedAmount: %v <  %v", outstandingValue, returnedAmount)
		if outstandingValue <= 0 {
			return
		}
		returnedAmount = outstandingValue
	}

	if err = transfer.Data.AddReturn(models.TransferReturnJson{
		TransferID: returnTransfer.ID,
		Time:       returnTransfer.Data.DtCreated, // TODO: Replace with DtActual?
		Amount:     returnedAmount,
	}); err != nil {
		return
	}

	transfer.Data.IsOutstanding = transfer.Data.GetOutstandingValue(time.Now()) > 0

	if err = Transfers.SaveTransfer(c, tx, transfer); err != nil {
		return
	}

	if dtdal.Reminder != nil {
		if err = dtdal.Reminder.DelayDiscardReminders(c, []string{transfer.ID}, returnTransfer.ID); err != nil {
			err = fmt.Errorf("failed to delay task to discard reminders: %w", err)
			return
		}
	}

	return
}
