package gaedal

import (
	"bytes"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"time"

	"context"
	"errors"
)

type TransferDalGae struct {
}

func NewTransferDalGae() TransferDalGae {
	return TransferDalGae{}
}

var _ dtdal.TransferDal = (*TransferDalGae)(nil)

func _loadDueOnTransfers(ctx context.Context, tx dal.ReadSession, userID string, limit int, filter func(q dal.QueryBuilder) dal.QueryBuilder) (transfers []models4debtus.TransferEntry, err error) {
	q := dal.From(models4debtus.TransfersCollection).
		WhereField("BothUserIDs", "=", userID).
		WhereField("IsOutstanding", "=", true).OrderBy(dal.AscendingField("DtDueOn"))
	q = filter(q).Limit(limit)
	query := q.SelectInto(models4debtus.NewTransferRecord)
	var (
		transferRecords []dal.Record
	)

	transferRecords, err = tx.QueryAllRecords(ctx, query)

	transfers = make([]models4debtus.TransferEntry, len(transferRecords))
	for i, transferRecord := range transferRecords {
		transfer := models4debtus.NewTransfer(transferRecord.Key().ID.(string), transferRecord.Data().(*models4debtus.TransferData))
		transfers[i] = transfer
	}
	return
}

func (transferDalGae TransferDalGae) LoadOverdueTransfers(ctx context.Context, tx dal.ReadSession, userID string, limit int) ([]models4debtus.TransferEntry, error) {
	return _loadDueOnTransfers(ctx, tx, userID, limit, func(q dal.QueryBuilder) dal.QueryBuilder {
		return q.WhereField("DtDueOn", dal.GreaterThen, time.Time{}).WhereField("DtDueOn", dal.LessThen, time.Now())
	})
}

func (transferDalGae TransferDalGae) LoadDueTransfers(ctx context.Context, tx dal.ReadSession, userID string, limit int) ([]models4debtus.TransferEntry, error) {
	return _loadDueOnTransfers(ctx, tx, userID, limit, func(q dal.QueryBuilder) dal.QueryBuilder {
		return q.WhereField("DtDueOn", dal.GreaterThen, time.Now())
	})
}

func (transferDalGae TransferDalGae) GetTransfersByID(ctx context.Context, tx dal.ReadSession, transferIDs []string) (transfers []models4debtus.TransferEntry, err error) {
	transfers = make([]models4debtus.TransferEntry, len(transferIDs))
	records := make([]dal.Record, len(transferIDs))
	for i, transferID := range transferIDs {
		transfers[i] = models4debtus.NewTransfer(transferID, nil)
		records[i] = transfers[i].Record
	}
	if err = tx.GetMulti(ctx, records); err != nil {
		return
	}
	return
}

