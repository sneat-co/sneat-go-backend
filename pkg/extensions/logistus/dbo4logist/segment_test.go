package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegmentDates_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       SegmentDates
		wantErr bool
	}{
		{"empty", SegmentDates{}, false},
		{"valid", SegmentDates{Departs: "2023-01-01", Arrives: "2023-01-02"}, false},
		{"equal", SegmentDates{Departs: "2023-01-01", Arrives: "2023-01-01"}, false},
		{"bad_departs", SegmentDates{Departs: "not-a-date"}, true},
		{"bad_arrives", SegmentDates{Arrives: "not-a-date"}, true},
		{"arrives_before_departs", SegmentDates{Departs: "2023-01-02", Arrives: "2023-01-01"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSegmentsFilter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       SegmentsFilter
		wantErr bool
	}{
		{"empty", SegmentsFilter{}, false},
		{"bad_contact_id", SegmentsFilter{ByContactID: " bad "}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContainerSegmentKey_String(t *testing.T) {
	k := ContainerSegmentKey{
		ContainerID: "c1",
		From:        SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "a", Role: CounterpartyRoleTrucker}},
		To:          SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "b", Role: CounterpartyRolePickPoint}},
	}
	s := k.String()
	assert.Contains(t, s, "container=c1")
	assert.Contains(t, s, "from=trucker:a")
	assert.Contains(t, s, "to=pick_point:b")
}

func TestContainerSegmentKey_Validate(t *testing.T) {
	valid := ContainerSegmentKey{
		ContainerID: "c1",
		From:        SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "a", Role: CounterpartyRoleTrucker}},
		To:          SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "b", Role: CounterpartyRoleTrucker}},
	}
	assert.NoError(t, valid.Validate())

	badFrom := valid
	badFrom.From = SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{Role: CounterpartyRoleTrucker}}
	assert.Error(t, badFrom.Validate())

	badTo := valid
	badTo.To = SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{Role: CounterpartyRoleTrucker}}
	assert.Error(t, badTo.Validate())

	sameContact := valid
	sameContact.To = sameContact.From
	assert.Error(t, sameContact.Validate())
}

func TestContainerSegment_Validate_table(t *testing.T) {
	base := func() ContainerSegment {
		return ContainerSegment{
			ContainerSegmentKey: ContainerSegmentKey{
				ContainerID: "c1",
				From:        SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "a", Role: CounterpartyRoleTrucker}},
				To:          SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "b", Role: CounterpartyRoleTrucker}},
			},
		}
	}
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, base().Validate())
	})
	t.Run("valid_with_dates", func(t *testing.T) {
		v := base()
		v.Dates = &SegmentDates{Departs: "2023-01-01", Arrives: "2023-01-02"}
		assert.NoError(t, v.Validate())
	})
	t.Run("bad_key", func(t *testing.T) {
		v := base()
		v.From = SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{Role: CounterpartyRoleTrucker}} // missing contactID
		assert.Error(t, v.Validate())
	})
	t.Run("bad_dates", func(t *testing.T) {
		v := base()
		v.Dates = &SegmentDates{Departs: "bad"}
		assert.Error(t, v.Validate())
	})
}

func TestWithSegments_Validate(t *testing.T) {
	seg := func(from, to string) *ContainerSegment {
		return &ContainerSegment{
			ContainerSegmentKey: ContainerSegmentKey{
				ContainerID: "c1",
				From:        SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: from, Role: CounterpartyRoleTrucker}},
				To:          SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: to, Role: CounterpartyRoleTrucker}},
			},
		}
	}
	assert.NoError(t, WithSegments{}.Validate())
	assert.NoError(t, WithSegments{Segments: []*ContainerSegment{seg("a", "b"), seg("c", "d")}}.Validate())

	// invalid child segment (missing from.contactID)
	bad := seg("a", "b")
	bad.From = SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{Role: CounterpartyRoleTrucker}}
	assert.Error(t, WithSegments{Segments: []*ContainerSegment{bad}}.Validate())

	// duplicate keys
	assert.Error(t, WithSegments{Segments: []*ContainerSegment{seg("a", "b"), seg("a", "b")}}.Validate())
}

