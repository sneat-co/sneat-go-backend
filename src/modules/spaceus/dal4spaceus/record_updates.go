package dal4spaceus

import "github.com/dal-go/dalgo/dal"

// RecordUpdates defines updates for a record
type RecordUpdates struct { // TODO: move to DALgo
	Key     *dal.Key
	Updates []dal.Update
}
