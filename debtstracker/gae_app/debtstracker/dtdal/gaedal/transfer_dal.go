package gaedal

import (
	"bytes"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/delaying"
	"time"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

type TransferDalGae struct {
}

func NewTransferDalGae() TransferDalGae {
	return TransferDalGae{}
}

var _ dtdal.TransferDal = (*TransferDalGae)(nil)

func _loadDueOnTransfers(c context.Context, tx dal.ReadSession, userID string, limit int, filter func(q dal.QueryBuilder) dal.QueryBuilder) (transfers []models.TransferEntry, err error) {
	q := dal.From(models.TransfersCollection).
		WhereField("BothUserIDs", "=", userID).
		WhereField("IsOutstanding", "=", true).OrderBy(dal.AscendingField("DtDueOn"))
	q = filter(q).Limit(limit)
	query := q.SelectInto(models.NewTransferRecord)
	var (
		transferRecords []dal.Record
	)

	transferRecords, err = tx.QueryAllRecords(c, query)

	transfers = make([]models.TransferEntry, len(transferRecords))
	for i, transferRecord := range transferRecords {
		transfer := models.NewTransfer(transferRecord.Key().ID.(string), transferRecord.Data().(*models.TransferData))
		transfers[i] = transfer
	}
	return
}

func (transferDalGae TransferDalGae) LoadOverdueTransfers(c context.Context, tx dal.ReadSession, userID string, limit int) ([]models.TransferEntry, error) {
	return _loadDueOnTransfers(c, tx, userID, limit, func(q dal.QueryBuilder) dal.QueryBuilder {
		return q.WhereField("DtDueOn", dal.GreaterThen, time.Time{}).WhereField("DtDueOn", dal.LessThen, time.Now())
	})
}

func (transferDalGae TransferDalGae) LoadDueTransfers(c context.Context, tx dal.ReadSession, userID string, limit int) ([]models.TransferEntry, error) {
	return _loadDueOnTransfers(c, tx, userID, limit, func(q dal.QueryBuilder) dal.QueryBuilder {
		return q.WhereField("DtDueOn", dal.GreaterThen, time.Now())
	})
}

func (transferDalGae TransferDalGae) GetTransfersByID(c context.Context, tx dal.ReadSession, transferIDs []string) (transfers []models.TransferEntry, err error) {
	transfers = make([]models.TransferEntry, len(transferIDs))
	records := make([]dal.Record, len(transferIDs))
	for i, transferID := range transferIDs {
		transfers[i] = models.NewTransfer(transferID, nil)
		records[i] = transfers[i].Record
	}
	if err = tx.GetMulti(c, records); err != nil {
		return
	}
	return
}

func (transferDalGae TransferDalGae) LoadOutstandingTransfers(c context.Context, tx dal.ReadSession, periodEnds time.Time, userID, contactID string, currency money.CurrencyCode, direction models.TransferDirection) (transfers []models.TransferEntry, err error) {
	log.Debugf(c, "TransferDalGae.LoadOutstandingTransfers(periodEnds=%v, userID=%v, contactID=%v currency=%v, direction=%v)", periodEnds, userID, contactID, currency, direction)
	const limit = 100

	// TODO: Load outstanding transfer just for the specific contact & specific direction
	q := dal.From(models.TransfersCollection).
		Where(
			dal.WhereField("BothUserIDs", dal.Equal, userID),
			dal.WhereField("Currency", dal.Equal, string(currency)),
			dal.WhereField("IsOutstanding", dal.Equal, true),
		).
		OrderBy(dal.AscendingField("DtCreated")).
		Limit(limit).
		SelectInto(models.NewTransferRecord)
	var transferRecords []dal.Record
	transferRecords, err = tx.QueryAllRecords(c, q)
	transfers = models.TransfersFromRecords(transferRecords)
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
		if err = delayFixTransfersIsOutstanding.EnqueueWork(c, delaying.With(common.QUEUE_TRANSFERS, "fix-transfers-is-outstanding", 0), transfersIDsToFixIsOutstanding); err != nil {
			log.Errorf(c, "failed to delay task to fix transfers IsOutstanding")
			err = nil
		}
	}
	if errorMessages.Len() > 0 {
		log.Errorf(c, errorMessages.String())
	}
	if warnings.Len() > 0 {
		log.Warningf(c, warnings.String())
	}
	if debugs.Len() > 0 {
		log.Debugf(c, debugs.String())
	}
	return
}

