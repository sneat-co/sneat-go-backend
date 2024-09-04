package facade4auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"time"
)

func CreateBotUserAndAppUserRecords(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	botUser BotUserEntry, params CreateUserWorkerParams, err error,
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

	started := time.Now()
	if params, err = getOrCreateAppUserRecordFromBotUser(ctx, tx, started, botUserData, remoteClientInfo); err != nil {
		return
	}
	telegramUserData.SetAppUserID(params.User.ID)
	if tx != nil {
		if err = tx.Insert(ctx, botUser.Record); err != nil {
			err = fmt.Errorf("failed to insert telegram user record: %w", err)
			return
		}
	}
	return
}

func getOrCreateAppUserRecordFromBotUser(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	started time.Time,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	params CreateUserWorkerParams, err error,
) {

	userToCreate := DataToCreateUser{
		AuthProvider: string(botUserData.PlatformID),
		Names: person.NameFields{
			FirstName: botUserData.FirstName,
			LastName:  botUserData.LastName,
		},
		PhotoURL:     botUserData.PhotoURL,
		LanguageCode: botUserData.LanguageCode,
		Account: appuser.AccountKey{
			Provider: string(botUserData.PlatformID),
			ID:       botUserData.BotUserID,
		},
		RemoteClient: remoteClientInfo,
	}

	var firebaseUserRecord *auth.UserRecord
	if firebaseUserRecord, err = createFirebaseUser(ctx, userToCreate); err != nil {
		return
	}

	providerUserInfo := &sneatauth.AuthProviderUserInfo{
		DisplayName: userToCreate.Names.GetFullName(),
		UID:         botUserData.BotUserID,
		ProviderID:  string(botUserData.PlatformID),
		PhotoURL:    userToCreate.PhotoURL,
	}
	userInfo := &sneatauth.AuthUserInfo{ // TODO: Does this duplicate userToCreate?
		AuthProviderUserInfo: providerUserInfo,
		ProviderUserInfo:     []*sneatauth.AuthProviderUserInfo{providerUserInfo}, // TODO: Why it duplicates AuthProviderUserInfo?
	}
	params = CreateUserWorkerParams{
		UserWorkerParams: &dal4userus.UserWorkerParams{
			Started: started,
			User:    dbo4userus.NewUserEntry(firebaseUserRecord.UID),
		},
	}
	if err = tx.Get(ctx, params.User.Record); err != nil && !dal.IsNotFound(err) {
		return
	}
	if err = createUserRecordsTxWorker(ctx, tx, userInfo, userToCreate, &params); err != nil {
		return
	}
	if err = params.ApplyChanges(ctx, tx); err != nil {
		err = fmt.Errorf("failed to apply changes generate by createUserRecordsTxWorker(): %w", err)
		return
	}
	return
}
