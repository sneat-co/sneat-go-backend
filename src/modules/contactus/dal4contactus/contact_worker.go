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

func (v ContactWorkerParams) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	return v.ContactusTeamWorkerParams.GetRecords(ctx, tx, append(records, v.Contact.Record)...)
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
	worker ContactWorker,
) error {
	contactWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, moduleWorkerParams *ContactusTeamWorkerParams) (err error) {
		params := NewContactWorkerParams(moduleWorkerParams, request.ContactID)
		if err = worker(ctx, tx, params); err != nil {
			return err
		}
		if err = applyContactUpdates(ctx, tx, params); err != nil {
			return err
		}
		return err
	}
	return RunContactusTeamWorker(ctx, user, request.TeamRequest, contactWorker)
}

func applyContactUpdates(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error) {
	if len(params.ContactUpdates) > 0 {
		if err = tx.Update(ctx, params.Contact.Record.Key(), params.ContactUpdates); err != nil {
			return err
		}
	}
	return err
}
