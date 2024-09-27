package facade4auth

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/bots-go-framework/bots-fw/botsfwconst"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/dto4auth"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"time"
)

func CreateBotUserAndAppUserRecords(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	appUserID string,
	botPlatformID botsfwconst.Platform,
	botUserData BotUserData,
	remoteClientInfo dbmodels.RemoteClientInfo,
) (
	botUser BotUserEntry, params CreateUserWorkerParams, err error,
) {
	started := time.Now()

	telegramUserData := new(models4bots.TelegramUserDbo)
	telegramUserData.AccessGranted = true // TODO: Do we really need this field?
	telegramUserData.DtCreated = started
	telegramUserData.DtUpdated = started
	telegramUserData.FirstName = botUserData.FirstName
	telegramUserData.LastName = botUserData.LastName
	telegramUserData.UserName = botUserData.Username
	botUser.Data = telegramUserData
	tgPlatformUserKey := botsdal.NewPlatformUserKey(telegram.PlatformID, botUserData.BotUserID)
	botUser = record.NewDataWithID[string, botsfwmodels.PlatformUserData](botUserData.BotUserID, tgPlatformUserKey, telegramUserData)

	botUser.ID = botUserData.BotUserID
	key := botsdal.NewPlatformUserKey(telegram.PlatformID, botUser.ID)
	botUser.Record = dal.NewRecordWithData(key, telegramUserData)

	if appUserID == "" {
		if params, err = getOrCreateAppUserRecordFromBotUser(ctx, tx, started, botUserData, remoteClientInfo); err != nil {
			err = fmt.Errorf("failed in getOrCreateAppUserRecordFromBotUser(): %w", err)
			return
		}
		_ = params.User.Data.AccountsOfUser.AddAccount(appuser.AccountKey{Provider: string(botPlatformID), ID: botUserData.BotUserID})
		botUser.Record.SetError(nil)
		params.QueueForInsert(botUser.Record)
	} else { // appUserID != ""
		params.UserWorkerParams = &dal4userus.UserWorkerParams{
			Started: started,
			User:    dbo4userus.NewUserEntry(appUserID),
		}
		if err = tx.Get(ctx, params.User.Record); err != nil {
			return
		}
		if updates := params.User.Data.AccountsOfUser.AddAccount(appuser.AccountKey{Provider: string(botPlatformID), ID: botUserData.BotUserID}); len(updates) > 0 {
			params.UserUpdates = append(params.UserUpdates, updates...)
			params.User.Record.MarkAsChanged()
		}
	}
	telegramUserData.SetAppUserID(params.User.ID)

	if err = telegramUserData.Validate(); err != nil {
		err = fmt.Errorf("newly created telegram user data is not valid: %w", err)
		botUser.Record.SetError(err)
		return
	}
	botUser.Record.SetError(dal.ErrRecordNotFound)
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

	userToCreate := dto4auth.DataToCreateUser{
		Names: person.NameFields{
			FirstName: botUserData.FirstName,
			LastName:  botUserData.LastName,
		},
		PhotoURL:     botUserData.PhotoURL,
		LanguageCode: botUserData.LanguageCode,
		AuthAccount: appuser.AccountKey{
			Provider: string(botUserData.PlatformID),
			ID:       botUserData.BotUserID,
		},
		RemoteClient: remoteClientInfo,
	}

	var uid string
	if uid, err = CreateUserInAuth(ctx, userToCreate); err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err2 := DeleteUserFromAuth(ctx, uid); err2 != nil {
				err = fmt.Errorf("failed to delete newly created firebase user: %v: ORIGINAL ERROR: %w", err2, err)
			}
		}
	}()

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
			User:    dbo4userus.NewUserEntry(uid),
		},
	}
	if err = tx.Get(ctx, params.User.Record); err != nil && !dal.IsNotFound(err) {
		return
	}
	if err = createUserRecordsTxWorker(ctx, tx, userInfo, userToCreate, &params); err != nil {
		return
	}
	return
}

var CreateUserInAuth func(ctx context.Context, userToCreate dto4auth.DataToCreateUser) (uid string, err error)
var DeleteUserFromAuth func(ctx context.Context, uid string) (err error)
