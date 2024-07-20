package facade4listus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// DeleteList deletes list
func DeleteList(ctx context.Context, user facade.User, request ListRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	uid := user.GetID()
	if uid == "" {
		return validation.NewErrRecordIsMissingRequiredField("user.ContactID()")
	}
	listType := request.ListType()
	id := dbo4listus.GetFullListID(listType, request.ListID)
	briefsAdapter := dal4spaceus.NewMapBriefsAdapter(
		func(teamModuleDbo *dbo4listus.ListusSpaceDbo) int {
			return len(teamModuleDbo.Lists)
		},
		func(teamModuleDbo *dbo4listus.ListusSpaceDbo, id string) ([]dal.Update, error) {
			delete(teamModuleDbo.Lists, id)
			return []dal.Update{{Field: "lists." + id, Value: dal.DeleteField}}, teamModuleDbo.Validate()
		},
	)
	spaceItemRequest := dal4spaceus.SpaceItemRequest{
		SpaceRequest: request.SpaceRequest,
		ID:           id,
	}
	err = dal4spaceus.DeleteSpaceItem(
		ctx,
		user,
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

func deleteListTxWorker(_ context.Context, _ dal.ReadwriteTransaction, _ *dal4spaceus.SpaceItemWorkerParams[*dbo4listus.ListusSpaceDbo, *dbo4listus.ListDbo]) (err error) {
	return errors.New("not implemented")
}
