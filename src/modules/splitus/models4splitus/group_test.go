package models4splitus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"strings"
	"testing"
)

func TestGroupEntity_ApplyBillBalanceDifference(t *testing.T) {
	splitusSpace := NewSplitusSpaceEntry("s1")

	{ // Try to apply empty difference
		if changed, err := splitusSpace.Data.ApplyBillBalanceDifference("EUR", briefs4splitus.BillBalanceDifference{}); err != nil {
			if !strings.HasSuffix(err.Error(), "not implemented yet") {
				t.Error(err)
			}
		} else if changed {
			t.Error("Should not return updates=true")
		}
	}

	splitusSpace.Data.SetGroupMembers([]briefs4splitus.SpaceSplitMember{
		{MemberBrief: briefs4splitus.MemberBrief{ID: "m1", UserID: "1", Name: "First member"}},
		//{MemberBrief: MemberBrief{ContactID: "m2", UserID: "2"}},
	})

	{ // Try to apply difference to empty balance
		if _, err := splitusSpace.Data.ApplyBillBalanceDifference("EUR", briefs4splitus.BillBalanceDifference{"m1": 100}); err == nil {
			t.Error("Shod return error")
		}
	}

	members := append(splitusSpace.Data.GetGroupMembers(), briefs4splitus.SpaceSplitMember{MemberBrief: briefs4splitus.MemberBrief{ID: "m2", UserID: "2", Name: "Second member"}})
	if updates := splitusSpace.Data.SetGroupMembers(members); len(updates) == 0 {
		t.Fatalf("Shor return updates=true")
	}

	//t.Log(splitusSpace.GetGroupMembers())

	{ // Try to add another member
		changed, err := splitusSpace.Data.ApplyBillBalanceDifference("EUR", briefs4splitus.BillBalanceDifference{
			"m1": -400,
			"m2": 400,
		})
		if err != nil {
			if !strings.HasSuffix(err.Error(), "not implemented yet") {
				t.Error(err)
			}
		} else if !changed {
			t.Error("Should return updates=true")
		}
	}
}
