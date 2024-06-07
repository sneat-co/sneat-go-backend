package inspector

import (
	"context"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"google.golang.org/appengine/v2"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/log"
	"google.golang.org/appengine/v2/datastore"
)

type contactWithBalances struct {
	models.ContactEntry
	transfersCount int
	balances       balances
}

type transfersInfo struct {
	count   int
	balance money.Balance
}

func newContactWithBalances(c context.Context, now time.Time, contact models.ContactEntry) contactWithBalances {
	balanceWithInterest, err := contact.Data.BalanceWithInterest(c, now)
	result := contactWithBalances{
		ContactEntry: contact,
		balances:     newBalances("contact", contact.Data.Balance(), balanceWithInterest),
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

func validateTransfers(c context.Context, userID string, userBalances balances) (
	byContactWithoutInterest map[string]transfersInfo, err error,
) {
	query := datastore.NewQuery(models.TransfersCollection).Filter("BothUserIDs=", userID)

	byContactWithoutInterest = make(map[string]transfersInfo)

	iterator := query.Run(c)

	for {
		transferEntity := new(models.TransferData)
		if _, err = iterator.Next(transferEntity); err != nil {
			if err == datastore.Done {
				break
			}
			panic(err)
		}
		userBalances.withoutInterest.Lock()
		row := userBalances.withoutInterest.byCurrency[transferEntity.Currency]
		contactID := transferEntity.To().ContactID
		var direction decimal.Decimal64p2
		switch {

		}
		switch userID {
		case transferEntity.From().UserID:
			direction = 1
		case transferEntity.To().UserID:
			direction = -1
		default:
			direction = 0
		}
		row.transfers += direction * transferEntity.AmountInCents
		if contactTransfersInfo, ok := byContactWithoutInterest[contactID]; ok {
			contactTransfersInfo.count += 1
			contactTransfersInfo.balance[transferEntity.Currency] += direction * transferEntity.AmountInCents
			byContactWithoutInterest[contactID] = contactTransfersInfo
		} else {
			byContactWithoutInterest[contactID] = transfersInfo{
				count:   1,
				balance: money.Balance{transferEntity.Currency: direction * transferEntity.AmountInCents},
			}
		}
		userBalances.withoutInterest.byCurrency[transferEntity.Currency] = row
		userBalances.withoutInterest.Unlock()
	}
	return
}

func validateContacts(c context.Context,
	now time.Time,
	user models.AppUser,
	userBalances balances,
) (
	contactsMissingInJson, contactsMissedByQuery, matchedContacts []contactWithBalances,
	contactInfosNotFoundInDb []models.UserContactJson,
	err error,
) {
	userContactsJson := user.Data.Contacts()
	contactInfos := make([]contactWithBalances, len(userContactsJson))
	contactInfosByID := make(map[string]contactWithBalances, len(contactInfos))

	contactsTotalWithoutInterest := make(money.Balance, len(userBalances.withoutInterest.byCurrency))
	contactsTotalWithInterest := make(money.Balance, len(userBalances.withInterest.byCurrency))

	updateBalance := func(contact models.ContactEntry) (ci contactWithBalances, err error) {
		contactBalanceWithoutInterest := contact.Data.Balance()
		contactBalanceWithInterest, err := contact.Data.BalanceWithInterest(c, now)
		if err == nil {
			return
		}
		ci = newContactWithBalances(c, now, contact)
		for currency, value := range contactBalanceWithoutInterest {
			contactsTotalWithoutInterest[currency] += value
		}
		for currency, value := range contactBalanceWithInterest {
			contactsTotalWithInterest[currency] += value
		}
		for _, userContactJson := range userContactsJson {
			if userContactJson.ID == contact.ID {
				for currency, value := range userContactJson.Balance() {
					row := ci.balances.withoutInterest.byCurrency[currency]
					row.user = value
					ci.balances.withoutInterest.byCurrency[currency] = row
				}

				if userContactBalanceWithInterest, err := userContactJson.BalanceWithInterest(c, now); err != nil {
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

	for i, userContactInfo := range userContactsJson {
		var contact models.ContactEntry
		if contact, err = facade.GetContactByID(c, nil, userContactInfo.ID); err != nil {
			if dal.IsNotFound(err) {
				contactInfosNotFoundInDb = append(contactInfosNotFoundInDb, userContactInfo)
			} else {
				panic(err)
			}
		}
		if contactInfos[i], err = updateBalance(contact); err != nil {
			return
		}
		contactInfosByID[contact.ID] = contactInfos[i]
	}

	for _, contact := range contactInfos {
		for _, userContact := range userContactsJson {
			if userContact.ID == contact.ID {
				goto foundInUserJson
			}
		}
		contactsMissingInJson = append(contactsMissingInJson, contact)
	foundInUserJson:
	}

	query := datastore.NewQuery(models.DebtusContactsCollection).Filter("UserID=", user.ID).KeysOnly()

	iterator := query.Run(c)

	for {
		var key *datastore.Key
		if key, err = iterator.Next(nil); err != nil {
			if err == datastore.Done {
				break
			}
			panic(err)
		}
		if contactInfo, ok := contactInfosByID[key.StringID()]; ok {
			matchedContacts = append(matchedContacts, contactInfo)
		} else {
			var contact models.ContactEntry
			if contact, err = facade.GetContactByID(c, nil, key.StringID()); err != nil {
				return
			}
			if contactInfo, err = updateBalance(contact); err != nil {
				return
			}
			contactInfos = append(contactInfos, contactInfo)
			contactsMissingInJson = append(contactInfos, contactInfo)
		}
	}

	defer func() {
		log.Debugf(c, "contactInfos: %v", contactInfos)
		log.Debugf(c, "contactsMissingInJson: %v", contactsMissingInJson)
		log.Debugf(c, "contactsMissedByQuery: %v", contactsMissedByQuery)
		log.Debugf(c, "matchedContacts: %v", matchedContacts)
	}()

	log.Debugf(c, "contactsTotalWithoutInterest: %v", contactsTotalWithoutInterest)
	log.Debugf(c, "contactsTotalWithInterest: %v", contactsTotalWithInterest)

	userBalances.withoutInterest.SetBalance(func(balances balancesByCurrency) {
		for currency, value := range contactsTotalWithoutInterest {
			row := balances.byCurrency[currency]
			row.contacts += value
			balances.byCurrency[currency] = row
		}
	})
	userBalances.withInterest.SetBalance(func(balances balancesByCurrency) {
		for currency, value := range contactsTotalWithInterest {
			row := balances.byCurrency[currency]
			row.contacts += value
			balances.byCurrency[currency] = row
		}
	})
	return
}

func userPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)

	userID := r.URL.Query().Get("id")
	if userID == "" {
		_, _ = w.Write([]byte("missing required parameter 'id'"))
		return
	}

	var (
		user                                                          models.AppUser
		contactsMissingInJson, contactsMissedByQuery, matchedContacts []contactWithBalances
		contactInfosNotFoundInDb                                      []models.UserContactJson
		err                                                           error
	)

	if user, err = facade.User.GetUserByID(c, nil, userID); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	wg := new(sync.WaitGroup)

	now := time.Now()

	userBalanceWithInterest, err := user.Data.BalanceWithInterest(c, now)
	userBalances := newBalances("user", user.Data.Balance(), userBalanceWithInterest)
	if err != nil {
		userBalances.withInterest.err = err
	}

	wg.Add(1)
	go func() { // TODO: Move to DAL?
		defer wg.Done()
		contactsMissingInJson, contactsMissedByQuery, matchedContacts, contactInfosNotFoundInDb, err = validateContacts(c, now, user, userBalances)
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

	log.Debugf(c, "matchedContacts: %v", matchedContacts)

	renderUserPage(now,
		user,
		userBalances,
		contactsMissingInJson, contactsMissedByQuery, matchedContacts, contactInfosNotFoundInDb,
		w)
}
