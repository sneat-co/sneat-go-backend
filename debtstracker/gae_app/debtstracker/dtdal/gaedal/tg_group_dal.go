package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type tgGroupDalGae struct {
}

func newTgGroupDalGae() tgGroupDalGae {
	return tgGroupDalGae{}
}

func (tgGroupDalGae) GetTgGroupByID(c context.Context, tx dal.ReadSession, id int64) (tgGroup models.TgGroup, err error) {
	tgGroup = models.NewTgGroup(id, nil)
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	return tgGroup, tx.Get(c, tgGroup.Record)
}

func (tgGroupDalGae) SaveTgGroup(c context.Context, tx dal.ReadwriteTransaction, tgGroup models.TgGroup) (err error) {
	return tx.Set(c, tgGroup.Record)
}
