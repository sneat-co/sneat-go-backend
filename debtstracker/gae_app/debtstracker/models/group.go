package models

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/random"
	"github.com/strongo/validation"
	"strings"

	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/decimal"
)

const GroupKind = "Group"

type Group struct {
	record.WithID[string]
	Data *GroupEntity
}

func NewGroup(id string, data *GroupEntity) Group {
	key := NewGroupKey(id)
	if data == nil {
		data = new(GroupEntity)
	}
	return Group{
		WithID: record.WithID[string]{
			ID:     id,
			Key:    key,
			Record: dal.NewRecordWithData(key, data),
		},
		Data: data,
	}
}

func NewGroupKey(id string) *dal.Key {
	if id == "" {
		key, err := dal.NewKeyWithOptions(GroupKind, dal.WithRandomStringID(dal.RandomLength(GroupIdLen)))
		if err != nil {
			panic(err.Error())
		}
		return key
	}
	return dal.NewKeyWithID(GroupKind, id)
}

type GroupEntity struct {
	CreatorUserID string
	//IsUser2User         bool   `datastore:",noindex"`
	Name                string             `datastore:",noindex"`
	Note                string             `datastore:",noindex,omitempty"`
	DefaultCurrency     money.CurrencyCode `datastore:",noindex,omitempty"`
	members             []GroupMemberJson
	MembersCount        int    `datastore:",noindex,omitempty"`
	MembersJson         string `datastore:",noindex,omitempty"`
	telegramGroups      []GroupTgChatJson
	TelegramGroupsCount int    `datastore:"TgGroupsCount,noindex,omitempty"`
	TelegramGroupsJson  string `datastore:"TgGroupsJson,noindex,omitempty"`
	billsHolder
}

func (entity *GroupEntity) ApplyBillBalanceDifference(currency money.CurrencyCode, diff BillBalanceDifference) (changed bool, err error) {
	if currency == "" {
		panic("currency parameter is required")
	}
	if strings.TrimSpace(string(currency)) != string(currency) {
		panic("currency parameter has leading ot closing spaces: " + currency)
	}

	groupMembers := entity.GetGroupMembers()

	var diffTotal, balanceTotal decimal.Decimal64p2
	diffCopy := make(BillBalanceDifference, len(diff))

	for i := range groupMembers {
		groupMemberID := groupMembers[i].ID

		if memberDifference, ok := diff[groupMemberID]; ok {
			delete(diff, groupMemberID)
			diffCopy[groupMemberID] = memberDifference
			if memberDifference == 0 {
				panic("memberDifference.Paid == 0 && memberDifference.Owes == 0, memberID: " + groupMemberID)
			}
			diffTotal += memberDifference
			if diffAmount := memberDifference; diffAmount != 0 {
				if groupMembers[i].Balance == nil || len(groupMembers) == 0 {
					groupMembers[i].Balance = money.Balance{currency: diffAmount}
					balanceTotal += diffAmount
				} else {
					groupMembers[i].Balance[currency] += diffAmount
					if len(groupMembers[i].Balance) == 0 {
						groupMembers[i].Balance = nil
					} else {
						balanceTotal += groupMembers[i].Balance[currency]
					}
				}
			}
		}
	}

	if len(diff) > 0 {
		err = fmt.Errorf("%w: %v", ErrNonGroupMember, diff)
		return
	}

	if diffTotal != 0 {
		err = fmt.Errorf("%w: diffTotal=%v, diff=%v", ErrBillOwesDiffTotalIsNotZero, diffTotal, diffCopy)
		return
	}

	if balanceTotal != 0 {
		err = fmt.Errorf("%wbalanceTotal=%v, diff=%v", ErrGroupTotalBalanceHasNonZeroValue, balanceTotal, diffCopy)
		return
	}
	return entity.SetGroupMembers(groupMembers), err
}

