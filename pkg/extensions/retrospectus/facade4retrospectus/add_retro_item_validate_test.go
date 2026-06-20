package facade4retrospectus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/facade4meetingus"
	"github.com/stretchr/testify/assert"
)

func newRetroItemRequest() RetroItemRequest {
	return RetroItemRequest{
		Request: facade4meetingus.Request{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: "space1"},
			MeetingID:    "meeting1",
		},
		Type: "good",
		Item: "item1",
	}
}

func TestRetroItemRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(r *RetroItemRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*RetroItemRequest) {}, wantErr: false},
		{name: "missing_space", mutate: func(r *RetroItemRequest) { r.SpaceID = "" }, wantErr: true},
		{name: "missing_meeting", mutate: func(r *RetroItemRequest) { r.MeetingID = "" }, wantErr: true},
		{name: "missing_type", mutate: func(r *RetroItemRequest) { r.Type = "" }, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRetroItemRequest()
			tt.mutate(&r)
			err := r.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddRetroItemRequest_Validate(t *testing.T) {
	newRequest := func() AddRetroItemRequest {
		return AddRetroItemRequest{
			RetroItemRequest: newRetroItemRequest(),
			Title:            "Good item",
		}
	}
	tests := []struct {
		name    string
		mutate  func(r *AddRetroItemRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*AddRetroItemRequest) {}, wantErr: false},
		{name: "missing_space", mutate: func(r *AddRetroItemRequest) { r.SpaceID = "" }, wantErr: true},
		{name: "missing_type", mutate: func(r *AddRetroItemRequest) { r.Type = "" }, wantErr: true},
		{name: "missing_title", mutate: func(r *AddRetroItemRequest) { r.Title = "" }, wantErr: true},
		{name: "title_with_leading_space", mutate: func(r *AddRetroItemRequest) { r.Title = " bad" }, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRequest()
			tt.mutate(&r)
			err := r.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
