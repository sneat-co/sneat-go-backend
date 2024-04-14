package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactusTeamWorkerParams = dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDbo]

func NewContactusTeamWorkerParams(userID, teamID string) *ContactusTeamWorkerParams {
	teamWorkerParams := dal4teamus.NewTeamWorkerParams(userID, teamID)
	return dal4teamus.NewTeamModuleWorkerParams(const4contactus.ModuleID, teamWorkerParams, new(models4contactus.ContactusTeamDbo))
}

func RunReadonlyContactusTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.TeamRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *ContactusTeamWorkerParams) (err error),
) error {
	return dal4teamus.RunReadonlyModuleTeamWorker(ctx, user, request, const4contactus.ModuleID, new(models4contactus.ContactusTeamDbo), worker)
}

func RunContactusTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.TeamRequest,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactusTeamWorkerParams) (err error),
) error {
	return dal4teamus.RunModuleTeamWorker(ctx, user, request, const4contactus.ModuleID, new(models4contactus.ContactusTeamDbo), worker)
}

func RunContactusTeamWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dto4teamus.TeamRequest,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactusTeamWorkerParams) (err error),
) error {
	return dal4teamus.RunModuleTeamWorkerTx(ctx, tx, user, request, const4contactus.ModuleID, new(models4contactus.ContactusTeamDbo), worker)
}
