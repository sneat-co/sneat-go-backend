package coretodo

import "github.com/dal-go/dalgo/dal"

// RecordUpdates defines updates for a record
type RecordUpdates struct { // TODO: move to DALgo
	Record  dal.Record
	Updates []dal.Update
}
