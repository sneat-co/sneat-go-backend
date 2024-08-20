package models4splitus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"time"

	"errors"
	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/decimal"
)

type BillsHistory = record.DataWithID[string, *BillsHistoryDbo]

//func (record BillsHistory) Entity() interface{} {
//	return record.BillsHistoryDbo
//}
//
//func (BillsHistory) NewEntity() interface{} {
//	return new(BillsHistoryDbo)
//}
//
//func (record *BillsHistory) SetEntity(entity interface{}) {
//	if entity == nil {
//		record.BillsHistoryDbo = nil
//	} else {
//		record.BillsHistoryDbo = entity.(*BillsHistoryDbo)
//	}
//}

//var _ db.EntityHolder = (*BillsHistory)(nil)

type BillsHistoryDbo struct {
	DtCreated            time.Time
	UserID               string
	StatusOld            string              `datastore:",noindex"`
	StatusNew            string              `datastore:",noindex"`
	Action               BillHistoryAction   `datastore:",noindex"`
	Currency             money.CurrencyCode  `datastore:",noindex"`
	TotalAmountDiff      decimal.Decimal64p2 `datastore:",noindex"`
	TotalAmountBefore    decimal.Decimal64p2 `datastore:",noindex"`
	TotalAmountAfter     decimal.Decimal64p2 `datastore:",noindex"`
	GroupIDs             []string
	BillIDs              []string
	BillsSettlementCount int                               `firestore:"billsSettlementCount,omitempty"`
	BillsSettlementJson  string                            `firestore:"billsSettlementJson,omitempty"`
	SplitMembersBefore   []briefs4splitus.SpaceSplitMember `firestore:"splitMembersBefore,omitempty"`
	SplitMembersAfter    []briefs4splitus.SpaceSplitMember `firestore:"members,omitempty"`
}

func (entity *BillsHistoryDbo) BillSettlements() (billSettlements []briefs4splitus.BillSettlementJson) {
	billSettlements = make([]briefs4splitus.BillSettlementJson, 0, entity.BillsSettlementCount)
	if err := ffjson.Unmarshal([]byte(entity.BillsSettlementJson), &billSettlements); err != nil {
		panic(err)
	}
	return
}

func (entity *BillsHistoryDbo) SetBillSettlements(groupID string, billSettlements []briefs4splitus.BillSettlementJson) { // TODO: Enable support for multiple groups
	if data, err := ffjson.Marshal(&billSettlements); err != nil {
		panic(err)
	} else {
		entity.BillsSettlementJson = string(data)
		entity.BillsSettlementCount = len(billSettlements)
		entity.BillIDs = make([]string, len(billSettlements))
		entity.GroupIDs = make([]string, 0, 1)
		for i, b := range billSettlements {
			entity.BillIDs[i] = b.BillID
			if b.GroupID != "" {
				for _, groupID := range entity.GroupIDs {
					if groupID == b.GroupID {
						goto groupFound
					}
				}
				entity.GroupIDs = append(entity.GroupIDs, b.GroupID)
			groupFound:
			}
		}
	}
}

func (entity *BillsHistoryDbo) Validate() (err error) {
	if entity.DtCreated.IsZero() {
		entity.DtCreated = time.Now()
	}
	if entity.Action == "" {
		err = errors.New("*BillsHistoryDbo.Action is empty")
		return
	}
	if entity.Action == BillHistoryActionSettled && entity.BillsSettlementJson == "" {
		err = errors.New("BillsSettlementJson is empty")
		return
	}
	// if entity.Currency == "" {
	// 	err = errors.New("Currency is empty")
	// 	return
	// }
	if len(entity.GroupIDs) == 0 {
		err = errors.New("len(entity.SpaceIDs) == 0")
		return
	}
	if entity.BillsSettlementJson == "" {
		if entity.BillsSettlementCount != 0 {
			err = errors.New("BillsSettlementJson is not empty && BillsSettlementCount !=  0")
			return
		}
	} else {
		bills := entity.BillSettlements()
		if entity.BillsSettlementCount != len(bills) {
			err = errors.New("BillsCount != len(bills)")
			return
		}
		var total decimal.Decimal64p2
		for i, b := range bills {
			total += b.Amount
			if entity.BillIDs[i] != b.BillID {
				err = fmt.Errorf("entity.BillIDs[%d]:%v != b.BillID:%v", i, entity.BillIDs[i], b.BillID)
			}
		}
		if entity.TotalAmountAfter != total {
			err = fmt.Errorf("entity.TotalAmount:%v != total:%v", entity.TotalAmountAfter, total)
			return
		}
	}
	//if properties, err = datastore.SaveStruct(entity); err != nil {
	//	return
	//}
	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
	//	"Currency":               gaedb.IsEmptyString,
	//	"TotalAmountDiff":        gaedb.IsZeroInt,
	//	"TotalAmountBefore":      gaedb.IsZeroInt,
	//	"TotalAmountAfter":       gaedb.IsZeroInt,
	//	"BillsSettlementCount":   gaedb.IsZeroInt,
	//	"BillsSettlementJson":    gaedb.IsEmptyJSON,
	//	"SplitMembersBefore": gaedb.IsEmptyJSON,
	//	"SplitMembersAfter":  gaedb.IsEmptyJSON,
	//}); err != nil {
	//	return
	//}
	return
}