func (transferDalGae TransferDalGae) LoadOutstandingTransfers(ctx context.Context, tx dal.ReadSession, periodEnds time.Time, userID, contactID string, currency money.CurrencyCode, direction models4debtus.TransferDirection) (transfers []models4debtus.TransferEntry, err error) {
	logus.Debugf(ctx, "TransferDalGae.LoadOutstandingTransfers(periodEnds=%v, userID=%v, contactID=%v currency=%v, direction=%v)", periodEnds, userID, contactID, currency, direction)
	const limit = 100

	// TODO: Load outstanding transfer just for the specific contact & specific direction
	q := dal.From(models4debtus.TransfersCollection).
		Where(
			dal.WhereField("BothUserIDs", dal.Equal, userID),
			dal.WhereField("Currency", dal.Equal, string(currency)),
			dal.WhereField("IsOutstanding", dal.Equal, true),
		).
		OrderBy(dal.AscendingField("DtCreated")).
		Limit(limit).
		SelectInto(models4debtus.NewTransferRecord)
	var transferRecords []dal.Record
	transferRecords, err = tx.QueryAllRecords(ctx, q)
	transfers = models4debtus.TransfersFromRecords(transferRecords)
	var errorMessages, warnings, debugs bytes.Buffer
	var transfersIDsToFixIsOutstanding []string
	for _, transfer := range transfers {
		if contactID != "" {
			if cpContactID := transfer.Data.CounterpartyInfoByUserID(userID).ContactID; cpContactID != contactID {
				debugs.WriteString(fmt.Sprintf("Skipped outstanding TransferEntry(id=%v) as counterpartyContactID != contactID: %v != %v\n", transfer.ID, cpContactID, contactID))
				continue
			}
		}
		if direction != "" {
			if d := transfer.Data.DirectionForUser(userID); d != direction {
				debugs.WriteString(fmt.Sprintf("Skipped outstanding TransferEntry(id=%v) as DirectionForUser(): %v\n", transfer.ID, d))
				continue
			}
		}

		if outstandingValue := transfer.Data.GetOutstandingValue(periodEnds); outstandingValue > 0 {
			transfers = append(transfers, transfer)
		} else if outstandingValue == 0 {
			_, _ = fmt.Fprintf(&warnings, "TransferEntry(id=%v) => GetOutstandingValue() == 0 && IsOutstanding==true\n", transfer.ID)
			transfersIDsToFixIsOutstanding = append(transfersIDsToFixIsOutstanding, transfer.ID)
		} else { // outstandingValue < 0
			_, _ = fmt.Fprintf(&errorMessages, "TransferEntry(id=%v) => IsOutstanding==true && GetOutstandingValue() < 0: %v\n", transfer.ID, outstandingValue)
		}
	}
	if len(transfersIDsToFixIsOutstanding) > 0 {
		if err = delayerFixTransfersIsOutstanding.EnqueueWork(ctx, delaying.With(const4debtus.QueueTransfers, "fix-api4transfers-is-outstanding", 0), transfersIDsToFixIsOutstanding); err != nil {
			logus.Errorf(ctx, "failed to delay task to fix api4transfers IsOutstanding")
			err = nil
		}
	}
	if errorMessages.Len() > 0 {
		logus.Errorf(ctx, errorMessages.String())
	}
	if warnings.Len() > 0 {
		logus.Warningf(ctx, warnings.String())
	}
	if debugs.Len() > 0 {
		logus.Debugf(ctx, debugs.String())
	}
	return
}

func delayedFixTransfersIsOutstanding(ctx context.Context, transferIDs []string) (err error) {
	logus.Debugf(ctx, "delayedFixTransfersIsOutstanding(%v)", transferIDs)
	for _, transferID := range transferIDs {
		if _, transferErr := fixTransferIsOutstanding(ctx, transferID); transferErr != nil {
			logus.Errorf(ctx, "Failed to fix transfer %v: %v", transferID, err)
			err = transferErr
		}
	}
	return
}

func fixTransferIsOutstanding(ctx context.Context, transferID string) (transfer models4debtus.TransferEntry, err error) {
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if transfer, err = facade4debtus.Transfers.GetTransferByID(ctx, tx, transferID); err != nil {
			return err
		}
		if transfer.Data.GetOutstandingValue(time.Now()) == 0 {
			transfer.Data.IsOutstanding = false
			return facade4debtus.Transfers.SaveTransfer(ctx, tx, transfer)
		}
		return nil
	})
	if err == nil {
		logus.Warningf(ctx, "Fixed IsOutstanding (set to false) for transfer %v", transferID)
	} else {
		logus.Errorf(ctx, "Failed to fix IsOutstanding for transfer %v", transferID)
	}
	return
}