func (entity *GroupEntity) GetTelegramGroups() (tgGroups []GroupTgChatJson, err error) {
	if entity.telegramGroups != nil {
		return entity.telegramGroups, nil
	}
	if entity.TelegramGroupsJson != "" {
		if err = ffjson.Unmarshal([]byte(entity.TelegramGroupsJson), &tgGroups); err != nil {
			return
		} else if len(tgGroups) != entity.TelegramGroupsCount {
			err = fmt.Errorf("%w: len([]GroupTgChatJson) != entity.TelegramGroupsCount", ErrJsonCountMismatch)
			return
		}
		entity.telegramGroups = tgGroups
	}
	return
}

func (entity *GroupEntity) SetTelegramGroups(tgGroups []GroupTgChatJson) (changed bool) {
	if data, err := ffjson.Marshal(tgGroups); err != nil {
		panic(err.Error())
	} else {
		if s := string(data); s != entity.TelegramGroupsJson {
			entity.TelegramGroupsJson = s
			changed = true
		}
		if l := len(tgGroups); l != entity.TelegramGroupsCount {
			entity.TelegramGroupsCount = l
			changed = true
		}
	}
	return
}

func (entity *GroupEntity) AddOrGetMember(userID, contactID, name string) (isNew, changed bool, index int, member GroupMemberJson, groupMembers []GroupMemberJson) {
	if userID == "" {
		panic(userID == "")
	}
	if name == "" {
		panic("name is empty string")
	}
	members := entity.GetMembers()
	groupMembers = entity.GetGroupMembers()
	var m MemberJson
	if index, m, isNew, changed = addOrGetMember(members, "", userID, contactID, name); isNew {
		member = GroupMemberJson{
			MemberJson: m,
		}
		groupMembers = append(groupMembers, member)
		if index != len(groupMembers)-1 {
			panic("index != len(groupMembers) - 1")
		}
		changed = true
	} else /* existing member */ if member = groupMembers[index]; member.ID != m.ID {
		panic("member.ID != m.ID")
	}
	if member.ID == "" {
		panic("member.ID is empty string")
	}
	return
}

func addOrGetMember(members []MemberJson, memberID, userID, contactID, name string) (index int, member MemberJson, isNew, changed bool) {
	if userID != "" || contactID != "" {
		for i, m := range members {
			if m.ID == memberID || m.UserID == userID {
				member = m
				if contactID != "" {
					for _, cID := range m.ContactIDs {
						if cID == contactID {
							goto contactFound
						}
					}
					m.ContactIDs = append(m.ContactIDs, contactID)
					changed = true
				contactFound:
				}
				member = m
				index = i
				return
			} else if contactID != "" {
				for _, cID := range m.ContactIDs {
					if cID == contactID {
						member = m
						index = i
						return
					}
				}
			}
		}
	}
	member = MemberJson{
		ID:     memberID,
		Name:   name,
		UserID: userID,
	}
	if member.ID == "" {
	randomID:
		for j := 0; j < 100; j++ {
			member.ID = random.ID(MemberIdLen)
			for _, m := range members {
				if m.ID == member.ID {
					continue randomID
				}
			}
			break
		}
		if member.ID == "" {
			panic("Failed to generate random member ID")
		}
	}
	return len(members), member, true, true
}

func (entity *GroupEntity) GetGroupMembers() []GroupMemberJson {
	members := make([]GroupMemberJson, entity.MembersCount)
	if entity.members != nil && len(entity.members) == entity.MembersCount {
		copy(members, entity.members)
		return members
	}
	if entity.MembersJson != "" {
		if err := ffjson.Unmarshal(([]byte)(entity.MembersJson), &members); err != nil {
			panic(err.Error())
		}
	}
	if len(members) != entity.MembersCount {
		panic("len(members) != entity.MembersCount")
	}
	entity.members = make([]GroupMemberJson, entity.MembersCount)
	copy(entity.members, members)
	return members
}

func (entity *GroupEntity) GetGroupMemberByID(id string) (GroupMemberJson, error) {
	if id == "" {
		return GroupMemberJson{}, fmt.Errorf("%w: empty id", dal.ErrRecordNotFound)
	}
	for _, m := range entity.GetGroupMembers() {
		if m.ID == id {
			return m, nil
		}
	}
	return GroupMemberJson{}, fmt.Errorf("%w: unknown id="+id, dal.ErrRecordNotFound)
}

