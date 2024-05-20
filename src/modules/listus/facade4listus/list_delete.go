package facade4listus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/models4listus"
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
	id := models4listus.GetFullListID(listType, request.ListID)
	key := dal4listus.NewTeamListKey(request.TeamID, id)
	input := dal4teamus.TeamItemRunnerInput[*models4listus.ListusTeamDto]{
		Counter:       "lists",
		TeamItem:      dal.NewRecord(key),
		BriefsAdapter: briefsAdapter(listType, request.ListID),
	}
	err = dal4teamus.DeleteTeamItem(ctx, user, input, const4listus.ModuleID, new(models4listus.ListusTeamDto), func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.TeamItemWorkerParams) (err error) {
		return errors.New("not implemented")
	})
	return
}
