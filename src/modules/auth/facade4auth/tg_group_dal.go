package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/models4auth"
	"github.com/sneat-co/sneat-go-core/facade"
)

type TgGroupDalGae struct {
}

func NewTgGroupDalGae() TgGroupDalGae {
	return TgGroupDalGae{}
}

func (TgGroupDalGae) GetTgGroupByID(ctx context.Context, tx dal.ReadSession, id int64) (tgGroup models4auth.TgGroup, err error) {
	tgGroup = models4auth.NewTgGroup(id, nil)
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
	}
	return tgGroup, tx.Get(ctx, tgGroup.Record)
}

func (TgGroupDalGae) SaveTgGroup(ctx context.Context, tx dal.ReadwriteTransaction, tgGroup models4auth.TgGroup) (err error) {
	return tx.Set(ctx, tgGroup.Record)
}
