package support

import (
	"context"
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2"
	"net/http"
	"reflect"
	"time"
)

func ValidateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	doFixes := r.URL.Query().Get("fix") == "all"
	spaceID := r.URL.Query().Get("spaceID")
	userID := r.URL.Query().Get("id")
	if userID == "" {
		logus.Errorf(ctx, "UserEntry ContactID is empty")
		return
	}
	user := dbo4userus.NewUserEntry(userID)
	var db dal.DB
	var err error
	if db, err = facade.GetDatabase(ctx); err != nil {
		logus.Errorf(ctx, "Failed to get database: %v", err)
		return
	}
	if err = db.Get(ctx, user.Record); err != nil {
		if dal.IsNotFound(err) {
			logus.Errorf(ctx, "UserEntry not found by key: %v", err)
		} else {
			logus.Errorf(ctx, "Failed to get user by key=%v: %v", user.Key, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	query := dal.From(const4contactus.ContactsCollection).WhereField("UserID", dal.Equal, userID).SelectInto(func() dal.Record {
		return dal.NewRecordWithIncompleteKey(dbo4userus.UsersCollection, reflect.Int64, new(dbo4userus.UserDbo))
	})
	userCounterpartyRecords, err := db.QueryAllRecords(ctx, query)
	if err != nil {
		logus.Errorf(ctx, "Failed to load user counterparties: %v", err)
		return
	}

	//slices.Sort(userCounterpartyIDs)

	counterpartyIDs := make([]string, len(userCounterpartyRecords))
	for i, v := range userCounterpartyRecords {
		counterpartyIDs[i] = v.Key().ID.(string)
	}
	//slices.Sort(counterpartyIDs)

	query = dal.From(models4debtus.TransfersCollection).WhereField("BothUserIDs", dal.Equal, userID).OrderBy(dal.AscendingField("DtCreated")).SelectInto(func() dal.Record {
		return dal.NewRecordWithIncompleteKey(models4debtus.AppUserKind, reflect.Int64, new(models4debtus.DebutsAppUserDataOBSOLETE))
	})

	transferRecords, err := db.QueryAllRecords(ctx, query)

	if err != nil {
		logus.Errorf(ctx, "Failed to load api4transfers: %v", err)
		return
	}

	type transfersInfo struct {
		Count  int
		LastID string
		LastAt time.Time
	}

	transfersInfoByCounterparty := make(map[string]transfersInfo, len(counterpartyIDs))

	for _, transferRecord := range transferRecords {
		transferEntity := transferRecord.Data().(*models4debtus.TransferData)
		counterpartyInfo := transferEntity.CounterpartyInfoByUserID(userID)
		counterpartyTransfersInfo := transfersInfoByCounterparty[counterpartyInfo.ContactID]
		counterpartyTransfersInfo.Count += 1
		if counterpartyTransfersInfo.LastAt.Before(transferEntity.DtCreated) {
			counterpartyTransfersInfo.LastAt = transferEntity.DtCreated
			counterpartyTransfersInfo.LastID = transferRecord.Key().ID.(string)
		}
		transfersInfoByCounterparty[counterpartyInfo.ContactID] = counterpartyTransfersInfo
	}

	//fixUserCounterparties := func() {
	//	var txUser models4debtus.AppUserOBSOLETE
	//	err := facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
	//		logus.Debugf(ctx, "Transaction started..")
	//		txUser = models4debtus.NewAppUserOBSOLETE(userID, nil)
	//		if err := tx.Get(ctx, txUser.Record); err != nil {
	//			return err
	//		}
	//		if txUser.Data.SavedCounter != user.Data.SavedCounter {
	//			return fmt.Errorf("user changed since last load: txUser.SavedCounter:%v != user.SavedCounter:%v", txUser.Data.SavedCounter, user.Data.SavedCounter)
	//		}
	//		txUser.Data.ContactsJson = ""
	//		for _, counterpartyRecord := range userCounterpartyRecords {
	//			counterpartyEntity := counterpartyRecord.Data().(*models4debtus.DebtusSpaceContactDbo)
	//			counterpartyID := counterpartyRecord.Key().ContactID.(string)
	//			if counterpartyTransfersInfo, ok := transfersInfoByCounterparty[counterpartyID]; ok {
	//				counterpartyEntity.LastTransferAt = counterpartyTransfersInfo.LastAt
	//				counterpartyEntity.LastTransferID = counterpartyTransfersInfo.LastID
	//				counterpartyEntity.CountOfTransfers = counterpartyTransfersInfo.Count
	//			} else {
	//				counterpartyEntity.CountOfTransfers = 0
	//				counterpartyEntity.LastTransferAt = time.Time{}
	//				counterpartyEntity.LastTransferID = ""
	//			}
	//			models4debtus.AddOrUpdateDebtusContact(&txUser, models4debtus.NewDebtusSpaceContactEntry(counterpartyID, counterpartyEntity))
	//		}
	//		if err = tx.Set(ctx, txUser.Record); err != nil {
	//			return fmt.Errorf("failed to save fixed user: %w", err)
	//		}
	//		return nil
	//	}, nil)
	//	if err != nil {
	//		logus.Errorf(ctx, "Failed to fix user.CounterpartyIDs: %v", err)
	//		return
	//	}
	//	logus.Infof(ctx, "Fixed user.ContactsJson\n\tfrom: %v\n\tto: %v", user.Data.ContactsJson, txUser.Data.ContactsJson)
	//	user = txUser
	//}

	//if len(userCounterpartyIDs) != len(counterpartyIDs) {
	//	logus.Warningf(ctx, "len(userCounterpartyIDs) != len(counterpartyIDs) => %v != %v", len(userCounterpartyIDs), len(counterpartyIDs))
	//	if doFixes {
	//		fixUserCounterparties()
	//	} else {
	//		return // Do not continue if counterparties are not in sync
	//	}
	//} else {
	//	for i, v := range userCounterpartyIDs {
	//		if counterpartyIDs[i] != v {
	//			logus.Warningf(ctx, "user.CounterpartyIDs != counterpartyKeys\n\tuserCounterpartyIDs:\n\t\t%v\n\tcounterpartyIDs:\n\t\t%v", userCounterpartyIDs, counterpartyIDs)
	//			if doFixes {
	//				fixUserCounterparties()
	//				break
	//			} else {
	//				return // Do not continue if counterparties are not in sync
	//			}
	//		}
	//	}
	//}
	//logus.Infof(ctx, "OK - UserEntry ContactsJson is OK")

	// We need counterparties by ContactID to check balance against api4transfers
	counterpartiesByID := make(map[int64]*models4debtus.DebtusSpaceContactDbo, len(counterpartyIDs))
	for _, counterpartyRecord := range userCounterpartyRecords {
		counterpartiesByID[counterpartyRecord.Key().ID.(int64)] = counterpartyRecord.Data().(*models4debtus.DebtusSpaceContactDbo)
	}

	debtusUser := models4debtus.NewDebtusUserEntry(userID)

	if err = db.Get(ctx, debtusUser.Record); err != nil {
		return
	}

	if len(transferRecords) > 0 && debtusUser.Data.LastTransferID == "" {
		if doFixes {
			var txUser dbo4userus.UserEntry
			err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
				if err = tx.Get(ctx, debtusUser.Record); err != nil {
					return err
				}
				if debtusUser.Data.LastTransferID == "" {
					i := len(transferRecords) - 1
					debtusUser.Data.LastTransferID = transferRecords[i].Key().ID.(string)
					debtusUser.Data.LastTransferAt = transferRecords[i].Data().(*models4debtus.TransferData).DtCreated
					err = tx.Set(ctx, txUser.Record)
					return err
				}
				return nil
			}, nil)
			if err != nil {
				logus.Errorf(ctx, "Failed to update user.LastTransferID")
			} else {
				logus.Infof(ctx, "Fixed user.LastTransferID")
				user = txUser
			}
		} else {
			logus.Warningf(ctx, "user.LastTransferID is not set")
		}
	}

	// Get stored user total balance
	transfersBalanceByCounterpartyID := make(map[string]money.Balance, len(counterpartyIDs))

	for i, transferRecord := range transferRecords {
		transferData := transferRecord.Data().(*models4debtus.TransferData)
		var counterpartyID string
		switch userID {
		case transferData.CreatorUserID:
			counterpartyID = transferData.Counterparty().ContactID
		case transferData.Counterparty().UserID:
			counterpartyID = transferData.Creator().ContactID
		default:
			logus.Errorf(ctx, "userID=%v is NOT equal to transferData.CreatorUserID=%v or transferData.DebtusSpaceContactEntry().UserID=%v", userID, transferData.CreatorUserID, transferData.Counterparty().UserID)
			return
		}
		transfersCounterpartyBalance, ok := transfersBalanceByCounterpartyID[counterpartyID]
		if !ok {
			transfersCounterpartyBalance = make(money.Balance)
			transfersBalanceByCounterpartyID[counterpartyID] = transfersCounterpartyBalance
		}
		value := transferData.GetAmount().Value
		currency := money.CurrencyCode(transferData.Currency)
		switch transferData.DirectionForUser(userID) {
		case models4debtus.TransferDirectionUser2Counterparty:
			transfersCounterpartyBalance[currency] += value
		case models4debtus.TransferDirectionCounterparty2User:
			transfersCounterpartyBalance[currency] -= value
		default:
			logus.Errorf(ctx, "TransferEntry %v has unknown direction: %v", transferRecords[i].Key().ID, transferData.DirectionForUser(userID))
			return
		}
	}

	//logus.Debugf(ctx, "transfersBalanceByCounterpartyID: %v", transfersBalanceByCounterpartyID)

	transfersTotalBalance := make(money.Balance)
	for _, transfersCounterpartyBalance := range transfersBalanceByCounterpartyID {
		for currency, value := range transfersCounterpartyBalance {
			if value == 0 {
				delete(transfersCounterpartyBalance, currency)
			} else {
				transfersTotalBalance[currency] += value
			}
		}
	}

	for currency, value := range transfersTotalBalance {
		if value == 0 {
			delete(transfersTotalBalance, currency)
		}
	}

	if len(debtusUser.Data.Balance) != len(transfersTotalBalance) {
		logus.Warningf(ctx, "len(userBalance) != len(transfersTotalBalance) =>\n\t%d: %v\n\t%d: %v", len(debtusUser.Data.Balance), debtusUser.Data.Balance, len(transfersTotalBalance), transfersTotalBalance)
	}

	userBalanceIsOK := true

	for currency, userVal := range debtusUser.Data.Balance {
		if transfersVal, ok := transfersTotalBalance[currency]; !ok {
			logus.Warningf(ctx, "UserEntry has %v=%v balance but no corresponding api4transfers' balance.", currency, userVal)
			userBalanceIsOK = false
		} else if transfersVal != userVal {
			logus.Warningf(ctx, "Currency(%v) UserEntry balance %v not equal to api4transfers' balance %v", currency, userVal, transfersVal)
			userBalanceIsOK = false
		}
	}

	for currency, transfersVal := range transfersTotalBalance {
		if _, ok := debtusUser.Data.Balance[currency]; !ok {
			logus.Warningf(ctx, "Transfers has %v=%v balance but no corresponding user balance.", currency, transfersVal)
			userBalanceIsOK = false
		}
	}

	if userBalanceIsOK {
		logus.Infof(ctx, "OK - UserEntry.Balance() is matching to %v api4transfers' balance.", len(transferRecords))
	} else {
		logus.Warningf(ctx, "Calculated balance for %v user api4transfers does not match user's total balance.", len(transferRecords))
		if !doFixes {
			logus.Debugf(ctx, "Pass fix=all to fix user balance")
		} else {
			err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
				txUser := dbo4userus.NewUserEntry(userID)
				spaceID := txUser.Data.GetFamilySpaceID()
				debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)
				if err := tx.Get(ctx, txUser.Record); err != nil {
					return err
				}
				if !debtusSpace.Data.Balance.Equal(debtusSpace.Data.Balance) {
					return errors.New("user changed: !reflect.DeepEqual(txUser.Balance(), user.Balance())")
				}

				debtusSpace.Data.Balance = transfersTotalBalance
				if err = tx.Set(ctx, txUser.Record); err != nil {
					return fmt.Errorf("failed to save user with fixed balance: %w", err)
				}
				return nil
			}, nil)
			if err != nil {
				err = fmt.Errorf("failed to fix user balance: %w", err)
				logus.Errorf(ctx, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			logus.Infof(ctx, "Fixed user balance")
		}
	}

	var counterpartyIDsWithMatchingBalance, counterpartyIDsWithNonMatchingBalance []string

	for _, counterpartyRecord := range userCounterpartyRecords {
		counterpartyKey := counterpartyRecord.Key()
		counterpartyID := counterpartyKey.ID.(string)
		counterparty := counterpartyRecord.Data().(*models4debtus.DebtusSpaceContactDbo)

		if transfersCounterpartyBalance := transfersBalanceByCounterpartyID[counterpartyID]; (len(transfersCounterpartyBalance) == 0 && len(counterparty.Balance) == 0) || reflect.DeepEqual(transfersCounterpartyBalance, counterparty.Balance) {
			counterpartyIDsWithMatchingBalance = append(counterpartyIDsWithMatchingBalance, counterpartyID)
		} else {
			counterpartyIDsWithNonMatchingBalance = append(counterpartyIDsWithNonMatchingBalance, counterpartyID)
			logus.Warningf(ctx, "DebtusSpaceContactEntry ContactID=%v has balance not matching api4transfers' balance:\n\tDebtusSpaceContactEntry: %v\n\tTransfers: %v", counterpartyID, counterparty.Balance, transfersCounterpartyBalance)
			if doFixes {
				//var txCounterparty models.DebtusSpaceContactEntry
				var db dal.DB
				if db, err = facade.GetDatabase(ctx); err != nil {
					logus.Errorf(ctx, "Failed to get database: %v", err)
					return
				}
				err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
					txCounterparty := models4debtus.NewDebtusSpaceContactEntry(spaceID, counterpartyKey.ID.(string), nil)
					if err := tx.Get(ctx, txCounterparty.Record); err != nil {
						return err
					}
					if !txCounterparty.Data.Balance.Equal(counterparty.Balance) {
						return errors.New("contact changed since check: !reflect.DeepEqual(txCounterparty.Balance(), counterparty.Balance())")
					}

					txCounterparty.Data.Balance = transfersCounterpartyBalance
					if err = tx.Set(ctx, txCounterparty.Record); err != nil {
						return fmt.Errorf("failed to save counterparty with ContactID=%v: %w", counterpartyID, err)
					}
					return nil
				}, nil)
				if err != nil {
					logus.Errorf(ctx, "Failed to fix counterparty with ContactID=%v: %v", counterpartyID, err)
				} else {
					logus.Infof(ctx, "Fixed counterparty with ContactID=%v", counterpartyID)
					//userCounterpartyRecords[i] = txCounterparty.Data
				}
			}
		}
	}
	if len(counterpartyIDsWithMatchingBalance) > 0 {
		logus.Infof(ctx, "There are %v counterparties with balance matching to api4transfers: %v", len(counterpartyIDsWithMatchingBalance), counterpartyIDsWithMatchingBalance)
	}
	if len(counterpartyIDsWithNonMatchingBalance) > 0 {
		logus.Warningf(ctx, "There are %v counterparties with balance NOT matching to api4transfers: %v", len(counterpartyIDsWithNonMatchingBalance), counterpartyIDsWithNonMatchingBalance)
	}
}
