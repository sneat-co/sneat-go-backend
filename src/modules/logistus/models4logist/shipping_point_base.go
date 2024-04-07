package models4logist

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// ShippingPointBase is used in OrderShippingPoint and ContainerPoint.
type ShippingPointBase struct {
	Status ShippingPointStatus `json:"status" firestore:"status"` // "pending", "completed"
	Notes  string              `json:"notes,omitempty" firestore:"notes,omitempty"`

	FreightPoint // This is the target numbers that if set should be matched by the sum of the FreightPoint's in the ContainerPoint.

	Started   *time.Time `json:"started,omitempty" firestore:"started,omitempty"`
	Completed *time.Time `json:"completed,omitempty" firestore:"completed,omitempty"`
}

// String representation of the ContainerSegmentKey.
func (v ShippingPointBase) String() string {
	return fmt.Sprintf("status=%s&load=%v&unload=%v", v.Status, v.ToLoad, v.ToUnload)
}

// Validate returns an error if the ShippingPointBase is invalid.
func (v ShippingPointBase) Validate() error {
	if err := validateShippingPointStatus("status", v.Status); err != nil {
		return err // Do not wrap error here
	}
	if err := v.FreightPoint.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Notes) != v.Notes {
		return validation.NewErrBadRecordFieldValue("notes", "leading or closing spaces")
	}
	return nil
}
