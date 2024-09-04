package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4spaceus"
)

func createDefaultUserSpacesTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *CreateUserWorkerParams,
) (
	spaces []dbo4spaceus.SpaceEntry,
	contactusSpaces []dal4contactus.ContactusSpaceEntry,
	err error,
) {
	for _, spaceType := range []core4spaceus.SpaceType{core4spaceus.SpaceTypeFamily, core4spaceus.SpaceTypePrivate} {
		if spaceID, _ := params.User.Data.GetFirstSpaceBriefBySpaceType(spaceType); spaceID == "" {
			var result facade4spaceus.CreateSpaceResult
			spaceRequest := dto4spaceus.CreateSpaceRequest{Type: spaceType}
			if result, err = facade4spaceus.CreateSpaceTxWorker(ctx, tx, params.Started, spaceRequest, params.UserWorkerParams); err != nil {
				return
			}
			spaces = append(spaces, result.Space)
			contactusSpaces = append(contactusSpaces, result.ContactusSpace)
			params.RecordsToInsert = append(params.RecordsToInsert, result.Records()...)
		}
	}
	return
}