func (entity *GroupEntity) GetGroupMemberByUserID(userID string) (GroupMemberJson, error) {
	if userID == "" {
		return GroupMemberJson{}, fmt.Errorf("%w: empty id", dal.ErrRecordNotFound)
	}
	for _, m := range entity.GetGroupMembers() {
		if m.UserID == userID {
			return m, nil
		}
	}
	return GroupMemberJson{}, fmt.Errorf("%w: unknown userID=%s", dal.ErrRecordNotFound, userID)
}

func (entity *GroupEntity) GetMembers() (members []MemberJson) {
	groupMembers := entity.GetGroupMembers()
	members = make([]MemberJson, len(groupMembers))
	for i, gm := range groupMembers {
		members[i] = gm.MemberJson
	}
	return
}

func (entity *GroupEntity) GetSplitMode() SplitMode {
	if entity.MembersCount == 0 {
		return SplitModeEqually
	}
	var min, max int
	for _, m := range entity.GetGroupMembers() {
		if m.Shares < min || min == 0 {
			min = m.Shares
		}
		if m.Shares > max {
			max = m.Shares
		}
	}
	if min == max {
		return SplitModeEqually
	}
	return SplitModeShare
}

func (entity *GroupEntity) TotalShares() (n int) {
	for _, m := range entity.GetGroupMembers() {
		n += m.Shares
	}
	return
}

func (entity *GroupEntity) UserIsMember(userID string) bool {
	for _, m := range entity.GetGroupMembers() {
		if m.UserID == userID {
			return true
		}
	}
	return false
}

func (entity *GroupEntity) SetGroupMembers(members []GroupMemberJson) (changed bool) {
	if len(members) == 0 {
		if changed = entity.MembersJson != ""; changed {
			entity.members = make([]GroupMemberJson, 0)
			entity.MembersJson = ""
			entity.MembersCount = 0
		}
		return
	}
	if err := entity.validateMembers(members, len(members)); err != nil {
		panic(err)
	}
	if data, err := ffjson.Marshal(members); err != nil {
		ffjson.Pool(data)
		panic(err)
	} else if membersJson := (string)(data); membersJson != entity.MembersJson {
		ffjson.Pool(data)
		if membersJson == "[]" {
			if entity.MembersJson == "" {
				return
			}
			membersJson = ""
		}
		changed = true
		entity.MembersJson = membersJson
		entity.members = make([]GroupMemberJson, len(members))
		copy(entity.members, members)
		entity.MembersCount = len(members)
	} else {
		ffjson.Pool(data)
	}
	return
}

func (entity *GroupEntity) validateMembers(members []GroupMemberJson, membersCount int) error {
	if membersCount != len(members) {
		return fmt.Errorf("entity.MembersCount != len(members), %d != %d", entity.MembersCount, len(members))
	}

	type Empty struct {
	}

	EMPTY := Empty{}

	totalBalance := make(money.Balance)

	userIDs := make(map[string]Empty, entity.MembersCount)
	contactIDs := make(map[string]Empty, entity.MembersCount)

	memberIDs := make(map[string]Empty, entity.MembersCount)

	for i, m := range members {
		if m.ID == "" {
			return fmt.Errorf("members[%d].ID is empty string", i)
		}
		if strings.TrimSpace(m.Name) == "" {
			return fmt.Errorf("members[%d].Name is empty string", i)
		}
		if _, ok := memberIDs[m.ID]; ok {
			return fmt.Errorf("members[%d]: Duplicate ID: %v", i, m.ID)
		}
		memberIDs[m.ID] = EMPTY
		if m.UserID == "" && len(m.ContactIDs) == 0 {
			return fmt.Errorf("members[%d]: m.UserID == 0 && len(m.ContactIDs) == 0", i)
		}
		if m.UserID != "" {
			if _, ok := userIDs[m.UserID]; ok {
				return fmt.Errorf("members[%d]: Duplicate UserID: %v", i, m.UserID)
			}
			userIDs[m.UserID] = EMPTY
		} else if len(m.ContactIDs) > 0 {
			for _, contactID := range m.ContactIDs {
				if _, ok := contactIDs[contactID]; ok {
					return fmt.Errorf("members[%d]: Duplicate ContactID: %v", i, contactID)
				}
				contactIDs[contactID] = EMPTY
			}
		}
		for currency, amount := range m.Balance {
			totalBalance[currency] += amount
		}
	}

	// Validate total balance is 0
	for currency, amount := range totalBalance {
		if amount != 0 {
			return fmt.Errorf("%w: %v=%v", ErrGroupTotalBalanceHasNonZeroValue, currency, amount)
		}
	}
	return nil
}

