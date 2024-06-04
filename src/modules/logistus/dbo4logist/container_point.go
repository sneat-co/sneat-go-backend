package dbo4logist

import (
	"fmt"
	"github.com/strongo/validation"
)

// ContainerPoint represents a shipping point of a container.
type ContainerPoint struct {
	ContainerID     string `json:"containerID" firestore:"containerID"`
	ShippingPointID string `json:"shippingPointID" firestore:"shippingPointID"`
	ContainerEndpoints
	ShippingPointBase
	RefNumber string `json:"refNumber,omitempty" firestore:"refNumber,omitempty"`
}

func (v ContainerPoint) String() string {
	return fmt.Sprintf("ContainerPoint(containerID=%s,shippingPointID=%s){Status=%s}", v.ContainerID, v.ShippingPointID, v.Status)
}

// Validate returns an error if the ShippingPointBase is invalid.
func (v ContainerPoint) Validate() error {
	if v.ShippingPointID == "" {
		return validation.NewErrRecordIsMissingRequiredField("shippingPointID")
	}
	if v.ContainerID == "" {
		return validation.NewErrRecordIsMissingRequiredField("containerID")
	}
	if err := v.ShippingPointBase.Validate(); err != nil {
		return err
	}
	if err := v.ContainerEndpoints.Validate(); err != nil {
		return err
	}
	if l := len(v.RefNumber); l > 50 {
		return validation.NewErrBadRecordFieldValue("refNumber", fmt.Sprintf("should be < 50 characters, got %d", l))
	}
	return nil
}
