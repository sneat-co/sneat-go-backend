package models4debtus

import (
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/anybot"
	"reflect"
)

func NewDebtusTelegramChatRecord() dal.Record {
	return dal.NewRecordWithIncompleteKey(botsfwtgmodels.TgChatCollection, reflect.String, new(anybot.SneatAppTgChatDbo))
}
