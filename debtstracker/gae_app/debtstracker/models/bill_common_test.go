package models

import (
	"testing"

	"github.com/strongo/decimal"
)

func TestBillCommon_GetBillMembers(t *testing.T) {
	billEntity := BillEntity{
		BillCommon: BillCommon{
			MembersJson:  `[{"ID": "2", "Name": "Second", "Paid": 2.47}, {"ID": "3", "Name": "Third"}]`,
			MembersCount: 2,
		},
	}

	billMembers := billEntity.GetBillMembers()

	if len(billMembers) != 2 {
		t.Errorf("Expected 2 memebers, got: %d", len(billMembers))
		return
	}

	verifyMember := func(m BillMemberJson, id, name string, paid decimal.Decimal64p2) {
		if m.ID != id {
			t.Errorf("Got ID: %v; expected: %v", m.ID, id)
		}
		if m.Name != name {
			t.Errorf("Got name: %v; expected: %v", m.Name, name)
		}
		if m.Paid != paid {
			t.Errorf("Got Paid: %v; expected: %v", m.Paid, paid)
		}
	}
	verifyMember(billMembers[0], "2", "Second", decimal.NewDecimal64p2(2, 47))
	verifyMember(billMembers[1], "3", "Third", 0)
}