func fixTransfersIsOutstanding(c context.Context, transferIDs []string) (err error) {
	log.Debugf(c, "fixTransfersIsOutstanding(%v)", transferIDs)
	for _, transferID := range transferIDs {
		if _, transferErr := fixTransferIsOutstanding(c, transferID); transferErr != nil {
			log.Errorf(c, "Failed to fix transfer %v: %v", transferID, err)
			err = transferErr
		}
	}
	return
}

func fixTransferIsOutstanding(c context.Context, transferID string) (transfer models.TransferEntry, err error) {
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if transfer, err = facade.Transfers.GetTransferByID(c, tx, transferID); err != nil {
			return err
		}
		if transfer.Data.GetOutstandingValue(time.Now()) == 0 {
			transfer.Data.IsOutstanding = false
			return facade.Transfers.SaveTransfer(c, tx, transfer)
		}
		return nil
	})
	if err == nil {
		log.Warningf(c, "Fixed IsOutstanding (set to false) for transfer %v", transferID)
	} else {
		log.Errorf(c, "Failed to fix IsOutstanding for transfer %v", transferID)
	}
	return
}

func (transferDalGae TransferDalGae) LoadTransfersByUserID(c context.Context, userID string, offset, limit int) (transfers []models.TransferEntry, hasMore bool, err error) {
	if limit == 0 {
		err = errors.New("limit == 0")
		return
	}
	if userID == "" {
		err = errors.New("userID == 0")
		return
	}
	q := dal.From(models.TransfersCollection).
		WhereField("BothUserIDs", dal.Equal, userID).
		OrderBy(dal.DescendingField("DtCreated")).
		SelectInto(models.NewTransferRecord)

	if transfers, err = transferDalGae.loadTransfers(c, q); err != nil {
		return
	}
	hasMore = len(transfers) > limit
	return
}

func (transferDalGae TransferDalGae) LoadTransferIDsByContactID(c context.Context, contactID string, limit int, startCursor string) (transferIDs []string, endCursor string, err error) {
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
	q := dal.From(models.TransfersCollection).
		WhereField("BothCounterpartyIDs", dal.Equal, contactID).
		Limit(limit).
		StartFrom(dal.Cursor(startCursor)).
		SelectInto(models.NewTransferRecord)

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
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	var reader dal.Reader
	if reader, err = db.QueryReader(c, q); err != nil {
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

func (transferDalGae TransferDalGae) LoadTransfersByContactID(c context.Context, contactID string, offset, limit int) (transfers []models.TransferEntry, hasMore bool, err error) {
	if limit == 0 {
		err = errors.New("LoadTransfersByContactID(): limit == 0")
		return
	}
	if contactID == "" {
		err = errors.New("LoadTransfersByContactID(): contactID == 0")
		return
	}
	q := dal.From(models.TransfersCollection).
		WhereField("BothCounterpartyIDs", dal.Equal, contactID).
		OrderBy(dal.DescendingField("DtCreated")).
		Limit(limit).
		Offset(offset).
		SelectInto(models.NewTransferRecord)

	if transfers, err = transferDalGae.loadTransfers(c, q); err != nil {
		return
	}
	hasMore = len(transfers) > limit
	return
}

func (transferDalGae TransferDalGae) LoadLatestTransfers(c context.Context, offset, limit int) ([]models.TransferEntry, error) {
	q := dal.From(models.TransfersCollection).
		OrderBy(dal.DescendingField("DtCreated")).
		Limit(limit).
		Offset(offset).
		SelectInto(models.NewTransferRecord)
	return transferDalGae.loadTransfers(c, q)
}

func (transferDalGae TransferDalGae) loadTransfers(c context.Context, q dal.Query) (transfers []models.TransferEntry, err error) {
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	//var reader dal.Reader
	//if reader, err = db.QueryReader(c, q); err != nil {
	//	return
	//}
	var records []dal.Record
	if records, err = db.QueryAllRecords(c, q); err != nil {
		return
	}
	transfers = make([]models.TransferEntry, len(records))
	for i, record := range records {
		transfers[i] = models.NewTransfer(record.Key().ID.(string), record.Data().(*models.TransferData))
	}
	return transfers, nil
}
