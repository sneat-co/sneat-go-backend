package facade4auth

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"strconv"
	"time"
)

func SignInWithTelegram(
	ctx context.Context, initData twainitdata.InitData, remoteClientInfo dbmodels.RemoteClientInfo,
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

	botUserID := strconv.FormatInt(initData.User.ID, 10)

	newBotUserData := func() BotUserData {
		return BotUserData{
			PlatformID: telegram.PlatformID,
			BotID:      "", // TODO: populate
			BotUserID:  botUserID,
			FirstName:  initData.User.FirstName,
			LastName:   initData.User.LastName,
			Username:   initData.User.Username,
		}
	}

	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if botUser, appUser, isNewUser, err = createBotUserAndAppUserRecordsTx(ctx, tx, botUserID, newBotUserData, remoteClientInfo); err != nil {
			return
		}
		return
	})
	return
}

func createBotUserAndAppUserRecordsTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserID string,
	newBotUserData func() BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	botUser BotUserEntry,
	appUser dbo4userus.UserEntry,
	isNewUser bool, // TODO: Document why needed or remove
	err error,
) {
	if botUser, err = botsdal.GetPlatformUser(ctx, tx, telegram.PlatformID, botUserID, new(models4bots.TelegramUserDbo)); err != nil {
		if dal.IsNotFound(err) {
			botUserData := newBotUserData()
			if botUser, appUser, err = CreateBotUserAndAppUserRecords(ctx, tx, botUserData, remoteClientInfo); err != nil {
				return
			}
			isNewUser = true
			return
		}
		return
	}

	if appUserID := botUser.Data.GetAppUserID(); appUserID == "" {
		botUserData := newBotUserData()
		if err = createAppUserRecordAndUpdateBotUserRecord(ctx, tx, botUserData, remoteClientInfo, botUser); err != nil {
			err = fmt.Errorf("failed in createAppUserRecordAndUpdateBotUserRecord(): %w", err)
			return
		}
	}
	return
}

func createAppUserRecordAndUpdateBotUserRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
	tgBotUser record.DataWithID[string, botsfwmodels.PlatformUserData],
) (err error) {
	started := time.Now()
	var user dbo4userus.UserEntry
	if user, err = getOrCreateAppUserRecordFromBotUser(ctx, tx, started, botUserData, remoteClientInfo); err != nil {
		return
	}
	tgBotUser.Data.SetAppUserID(user.ID)
	tgUserDbo := tgBotUser.Data.(*models4bots.TelegramUserDbo)
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
		{Field: "appUserID", Value: user.ID},
		{Field: "firstName", Value: tgUserDbo.FirstName},
		{Field: "lastName", Value: tgUserDbo.LastName},
		{Field: "userName", Value: tgUserDbo.UserName},
	}

	if err = tx.Update(ctx, tgBotUser.Key, updates); err != nil {
		err = fmt.Errorf("failed to update telegram user record: %w", err)
		return
	}
	return
}
