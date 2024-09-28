package facade4auth

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/record"
)

type BotUserEntry = record.DataWithID[string, botsfwmodels.PlatformUserData]
