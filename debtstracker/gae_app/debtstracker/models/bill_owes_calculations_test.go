package models

import (
	"testing"

	"github.com/strongo/decimal"
)

func TestUpdateMemberOwesForEqualSplit(t *testing.T) {
	members := []BillMemberJson{{}, {}, {}, {}}
	updateMemberOwesForEqualSplit(1001, "", members)
	verifyMemberOwes := func(i int, expecting decimal.Decimal64p2) {
		t.Helper()
		if members[i].Owes != expecting {
			t.Errorf("members[%d].Owes:%v != %v", i, members[i].Owes, expecting)
		}
	}
	verifyMemberOwes(0, 251)
	verifyMemberOwes(1, 250)
	verifyMemberOwes(2, 250)
	verifyMemberOwes(3, 250)
	//t.Logf("members +v: %+v", members)
}

func TestUpdateMemberOwesForEqualSplitWithAdjustment(t *testing.T) {
	members := []BillMemberJson{{}, {Adjustment: 200}, {}, {}}
	updateMemberOwesForEqualSplit(1001, "", members)
	verifyMemberOwes := func(i int, expecting decimal.Decimal64p2) {
		t.Helper()
		if members[i].Owes != expecting {
			t.Errorf("members[%d].Owes:%v != %v", i, members[i].Owes, expecting)
		}
	}
	verifyMemberOwes(0, 201)
	verifyMemberOwes(1, 400)
	verifyMemberOwes(2, 200)
	verifyMemberOwes(3, 200)
	//t.Logf("members +v: %+v", members)
}

func validateTotal(t *testing.T, members []BillMemberJson, expectedTotal decimal.Decimal64p2) {
	var total decimal.Decimal64p2
	for _, m := range members {
		total += m.Owes
	}
	if total != expectedTotal {
		t.Fatal("Wrong total", total)
	}
}

func TestUpdateMemberOwesForSplitByShares(t *testing.T) {
	var members []BillMemberJson

	members = []BillMemberJson{{}, {}, {}, {}}
	updateMemberOwesForSplitByShares(1000, "", members)
	if members[0].Owes != 250 || members[1].Owes != 250 || members[2].Owes != 250 || members[3].Owes != 250 {
		t.Fatal(members)
		return
	}
	validateTotal(t, members, 1000)

	members = []BillMemberJson{{}, {}, {}}
	updateMemberOwesForSplitByShares(1000, "", members)
	if members[0].Owes != 334 || members[1].Owes != 333 || members[2].Owes != 333 {
		t.Fatal(members)
		return
	}
	validateTotal(t, members, 1000)

	members = []BillMemberJson{{MemberJson: MemberJson{Shares: 3}}, {MemberJson: MemberJson{Shares: 2}}, {MemberJson: MemberJson{Shares: 1}}}
	updateMemberOwesForSplitByShares(1000, "", members)
	validateTotal(t, members, 1000)
}
