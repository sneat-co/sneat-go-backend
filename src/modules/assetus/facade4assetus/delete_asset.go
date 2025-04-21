package facade4assetus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteAsset deletes an asset
func DeleteAsset(ctx facade.ContextWithUser, request dto4spaceus.SpaceItemRequest) (err error) {
	if err = request.Validate(); err != nil {
		return fmt.Errorf("invalid request to facade4assetus.DeleteAsset: %w", err)
	}
	var getBriefsCount = func(teamModuleDbo *dbo4assetus.AssetusSpaceDbo) int {
		return len(teamModuleDbo.Assets)
	}
	briefsAdapter := dal4spaceus.NewMapBriefsAdapter[*dbo4assetus.AssetusSpaceDbo](
		getBriefsCount,
		func(teamModuleDbo *dbo4assetus.AssetusSpaceDbo, id string) ([]update.Update, error) {
			delete(teamModuleDbo.Assets, id)
			return []update.Update{update.ByFieldName("assets."+id, update.DeleteField)}, teamModuleDbo.Validate()
		},
	)

	return dal4spaceus.DeleteSpaceItem(ctx, ctx.User(), request,
		const4assetus.ModuleID, new(dbo4assetus.AssetusSpaceDbo),
		dal4assetus.AssetsCollection, new(dbo4assetus.AssetDbo),
		briefsAdapter,
		deleteAssetTxWorker,
	)
}

func deleteAssetTxWorker(_ context.Context, _ dal.ReadwriteTransaction, _ *dal4spaceus.SpaceItemWorkerParams[*dbo4assetus.AssetusSpaceDbo, *dbo4assetus.AssetDbo]) (err error) {
	return nil
}
