package debtustgbots

import (
	"context"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	dbo4userus2 "github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

func newBotChatData() botsfwmodels.BotChatData {
	return new(botsfwtgmodels.TgChatBaseData)
}

func newBotUserData() botsfwmodels.PlatformUserData {
	return new(botsfwtgmodels.TgPlatformUserBaseDbo)
}

func newAppUserData() botsfwmodels.AppUserData {
	return new(dbo4userus2.UserDbo)
}

func getAppUserByID(_ context.Context, tx dal.ReadSession, botID, appUserID string) (appUser record.DataWithID[string, botsfwmodels.AppUserData], err error) {
	appUserData := newAppUserData()
	key := dbo4userus2.NewUserKey(appUserID)
	appUser = record.NewDataWithID(appUserID, key, appUserData)
	return appUser, nil
}
