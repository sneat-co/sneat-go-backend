package models4splitus

import (
	"fmt"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"time"

	"errors"
	"github.com/strongo/decimal"
)

type BillDbo struct {
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

func NewBillEntity(data BillCommon) *BillDbo {
	return &BillDbo{
		BillCommon: data,
	}
}

type BillEntry = record.DataWithID[string, *BillDbo]

func NewBillEntry(id string, billCommon *BillCommon) BillEntry {
	var data *BillDbo
	if billCommon != nil {
		data = NewBillEntity(*billCommon)
	} else {
		data = new(BillDbo)
	}
	return BillEntry{
		WithID: record.WithID[string]{ID: id},
		Data:   data,
	}
}

//var _ db.EntityHolder = (*BillEntry)(nil)

//func (bill *BillEntry) Entity() interface{} {
//	return bill.Data
//}
//
//func (BillEntry) NewEntity() interface{} {
//	return new(BillDbo)
//}

//func (bill *BillEntry) SetEntity(entity interface{}) {
//	if entity == nil {
//		bill.BillDbo = nil
//	} else {
//		bill.BillDbo = entity.(*BillDbo)
//	}
//}

//func (entity *BillDbo) Load(ps []datastore.Property) error {
//	ps = entity.BillCommon.load(ps)
//	return datastore.LoadStruct(entity, ps)
//}

func (entity *BillDbo) Validate() (err error) {
	if err = entity.validateBalance(); err != nil {
		return
	}

	entity.DebtorIDs = make([]string, 0, len(entity.Members))
	entity.SponsorIDs = make([]string, 0, len(entity.Members))

	for _, m := range entity.Members {
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

func (entity *BillDbo) validateBalance() (err error) {
	if len(entity.Members) == 0 {
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

func (entity *BillDbo) GetBalance() (billBalanceByMember briefs4splitus.BillBalanceByMember) {
	members := entity.GetBillMembers()
	billBalanceByMember = make(briefs4splitus.BillBalanceByMember, len(members))

	for i, member := range members {

		if member.Owes < 0 {
			panic(fmt.Sprintf("member[%d].Owes < 0: %v", i, member.Owes))
		} else if member.Paid < 0 {
			panic(fmt.Sprintf("member[%d].Paid < 0: %v", i, member.Paid))
		}

		if member.Owes != 0 || member.Paid != 0 {
			billBalanceByMember[member.ID] = briefs4splitus.BillMemberBalance{
				Owes: member.Owes,
				Paid: member.Paid,
			}
		}
	}
	return
}

func (entity *BillDbo) SetBillMembers(members []*briefs4splitus.BillMemberBrief) (err error) {
	if err = entity.validateMembersForDuplicatesAndBasicChecks(members); err != nil {
		return
	}

	if err = entity.updateMemberOwes(members); err != nil {
		return
	}

	entity.Members = members

	entity.setUserIDs(members)

	return
}
