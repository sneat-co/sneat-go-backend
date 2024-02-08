package models

import (
	"fmt"
	"github.com/dal-go/dalgo/record"
	"time"

	"errors"
	"github.com/strongo/decimal"
)

const (
	BillKind = "Bill"
)

const (
	BillStatusDraft       = STATUS_DRAFT
	BillStatusDeleted     = STATUS_DELETED
	BillStatusOutstanding = "outstanding"
	BillStatusSettled     = "settled"
)

var (
	BillStatuses = [3]string{
		BillStatusDraft,
		BillStatusOutstanding,
		BillStatusSettled,
	}
	BillSplitModes = [4]SplitMode{
		// SplitModeAdjustment,
		SplitModeEqually,
		SplitModeExactAmount,
		SplitModePercentage,
		SplitModeShare,
	}
)

func IsValidBillSplit(split SplitMode) bool {
	for _, v := range BillSplitModes {
		if split == v {
			return true
		}
	}
	return false
}

func IsValidBillStatus(status string) bool {
	for _, v := range BillStatuses {
		if status == v {
			return true
		}
	}
	return false
}

type BillEntity struct {
	BillCommon
	DtDueToPay       time.Time `datastore:",noindex"` // TODO: Document diff between DtDueToPay & DtDueToCollect
	DtDueToCollect   time.Time `datastore:",noindex"`
	LocaleByMessage  []string  `datastore:",noindex"`
	TgChatMessageIDs []string  `datastore:",noindex"`
	DebtorIDs        []string
	SponsorIDs       []string
	SettlementIDs    []string
	//BalanceJson      string    `datastore:",noindex"`
	//BalanceVersion   int       `datastore:",noindex"`
	//balanceVersion   int       `datastore:"-"`
}

func NewBillEntity(data BillCommon) *BillEntity {
	return &BillEntity{
		BillCommon: data,
	}
}

type Bill = record.DataWithID[string, *BillEntity]

//type Bill struct {
//	record.WithID[string]
//	Data *BillEntity
//}

func NewBill(id string, billCommon *BillCommon) Bill {
	var data *BillEntity
	if billCommon != nil {
		data = NewBillEntity(*billCommon)
	} else {
		data = new(BillEntity)
	}
	return Bill{
		WithID: record.WithID[string]{ID: id},
		Data:   data,
	}
}

//var _ db.EntityHolder = (*Bill)(nil)

//func (bill *Bill) Entity() interface{} {
//	return bill.Data
//}
//
//func (Bill) NewEntity() interface{} {
//	return new(BillEntity)
//}

//func (bill *Bill) SetEntity(entity interface{}) {
//	if entity == nil {
//		bill.BillEntity = nil
//	} else {
//		bill.BillEntity = entity.(*BillEntity)
//	}
//}

//func (entity *BillEntity) Load(ps []datastore.Property) error {
//	ps = entity.BillCommon.load(ps)
//	return datastore.LoadStruct(entity, ps)
//}

func (entity *BillEntity) Validate() (err error) {
	if err = entity.validateBalance(); err != nil {
		return
	}

	entity.DebtorIDs = make([]string, 0, len(entity.members))
	entity.SponsorIDs = make([]string, 0, len(entity.members))

	for _, m := range entity.members {
		if m.Paid > m.Owes {
			entity.SponsorIDs = append(entity.SponsorIDs, m.ID)
		} else if m.Paid < m.Owes {
			entity.DebtorIDs = append(entity.DebtorIDs, m.ID)
		}
	}

	//if properties, err = datastore.SaveStruct(entity); err != nil {
	//	return
	//}
	if err = entity.BillCommon.Validate(); err != nil {
		return
	}
	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
	//	"DtDueToPay":     gaedb.IsZeroTime,
	//	"DtDueToCollect": gaedb.IsZeroTime,
	//}); err != nil {
	//	return
	//}
	return
}

var (
	ErrNegativeAmount                   = errors.New("negative amount")
	ErrTotalOwedIsNotMatchingBillAmount = errors.New("total owed is not matching bill amount")
	ErrTotalPaidIsGreaterThenBillAmount = errors.New("total paid is greater then bill amount")
	ErrBillTotalBalanceIsNotZero        = errors.New("total bill balance is not zero")
	ErrBillOwesDiffTotalIsNotZero       = errors.New("total bill difference of owes is not zero")
	ErrNonGroupMember                   = errors.New("non group member")
	ErrGroupTotalBalanceHasNonZeroValue = errors.New("group total balance has non zero value")
)

func (entity *BillEntity) validateBalance() (err error) {
	if entity.MembersCount == 0 {
		return
	}
	var (
		totalBalance decimal.Decimal64p2
		totalPaid    decimal.Decimal64p2
		totalOwed    decimal.Decimal64p2
	)

	members := entity.GetBillMembers()

	for i, member := range members {
		if member.Owes < 0 {
			err = fmt.Errorf("%w: members[%d] owes=%v", ErrNegativeAmount, i, member.Owes)
			return
		}
		if member.Paid < 0 {
			err = fmt.Errorf("%w: members[%d] paid=%v", ErrNegativeAmount, i, member.Paid)
			return
		}
		totalBalance += member.Paid - member.Owes
		totalPaid += member.Paid
		totalOwed += member.Owes
	}

	if totalOwed != entity.AmountTotal {
		err = fmt.Errorf("%w: totalOwed: %v, AmountTotal: %v", ErrTotalOwedIsNotMatchingBillAmount, totalOwed, entity.AmountTotal)
	}

	if totalPaid > entity.AmountTotal {
		err = fmt.Errorf("%w: totalPaid: %v, AmountTotal: %v", ErrTotalPaidIsGreaterThenBillAmount, totalPaid, entity.AmountTotal)
	}

	if totalBalance != 0 {
		err = fmt.Errorf("%w: totalBalance=%v, members: %+v", ErrBillTotalBalanceIsNotZero, totalBalance, members)
	}

	return
}

func (entity *BillEntity) GetBalance() (billBalanceByMember BillBalanceByMember) {
	members := entity.GetBillMembers()
	billBalanceByMember = make(BillBalanceByMember, len(members))

	for i, member := range members {

		if member.Owes < 0 {
			panic(fmt.Sprintf("member[%d].Owes < 0: %v", i, member.Owes))
		} else if member.Paid < 0 {
			panic(fmt.Sprintf("member[%d].Paid < 0: %v", i, member.Paid))
		}

		if member.Owes != 0 || member.Paid != 0 {
			billBalanceByMember[member.ID] = BillMemberBalance{
				Owes: member.Owes,
				Paid: member.Paid,
			}
		}
	}
	return
}

func (entity *BillEntity) SetBillMembers(members []BillMemberJson) (err error) {
	if err = entity.validateMembersForDuplicatesAndBasicChecks(members); err != nil {
		return
	}

	if err = entity.updateMemberOwes(members); err != nil {
		return
	}

	if err = entity.marshalMembersToJsonAndSetMembersCount(members); err != nil {
		return
	}

	entity.setUserIDs(members)

	return
}
