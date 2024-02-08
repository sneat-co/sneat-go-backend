package models

import "testing"

func TestBillBalanceByMember_BillDifference(t *testing.T) {
	previous := BillBalanceByMember{}
	current := BillBalanceByMember{}

	{ // Test empty
		if diff := current.BillBalanceDifference(previous); len(diff) != 0 {
			t.Error("Should be no difference", diff)
		}
	}

	{ // Test non empty current and empty previous
		previous = BillBalanceByMember{}
		current = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 1200, Owes: 400},
		}
		if diff := current.BillBalanceDifference(previous); len(diff) != 1 {
			t.Fatal("Should have single item", diff)
		} else if md, ok := diff["m1"]; !ok {
			t.Fatal("Item should be m1", diff)
		} else if md != 800 {
			t.Fatal("Wrong value", diff)
		}
	}

	{ // Test increase in Paid
		previous = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 10, Owes: 4},
		}
		current = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 12, Owes: 4},
		}
		if diff := current.BillBalanceDifference(previous); len(diff) != 1 {
			t.Fatal("Should have single item", diff)
		} else if md, ok := diff["m1"]; !ok {
			t.Fatal("Item should be m1", diff)
		} else if md != 2 {
			t.Fatal("Wrong value", diff)
		}
	}

	{ // Test increase in Owes
		previous = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 12, Owes: 1},
		}
		current = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 12, Owes: 4},
		}
		if diff := current.BillBalanceDifference(previous); len(diff) != 1 {
			t.Fatal("Should have single item", diff)
		} else if md, ok := diff["m1"]; !ok {
			t.Fatal("Item should be m1", diff)
		} else if md != -3 {
			t.Fatal("Wrong value", diff)
		}
	}

	{ // Test decrease in Paid & Owes
		previous = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 1500, Owes: 900},
		}
		current = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 1200, Owes: 400},
		}
		if diff := current.BillBalanceDifference(previous); len(diff) != 1 {
			t.Fatal("Should have single item", diff)
		} else if md, ok := diff["m1"]; !ok {
			t.Fatal("Item should be m1", diff)
		} else if md != 200 {
			t.Fatal("Wrong value", diff)
		}
	}

	{ // Test in member added
		previous = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 1200, Owes: 1200},
		}
		current = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 1200, Owes: 600},
			"m2": BillMemberBalance{Paid: 0, Owes: 600},
		}
		if diff := current.BillBalanceDifference(previous); len(diff) != 2 {
			t.Fatal("Should have 2 items", diff)
		} else if m1, ok := diff["m1"]; !ok {
			t.Fatal("Item should be m1", diff)
		} else if m2, ok := diff["m2"]; !ok {
			t.Fatal("Item should be m1", diff)
		} else if m1 != 600 {
			t.Fatal("Wrong m1 diff", m1)
		} else if m2 != -600 {
			t.Fatal("Wrong m2 diff", m2)
		}
	}

	{ // Test in member swapped
		previous = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 1200, Owes: 600},
			"m2": BillMemberBalance{Paid: 0, Owes: 600},
		}
		current = BillBalanceByMember{
			"m1": BillMemberBalance{Paid: 1200, Owes: 600},
			"m3": BillMemberBalance{Paid: 0, Owes: 600},
		}
		if diff := current.BillBalanceDifference(previous); len(diff) != 2 {
			t.Fatal("Should have 2 items", diff)
		} else if m2, ok := diff["m2"]; !ok {
			t.Fatal("Item should be m2", diff)
		} else if m3, ok := diff["m3"]; !ok {
			t.Fatal("Item should be m3", diff)
		} else if m2 != 600 {
			t.Fatal("Wrong m2 diff", m2)
		} else if m3 != -600 {
			t.Fatal("Wrong m3 diff", m3)
		}
	}
}

func TestBillBalanceDifference_IsAffectingGroupBalance(t *testing.T) {
	var diff BillBalanceDifference

	{ // verify empty
		diff = BillBalanceDifference{}
		if diff.IsNoDifference() {
			t.Fatal("should be false for empty map")
		}
	}

	{ // verify paid=owes for single member
		diff = BillBalanceDifference{
			"m1": 0,
		}
		if diff.IsNoDifference() {
			t.Fatal("should be false for empty map")
		}
	}

	{ // verify paid=owes for 2 members
		diff = BillBalanceDifference{
			"m1": 0,
			"m2": 0,
		}
		if diff.IsNoDifference() {
			t.Fatal("should be false for empty map")
		}
	}
}
