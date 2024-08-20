package debtustgbots

import (
	"context"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
)

func newBotChatData() botsfwmodels.BotChatData {
	return nil
}

func newBotUserData() botsfwmodels.BotUserData {
	return nil
}
func newAppUserData() botsfwmodels.AppUserData {
	return nil
}
func getAppUserByID(c context.Context, tx dal.ReadSession, botID, appUserID string) (appUser record.DataWithID[string, botsfwmodels.AppUserData], err error) {
	appUserData := newAppUserData()
	key := dbo4userus.NewUserKey(appUserID)
	appUser = record.NewDataWithID(appUserID, key, appUserData)
	return appUser, nil
}
