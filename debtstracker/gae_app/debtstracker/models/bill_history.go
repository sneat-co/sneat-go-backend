package models

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"time"

	"errors"
	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/decimal"
)

const (
	BillsHistoryKind = "BillH"
)

type BillsHistory struct {
	record.WithID[string]
	Data *BillsHistoryEntity
}

func (BillsHistory) Kind() string {
	return BillsHistoryKind
}

//func (record BillsHistory) Entity() interface{} {
//	return record.BillsHistoryEntity
//}
//
//func (BillsHistory) NewEntity() interface{} {
//	return new(BillsHistoryEntity)
//}
//
//func (record *BillsHistory) SetEntity(entity interface{}) {
//	if entity == nil {
//		record.BillsHistoryEntity = nil
//	} else {
//		record.BillsHistoryEntity = entity.(*BillsHistoryEntity)
//	}
//}

//var _ db.EntityHolder = (*BillsHistory)(nil)

type BillsHistoryEntity struct {
	DtCreated              time.Time
	UserID                 string
	StatusOld              string              `datastore:",noindex"`
	StatusNew              string              `datastore:",noindex"`
	Action                 BillHistoryAction   `datastore:",noindex"`
	Currency               money.CurrencyCode  `datastore:",noindex"`
	TotalAmountDiff        decimal.Decimal64p2 `datastore:",noindex"`
	TotalAmountBefore      decimal.Decimal64p2 `datastore:",noindex"`
	TotalAmountAfter       decimal.Decimal64p2 `datastore:",noindex"`
	GroupIDs               []string
	BillIDs                []string
	BillsSettlementCount   int    `datastore:",noindex"`
	BillsSettlementJson    string `datastore:",noindex"`
	GroupMembersJsonBefore string `datastore:",noindex"`
	GroupMembersJsonAfter  string `datastore:",noindex"`
}

func (entity *BillsHistoryEntity) BillSettlements() (billSettlements []BillSettlementJson) {
	billSettlements = make([]BillSettlementJson, 0, entity.BillsSettlementCount)
	if err := ffjson.Unmarshal([]byte(entity.BillsSettlementJson), &billSettlements); err != nil {
		panic(err)
	}
	return
}

func (entity *BillsHistoryEntity) SetBillSettlements(groupID string, billSettlements []BillSettlementJson) { // TODO: Enable support for multiple groups
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

func (entity *BillsHistoryEntity) Validate() (err error) {
	if entity.DtCreated.IsZero() {
		entity.DtCreated = time.Now()
	}
	if entity.Action == "" {
		err = errors.New("*BillsHistoryEntity.Action is empty")
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
		err = errors.New("len(entity.GroupIDs) == 0")
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
	//	"GroupMembersJsonBefore": gaedb.IsEmptyJSON,
	//	"GroupMembersJsonAfter":  gaedb.IsEmptyJSON,
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

func NewBillHistoryBillCreated(bill Bill, groupEntity *GroupEntity) (bh BillsHistory) {
	key, err := dal.NewKeyWithOptions(BillsHistoryKind, dal.WithRandomStringID(dal.RandomLength(BillsHistoryIdLen)))
	if err != nil {
		panic(err)
	}
	data := &BillsHistoryEntity{
		Currency:         bill.Data.Currency,
		UserID:           bill.Data.CreatorUserID,
		TotalAmountAfter: bill.Data.AmountTotal,
		Action:           BillHistoryActionCreated,
		BillIDs:          []string{bill.ID},
		GroupIDs:         []string{bill.Data.UserGroupID},
	}
	bh = BillsHistory{
		WithID: record.WithID[string]{
			Key:    key,
			Record: dal.NewRecordWithData(key, data),
		},
		Data: data,
	}
	if groupEntity != nil {
		bh.Data.GroupMembersJsonAfter = groupEntity.MembersJson
	}
	return
}

func NewBillHistoryMemberAdded(userID string, bill Bill, totalAboutBefore decimal.Decimal64p2, groupMemberJsonBefore, groupMemberJsonAfter string) (bh BillsHistory) {
	bh = BillsHistory{
		Data: &BillsHistoryEntity{
			UserID:            userID,
			Currency:          bill.Data.Currency,
			TotalAmountBefore: totalAboutBefore,
			TotalAmountAfter:  bill.Data.AmountTotal,
			Action:            BillHistoryActionMemberAdded,
			BillIDs:           []string{bill.ID},
			GroupIDs:          []string{bill.Data.UserGroupID},
		},
	}
	bh.Data.GroupMembersJsonBefore = groupMemberJsonBefore
	bh.Data.GroupMembersJsonAfter = groupMemberJsonAfter
	return
}

func NewBillHistoryBillDeleted(userID string, bill Bill) (record BillsHistory) {
	panic("TODO: create key with random ID using dalgo insert options")
	//return BillsHistory{
	//	Data: &BillsHistoryEntity{
	//		StatusOld:         bill.Data.Status,
	//		StatusNew:         BillStatusDeleted,
	//		UserID:            userID,
	//		Currency:          bill.Data.Currency,
	//		TotalAmountBefore: bill.Data.AmountTotal,
	//		TotalAmountAfter:  bill.Data.AmountTotal,
	//		Action:            BillHistoryActionMemberAdded,
	//		BillIDs:           []string{bill.ID},
	//		GroupIDs:          []string{bill.Data.UserGroupID},
	//	},
	//}
}

func NewBillHistoryBillRestored(userID string, bill Bill) (record BillsHistory) {
	return BillsHistory{
		Data: &BillsHistoryEntity{
			StatusOld:         BillStatusDeleted,
			StatusNew:         bill.Data.Status,
			UserID:            userID,
			Currency:          bill.Data.Currency,
			TotalAmountBefore: bill.Data.AmountTotal,
			TotalAmountAfter:  bill.Data.AmountTotal,
			Action:            BillHistoryActionMemberAdded,
			BillIDs:           []string{bill.ID},
			GroupIDs:          []string{bill.Data.UserGroupID},
		},
	}
}
