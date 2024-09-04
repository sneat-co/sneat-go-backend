package coretodo

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
)

type WithRecordChanges struct { // TODO: move to github.com/dal-go/dalgo/dal?
	recordsToInsert []dal.Record // We might need to consider to use []*dal.Record to simplify updating dal.Record.ID
	RecordsToUpdate []RecordUpdates
	RecordsToDelete []*dal.Key
}

func (v *WithRecordChanges) RecordsToInsert() (records []dal.Record) {
	records = make([]dal.Record, len(v.recordsToInsert))
	copy(records, v.recordsToInsert)
	return
}

func (v *WithRecordChanges) QueueForInsert(records ...dal.Record) {
	for i, record := range records {
		if record == nil {
			panic(fmt.Sprintf("record #%d is required", i))
		}
		if record.Key() == nil {
			panic(fmt.Sprintf("record #%d.Key() is required", i))
		}
		if record.Data() == nil {
			panic(fmt.Sprintf("record #%d.Data() is required", i))
		}
		v.recordsToInsert = append(v.recordsToInsert, record)
	}
}

func (v *WithRecordChanges) ApplyChanges(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
	if len(v.recordsToInsert) > 0 {
		if err = tx.InsertMulti(ctx, v.recordsToInsert); err != nil {
			err = fmt.Errorf("failed to insert records: %w", err)
			return
		}
	}
	if len(v.RecordsToUpdate) > 0 {
		for _, record2update := range v.RecordsToUpdate {
			key := record2update.Record.Key()
			if err = tx.Update(ctx, key, record2update.Updates); err != nil {
				return fmt.Errorf("failed to update record %s: %w", key, err)
			}
		}
	}
	if len(v.RecordsToDelete) > 0 {
		if err = tx.DeleteMulti(ctx, v.RecordsToDelete); err != nil {
			err = fmt.Errorf("failed to delete records: %w", err)
			return
		}
	}
	v.recordsToInsert = nil
	v.RecordsToUpdate = nil
	v.RecordsToDelete = nil
	return
}
