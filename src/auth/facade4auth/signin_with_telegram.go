package facade4auth

import (
	"context"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"strconv"
)

func SignInWithTelegram(
	ctx context.Context, botID string, initData twainitdata.InitData, remoteClientInfo dbmodels.RemoteClientInfo,
) (
	botUser BotUserEntry,
	appUser dbo4userus.UserEntry,
	isNewUser bool, // TODO: Document why needed or remove
	err error,
) {
	botUserID := strconv.FormatInt(initData.User.ID, 10)

	newBotUserData := func() BotUserData {
		return BotUserData{
			PlatformID: telegram.PlatformID,
			BotID:      botID,
			BotUserID:  botUserID,
			FirstName:  initData.User.FirstName,
			LastName:   initData.User.LastName,
			Username:   initData.User.Username,
		}
	}
	return SignInWithBot(ctx, remoteClientInfo, telegram.PlatformID, botUserID, newBotUserData)
}
