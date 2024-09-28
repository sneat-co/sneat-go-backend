package facade4auth

import (
	"context"
	"errors"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
	"reflect"
)

type TgUserDalGae struct {
}

func NewTgUserDalGae() TgUserDalGae {
	return TgUserDalGae{}
}

func (TgUserDalGae) FindByUserName(ctx context.Context, tx dal.ReadSession, userName string) (tgUsers []botsfwtgmodels.TgPlatformUser, err error) {
	if tx == nil {
		tx, err = facade.GetSneatDB(ctx)
		if err != nil {
			return
		}
	}
	q := dal.From(botsfwtgmodels.BotUserCollection).
		WhereField("UserName", dal.Equal, userName)

	query := q.SelectInto(func() dal.Record {
		return dal.NewRecordWithIncompleteKey(botsfwtgmodels.BotUserCollection, reflect.Int, new(botsfwtgmodels.TgPlatformUser))
	})
	var records []dal.Record

	if records, err = tx.QueryAllRecords(ctx, query); err != nil {
		return
	}
	tgUsers = make([]botsfwtgmodels.TgPlatformUser, len(records))
	//for i, r := range records {
	//	tgUsers[i] = botsfwtgmodels.TgBotUserBaseData{
	//		WithID: record.NewWithID(r.Key().ContactID.(int64), r.Key(), r.Data),
	//		Data:   r.Data().(*botsfwtgmodels.TgBotUserData),
	//	}
	//}
	return tgUsers, errors.New("not implemented")
}
