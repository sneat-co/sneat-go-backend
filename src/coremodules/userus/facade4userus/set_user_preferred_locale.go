package facade4userus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
)

func SetUserPreferredLocale(ctx context.Context, userCtx facade.UserContext, localeCode5 string) (err error) {
	err = dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) (err error) {
		if dal.IsNotFound(err) {
			logus.Errorf(ctx, "User not found by ContactID: %v", err)
			return nil
		}
		user := userWorkerParams.User
		if err == nil && user.Data.PreferredLocale != localeCode5 {
			if userWorkerParams.UserUpdates, err = user.Data.SetPreferredLocale(localeCode5); err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}
	return err
}
