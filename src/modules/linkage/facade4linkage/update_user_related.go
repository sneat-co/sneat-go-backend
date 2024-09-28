package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
)

func updateUserRelated(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	objectRef dbo4linkage.SpaceModuleItemRef,
	item record.DataWithID[string, *dbo4linkage.WithRelated],
) (userUpdates record.Updates, err error) {
	user := dbo4userus.NewUserEntry(userID)
	if err = tx.Get(ctx, user.Record); err != nil {
		return
	}

	return
}