func (transferDalGae TransferDalGae) LoadTransfersByUserID(ctx context.Context, userID string, offset, limit int) (transfers []models4debtus.TransferEntry, hasMore bool, err error) {
	if limit == 0 {
		err = errors.New("limit == 0")
		return
	}
	if userID == "" {
		err = errors.New("userID == 0")
		return
	}
	q := dal.From(models4debtus.TransfersCollection).
		WhereField("BothUserIDs", dal.Equal, userID).
		OrderBy(dal.DescendingField("DtCreated")).
		SelectInto(models4debtus.NewTransferRecord)

	if transfers, err = transferDalGae.loadTransfers(ctx, q); err != nil {
		return
	}
	hasMore = len(transfers) > limit
	return
}

func (transferDalGae TransferDalGae) LoadTransferIDsByContactID(ctx context.Context, contactID string, limit int, startCursor string) (transferIDs []string, endCursor string, err error) {
	if limit == 0 {
		err = errors.New("LoadTransferIDsByContactID(): limit == 0")
		return
	} else if limit > 1000 {
		err = errors.New("LoadTransferIDsByContactID(): limit > 1000")
		return
	}
	if contactID == "" {
		err = errors.New("LoadTransferIDsByContactID(): contactID == 0")
		return
	}
	q := dal.From(models4debtus.TransfersCollection).
		WhereField("BothCounterpartyIDs", dal.Equal, contactID).
		Limit(limit).
		StartFrom(dal.Cursor(startCursor)).
		SelectInto(models4debtus.NewTransferRecord)

	//if startCursor != "" {
	//	var decodedCursor datastore.Cursor
	//	if decodedCursor, err = datastore.DecodeCursor(startCursor); err != nil {
	//		return
	//	} else {
	//		q = q.Start(decodedCursor)
	//	}
	//}

	transferIDs = make([]string, 0, limit)
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	var reader dal.Reader
	if reader, err = db.QueryReader(ctx, q); err != nil {
		return
	}
	var record dal.Record
	for record, err = reader.Next(); err != nil; {
		if dal.ErrNoMoreRecords == err {
			if endCursor, err = reader.Cursor(); err != nil {
				return
			}
			return
		} else if err != nil {
			return
		}
		transferIDs = append(transferIDs, record.Key().ID.(string))
	}
	return
}

func (transferDalGae TransferDalGae) LoadTransfersByContactID(ctx context.Context, contactID string, offset, limit int) (transfers []models4debtus.TransferEntry, hasMore bool, err error) {
	if limit == 0 {
		err = errors.New("LoadTransfersByContactID(): limit == 0")
		return
	}
	if contactID == "" {
		err = errors.New("LoadTransfersByContactID(): contactID == 0")
		return
	}
	q := dal.From(models4debtus.TransfersCollection).
		WhereField("BothCounterpartyIDs", dal.Equal, contactID).
		OrderBy(dal.DescendingField("DtCreated")).
		Limit(limit).
		Offset(offset).
		SelectInto(models4debtus.NewTransferRecord)

	if transfers, err = transferDalGae.loadTransfers(ctx, q); err != nil {
		return
	}
	hasMore = len(transfers) > limit
	return
}

func (transferDalGae TransferDalGae) LoadLatestTransfers(ctx context.Context, offset, limit int) ([]models4debtus.TransferEntry, error) {
	q := dal.From(models4debtus.TransfersCollection).
		OrderBy(dal.DescendingField("DtCreated")).
		Limit(limit).
		Offset(offset).
		SelectInto(models4debtus.NewTransferRecord)
	return transferDalGae.loadTransfers(ctx, q)
}

func (transferDalGae TransferDalGae) loadTransfers(ctx context.Context, q dal.Query) (transfers []models4debtus.TransferEntry, err error) {
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	//var reader dal.Reader
	//if reader, err = db.QueryReader(ctx, q); err != nil {
	//	return
	//}
	var records []dal.Record
	if records, err = db.QueryAllRecords(ctx, q); err != nil {
		return
	}
	transfers = make([]models4debtus.TransferEntry, len(records))
	for i, record := range records {
		transfers[i] = models4debtus.NewTransfer(record.Key().ID.(string), record.Data().(*models4debtus.TransferData))
	}
	return transfers, nil
}
