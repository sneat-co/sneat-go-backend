package inspector

import (
	"context"
	"errors"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/logus"
	"net/http"
	"sync"
	"time"
)

type contactWithBalances struct {
	models4debtus.DebtusSpaceContactEntry
	transfersCount int
	balances       balances
}

type transfersInfo struct {
	count   int
	balance money.Balance
}

func newContactWithBalances(ctx context.Context, now time.Time, contact models4debtus.DebtusSpaceContactEntry) contactWithBalances {
	balanceWithInterest, err := contact.Data.BalanceWithInterest(ctx, now)
	result := contactWithBalances{
		DebtusSpaceContactEntry: contact,
		balances:                newBalances("contact", contact.Data.Balance, balanceWithInterest),
	}
	result.balances.withInterest.err = err
	return result
}

func newBalanceSummary(who string, balance money.Balance) (balances balancesByCurrency) {
	balances = balancesByCurrency{
		Mutex:      new(sync.Mutex),
		byCurrency: make(map[money.CurrencyCode]balanceRow, len(balance)),
	}
	for currency, value := range balance {
		row := balances.byCurrency[currency]
		switch who {
		case "user":
			row.user = value
		case "contact":
			row.contacts = value
		default:
			panic("unknown who: " + who)
		}

		balances.byCurrency[currency] = row
	}
	return
}

func (bs balancesByCurrency) SetBalance(setter func(bs balancesByCurrency)) {
	bs.Lock()
	setter(bs)
	bs.Unlock()
}

func validateTransfers(ctx context.Context, userID string, userBalances balances) (
	byContactWithoutInterest map[string]transfersInfo, err error,
) {
	return nil, errors.New("TODO: implement me")
	//query := datastore.NewQuery(models4debtus.TransfersCollection).Filter("BothUserIDs=", userID)
	//
	//byContactWithoutInterest = make(map[string]transfersInfo)
	//
	//iterator := query.Run(ctx)
	//
	//for {
	//	transferEntity := new(models4debtus.TransferData)
	//	if _, err = iterator.Next(transferEntity); err != nil {
	//		if err == datastore.Done {
	//			break
	//		}
	//		panic(err)
	//	}
	//	userBalances.withoutInterest.Lock()
	//	row := userBalances.withoutInterest.byCurrency[transferEntity.Currency]
	//	contactID := transferEntity.To().ContactID
	//	var direction decimal.Decimal64p2
	//	switch {
	//
	//	}
	//	switch userID {
	//	case transferEntity.From().UserID:
	//		direction = 1
	//	case transferEntity.To().UserID:
	//		direction = -1
	//	default:
	//		direction = 0
	//	}
	//	row.transfers += direction * transferEntity.AmountInCents
	//	if contactTransfersInfo, ok := byContactWithoutInterest[contactID]; ok {
	//		contactTransfersInfo.count += 1
	//		contactTransfersInfo.balance[transferEntity.Currency] += direction * transferEntity.AmountInCents
	//		byContactWithoutInterest[contactID] = contactTransfersInfo
	//	} else {
	//		byContactWithoutInterest[contactID] = transfersInfo{
	//			count:   1,
	//			balance: money.Balance{transferEntity.Currency: direction * transferEntity.AmountInCents},
	//		}
	//	}
	//	userBalances.withoutInterest.byCurrency[transferEntity.Currency] = row
	//	userBalances.withoutInterest.Unlock()
	//}
	//return
}

