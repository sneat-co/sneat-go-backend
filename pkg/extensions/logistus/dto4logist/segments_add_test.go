package dto4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/stretchr/testify/assert"
)

func validSegmentContainerData() SegmentContainerData {
	return SegmentContainerData{ID: "c1"}
}

func TestSegmentContainerData_Validate(t *testing.T) {
	badLoad := &dbo4logist.FreightLoad{NumberOfPallets: -1}
	tests := []struct {
		name    string
		v       SegmentContainerData
		wantErr bool
	}{
		{"valid", validSegmentContainerData(), false},
		{"missing_id", SegmentContainerData{}, true},
		{"valid_with_loads", SegmentContainerData{
			ID: "c1",
			FreightPoint: dbo4logist.FreightPoint{
				Tasks:    []dbo4logist.ShippingPointTask{"load", "unload"},
				ToLoad:   &dbo4logist.FreightLoad{NumberOfPallets: 1},
				ToUnload: &dbo4logist.FreightLoad{NumberOfPallets: 1},
			},
		}, false},
		{"bad_to_load", SegmentContainerData{
			ID:           "c1",
			FreightPoint: dbo4logist.FreightPoint{ToLoad: badLoad},
		}, true},
		{"bad_to_unload", SegmentContainerData{
			ID:           "c1",
			FreightPoint: dbo4logist.FreightPoint{ToUnload: badLoad},
		}, true},
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

func validAddSegmentParty() AddSegmentParty {
	return AddSegmentParty{
		Counterparty: dbo4logist.SegmentCounterparty{ContactID: "c1", Role: dbo4logist.CounterpartyRoleTrucker},
	}
}

func TestAddSegmentParty_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       AddSegmentParty
		wantErr bool
	}{
		{"valid", validAddSegmentParty(), false},
		{"valid_ref", AddSegmentParty{
			Counterparty: dbo4logist.SegmentCounterparty{ContactID: "c1", Role: dbo4logist.CounterpartyRoleTrucker},
			RefNumber:    "REF1",
		}, false},
		{"bad_counterparty", AddSegmentParty{}, true},
		{"untrimmed_ref", AddSegmentParty{
			Counterparty: dbo4logist.SegmentCounterparty{ContactID: "c1", Role: dbo4logist.CounterpartyRoleTrucker},
			RefNumber:    " REF1 ",
		}, true},
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

func validAddSegmentEndpoint() AddSegmentEndpoint {
	return AddSegmentEndpoint{AddSegmentParty: validAddSegmentParty()}
}

func TestAddSegmentEndpoint_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       AddSegmentEndpoint
		wantErr bool
	}{
		{"valid", validAddSegmentEndpoint(), false},
		{"valid_date", AddSegmentEndpoint{AddSegmentParty: validAddSegmentParty(), Date: "2023-01-01"}, false},
		{"bad_party", AddSegmentEndpoint{}, true},
		{"bad_date", AddSegmentEndpoint{AddSegmentParty: validAddSegmentParty(), Date: "not-a-date"}, true},
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

func TestAddSegmentsRequest_Validate(t *testing.T) {
	trucker := &AddSegmentParty{
		Counterparty: dbo4logist.SegmentCounterparty{ContactID: "c1", Role: dbo4logist.CounterpartyRoleTrucker},
	}
	nonTrucker := &AddSegmentParty{
		Counterparty: dbo4logist.SegmentCounterparty{ContactID: "c1", Role: dbo4logist.CounterpartyRoleDispatcher},
	}
	tests := []struct {
		name    string
		v       AddSegmentsRequest
		wantErr bool
	}{
		{"valid", AddSegmentsRequest{
			From:       validAddSegmentEndpoint(),
			To:         validAddSegmentEndpoint(),
			Containers: []SegmentContainerData{validSegmentContainerData()},
		}, false},
		{"valid_with_trucker_by", AddSegmentsRequest{
			From:       validAddSegmentEndpoint(),
			To:         validAddSegmentEndpoint(),
			By:         trucker,
			Containers: []SegmentContainerData{validSegmentContainerData()},
		}, false},
		{"bad_from", AddSegmentsRequest{
			From:       AddSegmentEndpoint{},
			To:         validAddSegmentEndpoint(),
			Containers: []SegmentContainerData{validSegmentContainerData()},
		}, true},
		{"bad_to", AddSegmentsRequest{
			From:       validAddSegmentEndpoint(),
			To:         AddSegmentEndpoint{},
			Containers: []SegmentContainerData{validSegmentContainerData()},
		}, true},
		{"bad_by", AddSegmentsRequest{
			From:       validAddSegmentEndpoint(),
			To:         validAddSegmentEndpoint(),
			By:         &AddSegmentParty{},
			Containers: []SegmentContainerData{validSegmentContainerData()},
		}, true},
		{"by_not_trucker", AddSegmentsRequest{
			From:       validAddSegmentEndpoint(),
			To:         validAddSegmentEndpoint(),
			By:         nonTrucker,
			Containers: []SegmentContainerData{validSegmentContainerData()},
		}, true},
		{"no_containers", AddSegmentsRequest{
			From: validAddSegmentEndpoint(),
			To:   validAddSegmentEndpoint(),
		}, true},
		{"bad_container", AddSegmentsRequest{
			From:       validAddSegmentEndpoint(),
			To:         validAddSegmentEndpoint(),
			Containers: []SegmentContainerData{{}},
		}, true},
		{"duplicate_container", AddSegmentsRequest{
			From: validAddSegmentEndpoint(),
			To:   validAddSegmentEndpoint(),
			Containers: []SegmentContainerData{
				{ID: "c1"},
				{ID: "c1"},
			},
		}, true},
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
