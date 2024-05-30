package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
)

func updateUserRelated(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	objectRef models4linkage.TeamModuleItemRef,
	item record.DataWithID[string, *models4linkage.WithRelated],
) (userUpdates dal4teamus.RecordUpdates, err error) {
	user := models4userus.NewUser(userID)
	if err = tx.Get(ctx, user.Record); err != nil {
		return
	}

	return
}