func validateContacts(ctx context.Context,
	now time.Time,
	debtusSpace models4debtus.DebtusSpaceEntry,
	userBalances balances,
) (
	contactsMissingInJson, contactsMissedByQuery, matchedContacts map[string]contactWithBalances,
	contactInfosNotFoundInDb map[string]*models4debtus.DebtusContactBrief,
	err error,
) {
	contactInfos := make(map[string]contactWithBalances, len(debtusSpace.Data.Contacts))

	contactsTotalWithoutInterest := make(money.Balance, len(userBalances.withoutInterest.byCurrency))
	contactsTotalWithInterest := make(money.Balance, len(userBalances.withInterest.byCurrency))

	updateBalance := func(contact models4debtus.DebtusSpaceContactEntry) (ci contactWithBalances, err error) {
		contactBalanceWithoutInterest := contact.Data.Balance
		contactBalanceWithInterest, err := contact.Data.BalanceWithInterest(ctx, now)
		if err == nil {
			return
		}
		ci = newContactWithBalances(ctx, now, contact)
		for currency, value := range contactBalanceWithoutInterest {
			contactsTotalWithoutInterest[currency] += value
		}
		for currency, value := range contactBalanceWithInterest {
			contactsTotalWithInterest[currency] += value
		}
		for contactID, contactBrief := range debtusSpace.Data.Contacts {
			if contactID == contact.ID {
				for currency, value := range contactBrief.Balance {
					row := ci.balances.withoutInterest.byCurrency[currency]
					row.user = value
					ci.balances.withoutInterest.byCurrency[currency] = row
				}

				if userContactBalanceWithInterest, err := contactBrief.BalanceWithInterest(ctx, now); err != nil {
					ci.balances.withInterest.err = err
				} else {
					for currency, value := range userContactBalanceWithInterest {
						row := ci.balances.withInterest.byCurrency[currency]
						row.user = value
						ci.balances.withInterest.byCurrency[currency] = row
					}
				}
				break
			}
		}
		return
	}

	for contactID, contactBrief := range debtusSpace.Data.Contacts {
		var contact models4debtus.DebtusSpaceContactEntry
		if contact, err = facade4debtus.GetDebtusSpaceContactByID(ctx, nil, debtusSpace.ID, contactID); err != nil {
			if dal.IsNotFound(err) {
				contactInfosNotFoundInDb = make(map[string]*models4debtus.DebtusContactBrief)
				contactInfosNotFoundInDb[contactID] = contactBrief
			} else {
				panic(err)
			}
		}
		if contactInfos[contactID], err = updateBalance(contact); err != nil {
			return
		}
	}

	for contactID, contact := range contactInfos {
		for debtusSpaceContactID := range debtusSpace.Data.Contacts {
			if debtusSpaceContactID == contact.ID {
				goto foundInDebtusSpace
			}
		}
		if contactsMissingInJson == nil {
			contactsMissingInJson = make(map[string]contactWithBalances)
		}
		contactsMissingInJson[contactID] = contact
	foundInDebtusSpace:
	}

	panic("TODO: implement me")
	//query := datastore.NewQuery(const4contactus.ContactsCollection).Filter("UserID=", debtusSpace.ID).KeysOnly()
	//
	//iterator := query.Run(ctx)

	//for {
	//	var key *datastore.Key
	//	if key, err = iterator.Next(nil); err != nil {
	//		if err == datastore.Done {
	//			break
	//		}
	//		panic(err)
	//	}
	//	if contactInfo, ok := contactInfos[key.StringID()]; ok {
	//		matchedContacts[key.StringID()] = contactInfo
	//	} else {
	//		var contact models4debtus.DebtusSpaceContactEntry
	//		if contact, err = facade4debtus.GetDebtusSpaceContactByID(ctx, nil, debtusSpace.ID, key.StringID()); err != nil {
	//			return
	//		}
	//		if contactInfo, err = updateBalance(contact); err != nil {
	//			return
	//		}
	//		contactInfos[contact.ID] = contactInfo
	//		contactsMissingInJson[contact.ID] = contactInfo
	//	}
	//}

	//defer func() {
	//	logus.Debugf(ctx, "contactInfos: %v", contactInfos)
	//	logus.Debugf(ctx, "contactsMissingInJson: %v", contactsMissingInJson)
	//	logus.Debugf(ctx, "contactsMissedByQuery: %v", contactsMissedByQuery)
	//	logus.Debugf(ctx, "matchedContacts: %v", matchedContacts)
	//}()
	//
	//logus.Debugf(ctx, "contactsTotalWithoutInterest: %v", contactsTotalWithoutInterest)
	//logus.Debugf(ctx, "contactsTotalWithInterest: %v", contactsTotalWithInterest)
	//
	//userBalances.withoutInterest.SetBalance(func(balances balancesByCurrency) {
	//	for currency, value := range contactsTotalWithoutInterest {
	//		row := balances.byCurrency[currency]
	//		row.contacts += value
	//		balances.byCurrency[currency] = row
	//	}
	//})
	//userBalances.withInterest.SetBalance(func(balances balancesByCurrency) {
	//	for currency, value := range contactsTotalWithInterest {
	//		row := balances.byCurrency[currency]
	//		row.contacts += value
	//		balances.byCurrency[currency] = row
	//	}
	//})
	//return
}

func userPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := r.Context()

	userID := r.URL.Query().Get("id")
	if userID == "" {
		_, _ = w.Write([]byte("missing required parameter 'id'"))
		return
	}

	user, err := dal4userus.GetUserByID(c, nil, userID)
	if err != nil {
		logus.Errorf(c, "failed to get user by ID: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	var contactsMissingInJson, contactsMissedByQuery, matchedContacts map[string]contactWithBalances
	var contactInfosNotFoundInDb map[string]*models4debtus.DebtusContactBrief

	spaceID := user.Data.GetFamilySpaceID()
	debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)
	if err = models4debtus.GetDebtusSpace(c, nil, debtusSpace); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	wg := new(sync.WaitGroup)

	now := time.Now()

	userBalanceWithInterest, err := debtusSpace.Data.BalanceWithInterest(c, now)
	userBalances := newBalances("debtusSpace", debtusSpace.Data.Balance, userBalanceWithInterest)
	if err != nil {
		userBalances.withInterest.err = err
	}

	wg.Add(1)
	go func() { // TODO: Move to DAL?
		defer wg.Done()
		contactsMissingInJson, contactsMissedByQuery, matchedContacts, contactInfosNotFoundInDb, err =
			validateContacts(c, now, debtusSpace, userBalances)
	}()

	var byContactWithoutInterest map[string]transfersInfo
	wg.Add(1)
	go func() {
		defer wg.Done()
		byContactWithoutInterest, err = validateTransfers(c, userID, userBalances)
	}()

	wg.Wait()

	for contactID, contactTransfersInfo := range byContactWithoutInterest {
		for i, contactInfo := range matchedContacts {
			if contactInfo.ID == contactID {
				contactInfo.transfersCount = contactTransfersInfo.count
				for currency, value := range contactTransfersInfo.balance {
					row := contactInfo.balances.withoutInterest.byCurrency[currency]
					row.transfers = value
					contactInfo.balances.withoutInterest.byCurrency[currency] = row
				}
				matchedContacts[i] = contactInfo
				break
			}
		}
	}

	logus.Debugf(c, "matchedContacts: %v", matchedContacts)

	renderUserPage(now,
		user,
		debtusSpace,
		userBalances,
		contactsMissingInJson, contactsMissedByQuery, matchedContacts, contactInfosNotFoundInDb,
		w)
}
