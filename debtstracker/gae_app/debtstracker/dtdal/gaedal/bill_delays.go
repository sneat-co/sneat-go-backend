package gaedal

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"sync"
)

const updateUsersWithBillKeyName = "update-users-with-bill"

func DelayUpdateUsersWithBill(c context.Context, billID string, userIDs []string) (err error) {
	return delayUpdateUsersWithBill.EnqueueWork(c, delaying.With(common.QUEUE_BILLS, updateUsersWithBillKeyName, 0), billID, userIDs)
}

func updateUsersWithBill(c context.Context, billID string, userIDs []string) (err error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(userIDs))
	for i := range userIDs {
		go func(i int) {
			defer wg.Done()
			if err2 := delayUpdateUserWithBill.EnqueueWork(c, delaying.With(common.QUEUE_BILLS, updateUserWithBillKeyName, 0), billID, userIDs[i]); err != nil {
				err = err2
			}
		}(i)
	}
	wg.Wait()
	return
}

const updateUserWithBillKeyName = "delayedUpdateUserWithBill"

func delayedUpdateUserWithBill(c context.Context, billID, userID string) (err error) {
	logus.Debugf(c, "delayedUpdateUserWithBill(billID=%v, userID=%v)", billID, userID)
	var (
		bill             models.Bill
		wg               sync.WaitGroup
		billErr          error
		userChanged      bool
		userIsBillMember bool
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if bill, billErr = facade2debtus.GetBillByID(c, nil, billID); err != nil {
			return
		}
	}()
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var user models.AppUser
		if user, err = dtdal.User.GetUserByStrID(c, userID); err != nil {
			return
		}
		wg.Wait()
		if billErr != nil {
			return fmt.Errorf("%w: failed to get bill", billErr)
		} else if bill.Data == nil {
			return errors.New("bill.BillEntity == nil")
		}
		var userBillBalance decimal.Decimal64p2
		if bill.Data.Status != models.BillStatusDeleted {
			for _, billMember := range bill.Data.GetBillMembers() {
				if billMember.UserID == userID {
					userBillBalance = billMember.Balance()
					userIsBillMember = true
					logus.Debugf(c, "userBillBalance: %v; billMember.Owes: %v; billMember.Paid: %v",
						userBillBalance, billMember.Owes, billMember.Paid)
					break
				}
			}
		}

		logus.Debugf(c, "userIsBillMember: %v", userIsBillMember)

		shouldBeInOutstanding := userIsBillMember && (bill.Data.Status == models.BillStatusOutstanding || bill.Data.Status == models.BillStatusDraft)
		userOutstandingBills := user.Data.GetOutstandingBills()
		for i, userOutstandingBill := range userOutstandingBills {
			if userOutstandingBill.ID == billID {
				if !shouldBeInOutstanding {
					// Remove bill info from the user
					userOutstandingBills = append(userOutstandingBills[:i], userOutstandingBills[i+1:]...)
				} else {
					if billUserGroupID := bill.Data.GetUserGroupID(); userOutstandingBill.GroupID != billUserGroupID {
						userOutstandingBill.GroupID = billUserGroupID
						userChanged = true
					}
					if userOutstandingBill.UserBalance != userBillBalance {
						userOutstandingBill.UserBalance = userBillBalance
						userChanged = true
					}
					if userOutstandingBill.Total != bill.Data.AmountTotal {
						userOutstandingBill.Total = bill.Data.AmountTotal
						userChanged = true
					}
					if userOutstandingBill.Currency != bill.Data.Currency {
						userOutstandingBill.Currency = bill.Data.Currency
						userChanged = true
					}
					if userOutstandingBill.Name != bill.Data.Name {
						userOutstandingBill.Name = bill.Data.Name
						userChanged = true
					}
					userOutstandingBills[i] = userOutstandingBill
				}
				goto doneWithChanges
			}
		}
		if shouldBeInOutstanding {
			userOutstandingBills = append(userOutstandingBills, models.BillJson{
				ID:           bill.ID,
				Name:         bill.Data.Name,
				MembersCount: bill.Data.MembersCount,
				Total:        bill.Data.AmountTotal,
				Currency:     bill.Data.Currency,
				UserBalance:  userBillBalance,
				GroupID:      bill.Data.GetUserGroupID(),
			})
			userChanged = true
		}
	doneWithChanges:
		if userChanged {
			if _, err = user.Data.SetOutstandingBills(userOutstandingBills); err != nil {
				return
			}
			if err = facade2debtus.User.SaveUser(c, tx, user); err != nil {
				return
			}
		} else {
			logus.Debugf(c, "User not changed, ID: %v", user.ID)
		}
		return
	}); err != nil {
		if dal.IsNotFound(err) {
			logus.Errorf(c, err.Error())
			err = nil
		}
		return
	}
	if userChanged {
		logus.Infof(c, "User %v updated with info for bill %v", userID, billID)
	}
	return
}
