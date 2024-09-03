package facade4auth

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
)

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

	if appUser, err = createUserFromBotUser(ctx, tx, botUserData, remoteClientInfo); err != nil {
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

func createUserFromBotUser(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	user dbo4userus.UserEntry, err error,
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
	return createUserFromBot(ctx, tx, userToCreate)
}
