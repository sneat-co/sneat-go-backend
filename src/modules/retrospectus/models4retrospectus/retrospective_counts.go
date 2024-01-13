package models4retrospectus

// Kept here for now to get rid of cyclic imports

import "github.com/strongo/validation"

// RetrospectiveCounts record // TODO: move to retrospectives module
type RetrospectiveCounts struct {
	ItemsByUserAndType map[string]map[string]int `json:"itemsByUserAndType,omitempty" firestore:"itemsByUserAndType,omitempty"`
}

// Validate validates RetrospectiveCounts record
func (v *RetrospectiveCounts) Validate() error {
	if len(v.ItemsByUserAndType) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("ItemsByUserAndType")
	}
	return nil
}
