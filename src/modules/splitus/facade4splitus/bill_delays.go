package facade4splitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/const4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/decimal"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"sync"
)

const updateUsersWithBillKeyName = "update-users-with-bill"

func DelayUpdateUsersWithBill(ctx context.Context, billID string, userIDs []string) (err error) {
	return delayerUpdateUsersWithBill.EnqueueWork(ctx, delaying.With(const4splitus.QueueSplitus, updateUsersWithBillKeyName, 0), billID, userIDs)
}

func delayedUpdateUsersWithBill(ctx context.Context, billID string, userIDs []string) (err error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(userIDs))
	for i := range userIDs {
		go func(i int) {
			defer wg.Done()
			if err2 := delayerUpdateUserWithBill.EnqueueWork(ctx, delaying.With(const4splitus.QueueSplitus, updateUserWithBillKeyName, 0), billID, userIDs[i]); err != nil {
				err = err2
			}
		}(i)
	}
	wg.Wait()
	return
}

const updateUserWithBillKeyName = "delayedUpdateUserWithBill"

func delayedUpdateUserWithBill(ctx context.Context, billID, userID string) (err error) {
	logus.Debugf(ctx, "delayedUpdateUserWithBill(billID=%v, userID=%v)", billID, userID)
	var (
		bill                 models4splitus.BillEntry
		wg                   sync.WaitGroup
		billErr              error
		userModuleDboChanged bool
		userIsBillMember     bool
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if bill, billErr = GetBillByID(ctx, nil, billID); err != nil {
			return
		}
	}()

	if err = dal4userus.RunUserModuleWorker(ctx, userID, const4debtus.ModuleID, new(models4splitus.SplitusUserDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserModuleWorkerParams[models4splitus.SplitusUserDbo]) (err error) {
			if err = tx.Get(ctx, params.UserModule.Record); err != nil {
				return err
			}
			wg.Wait()
			if billErr != nil {
				return fmt.Errorf("%w: failed to get bill", billErr)
			} else if bill.Data == nil {
				return errors.New("bill.BillDbo == nil")
			}
			var userBillBalance decimal.Decimal64p2
			if bill.Data.Status != models4splitus.BillStatusDeleted {
				for _, billMember := range bill.Data.GetBillMembers() {
					if billMember.UserID == userID {
						userBillBalance = billMember.Balance()
						userIsBillMember = true
						logus.Debugf(ctx, "userBillBalance: %v; billMember.Owes: %v; billMember.Paid: %v",
							userBillBalance, billMember.Owes, billMember.Paid)
						break
					}
				}
			}

			logus.Debugf(ctx, "userIsBillMember: %v", userIsBillMember)

			shouldBeInOutstanding := userIsBillMember && (bill.Data.Status == models4splitus.BillStatusOutstanding || bill.Data.Status == models4splitus.BillStatusDraft)
			userOutstandingBills := params.UserModule.Data.GetOutstandingBills()
			if billBrief, ok := userOutstandingBills[billID]; ok && !shouldBeInOutstanding {
				delete(userOutstandingBills, billID)
			} else if ok && shouldBeInOutstanding {
				if billUserGroupID := bill.Data.GetUserGroupID(); billBrief.GroupID != billUserGroupID {
					billBrief.GroupID = billUserGroupID
					userModuleDboChanged = true
				}
				if billBrief.UserBalance != userBillBalance {
					billBrief.UserBalance = userBillBalance
					userModuleDboChanged = true
				}
				if billBrief.Total != bill.Data.AmountTotal {
					billBrief.Total = bill.Data.AmountTotal
					userModuleDboChanged = true
				}
				if billBrief.Currency != bill.Data.Currency {
					billBrief.Currency = bill.Data.Currency
					userModuleDboChanged = true
				}
				if billBrief.Name != bill.Data.Name {
					billBrief.Name = bill.Data.Name
					userModuleDboChanged = true
				}
				userOutstandingBills[billID] = billBrief
			} else if !ok && shouldBeInOutstanding {
				userOutstandingBills[billID] = briefs4splitus.BillBrief{
					Name:         bill.Data.Name,
					MembersCount: len(bill.Data.Members),
					Total:        bill.Data.AmountTotal,
					Currency:     bill.Data.Currency,
					UserBalance:  userBillBalance,
					GroupID:      bill.Data.GetUserGroupID(),
				}
				userModuleDboChanged = true
			}
			if userModuleDboChanged {
				if err = params.UserModule.Data.SetOutstandingBills(userOutstandingBills); err != nil {
					return
				}

				if err = tx.Set(ctx, params.UserModule.Record); err != nil {
					return
				}
			} else {
				logus.Debugf(ctx, "User not changed, ContactID: %v", userID)
			}
			return
		}); err != nil {
		if dal.IsNotFound(err) {
			logus.Errorf(ctx, err.Error())
			err = nil
		}
		return
	}
	if userModuleDboChanged {
		logus.Infof(ctx, "User %v updated with info for bill %v", userID, billID)
	}
	return
}
