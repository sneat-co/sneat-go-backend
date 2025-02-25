package dbo4retrospectus

import (
	"github.com/dal-go/dalgo/dal"
)

// NewRetrospectiveKey creates a new retrospective key
func NewRetrospectiveKey(retrospectiveID string, parent *dal.Key) (retrospectiveKey *dal.Key) {
	retrospectiveKey = dal.NewKeyWithParentAndID(parent, "api4meetingus", retrospectiveID)
	return retrospectiveKey
}
