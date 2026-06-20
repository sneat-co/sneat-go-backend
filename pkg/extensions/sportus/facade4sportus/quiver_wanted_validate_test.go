package facade4sportus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/sportus/dbo4sportus"
	"github.com/stretchr/testify/assert"
)

func TestCreateWantedRequest_Validate(t *testing.T) {
	now := time.Now()
	newRequest := func() CreateWantedRequest {
		return CreateWantedRequest{
			Wanted: dbo4sportus.Wanted{
				Item: dbo4sportus.Item{
					UserID:    "user1",
					Status:    "active",
					DtCreated: now,
					DtUpdated: now,
				},
			},
		}
	}
	tests := []struct {
		name    string
		mutate  func(r *CreateWantedRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*CreateWantedRequest) {}, wantErr: false},
		{name: "missing_user", mutate: func(r *CreateWantedRequest) { r.UserID = "" }, wantErr: true},
		{name: "missing_status", mutate: func(r *CreateWantedRequest) { r.Status = "" }, wantErr: true},
		{name: "unknown_status", mutate: func(r *CreateWantedRequest) { r.Status = "bogus" }, wantErr: true},
		{name: "missing_dt_created", mutate: func(r *CreateWantedRequest) { r.DtCreated = time.Time{} }, wantErr: true},
		{name: "missing_dt_updated", mutate: func(r *CreateWantedRequest) { r.DtUpdated = time.Time{} }, wantErr: true},
		{name: "updated_before_created", mutate: func(r *CreateWantedRequest) { r.DtUpdated = now.Add(-time.Hour) }, wantErr: true},
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
