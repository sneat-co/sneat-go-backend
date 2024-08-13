package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-core/facade"
)

type TgGroupDalGae struct {
}

func NewTgGroupDalGae() TgGroupDalGae {
	return TgGroupDalGae{}
}

func (TgGroupDalGae) GetTgGroupByID(c context.Context, tx dal.ReadSession, id int64) (tgGroup models4auth.TgGroup, err error) {
	tgGroup = models4auth.NewTgGroup(id, nil)
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	return tgGroup, tx.Get(c, tgGroup.Record)
}

func (TgGroupDalGae) SaveTgGroup(c context.Context, tx dal.ReadwriteTransaction, tgGroup models4auth.TgGroup) (err error) {
	return tx.Set(c, tgGroup.Record)
}
