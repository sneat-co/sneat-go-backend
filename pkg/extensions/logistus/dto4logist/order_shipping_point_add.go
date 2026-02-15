package dto4logist

import (
	"fmt"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/strongo/validation"
)

type task = dbo4logist.ShippingPointTask

type AddContainerPoint struct {
	ID    string `json:"id"`    // container ID
	Tasks []task `json:"tasks"` // e.g. "load", "unload", "pick", "drop"
}

func (v AddContainerPoint) Validate() error {
	if err := validateContainerID("id", v.ID); err != nil {
		return err
	}
	return dbo4logist.ValidateShippingPointTasksRequest(v.Tasks, true)
}

// AddOrderShippingPointRequest represents a request to create a new order shipping point.
type AddOrderShippingPointRequest struct {
	OrderRequest
	Tasks             []task `json:"tasks,omitempty"` // e.g. "load", "unload", "pick", "drop"
	LocationContactID string `json:"locationContactID"`

	Containers []AddContainerPoint `json:"containers,omitempty"`
}

// Validate returns an error if AddOrderShippingPointRequest is invalid.
func (v AddOrderShippingPointRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if err := dbo4logist.ValidateShippingPointTasksRequest(v.Tasks, len(v.Containers) == 0); err != nil {
		return err
	}
	if v.LocationContactID == "" {
		return validation.NewErrRequestIsMissingRequiredField("locationContactID")
	}
	for i, container := range v.Containers {
		if err := container.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("containers[%d]", i), err.Error())
		}
		for j, container2 := range v.Containers {
			if j != i && container2.ID == container.ID {
				return validation.NewErrBadRequestFieldValue("containers", fmt.Sprintf(`duplicate container IDs at indexes %d & %d: ContactID="%v"`, j, i, container.ID))
			}
		}
	}
	return nil
}
