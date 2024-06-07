package facade4listus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
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
	briefsAdapter := dal4teamus.NewMapBriefsAdapter(
		func(teamModuleDbo *dbo4listus.ListusTeamDbo) int {
			return len(teamModuleDbo.Lists)
		},
		func(teamModuleDbo *dbo4listus.ListusTeamDbo, id string) ([]dal.Update, error) {
			delete(teamModuleDbo.Lists, id)
			return []dal.Update{{Field: "lists." + id, Value: dal.DeleteField}}, teamModuleDbo.Validate()
		},
	)
	teamItemRequest := dal4teamus.TeamItemRequest{
		TeamRequest: request.TeamRequest,
		ID:          id,
	}
	err = dal4teamus.DeleteTeamItem(
		ctx,
		user,
		teamItemRequest,
		const4listus.ModuleID,
		new(dbo4listus.ListusTeamDbo),
		dbo4listus.ListsCollection,
		new(dbo4listus.ListDbo),
		briefsAdapter,
		deleteListTxWorker,
	)

	return
}

func deleteListTxWorker(_ context.Context, _ dal.ReadwriteTransaction, _ *dal4teamus.TeamItemWorkerParams[*dbo4listus.ListusTeamDbo, *dbo4listus.ListDbo]) (err error) {
	return errors.New("not implemented")
}
