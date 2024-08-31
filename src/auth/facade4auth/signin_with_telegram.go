package facade4auth

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/bots-go-framework/bots-fw/botsfwconst"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
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

type BotUserData struct {
	PlatformID   botsfwconst.Platform
	BotID        string
	BotUserID    string
	FirstName    string
	LastName     string
	Username     string
	PhotoURL     string
	LanguageCode string
}

func CreateBotUserAndAppUserRecords(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	botUser BotUserEntry, appUser dbo4userus.UserEntry, err error,
) {

	telegramUserData := new(models4bots.TelegramUserDbo)
	telegramUserData.FirstName = botUserData.FirstName
	telegramUserData.LastName = botUserData.LastName
	telegramUserData.UserName = botUserData.Username
	botUser.Data = telegramUserData
	tgPlatformUserKey := botsdal.NewPlatformUserKey(telegram.PlatformID, botUserData.BotUserID)
	botUser = record.NewDataWithID[string, botsfwmodels.PlatformUserData](botUserData.BotUserID, tgPlatformUserKey, telegramUserData)

	botUser.ID = botUserData.BotUserID
	key := botsdal.NewPlatformUserKey(telegram.PlatformID, botUser.ID)
	botUser.Record = dal.NewRecordWithData(key, telegramUserData)

	if appUser, err = createUserFromTelegramUser(ctx, tx, botUserData, remoteClientInfo); err != nil {
		return
	}
	telegramUserData.SetAppUserID(appUser.ID)
	if tx != nil {
		if err = tx.Insert(ctx, botUser.Record); err != nil {
			err = fmt.Errorf("failed to insert telegram user record: %w", err)
			return
		}
	}
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
	if user, err = createUserFromTelegramUser(ctx, tx, botUserData, remoteClientInfo); err != nil {
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

func createUserFromTelegramUser(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	user dbo4userus.UserEntry, err error,
) {
	userToCreate := &UserToCreate{
		Names: person.NameFields{
			FirstName: botUserData.FirstName,
			LastName:  botUserData.LastName,
		},
		PhotoURL:     botUserData.PhotoURL,
		LanguageCode: botUserData.LanguageCode,
		Account: appuser.AccountKey{
			Provider: telegram.PlatformID,
			ID:       botUserData.BotUserID,
		},
		RemoteClientInfo: remoteClientInfo,
	}
	return createUser(ctx, tx, userToCreate)
}
