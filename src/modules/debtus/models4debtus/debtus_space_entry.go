package models4debtus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
)

type DebtusSpaceEntry = record.DataWithID[string, *DebtusSpaceDbo]

func NewDebtusSpaceEntry(spaceID string) DebtusSpaceEntry {
	key := dbo4spaceus.NewSpaceModuleKey(spaceID, const4debtus.ModuleID)
	return record.NewDataWithID(spaceID, key, new(DebtusSpaceDbo))
}

func GetDebtusSpace(ctx context.Context, tx dal.ReadSession, space DebtusSpaceEntry) error {
	return tx.Get(ctx, space.Record)
}
