package models

import (
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"reflect"
)

type DebtusTelegramChat struct {
	record.WithID[string]
	//botsfwtgmodels.ChatEntity
	Data *DebtusTelegramChatData
}

var _ botsfwtgmodels.TgChatData = (*DebtusTelegramChatData)(nil)

// DebtusTelegramChatData is a data structure for storing debtus data related to specific telegram chat
type DebtusTelegramChatData struct {
	botsfwtgmodels.TgChatBaseData
	DebtusChatData
}

func NewDebtusTelegramChatRecord() dal.Record {
	return dal.NewRecordWithIncompleteKey(botsfwtgmodels.TgChatCollection, reflect.String, new(DebtusTelegramChatData))
}

func (v *DebtusTelegramChatData) BaseChatData() *botsfwtgmodels.TgChatBaseData {
	return &v.TgChatBaseData
}

func (v *DebtusTelegramChatData) Validate() (err error) {
	//if properties, err = datastore.SaveStruct(entity); err != nil {
	//	return properties, err
	//}
	//if properties, err = entity.TgChatEntityBase.CleanProperties(properties); err != nil {
	//	return
	//}
	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
	//	"GetUserGroupID":   gaedb.IsEmptyString,
	//	"TgChatInstanceID": gaedb.IsEmptyString,
	//}); err != nil {
	//	return
	//}
	return
}
