package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactWorkerParams struct {
	*ContactusTeamWorkerParams
	Contact        ContactEntry
	ContactUpdates []dal.Update
}

func NewContactWorkerParams(moduleParams *ContactusTeamWorkerParams, contactID string) *ContactWorkerParams {
	return &ContactWorkerParams{
		ContactusTeamWorkerParams: moduleParams,
		Contact:                   NewContactEntry(moduleParams.Team.ID, contactID),
	}
}

type ContactWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error)

func RunContactWorker(
	ctx context.Context,
	user facade.User,
	request dto4contactus.ContactRequest,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error),
) error {
	moduleWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, moduleWorkerParams *ContactusTeamWorkerParams) (err error) {
		params := NewContactWorkerParams(moduleWorkerParams, request.ContactID)
		return worker(ctx, tx, params)
	}
	return RunContactusTeamWorker(ctx, user, request.TeamRequest, moduleWorker)
}

func RunContactWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dto4contactus.ContactRequest,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error),
) error {
	moduleWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, moduleWorkerParams *ContactusTeamWorkerParams) (err error) {
		params := NewContactWorkerParams(moduleWorkerParams, request.ContactID)
		return worker(ctx, tx, params)
	}
	return RunContactusTeamWorkerTx(ctx, tx, user, request.TeamRequest, moduleWorker)
}
