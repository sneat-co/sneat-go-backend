package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactusSpaceWorkerParams = dal4teamus.ModuleSpaceWorkerParams[*models4contactus.ContactusSpaceDbo]

func NewContactusSpaceWorkerParams(userID, spaceID string) *ContactusSpaceWorkerParams {
	teamWorkerParams := dal4teamus.NewSpaceWorkerParams(userID, spaceID)
	return dal4teamus.NewSpaceModuleWorkerParams(const4contactus.ModuleID, teamWorkerParams, new(models4contactus.ContactusSpaceDbo))
}

func RunReadonlyContactusSpaceWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *ContactusSpaceWorkerParams) (err error),
) error {
	return dal4teamus.RunReadonlyModuleSpaceWorker(ctx, user, request, const4contactus.ModuleID, new(models4contactus.ContactusSpaceDbo), worker)
}

type ContactusModuleWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactusSpaceWorkerParams) (err error)

func RunContactusSpaceWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.SpaceRequest,
	worker ContactusModuleWorker,
) error {
	return dal4teamus.RunModuleSpaceWorker(ctx, user, request, const4contactus.ModuleID, new(models4contactus.ContactusSpaceDbo), worker)
}

func RunContactusSpaceWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dto4teamus.SpaceRequest,
	worker ContactusModuleWorker,
) error {
	return dal4teamus.RunModuleSpaceWorkerTx(ctx, tx, user, request, const4contactus.ModuleID, new(models4contactus.ContactusSpaceDbo), worker)
}
