package facade4auth

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/bots-go-framework/bots-fw/botsfwconst"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/coretodo"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"time"
)

func SignInWithBot(
	ctx context.Context,
	remoteClientInfo dbmodels.RemoteClientInfo,
	botPlatformID botsfwconst.Platform,
	botUserID string,
	newBotUserData func() BotUserData,
) (
	botUser BotUserEntry,
	appUser dbo4userus.UserEntry,
	isNewUser bool, // TODO: Document why needed or remove
	err error,
) {
	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		var params CreateUserWorkerParams
		if botUser, params, isNewUser, err = createBotUserAndAppUserRecordsTx(ctx, tx, botPlatformID, botUserID, newBotUserData, remoteClientInfo); err != nil {
			return
		}
		now := time.Now() // TODO: Should be in sync with one used in createBotUserAndAppUserRecordsTx()
		if params.User.Record.Exists() {
			params.UserUpdates = append(params.UserUpdates, params.User.Data.SetLastLoginAt(now))
			params.RecordsToUpdate = append(params.RecordsToUpdate, coretodo.RecordUpdates{
				Record:  params.User.Record,
				Updates: params.UserUpdates,
			})
		} else {
			params.QueueForInsert(params.User.Record)
		}
		if botUser.Record.Exists() {
			params.RecordsToUpdate = append(params.RecordsToUpdate, coretodo.RecordUpdates{
				Record: botUser.Record,
				Updates: []dal.Update{ // TODO: Should be populated inside of createBotUserAndAppUserRecordsTx()?
					{Field: "appUserID", Value: params.User.ID},
					{Field: "dtUpdated", Value: now},
				},
			})
		} else {
			params.QueueForInsert(botUser.Record)
		}
		if err = params.ApplyChanges(ctx, tx); err != nil {
			err = fmt.Errorf("failed to apply changes returned by createBotUserAndAppUserRecordsTx(): %w", err)
		}
		return err
	})
	return

}

func createBotUserAndAppUserRecordsTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botPlatformID botsfwconst.Platform,
	botUserID string,
	newBotUserData func() BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	botUser BotUserEntry,
	params CreateUserWorkerParams,
	isNewUser bool, // TODO: Document why needed or remove
	err error,
) {
	var appUserID string
	authToken := sneatauth.AuthTokenFromContext(ctx)
	if authToken != nil {
		appUserID = authToken.UID
	}
	if botUser, err = botsdal.GetPlatformUser(ctx, tx, telegram.PlatformID, botUserID, new(models4bots.TelegramUserDbo)); err != nil {
		if dal.IsNotFound(err) {
			botUserData := newBotUserData()
			if botUser, params, err = CreateBotUserAndAppUserRecords(ctx, tx, appUserID, botPlatformID, botUserData, remoteClientInfo); err != nil {
				return
			}
			isNewUser = true
			return
		}
		return
	} else if botAppUserID := botUser.Data.GetAppUserID(); botAppUserID == "" {
		botUserData := newBotUserData()
		if appUserID == "" {
			if params, err = createAppUserRecordsAndUpdateBotUserRecord(ctx, tx, botUserData, remoteClientInfo, botUser); err != nil {
				err = fmt.Errorf("failed in createAppUserRecordsAndUpdateBotUserRecord(): %w", err)
				return
			}
		} else {
			botUser.Data.SetAppUserID(appUserID)
		}
	} else if appUserID != "" && botAppUserID != appUserID {
		err = fmt.Errorf("bot user is already linked to another app user")
		return
	}
	return
}

func createAppUserRecordsAndUpdateBotUserRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
	botUser record.DataWithID[string, botsfwmodels.PlatformUserData],
) (params CreateUserWorkerParams, err error) {
	started := time.Now()
	if params, err = getOrCreateAppUserRecordFromBotUser(ctx, tx, started, botUserData, remoteClientInfo); err != nil {
		return
	}
	botUser.Data.SetAppUserID(params.User.ID)
	tgUserDbo := botUser.Data.(*models4bots.TelegramUserDbo)
	tgUserDbo.DtCreated = started
	tgUserDbo.DtUpdated = started
	tgUserDbo.FirstName = botUserData.FirstName
	tgUserDbo.LastName = botUserData.LastName
	tgUserDbo.UserName = botUserData.Username

	if err = tgUserDbo.Validate(); err != nil {
		err = fmt.Errorf("failed to validate telegram user data: %w", err)
		return
	}

	updates := []dal.Update{
		{Field: "appUserID", Value: params.User.ID},
		{Field: "firstName", Value: tgUserDbo.FirstName},
		{Field: "lastName", Value: tgUserDbo.LastName},
		{Field: "userName", Value: tgUserDbo.UserName},
	}

	if err = tx.Update(ctx, botUser.Key, updates); err != nil {
		err = fmt.Errorf("failed to update telegram user record: %w", err)
		return
	}
	return
}
