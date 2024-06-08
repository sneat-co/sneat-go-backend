package facade4assetus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteAsset deletes an asset
func DeleteAsset(ctx context.Context, user facade.User, request dal4teamus.TeamItemRequest) (err error) {
	if err = request.Validate(); err != nil {
		return fmt.Errorf("invalid request to facade4assetus.DeleteAsset: %w", err)
	}
	var getBriefsCount = func(teamModuleDbo *dbo4assetus.AssetusTeamDbo) int {
		return len(teamModuleDbo.Assets)
	}
	briefsAdapter := dal4teamus.NewMapBriefsAdapter[*dbo4assetus.AssetusTeamDbo](
		getBriefsCount,
		func(teamModuleDbo *dbo4assetus.AssetusTeamDbo, id string) ([]dal.Update, error) {
			delete(teamModuleDbo.Assets, id)
			return []dal.Update{{Field: "assets." + id, Value: dal.DeleteField}}, teamModuleDbo.Validate()
		},
	)

	return dal4teamus.DeleteTeamItem(ctx, user, request,
		const4assetus.ModuleID, new(dbo4assetus.AssetusTeamDbo),
		dal4assetus.AssetsCollection, new(dbo4assetus.AssetDbo),
		briefsAdapter,
		deleteAssetTxWorker,
	)
}

func deleteAssetTxWorker(_ context.Context, _ dal.ReadwriteTransaction, _ *dal4teamus.TeamItemWorkerParams[*dbo4assetus.AssetusTeamDbo, *dbo4assetus.AssetDbo]) (err error) {
	return nil
}
