package facade4listus

import (
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// DeleteList deletes list
func DeleteList(ctx facade.ContextWithUser, request dto4listus.ListRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	userCtx := ctx.User()
	uid := userCtx.GetUserID()
	if uid == "" {
		return validation.NewErrRecordIsMissingRequiredField("userCtx.ContactID()")
	}
	briefsAdapter := dal4spaceus.NewMapBriefsAdapter(
		func(teamModuleDbo *dbo4listus.ListusSpaceDbo) int {
			return len(teamModuleDbo.Lists)
		},
		func(teamModuleDbo *dbo4listus.ListusSpaceDbo, id string) ([]update.Update, error) {
			delete(teamModuleDbo.Lists, id)
			return []update.Update{update.ByFieldPath([]string{"lists", id}, update.DeleteField)}, teamModuleDbo.Validate()
		},
	)
	spaceItemRequest := dto4spaceus.SpaceItemRequest{
		SpaceRequest: request.SpaceRequest,
		ID:           string(request.ListID),
	}
	err = dal4spaceus.DeleteSpaceItem(
		ctx,
		spaceItemRequest,
		const4listus.ModuleID,
		new(dbo4listus.ListusSpaceDbo),
		dbo4listus.ListsCollection,
		new(dbo4listus.ListDbo),
		briefsAdapter,
		deleteListTxWorker,
	)

	return
}

func deleteListTxWorker(_ facade.ContextWithUser, _ dal.ReadwriteTransaction, _ *dal4spaceus.SpaceItemWorkerParams[*dbo4listus.ListusSpaceDbo, *dbo4listus.ListDbo]) (err error) {
	return errors.New("not implemented")
}
