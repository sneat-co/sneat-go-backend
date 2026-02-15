package dto4listus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestListRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     ListRequest
		wantErr bool
	}{
		{"valid", ListRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			ListID:       dbo4listus.NewListKey(dbo4listus.ListTypeToDo, "123"),
		}, false},
		{"invalid_id", ListRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			ListID:       "invalid",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ListRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateListRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateListRequest
		wantErr bool
	}{
		{"valid", CreateListRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			Type:         dbo4listus.ListTypeToDo,
			Title:        "My List",
		}, false},
		{"missing_type", CreateListRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			Title:        "My List",
		}, true},
		{"missing_title", CreateListRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			Type:         dbo4listus.ListTypeToDo,
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateListRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
