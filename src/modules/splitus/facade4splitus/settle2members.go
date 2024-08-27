package facade4splitus

import (
	"context"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
	"github.com/strongo/logus"
	"reflect"
)

func Settle2members(ctx context.Context, spaceID, debtorID, sponsorID string, currency money.CurrencyCode, amount decimal.Decimal64p2) (err error) {
	logus.Debugf(ctx, "Settle2members(spaceID=%v, debtorID=%v, sponsorID=%v, currency=%v, amount=%v)", spaceID, debtorID, sponsorID, currency, amount)
	query := dal.From(models4splitus.BillKind).
		WhereField("GetUserGroupID", dal.Equal, spaceID).
		WhereField("Currency", dal.Equal, string(currency)).
		WhereField("DebtorIDs", dal.Equal, debtorID).
		WhereField("SponsorIDs", dal.Equal, sponsorID).
		OrderBy(dal.AscendingField("DtCreated")).
		Limit(20).
		SelectKeysOnly(reflect.String)

	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return
	}
	var reader dal.Reader

	if reader, err = db.QueryReader(ctx, query); err != nil {
		return err
	}
	var ids []string
	if ids, err = dal.SelectAllIDs[string](reader, dal.WithLimit(query.Limit())); err != nil {
		return err
	}

	if len(ids) == 0 {
		logus.Errorf(ctx, "No bills found to settle")
		return
	} else {
		logus.Debugf(ctx, "ids: %+v", ids)
	}

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		splitusSpace := models4splitus.NewSplitusSpaceEntry(spaceID)
		var groupDebtor, groupSponsor briefs4splitus.SpaceSplitMember

		if err = tx.Get(ctx, splitusSpace.Record); err != nil {
			return
		}

		billsSettlement := models4splitus.BillsHistory{
			Data: &models4splitus.BillsHistoryDbo{
				Action:             models4splitus.BillHistoryActionSettled,
				Currency:           currency,
				SplitMembersBefore: splitusSpace.Data.Members,
			},
		}

		if groupDebtor, err = splitusSpace.Data.GetGroupMemberByID(debtorID); err != nil {
			return fmt.Errorf("unknown debtorID=%s: %w", debtorID, err)
		}
		if groupSponsor, err = splitusSpace.Data.GetGroupMemberByID(sponsorID); err != nil {
			return fmt.Errorf("unknown sponsorID=%s: %w", sponsorID, err)
		}

		if v, ok := groupDebtor.Balance[currency]; !ok {
			return fmt.Errorf("splitusSpace debtor has no balance in currency=%v", currency)
		} else if -v < amount {
			logus.Warningf(ctx, "Debtor balance is less then settling amount")
			amount = -v
		}
		if v, ok := groupSponsor.Balance[currency]; !ok {
			return fmt.Errorf("splitusSpace sponsor has no balance in currency=%v", currency)
		} else if v < amount {
			logus.Warningf(ctx, "sponsor balance is less then settling amount")
			amount = v
		}

		billsToSave := make([]models4splitus.BillEntry, 0, len(ids))

		settlementBills := make([]briefs4splitus.BillSettlementJson, 0, len(ids))

		for _, id := range ids {
			if amount == 0 {
				break
			} else if amount < 0 {
				panic(fmt.Sprintf("amount < 0: %v", amount))
			}
			var bill models4splitus.BillEntry
			if bill, err = GetBillByID(ctx, tx, id); err != nil {
				return
			}
			billMembers := bill.Data.GetBillMembers()
			var debtor, sponsor *briefs4splitus.BillMemberBrief
			var debtorInvertedBalance, diff decimal.Decimal64p2
			for i := range billMembers {
				switch billMembers[i].ID {
				case debtorID:
					if debtor = billMembers[i]; debtor.Balance() >= 0 {
						logus.Warningf(ctx, "Got debtor %v with positive balance = %v", debtor.ID, debtor.Balance())
						goto nextBill
					}
					if sponsor != nil {
						break
					}
				case sponsorID:
					if sponsor = billMembers[i]; sponsor.Balance() <= 0 {
						logus.Warningf(ctx, "Got sponsor %v with negative balance = %v", sponsor.ID, sponsor.Balance())
						goto nextBill
					}
					if debtor != nil {
						break
					}
				}
			}
			if debtor == nil {
				logus.Warningf(ctx, "Debtor not found by ContactID="+debtorID)
				goto nextBill
			}
			if sponsor == nil {
				logus.Warningf(ctx, "Sponsor not found by ContactID="+sponsorID)
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

			logus.Debugf(ctx, "diff: %v", diff)
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

			logus.Debugf(ctx, "groupDebtor.Balance: %v", groupDebtor.Balance)
			logus.Debugf(ctx, "groupSponsor.Balance: %v", groupSponsor.Balance)

			billsToSave = append(billsToSave, bill)
			settlementBills = append(settlementBills, briefs4splitus.BillSettlementJson{
				BillID:    bill.ID,
				GroupID:   spaceID,
				DebtorID:  debtorID,
				SponsorID: sponsorID,
				Amount:    diff,
			})

		nextBill:
		}

		if len(billsToSave) > 0 {
			billsSettlement.Data.SetBillSettlements(spaceID, settlementBills)
			if err = dtdal.InsertWithRandomStringID(ctx, tx, billsSettlement.Record); err != nil {
				return
			}
			toSave := make([]dal.Record, len(billsToSave)+1)
			toSave[0] = splitusSpace.Record
			for i, bill := range billsToSave {
				bill.Data.SettlementIDs = append(bill.Data.SettlementIDs, billsSettlement.ID)
				toSave[i+1] = bill.Record
			}

			groupMembers := splitusSpace.Data.GetGroupMembers()
			for i, m := range groupMembers {
				switch m.ID {
				case debtorID:
					groupMembers[i] = groupDebtor
				case sponsorID:
					groupMembers[i] = groupSponsor
				}
			}
			if updates := splitusSpace.Data.SetGroupMembers(groupMembers); len(updates) == 0 {
				panic("GroupEntry members not changed - something wrong")
			}
			if err = tx.SetMulti(ctx, toSave); err != nil {
				return
			}
			billsSettlement.Data.SplitMembersAfter = splitusSpace.Data.Members
			if err = tx.Set(ctx, billsSettlement.Record); err != nil {
				return
			}
		} else {
			logus.Errorf(ctx, "No bills found to settle")
		}

		return
	})

	return
}
