package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactWorkerParams struct {
	*dal4contactus.ContactusTeamWorkerParams
	Contact dal4contactus.ContactEntry
}

type ContactWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error)

func RunContactWorker(ctx context.Context, user facade.User, request dto4contactus.ContactRequest, worker ContactWorker) (err error) {
	return dal4contactus.RunContactusTeamWorker(ctx, user, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, contactusTeamParams *dal4contactus.ContactusTeamWorkerParams) (err error) {
			params := &ContactWorkerParams{
				ContactusTeamWorkerParams: contactusTeamParams,
				Contact:                   dal4contactus.NewContactEntry(contactusTeamParams.Team.ID, request.ContactID),
			}
			return worker(ctx, tx, params)
		},
	)
}
