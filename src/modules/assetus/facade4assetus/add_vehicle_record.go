package facade4assetus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type CreateVehicleRecordResponse struct {
	ID string `json:"id"`
}

func AddVehicleRecord(ctx context.Context, user facade.UserContext, request dto4assetus.AddVehicleRecordRequest) (response CreateVehicleRecordResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4assetus.RunAssetusSpaceWorker(ctx, user,
		request.SpaceRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4assetus.AssetusSpaceWorkerParams) (err error) {
			response, err = addVehicleRecordTx(ctx, tx, user, request, params)
			return err
		},
	)
	return
}

// addVehicleRecordTx creates dbo4assetus.VehicleRecordDbo in n /teams/{teamID}/modules/assetus/{assetID}/mileage/{randomRecordID}
func addVehicleRecordTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.UserContext,
	request dto4assetus.AddVehicleRecordRequest,
	params *dal4assetus.AssetusSpaceWorkerParams,
	// params *dal4teamus.ModuleTeamWorkerParams[*dal4assetus.Mileage],
) (
	response CreateVehicleRecordResponse, err error,
) {
	_ = fmt.Sprintf("%v, %v, %v, %v, %v", ctx, tx, user, request, params) // TODO: remove this temp line

	// TODO:
	// 1. Get asset record by ID using tx.Get()
	// 2. Verify asset exists by checking if (dal.IsErrNotFound(err))
	// 3. Create dbo4assetus.VehicleRecordDbo in /teams/{teamID}/modules/assetus/{assetID}/mileage/{randomRecordID} using VehicleRecordDbo tx.Insert()
	// 4. Update asset.extra.Mileages with mileage record ID
	// 4.1 Update asset record
	// 4.2 Validate asset record
	// 4.3 update asset record using tx.Update()

	return response, err
}
