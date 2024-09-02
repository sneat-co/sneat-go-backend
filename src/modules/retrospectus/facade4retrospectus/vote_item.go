package facade4retrospectus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// VoteItemRequest request
type VoteItemRequest struct {
	RetroItemRequest
	Points int `json:"points"`
}

// Validate validates request
func (v *VoteItemRequest) Validate() error {
	if err := v.RetroItemRequest.Validate(); err != nil {
		return err
	}
	if v.Points == 0 {
		return validation.NewErrRecordIsMissingRequiredField("points")
	}
	return nil
}

// VoteItem votes an item
func VoteItem(ctx context.Context, userCtx facade.UserContext, request VoteItemRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	uid := userCtx.GetUserID()
	err := runRetroWorker(ctx, userCtx, request.Request,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) error {
			retrospective := params.Meeting.Record.Data().(*dbo4retrospectus.Retrospective)
			nodesByID, err := retrospective.GetMapOfRetroItemsByID()
			if err != nil {
				return err
			}
			itemNode := nodesByID[request.Item]
			item := itemNode.Item()
			points := item.VotesByUser[uid]
			if points == request.Points {
				return nil
			}
			updates := []dal.Update{{
				Field: fmt.Sprintf("%v.votesByUser.%v", itemNode.GetUpdatePath(nodesByID), uid),
			}}
			if request.Points == 0 {
				updates[0].Value = dal.DeleteField
			} else {
				updates[0].Value = request.Points
			}
			item.VotesByUser[uid] = request.Points
			if err = txUpdateRetrospective(ctx, tx, params.Meeting.Key, retrospective, updates); err != nil {
				return err
			}
			return err
		})
	return err
}
