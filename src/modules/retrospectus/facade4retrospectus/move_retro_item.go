package facade4retrospectus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// MoveRetroItem moves item
func MoveRetroItem(ctx context.Context, userCtx facade.UserContext, request MoveRetroItemRequest) (err error) {
	uid := userCtx.GetUserID()
	var retrospectiveKey *dal.Key
	if request.MeetingID == UpcomingRetrospectiveID {
		retrospectiveKey = dbo4retrospectus.NewRetrospectiveKey(string(request.SpaceID), dbo4userus.NewUserKey(uid))
	} else {
		retrospectiveKey = dbo4retrospectus.NewRetrospectiveKey(request.MeetingID, newSpaceKey(request.SpaceID))
	}

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		retrospective := new(dbo4retrospectus.Retrospective)
		retrospectiveRecord := dal.NewRecordWithData(retrospectiveKey, retrospective)

		err = tx.Get(ctx, retrospectiveRecord)
		//err = txGetRetrospective(ctx, tx, retrospectiveRecord)
		//panic(fmt.Sprintf("err: %v\n retrospectiveRecord: %+v", err, retrospectiveRecord.Data()))
		if err != nil {
			return err
		} else if err := retrospectiveRecord.Error(); err != nil {
			return fmt.Errorf("retrospectiveRecord.Error(): %w", err)
		} else if !retrospectiveRecord.Exists() {
			return fmt.Errorf("retrospective not found by id: %v-%v", request.SpaceID, request.MeetingID)
		}
		if err = dbo4retrospectus.MoveRetroItem(retrospective.Items, request.Item, request.From, request.To); err != nil {
			return err
		}
		if err = txUpdateRetrospective(ctx, tx, retrospectiveKey, retrospective, []update.Update{
			update.ByFieldName("items", retrospective.Items),
		}); err != nil {
			return err
		}
		return nil
	})

	return err
}
