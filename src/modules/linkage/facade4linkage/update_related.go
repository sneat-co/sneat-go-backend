package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateItemRelationships(ctx context.Context, userCtx facade.User, request dto4linkage.UpdateItemRequest) (err error) {
	err = dal4teamus.RunTeamWorker(ctx, userCtx, request.TeamID, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.TeamWorkerParams) (err error) {
		return nil
	})
	return err
}
