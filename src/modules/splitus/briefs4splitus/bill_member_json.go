package briefs4splitus

//go:generate ffjson $GOFILE

import (
	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/decimal"
)

type MemberBrief struct {
	// We use string IDs as it's faster to marshal and will be more compact in future
	ID            string                      // required
	Name          string                      // required
	AddedByUserID string                      `json:",omitempty"`
	UserID        string                      `json:",omitempty"`
	TgUserID      string                      `json:",omitempty"`
	ContactIDs    []string                    `json:",omitempty"`
	ContactByUser MemberContactBriefsByUserID `json:",omitempty"`
	Shares        int                         `json:",omitempty"`
}

var _ SplitMember = (*MemberBrief)(nil)

func (m MemberBrief) GetID() string {
	return m.ID
}

func (m MemberBrief) GetName() string {
	return m.Name
}

func (m MemberBrief) GetShares() int {
	return m.Shares
}

type SpaceSplitMember struct {
	MemberBrief
	Balance money.Balance `json:",omitempty"`
}

var _ SplitMember = (*SpaceSplitMember)(nil)

func (m *SpaceSplitMember) String() string {
	buffer, _ := ffjson.MarshalFast(m)
	return string(buffer)
}

type GroupBalanceByCurrencyAndMember map[money.CurrencyCode]map[string]decimal.Decimal64p2

type BillMemberBrief struct {
	MemberBrief
	Paid       decimal.Decimal64p2 `json:",omitempty"`
	Owes       decimal.Decimal64p2 `json:",omitempty"`
	Percent    decimal.Decimal64p2 `json:",omitempty"`
	Adjustment decimal.Decimal64p2 `json:",omitempty"`
	//transferIDs []int             `json:",omitempty"`
}

func (m BillMemberBrief) Balance() decimal.Decimal64p2 {
	return m.Paid - m.Owes
}

func (m *BillMemberBrief) String() string {
	buffer, _ := ffjson.MarshalFast(m)
	return string(buffer)
}

type MemberContactBrief struct {
	ContactID   string
	ContactName string
}

type MemberContactBriefsByUserID map[string]MemberContactBrief
