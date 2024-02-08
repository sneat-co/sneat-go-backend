package models

import "testing"

func TestGroupEntity_ApplyBillBalanceDifference(t *testing.T) {
	groupEntity := GroupEntity{}

	{ // Try to apply empty difference
		if changed, err := groupEntity.ApplyBillBalanceDifference("EUR", BillBalanceDifference{}); err != nil {
			t.Error(err)
		} else if changed {
			t.Error("Should not return changed=true")
		}
	}

	groupEntity.SetGroupMembers([]GroupMemberJson{
		{MemberJson: MemberJson{ID: "m1", UserID: "1", Name: "First member"}},
		//{MemberJson: MemberJson{ID: "m2", UserID: "2"}},
	})

	{ // Try to apply difference to empty balance
		if _, err := groupEntity.ApplyBillBalanceDifference("EUR", BillBalanceDifference{"m1": 100}); err == nil {
			t.Error("Shod return error")
		}
	}

	members := append(groupEntity.GetGroupMembers(), GroupMemberJson{MemberJson: MemberJson{ID: "m2", UserID: "2", Name: "Second member"}})
	if changed := groupEntity.SetGroupMembers(members); !changed {
		t.Fatalf("Shor return changed=true")
	}

	//t.Log(groupEntity.GetGroupMembers())

	{ // Try to add another member
		changed, err := groupEntity.ApplyBillBalanceDifference("EUR", BillBalanceDifference{
			"m1": -400,
			"m2": 400,
		})
		if err != nil {
			t.Error(err)
		} else if !changed {
			t.Error("Should return changed=true")
		}
	}

	{ // test splti first + paid then
		groupEntity = GroupEntity{}
	}
}
