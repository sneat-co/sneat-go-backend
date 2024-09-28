package facade4auth

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/appuser"
)

func Disconnect(ctx context.Context, userCtx facade.UserContext, provider string) (err error) {
	disconnectTx := func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
		var ua *appuser.AccountKey
		if ua, err = params.User.Data.AccountsOfUser.GetAccount(provider, ""); err != nil {
			return
		}
		if params.User.Data.AccountsOfUser.RemoveAccount(*ua) {
			params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "accounts", Value: params.User.Data.Accounts})
			params.User.Record.MarkAsChanged()
		}
		botPlatformKey := dal.NewKeyWithID("botPlatforms", provider)
		botUserKey := dal.NewKeyWithParentAndID(botPlatformKey, "botUsers", ua.ID)
		botUser := record.NewDataWithID(ua.ID, botUserKey, new(models4bots.TelegramUserDbo))
		if err = tx.Get(ctx, botUser.Record); err != nil {
			return
		}
		botUser.Data.SetAppUserID("")
		if err = botUser.Data.Validate(); err != nil {
			return fmt.Errorf("bot user data is invalid after cleaning appUserID: %w", err)
		}
		if err = tx.Update(ctx, botUserKey, []dal.Update{{Field: "appUserID", Value: ""}}); err != nil {
			return
		}
		return
	}
	return dal4userus.RunUserWorker(ctx, userCtx, true, disconnectTx)
}
