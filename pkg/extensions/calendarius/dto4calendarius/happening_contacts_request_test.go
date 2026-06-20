package dto4calendarius

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/stretchr/testify/assert"
)

func validContactRef() dbo4linkage.ShortSpaceModuleItemRef {
	return dbo4linkage.ShortSpaceModuleItemRef{ID: "contact1"}
}

func TestHappeningContactRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     HappeningContactRequest
		wantErr bool
	}{
		{"valid", HappeningContactRequest{
			HappeningRequest: validHappeningRequest(),
			Contact:          validContactRef(),
		}, false},
		{"invalid_happening_request", HappeningContactRequest{
			HappeningRequest: HappeningRequest{},
			Contact:          validContactRef(),
		}, true},
		{"missing_contact_id", HappeningContactRequest{
			HappeningRequest: validHappeningRequest(),
			Contact:          dbo4linkage.ShortSpaceModuleItemRef{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHappeningContactsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     HappeningContactsRequest
		wantErr bool
	}{
		{"valid", HappeningContactsRequest{
			HappeningRequest: validHappeningRequest(),
			Contacts:         []dbo4linkage.ShortSpaceModuleItemRef{validContactRef()},
		}, false},
		{"invalid_happening_request", HappeningContactsRequest{
			HappeningRequest: HappeningRequest{},
			Contacts:         []dbo4linkage.ShortSpaceModuleItemRef{validContactRef()},
		}, true},
		{"empty_contacts", HappeningContactsRequest{
			HappeningRequest: validHappeningRequest(),
		}, true},
		{"invalid_contact", HappeningContactsRequest{
			HappeningRequest: validHappeningRequest(),
			Contacts:         []dbo4linkage.ShortSpaceModuleItemRef{{}},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
