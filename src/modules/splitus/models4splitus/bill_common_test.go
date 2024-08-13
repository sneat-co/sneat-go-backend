package models4splitus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"testing"

	"github.com/strongo/decimal"
)

func TestBillCommon_GetBillMembers(t *testing.T) {
	billEntity := BillDbo{
		BillCommon: BillCommon{
			Members: []*briefs4splitus.BillMemberBrief{
				{
					Paid: decimal.NewDecimal64p2(2, 47),
					MemberBrief: briefs4splitus.MemberBrief{
						ID:   "2",
						Name: "Second",
					},
				},
				{
					MemberBrief: briefs4splitus.MemberBrief{
						ID:   "3",
						Name: "Third",
					},
				},
			},
		},
	}

	billMembers := billEntity.GetBillMembers()

	if len(billMembers) != 2 {
		t.Errorf("Expected 2 memebers, got: %d", len(billMembers))
		return
	}

	verifyMember := func(m *briefs4splitus.BillMemberBrief, id, name string, paid decimal.Decimal64p2) {
		t.Helper()
		if m.ID != id {
			t.Errorf("Got ContactID: %v; expected: %v", m.ID, id)
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