func TestWithSegments_GetSegmentsByFilter(t *testing.T) {
	s1 := &ContainerSegment{ContainerSegmentKey: ContainerSegmentKey{ContainerID: "c1",
		From: SegmentEndpoint{ShippingPointID: "sp1"}, To: SegmentEndpoint{ShippingPointID: "sp2"}}, ByContactID: "x"}
	s2 := &ContainerSegment{ContainerSegmentKey: ContainerSegmentKey{ContainerID: "c2",
		From: SegmentEndpoint{ShippingPointID: "sp3"}, To: SegmentEndpoint{ShippingPointID: "sp4"}}, ByContactID: "y"}
	ws := WithSegments{Segments: []*ContainerSegment{s1, s2}}

	assert.Len(t, ws.GetSegmentsByFilter(SegmentsFilter{}), 2)
	assert.Equal(t, []*ContainerSegment{s1}, ws.GetSegmentsByFilter(SegmentsFilter{ContainerIDs: []string{"c1"}}))
	assert.Equal(t, []*ContainerSegment{s1}, ws.GetSegmentsByFilter(SegmentsFilter{FromShippingPointID: "sp1"}))
	assert.Equal(t, []*ContainerSegment{s2}, ws.GetSegmentsByFilter(SegmentsFilter{ToShippingPointID: "sp4"}))
	assert.Equal(t, []*ContainerSegment{s2}, ws.GetSegmentsByFilter(SegmentsFilter{ByContactID: "y"}))
	assert.Empty(t, ws.GetSegmentsByFilter(SegmentsFilter{ContainerIDs: []string{"none"}}))
}

func TestWithSegments_GetSegmentByKey(t *testing.T) {
	s1 := &ContainerSegment{ContainerSegmentKey: ContainerSegmentKey{ContainerID: "c1",
		From: SegmentEndpoint{ShippingPointID: "sp1"}, To: SegmentEndpoint{ShippingPointID: "sp2"}}}
	ws := WithSegments{Segments: []*ContainerSegment{s1}}

	got := ws.GetSegmentByKey(ContainerSegmentKey{ContainerID: "c1", From: SegmentEndpoint{ShippingPointID: "sp1"}})
	assert.Equal(t, s1, got)

	got = ws.GetSegmentByKey(ContainerSegmentKey{ContainerID: "c1", To: SegmentEndpoint{ShippingPointID: "sp2"}})
	assert.Equal(t, s1, got)

	assert.Nil(t, ws.GetSegmentByKey(ContainerSegmentKey{ContainerID: "c1", From: SegmentEndpoint{ShippingPointID: "other"}}))
	assert.Nil(t, ws.GetSegmentByKey(ContainerSegmentKey{ContainerID: "missing"}))

	assert.Panics(t, func() { ws.GetSegmentByKey(ContainerSegmentKey{}) })
}

func TestWithSegments_DeleteSegments(t *testing.T) {
	s1 := &ContainerSegment{ContainerSegmentKey: ContainerSegmentKey{ContainerID: "c1",
		From: SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "a"}},
		To:   SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "b"}}}}
	s2 := &ContainerSegment{ContainerSegmentKey: ContainerSegmentKey{ContainerID: "c2",
		From: SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "x"}},
		To:   SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "y"}}}}
	ws := WithSegments{Segments: []*ContainerSegment{s1, s2}}
	remaining := ws.DeleteSegments([]*ContainerSegment{s1})
	// s1 is dropped where it matches; s2 stays
	assert.Contains(t, remaining, s2)
}

func TestWithSegments_Updates(t *testing.T) {
	empty := WithSegments{}.Updates()
	assert.Len(t, empty, 1)
	assert.Equal(t, "segments", empty[0].FieldName())

	withData := WithSegments{Segments: []*ContainerSegment{{}}}.Updates()
	assert.Len(t, withData, 1)
	assert.Equal(t, "segments", withData[0].FieldName())
}
