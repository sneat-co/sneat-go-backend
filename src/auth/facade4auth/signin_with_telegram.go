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

	botUserData := BotUserData{
		PlatformID: telegram.PlatformID,
		BotID:      "", // TODO: populate
		BotUserID:  strconv.FormatInt(initData.User.ID, 10),
		FirstName:  initData.User.FirstName,
		LastName:   initData.User.LastName,
		Username:   initData.User.Username,
	}

	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if botUser, err = botsdal.GetPlatformUser(ctx, db, telegram.PlatformID, botUserData.BotUserID, new(models4bots.TelegramUserDbo)); err != nil {
			if !dal.IsNotFound(err) {
				return
			}
			if botUser, appUser, err = CreateBotUserAndAppUserRecords(ctx, tx, botUserData, remoteClientInfo); err != nil {
				return
			}
			isNewUser = true
			return
		}

		if appUserID := botUser.Data.GetAppUserID(); appUserID == "" {

			if err = createAppUserRecordAndUpdateTelegramUserRecord(ctx, tx, botUserData, remoteClientInfo, botUser); err != nil {
				err = fmt.Errorf("failed in createAppUserRecordAndUpdateTelegramUserRecord(): %w", err)
				return
			}
		}
		return
	})
	return
}

func createAppUserRecordAndUpdateTelegramUserRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
	tgBotUser record.DataWithID[string, botsfwmodels.PlatformUserData],
) (err error) {
	var user dbo4userus.UserEntry
	if user, err = createUserFromBotUser(ctx, tx, botUserData, remoteClientInfo); err != nil {
		return
	}
	tgBotUser.Data.SetAppUserID(user.ID)
	tgUserDbo := tgBotUser.Data.(*models4bots.TelegramUserDbo)
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
