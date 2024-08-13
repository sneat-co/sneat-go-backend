package facade4debtus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"reflect"
	"time"
)

func delayedUpdateSpaceHasDueTransfers(ctx context.Context, userID, spaceID string) (err error) {
	logus.Infof(ctx, "delayedUpdateSpaceHasDueTransfers(userID=%v)", userID)
	userCtx := facade.NewUserContext(userID)
	err = dal4spaceus.RunModuleSpaceWorker(ctx, userCtx, spaceID, const4debtus.ModuleID, new(models4debtus.DebtusSpaceDbo),
		func(c context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*models4debtus.DebtusSpaceDbo]) error {
			if !params.SpaceModuleEntry.Data.HasDueTransfers {
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, params.SpaceModuleEntry.Data.SetHasDueTransfers(true))
				params.SpaceModuleEntry.Record.MarkAsChanged()
			}
			return nil
		})
	return
}

func checkHasDueTransfers(ctx context.Context, db dal.ReadSession, userID, spaceID string) (hasDueTransfer bool, err error) {
	q := dal.From(models4debtus.TransfersCollection).
		WhereField("BothUserIDs", dal.Equal, userID).
		WhereField("IsOutstanding", dal.Equal, true).
		WhereField("DtDueOn", dal.GreaterThen, time.Time{}).
		Limit(1).
		SelectKeysOnly(reflect.Int)

	var reader dal.Reader
	if reader, err = db.QueryReader(ctx, q); err != nil {
		return
	}

	var transferIDs []int
	if transferIDs, err = dal.SelectAllIDs[int](reader, dal.WithLimit(q.Limit())); err != nil {
		return
	}

	return len(transferIDs) > 0, nil
}

func delayedUpdateUserHasDueTransfers(ctx context.Context, userID, spaceID string) (err error) {
	logus.Infof(ctx, "delayerUpdateUserHasDueTransfers(userID=%v)", userID)
	if userID == "" {
		logus.Errorf(ctx, "userID == 0")
		return nil
	}

	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return err
	}

	debtusUser := models4debtus.NewDebtusUserEntry(userID)
	if err = db.Get(ctx, debtusUser.Record); err != nil && !dal.IsNotFound(err) {
		return err
	}

	if debtusUser.Data.HasDueTransfers {
		logus.Infof(ctx, "Already debtusUser.HasDueTransfers == %v", true)
		return nil
	}

	var hasDueTransfers bool
	if hasDueTransfers, err = checkHasDueTransfers(ctx, db, userID, spaceID); err != nil {
		return err
	}

	if !hasDueTransfers {
		logus.Infof(ctx, "No due transfers found")
		return nil
	}

	err = dal4userus.RunUserModuleWorker[models4debtus.DebtusUserDbo](ctx, userID, const4debtus.ModuleID, new(models4debtus.DebtusUserDbo),
		func(tc context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserModuleWorkerParams[models4debtus.DebtusUserDbo]) error {
			if !params.UserModule.Data.HasDueTransfers {
				params.UserModuleUpdates = append(params.UserModuleUpdates, params.UserModule.Data.SetHasDueTransfers(true))
				logus.Infof(ctx, "User updated & saved to datastore")
			}
			return nil
		})
	return err
}
