package models4logist

import (
	"fmt"
	"github.com/strongo/validation"
	"time"
)

const (
	// OrderStepStatusPending = "pending"
	OrderStepStatusPending = "pending"

	// OrderStepStatusCompleted = "completed"
	OrderStepStatusCompleted = "completed"
)

// OrderStep - TODO: document intended usage
type OrderStep struct {
	ID          string     `json:"id" firestore:"id"`
	Status      string     `json:"status" firestore:"status"`                               // "pending", "completed"
	PlannedDate string     `json:"plannedDate,omitempty" firestore:"plannedDate,omitempty"` // ISO8601
	Completed   *time.Time `json:"completed,omitempty" firestore:"completed,omitempty"`
}

// Validate validates OrderStep
func (v OrderStep) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	switch v.Status {
	case OrderStepStatusPending, OrderStepStatusCompleted: // OK
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	default:
		return validation.NewErrBadRecordFieldValue("status", fmt.Sprintf("unknown status: [%v]", v.Status))
	}
	return nil
}