//func (entity *GroupEntity) Load(ps []datastore.Property) (err error) {
//	if ps, err = gaedb.CleanProperties(ps, map[string]gaedb.IsOkToRemove{
//		"Status": gaedb.IsObsolete,
//	}); err != nil {
//		return
//	}
//	if err = datastore.LoadStruct(entity, ps); err != nil {
//		return err
//	}
//	return nil
//}

//var groupPropertiesToClean = map[string]gaedb.IsOkToRemove{
//	"DefaultCurrency":       gaedb.IsEmptyString,
//	"MembersCount":          gaedb.IsZeroInt,
//	"MemberLastID":          gaedb.IsZeroInt,
//	"MembersJson":           gaedb.IsEmptyJSON,
//	"Note":                  gaedb.IsEmptyString,
//	"OutstandingBillsJson":  gaedb.IsEmptyJSON,
//	"OutstandingBillsCount": gaedb.IsZeroInt,
//	"TgGroupsCount":         gaedb.IsZeroInt,
//	"TgGroupsJson":          gaedb.IsEmptyJSON,
//}

func (entity *GroupEntity) Validate() error {
	if entity.CreatorUserID == "" {
		return validation.NewErrRecordIsMissingRequiredField("CreatorUserID")
	}
	if strings.TrimSpace(entity.Name) == "" {
		return validation.NewErrRecordIsMissingRequiredField("Name")
	}
	if err := entity.validateMembers(entity.GetGroupMembers(), entity.MembersCount); err != nil {
		return err
	}
	//ps, err := datastore.SaveStruct(entity)
	//if ps, err = gaedb.CleanProperties(ps, groupPropertiesToClean); err != nil {
	//	return ps, fmt.Errorf("%w: failed to clean properties for *GroupEntity", err)
	//}
	//checkHasProperties(GroupKind, ps)
	return nil
}

func (entity *GroupEntity) AddBill(bill Bill) (changed bool, err error) {
	outstandingBills := entity.GetOutstandingBills()

	for i, b := range outstandingBills {
		if b.ID == bill.ID {
			if b.Name != bill.Data.Name {
				outstandingBills[i].Name = bill.Data.Name
				changed = true
			}
			if b.MembersCount != bill.Data.MembersCount {
				outstandingBills[i].MembersCount = bill.Data.MembersCount
				changed = true
			}
			if b.Total != bill.Data.AmountTotal {
				outstandingBills[i].Total = bill.Data.AmountTotal
				changed = true
			}
			goto addedOrUpdatedOrNotChanged
		}
	}
	outstandingBills = append(outstandingBills, BillJson{
		ID:           bill.ID,
		Name:         bill.Data.Name,
		MembersCount: bill.Data.MembersCount,
		Total:        bill.Data.AmountTotal,
		Currency:     bill.Data.Currency,
	})
addedOrUpdatedOrNotChanged:
	if changed {
		if _, err = entity.SetOutstandingBills(outstandingBills); err != nil {
			return
		}
		groupMembers := entity.GetGroupMembers()
		billMembers := bill.Data.GetBillMembers()
		for j, groupMember := range groupMembers {
			for _, billMember := range billMembers {
				if billMember.ID == groupMember.ID {
					groupMember.Balance[bill.Data.Currency] += billMember.Balance()
					groupMembers[j] = groupMember
					break
				}
			}
		}
		entity.SetGroupMembers(groupMembers)
	}
	return
}
