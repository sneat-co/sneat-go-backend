package facade4retrospectus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/models4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// MoveRetroItem moves item
func MoveRetroItem(ctx context.Context, userFacade facade.User, request MoveRetroItemRequest) (err error) {
	uid := userFacade.GetID()
	var retrospectiveKey *dal.Key
	if request.MeetingID == UpcomingRetrospectiveID {
		retrospectiveKey = models4retrospectus.NewRetrospectiveKey(request.TeamID, models4userus.NewUserKey(uid))
	} else {
		retrospectiveKey = models4retrospectus.NewRetrospectiveKey(request.MeetingID, newTeamKey(request.TeamID))
	}

	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		retrospective := new(models4retrospectus.Retrospective)
		retrospectiveRecord := dal.NewRecordWithData(retrospectiveKey, retrospective)

		err = tx.Get(ctx, retrospectiveRecord)
		//err = txGetRetrospective(ctx, tx, retrospectiveRecord)
		//panic(fmt.Sprintf("err: %v\n retrospectiveRecord: %+v", err, retrospectiveRecord.Data()))
		if err != nil {
			return err
		} else if err := retrospectiveRecord.Error(); err != nil {
			return fmt.Errorf("retrospectiveRecord.Error(): %w", err)
		} else if !retrospectiveRecord.Exists() {
			return fmt.Errorf("retrospective not found by id: %v-%v", request.TeamID, request.MeetingID)
		}
		if err = models4retrospectus.MoveRetroItem(retrospective.Items, request.Item, request.From, request.To); err != nil {
			return err
		}
		if err = txUpdateRetrospective(ctx, tx, retrospectiveKey, retrospective, []dal.Update{
			{
				Field: "items",
				Value: retrospective.Items,
			},
		}); err != nil {
			return err
		}
		return nil
	})

	return err
}
