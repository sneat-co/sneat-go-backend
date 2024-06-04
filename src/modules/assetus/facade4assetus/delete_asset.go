package facade4assetus

import (
	"context"
	"errors"
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
	if user == nil || user.GetID() == "" {
		return errors.New("no user context")
	}
	asset := dal4assetus.NewAssetEntry(request.TeamID, request.ID)
	input := dal4teamus.TeamItemRunnerInput[*dbo4assetus.AssetusTeamDbo]{
		TeamRequest:        request.TeamRequest,
		BriefsAdapter:      nil,
		IsTeamRecordNeeded: true,
		Counter:            "assets",
		ShortID:            request.ID,
		TeamItem:           asset.Record,
	}
	err = dal4teamus.DeleteTeamItem(ctx, user, input, const4assetus.ModuleID, new(dbo4assetus.AssetusTeamDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.TeamItemWorkerParams) (err error) {
			return nil //errors.New("not implemented")
		})
	return
}
