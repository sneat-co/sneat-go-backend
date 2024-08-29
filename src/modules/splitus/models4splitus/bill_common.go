package models4splitus

import (
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/strongo/decimal"
	"github.com/strongo/strongoapp/with"
	"strconv"
)

type BillCommon struct {
	PayMode PayMode `firestore:"payMode"`
	SpaceID string  `firestore:"spaceID"`
	Status  string  `firestore:"status"`
	with.CreatedFields
	AmountTotal        decimal.Decimal64p2
	Currency           money.CurrencyCode
	UserIDs            []string
	ContactIDs         []string                          // Holds contact IDs, so we can update names in MembersJson on contact changed
	TgInlineMessageIDs []string                          `firestore:"tgInlineMessageIDs,omitempty"`
	CreatorUserID      string                            `firestore:"creatorUserID,omitempty"`
	Name               string                            `firestore:"name,omitempty"`
	SplitMode          SplitMode                         `firestore:"splitMode,omitempty"`
	Members            []*briefs4splitus.BillMemberBrief `firestore:"members,omitempty"`
	LastMemberID       int                               `firestore:"lastMemberID,omitempty"`
	Shares             int                               `firestore:"shares,omitempty"`
}

// GetUserGroupID returns user group ContactID
func (entity *BillCommon) GetUserGroupID() string {
	return entity.SpaceID
}

var (
	ErrBillAlreadyAssignedToAnotherGroup = errors.New("bill already assigned to another group ")
)

// AssignToGroup assigns bill to group
func (entity *BillCommon) AssignToGroup(groupID string) (err error) {
	if groupID == "" {
		err = errors.New("*BillCommon.AssignToGroup(): parameter groupID is required")
		return
	}
	if entity.SpaceID == "" {
		entity.SpaceID = groupID
	} else if entity.SpaceID != groupID {
		err = fmt.Errorf("%w: %s", ErrBillAlreadyAssignedToAnotherGroup, entity.SpaceID)
	}
	return
}

func (entity *BillCommon) AddOrGetMember(
	groupMemberID, userID, contactID, name string,
) (
	isNew, changed bool, index int, member *briefs4splitus.BillMemberBrief, billMembers []*briefs4splitus.BillMemberBrief,
) {
	members := entity.GetMembers()
	var m briefs4splitus.MemberBrief
	if index, m, isNew, changed = AddOrGetMember(members, groupMemberID, userID, contactID, name); isNew {
		member = &briefs4splitus.BillMemberBrief{
			MemberBrief: m,
		}
		entity.Members = append(entity.Members, member)
		if index != len(billMembers)-1 {
			panic("index != len(billMembers) - 1")
		}
		changed = true
	} else /* existing member */ if member = billMembers[index]; member.ID != m.ID {
		panic("member.ContactID != m.ContactID")
	}
	if member.ID == "" {
		panic("member.ContactID is empty string")
	}
	return
}

func (entity *BillCommon) IsOkToSplit() bool {
	if len(entity.Members) <= 1 {
		return false
	}

	var paidByMembers decimal.Decimal64p2
	for _, m := range entity.GetBillMembers() {
		paidByMembers += m.Paid
		// owedByMembers += m.Owes
	}
	return paidByMembers == entity.AmountTotal
}

func (entity *BillCommon) TotalAmount() money.Amount {
	return money.NewAmount(money.CurrencyCode(entity.Currency), entity.AmountTotal)
}

func (entity *BillCommon) GetBillMembers() (members []*briefs4splitus.BillMemberBrief) {
	members = make([]*briefs4splitus.BillMemberBrief, len(entity.Members))
	copy(members, entity.Members)
	return
}

func (entity *BillCommon) GetMembers() (members []briefs4splitus.MemberBrief) {
	billMembers := entity.GetBillMembers()
	members = make([]briefs4splitus.MemberBrief, len(billMembers))
	for i, bm := range billMembers {
		members[i] = bm.MemberBrief
	}
	return
}

func (entity *BillCommon) validateMembersForDuplicatesAndBasicChecks(members []*briefs4splitus.BillMemberBrief) error {
	isEquallySplit := true
	// maxShares := 0

	uniqueUserIDs := make(map[string]int, len(members))
	for i, member := range members {
		if member.ID == "" {
			entity.LastMemberID++
			member.ID = strconv.Itoa(entity.LastMemberID)
		}
		if isEquallySplit {
			// if member.Shares > maxShares {
			// 	maxShares = member.Shares
			// }
			if member.Adjustment != 0 || (i > 0 && member.Shares != members[i-1].Shares) {
				isEquallySplit = false
			}
		}
		if member.UserID != "" {
			for _, uniqueUserID := range uniqueUserIDs {
				if i0, ok := uniqueUserIDs[member.UserID]; ok {
					return fmt.Errorf("duplicate members with same UserID=%d: members[%d].UserID == members[%d].UserID", uniqueUserID, i, i0)
				}
			}
			uniqueUserIDs[member.UserID] = i
		}
		if member.Name == "" {
			return fmt.Errorf("no name for the members[%d]", i)
		}
		if member.Owes > entity.AmountTotal {
			return fmt.Errorf("members[%d].Owes > entity.AmountTotal", i)
		}
		if member.Adjustment > entity.AmountTotal || (member.Adjustment < 0 && -1*member.Adjustment > entity.AmountTotal) {
			return fmt.Errorf("members[%d].Adjustment is too big", i)
		}
	}
	return nil
}

func (entity *BillCommon) setUserIDs(members []*briefs4splitus.BillMemberBrief) {
	entity.UserIDs = make([]string, 0, len(members))
Members:
	for _, m := range members {
		if m.UserID != "" {
			for _, userID := range entity.UserIDs {
				if userID == m.UserID {
					continue Members
				}
			}
			entity.UserIDs = append(entity.UserIDs, m.UserID)
		}
	}
}

//func (entity *BillCommon) load(ps []datastore.Property) []datastore.Property {
//	for i, p := range ps {
//		if p.Name == "GetUserGroupID" {
//			entity.SpaceRef = p.Value.(string)
//			return append(ps[:i], ps[i+1:]...)
//		}
//	}
//	return ps
//}

func (entity *BillCommon) Validate() (err error) {
	if entity.CreatorUserID == "" {
		panic("entity.CreatorUserID is empty string")
	}
	if entity.SplitMode == "" {
		return errors.New("entity.SplitMode is empty string")
	}
	if entity.Status == "" {
		return errors.New("entity.Status is empty string")
	}
	if entity.CreatedAt.IsZero() {
		return errors.New("entity.DtCreated is zero")
	}
	//if filtered, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
	//	"MembersCount": gaedb.IsZeroInt,
	//	"MembersJson":  gaedb.IsEmptyJSON,
	//	"PayMode":      gaedb.IsEmptyString,
	//	"ContactName":  gaedb.IsEmptyString,
	//	"SplitMode":    gaedb.IsEmptyString,
	//	"Shares":       gaedb.IsZeroInt,
	//}); err != nil {
	//	return
	//}
	//if entity.SpaceRef != "" {
	//	filtered = append(filtered, datastore.Property{Name: "GetUserGroupID", Value: entity.SpaceRef, NoIndex: false})
	//}
	return
}
