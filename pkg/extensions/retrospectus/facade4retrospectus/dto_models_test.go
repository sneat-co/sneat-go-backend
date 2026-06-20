package facade4retrospectus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus/dbo4retrospectus"
	"github.com/stretchr/testify/assert"
)

func newRetroRequest() RetroRequest {
	return RetroRequest{
		SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: "space1"},
		MeetingID:    "meeting1",
	}
}

func TestMoveRetroItemRequest_Validate(t *testing.T) {
	newRequest := func() MoveRetroItemRequest {
		return MoveRetroItemRequest{
			Request: facade4meetingus.Request{
				SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: "space1"},
				MeetingID:    "meeting1",
			},
			Item: "item1",
			From: dbo4retrospectus.TreePosition{Parent: "good", Index: 0},
			To:   dbo4retrospectus.TreePosition{Parent: "good", Index: 1},
		}
	}
	tests := []struct {
		name    string
		mutate  func(r *MoveRetroItemRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*MoveRetroItemRequest) {}, wantErr: false},
		{name: "missing_space", mutate: func(r *MoveRetroItemRequest) { r.SpaceID = "" }, wantErr: true},
		{name: "missing_meeting", mutate: func(r *MoveRetroItemRequest) { r.MeetingID = "" }, wantErr: true},
		{name: "negative_from_index", mutate: func(r *MoveRetroItemRequest) { r.From.Index = -1 }, wantErr: true},
		{name: "negative_to_index", mutate: func(r *MoveRetroItemRequest) { r.To.Index = -1 }, wantErr: true},
		{name: "same_position", mutate: func(r *MoveRetroItemRequest) { r.To = r.From }, wantErr: true},
		{name: "to_index_too_large", mutate: func(r *MoveRetroItemRequest) { r.To.Index = 101 }, wantErr: true},
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

func TestStartRetrospectiveRequest_Validate(t *testing.T) {
	newRequest := func() StartRetrospectiveRequest {
		return StartRetrospectiveRequest{
			RetroRequest:       newRetroRequest(),
			DurationsInMinutes: RetroDurations{Feedback: 10, Review: 5},
		}
	}
	tests := []struct {
		name    string
		mutate  func(r *StartRetrospectiveRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*StartRetrospectiveRequest) {}, wantErr: false},
		{name: "zero_durations", mutate: func(r *StartRetrospectiveRequest) { r.DurationsInMinutes = RetroDurations{} }, wantErr: false},
		{name: "missing_space", mutate: func(r *StartRetrospectiveRequest) { r.SpaceID = "" }, wantErr: true},
		{name: "missing_meeting", mutate: func(r *StartRetrospectiveRequest) { r.MeetingID = "" }, wantErr: true},
		{name: "negative_feedback", mutate: func(r *StartRetrospectiveRequest) { r.DurationsInMinutes.Feedback = -1 }, wantErr: true},
		{name: "negative_review", mutate: func(r *StartRetrospectiveRequest) { r.DurationsInMinutes.Review = -1 }, wantErr: true},
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
