package dal4bots

import (
	"context"
	"github.com/bots-go-framework/bots-fw/botsfw/botsdal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/botscore/models4bots"
)

func GetBotUserRecord(ctx context.Context, botPlatformID, botID, userID string) (
	tgBotUser record.DataWithID[string, *models4bots.TelegramUserDbo], err error,
) {
	key := botsdal.NewPlatformUserKey("botUsers", userID)
	tgBotUser = record.NewDataWithID(userID, key, new(models4bots.TelegramUserDbo))
	return
}