type BillHistoryAction string

const (
	BillHistoryActionCreated     BillHistoryAction = "created"
	BillHistoryActionMemberAdded BillHistoryAction = "member-added"
	BillHistoryActionSettled     BillHistoryAction = "settled"
	//BillHistoryActionDeleted     BillHistoryAction = "deleted"
	//BillHistoryActionRestored    BillHistoryAction = "restored"
)

func NewBillHistoryBillCreated(bill BillEntry, splitusSpaceDbo *SplitusSpaceDbo) (bh BillsHistory) {
	key, err := dal.NewKeyWithOptions(BillsHistoryCollection, dal.WithRandomStringID(dal.RandomLength(const4debtus.BillsHistoryIdLen)))
	if err != nil {
		panic(err)
	}
	data := &BillsHistoryDbo{
		Currency:         bill.Data.Currency,
		UserID:           bill.Data.CreatorUserID,
		TotalAmountAfter: bill.Data.AmountTotal,
		Action:           BillHistoryActionCreated,
		BillIDs:          []string{bill.ID},
		GroupIDs:         []string{bill.Data.SpaceID},
	}
	bh = BillsHistory{
		WithID: record.WithID[string]{
			Key:    key,
			Record: dal.NewRecordWithData(key, data),
		},
		Data: data,
	}
	if splitusSpaceDbo != nil {
		bh.Data.SplitMembersAfter = splitusSpaceDbo.Members
	}
	return
}

func NewBillHistoryMemberAdded(
	userID string,
	bill BillEntry,
	totalAboutBefore decimal.Decimal64p2,
	splitMembersBefore, splitMembersAfter []briefs4splitus.SpaceSplitMember,
) (bh BillsHistory) {
	bh = BillsHistory{
		Data: &BillsHistoryDbo{
			UserID:            userID,
			Currency:          bill.Data.Currency,
			TotalAmountBefore: totalAboutBefore,
			TotalAmountAfter:  bill.Data.AmountTotal,
			Action:            BillHistoryActionMemberAdded,
			BillIDs:           []string{bill.ID},
			GroupIDs:          []string{bill.Data.SpaceID},
		},
	}
	bh.Data.SplitMembersBefore = splitMembersBefore
	bh.Data.SplitMembersAfter = splitMembersAfter
	return
}

func NewBillHistoryBillDeleted(userID string, bill BillEntry) (record BillsHistory) {
	panic("TODO: create key with random ContactID using dalgo insert options")
	//return BillsHistory{
	//	Data: &BillsHistoryDbo{
	//		StatusOld:         bill.Data.Status,
	//		StatusNew:         BillStatusDeleted,
	//		UserID:            userID,
	//		Currency:          bill.Data.Currency,
	//		TotalAmountBefore: bill.Data.AmountTotal,
	//		TotalAmountAfter:  bill.Data.AmountTotal,
	//		Action:            BillHistoryActionMemberAdded,
	//		BillIDs:           []string{bill.ContactID},
	//		SpaceIDs:          []string{bill.Data.SpaceID},
	//	},
	//}
}

func NewBillHistoryBillRestored(userID string, bill BillEntry) (record BillsHistory) {
	return BillsHistory{
		Data: &BillsHistoryDbo{
			StatusOld:         BillStatusDeleted,
			StatusNew:         bill.Data.Status,
			UserID:            userID,
			Currency:          bill.Data.Currency,
			TotalAmountBefore: bill.Data.AmountTotal,
			TotalAmountAfter:  bill.Data.AmountTotal,
			Action:            BillHistoryActionMemberAdded,
			BillIDs:           []string{bill.ID},
			GroupIDs:          []string{bill.Data.SpaceID},
		},
	}
}
