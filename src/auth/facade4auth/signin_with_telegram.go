package facade4auth

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/bots-go-framework/bots-fw/botsfw/botsdal"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/bots/models4bots"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"strconv"
)

type TelegramPlatformUserEntry = record.DataWithID[string, botsfwmodels.PlatformUserData]

func SignInWithTelegram(
	ctx context.Context, initData twainitdata.InitData, remoteClientInfo dbmodels.RemoteClientInfo,
) (
	tgPlatformUser TelegramPlatformUserEntry,
	isNewUser bool, // TODO: Document why needed or remove
	err error,
) {
	tgUserID := strconv.FormatInt(initData.User.ID, 10)

	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return
	}

	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if tgPlatformUser, err = botsdal.GetPlatformUser(ctx, db, telegram.PlatformID, tgUserID, new(models4bots.TelegramUserDbo)); err != nil {
			if !dal.IsNotFound(err) {
				return
			}
			if tgPlatformUser, err = createTelegramUserAndAppUserRecords(ctx, tx, initData, remoteClientInfo); err != nil {
				return
			}
			isNewUser = true
			return
		}

		if appUserID := tgPlatformUser.Data.GetAppUserID(); appUserID == "" {

			if err = createAppUserRecordAndUpdateTelegramUserRecord(ctx, tx, initData, remoteClientInfo, tgPlatformUser); err != nil {
				err = fmt.Errorf("failed in createAppUserRecordAndUpdateTelegramUserRecord(): %w", err)
				return
			}
		}
		return
	})
	return
}

func createTelegramUserAndAppUserRecords(
	ctx context.Context, tx dal.ReadwriteTransaction,
	initData twainitdata.InitData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	tgPlatformUser TelegramPlatformUserEntry, err error,
) {
	telegramUserData := new(models4bots.TelegramUserDbo)
	telegramUserData.FirstName = initData.User.FirstName
	telegramUserData.LastName = initData.User.LastName
	telegramUserData.UserName = initData.User.Username
	tgPlatformUser.Data = telegramUserData
	tgUserID := strconv.FormatInt(initData.User.ID, 10)
	tgPlatformUserKey := botsdal.NewPlatformUserKey(telegram.PlatformID, tgUserID)
	tgPlatformUser = record.NewDataWithID[string, botsfwmodels.PlatformUserData](tgUserID, tgPlatformUserKey, telegramUserData)

	tgPlatformUser.ID = strconv.FormatInt(initData.User.ID, 10)
	key := botsdal.NewPlatformUserKey(telegram.PlatformID, tgPlatformUser.ID)
	tgPlatformUser.Record = dal.NewRecordWithData(key, telegramUserData)

	var uid string
	if uid, err = createUserFromTelegramUser(ctx, tx, initData, remoteClientInfo); err != nil {
		return
	}
	telegramUserData.SetAppUserID(uid)
	if err = tx.Insert(ctx, tgPlatformUser.Record); err != nil {
		err = fmt.Errorf("failed to insert telegram user record: %w", err)
		return
	}
	return
}

func createAppUserRecordAndUpdateTelegramUserRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	initData twainitdata.InitData,
	remoteClientInfo dbmodels.RemoteClientInfo,
	tgBotUser record.DataWithID[string, botsfwmodels.PlatformUserData],
) (err error) {
	var uid string
	if uid, err = createUserFromTelegramUser(ctx, tx, initData, remoteClientInfo); err != nil {
		return
	}
	tgBotUser.Data.SetAppUserID(uid)
	tgUserDbo := tgBotUser.Data.(*models4bots.TelegramUserDbo)
	tgUserDbo.FirstName = initData.User.FirstName
	tgUserDbo.LastName = initData.User.LastName
	tgUserDbo.UserName = initData.User.Username

	if err = tgUserDbo.Validate(); err != nil {
		err = fmt.Errorf("failed to validate telegram user data: %w", err)
		return
	}

	updates := []dal.Update{
		{Field: "appUserID", Value: uid},
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

func createUserFromTelegramUser(ctx context.Context, tx dal.ReadwriteTransaction, initData twainitdata.InitData, remoteClientInfo dbmodels.RemoteClientInfo) (uid string, err error) {
	userToCreate := &UserToCreate{
		Names: person.NameFields{
			FirstName: initData.User.FirstName,
			LastName:  initData.User.LastName,
		},
		PhotoURL:     initData.User.PhotoURL,
		LanguageCode: initData.User.LanguageCode,
		Account: appuser.AccountKey{
			Provider: telegram.PlatformID,
			ID:       strconv.FormatInt(initData.User.ID, 10),
		},
		RemoteClientInfo: remoteClientInfo,
	}
	return createUser(ctx, tx, userToCreate)
}
