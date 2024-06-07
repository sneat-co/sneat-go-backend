package gaedal

import (
	"context"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/log"
	"reflect"
)

func Settle2members(c context.Context, groupID, debtorID, sponsorID string, currency money.CurrencyCode, amount decimal.Decimal64p2) (err error) {
	log.Debugf(c, "Settle2members(groupID=%v, debtorID=%v, sponsorID=%v, currency=%v, amount=%v)", groupID, debtorID, sponsorID, currency, amount)
	query := dal.From(models.BillKind).
		WhereField("GetUserGroupID", dal.Equal, groupID).
		WhereField("Currency", dal.Equal, string(currency)).
		WhereField("DebtorIDs", dal.Equal, debtorID).
		WhereField("SponsorIDs", dal.Equal, sponsorID).
		OrderBy(dal.AscendingField("DtCreated")).
		Limit(20).
		SelectKeysOnly(reflect.String)

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	var reader dal.Reader

	if reader, err = db.QueryReader(c, query); err != nil {
		return err
	}
	var ids []string
	if ids, err = dal.SelectAllIDs[string](reader, dal.WithLimit(query.Limit())); err != nil {
		return err
	}

	if len(ids) == 0 {
		log.Errorf(c, "No bills found to settle")
		return
	} else {
		log.Debugf(c, "ids: %+v", ids)
	}

	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var (
			group                     models.GroupEntry
			groupDebtor, groupSponsor models.GroupMemberJson
		)
		if group, err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
			return
		}

		billsSettlement := models.BillsHistory{
			Data: &models.BillsHistoryDbo{
				Action:                 models.BillHistoryActionSettled,
				Currency:               currency,
				GroupMembersJsonBefore: group.Data.MembersJson,
			},
		}

		if groupDebtor, err = group.Data.GetGroupMemberByID(debtorID); err != nil {
			return fmt.Errorf("unknown debtorID=%s: %w", debtorID, err)
		}
		if groupSponsor, err = group.Data.GetGroupMemberByID(sponsorID); err != nil {
			return fmt.Errorf("unknown sponsorID=%s: %w", sponsorID, err)
		}

		if v, ok := groupDebtor.Balance[currency]; !ok {
			return fmt.Errorf("group debtor has no balance in currency=%v", currency)
		} else if -v < amount {
			log.Warningf(c, "Debtor balance is less then settling amount")
			amount = -v
		}
		if v, ok := groupSponsor.Balance[currency]; !ok {
			return fmt.Errorf("group sponsor has no balance in currency=%v", currency)
		} else if v < amount {
			log.Warningf(c, "sponsor balance is less then settling amount")
			amount = v
		}

		billsToSave := make([]models.Bill, 0, len(ids))

		settlementBills := make([]models.BillSettlementJson, 0, len(ids))

		for _, id := range ids {
			if amount == 0 {
				break
			} else if amount < 0 {
				panic(fmt.Sprintf("amount < 0: %v", amount))
			}
			var bill models.Bill
			if bill, err = facade.GetBillByID(c, tx, id); err != nil {
				return
			}
			billMembers := bill.Data.GetBillMembers()
			var debtor, sponsor *models.BillMemberJson
			var debtorInvertedBalance, diff decimal.Decimal64p2
			for i := range billMembers {
				switch billMembers[i].ID {
				case debtorID:
					if debtor = &billMembers[i]; debtor.Balance() >= 0 {
						log.Warningf(c, "Got debtor %v with positive balance = %v", debtor.ID, debtor.Balance())
						goto nextBill
					}
					if sponsor != nil {
						break
					}
				case sponsorID:
					if sponsor = &billMembers[i]; sponsor.Balance() <= 0 {
						log.Warningf(c, "Got sponsor %v with negative balance = %v", sponsor.ID, sponsor.Balance())
						goto nextBill
					}
					if debtor != nil {
						break
					}
				}
			}
			if debtor == nil {
				log.Warningf(c, "Debtor not found by ID="+debtorID)
				goto nextBill
			}
			if sponsor == nil {
				log.Warningf(c, "Sponsor not found by ID="+sponsorID)
				goto nextBill
			}
			debtorInvertedBalance = -1 * debtor.Balance()
			if debtorInvertedBalance <= sponsor.Balance() {
				diff = debtorInvertedBalance
			} else {
				diff = sponsor.Balance()
			}

			if diff > amount {
				diff = amount
			}

			log.Debugf(c, "diff: %v", diff)
			amount -= diff
			billsSettlement.Data.TotalAmountDiff += diff

			debtor.Paid += diff
			sponsor.Paid -= diff

			groupDebtor.Balance[currency] += diff
			if groupDebtor.Balance[currency] == 0 {
				delete(groupDebtor.Balance, currency)
			}
			groupSponsor.Balance[currency] -= diff
			if groupSponsor.Balance[currency] == 0 {
				delete(groupSponsor.Balance, currency)
			}

			if err = bill.Data.SetBillMembers(billMembers); err != nil {
				return
			}

			log.Debugf(c, "groupDebtor.Balance: %v", groupDebtor.Balance)
			log.Debugf(c, "groupSponsor.Balance: %v", groupSponsor.Balance)

			billsToSave = append(billsToSave, bill)
			settlementBills = append(settlementBills, models.BillSettlementJson{
				BillID:    bill.ID,
				GroupID:   groupID,
				DebtorID:  debtorID,
				SponsorID: sponsorID,
				Amount:    diff,
			})

		nextBill:
		}

		if len(billsToSave) > 0 {
			billsSettlement.Data.SetBillSettlements(groupID, settlementBills)
			if err = dtdal.InsertWithRandomStringID(c, tx, billsSettlement.Record); err != nil {
				return
			}
			toSave := make([]dal.Record, len(billsToSave)+1)
			toSave[0] = group.Record
			for i, bill := range billsToSave {
				bill.Data.SettlementIDs = append(bill.Data.SettlementIDs, billsSettlement.ID)
				toSave[i+1] = bill.Record
			}

			groupMembers := group.Data.GetGroupMembers()
			for i, m := range groupMembers {
				switch m.ID {
				case debtorID:
					groupMembers[i] = groupDebtor
				case sponsorID:
					groupMembers[i] = groupSponsor
				}
			}
			if changed := group.Data.SetGroupMembers(groupMembers); !changed {
				panic("GroupEntry members not changed - something wrong")
			}
			if err = tx.SetMulti(c, toSave); err != nil {
				return
			}
			billsSettlement.Data.GroupMembersJsonAfter = group.Data.MembersJson
			if err = tx.Set(c, billsSettlement.Record); err != nil {
				return
			}
		} else {
			log.Errorf(c, "No bills found to settle")
		}

		return
	})

	return
}
