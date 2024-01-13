package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"time"
)

type Canceled struct {
	At     time.Time       `json:"at" firestore:"at"`
	By     dbmodels.ByUser `json:"by" firestore:"by"`
	Reason string          `json:"reason,omitempty" firestore:"reason,omitempty"`
}

func (v Canceled) Validate() error {
	if v.At.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("at")
	}
	if v.By.UID == "" {
		return validation.NewErrRecordIsMissingRequiredField("by")
	}
	if len(v.Reason) > ReasonMaxLen {
		return validation.NewErrBadRecordFieldValue("reason",
			fmt.Sprintf("maximum length of reason is %v, got %v", ReasonMaxLen, len(v.Reason)))
	}
	return nil
}

type SingleHappeningSlotCancellation struct {
	SlotIDs  []string `json:"slotIDs" firestore:"slotIDs"`
	Canceled Canceled `json:"canceled" firestore:"canceled"`
	Reason   string   `json:"reason,omitempty" firestore:"reason,omitempty"`
}
