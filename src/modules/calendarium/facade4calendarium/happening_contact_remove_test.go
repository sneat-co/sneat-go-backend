package facade4calendarium

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"testing"
)

func TestRemoveParticipantsFromHappening(t *testing.T) {
	type args struct {
		ctx     facade.ContextWithUser
		request dto4calendarium.HappeningContactsRequest
	}
	tests := []struct { // TODO(help-wanted): Add tests cases
		name     string
		args     args
		checkErr func(t *testing.T, err error)
	}{
		{
			name: "remove participants from happening without participants",
			args: args{
				ctx: facade.NewContextWithUserID(context.Background(), "user1"),
				request: dto4calendarium.HappeningContactsRequest{
					HappeningRequest: dto4calendarium.HappeningRequest{
						SpaceRequest: dto4spaceus.SpaceRequest{
							SpaceID: "space1",
						},
					},
				},
			},
			checkErr: func(t *testing.T, err error) {
				if !validation.IsBadRequestError(err) {
					t.Errorf("RemoveParticipantsFromHappening() expected BadRequestError but returned  = %T: %v", err, err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RemoveParticipantsFromHappening(tt.args.ctx, tt.args.request)
			if err != nil {
				if tt.checkErr == nil {
					t.Errorf("RemoveParticipantsFromHappening() returned unexpected error = %v", err)
				} else {
					tt.checkErr(t, err)
				}
			}
		})
	}
}
