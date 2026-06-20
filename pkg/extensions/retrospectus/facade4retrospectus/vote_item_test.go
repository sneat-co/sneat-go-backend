package facade4retrospectus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoteItemRequest_Validate(t *testing.T) {
	newRequest := func() VoteItemRequest {
		return VoteItemRequest{
			RetroItemRequest: newRetroItemRequest(),
			Points:           3,
		}
	}
	tests := []struct {
		name    string
		mutate  func(r *VoteItemRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*VoteItemRequest) {}, wantErr: false},
		{name: "negative_points", mutate: func(r *VoteItemRequest) { r.Points = -1 }, wantErr: false},
		{name: "zero_points", mutate: func(r *VoteItemRequest) { r.Points = 0 }, wantErr: true},
		{name: "missing_space", mutate: func(r *VoteItemRequest) { r.SpaceID = "" }, wantErr: true},
		{name: "missing_type", mutate: func(r *VoteItemRequest) { r.Type = "" }, wantErr: true},
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
