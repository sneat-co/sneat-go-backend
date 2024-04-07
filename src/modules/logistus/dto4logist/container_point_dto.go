package dto4logist

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

// ContainerPointsRequest is a request regards container points
type ContainerPointsRequest struct {
	OrderRequest
	ContainerID      string   `json:"containerID"`
	ShippingPointIDs []string `json:"shippingPointIDs"`
}

// ContainerPointRequest is a request regards container point
type ContainerPointRequest struct {
	OrderRequest
	ShippingPointID string `json:"shippingPointID"`
	ContainerID     string `json:"containerID"`
}

// Validate returns error if request is invalid
func (v ContainerPointsRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if err := validateContainerID("containerID", v.ContainerID); err != nil {
		return validation.NewBadRequestError(err)
	}
	if len(v.ShippingPointIDs) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("shippingPointIDs")
	}
	for i, shippingPointID := range v.ShippingPointIDs {
		if err := models4teamus.ValidateShippingPointID(shippingPointID); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("shippingPointIDs[%v]", i), err.Error())
		}
	}
	return nil
}

// Validate returns error if request is invalid
func (v ContainerPointRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if err := validateContainerID("containerID", v.ContainerID); err != nil {
		return validation.NewBadRequestError(err)
	}
	if err := models4teamus.ValidateShippingPointID(v.ShippingPointID); err != nil {
		return validation.NewErrBadRequestFieldValue("shippingPointID", err.Error())
	}
	return nil
}

// AddContainerPointsRequest is a request to add a container point to an order
type AddContainerPointsRequest struct {
	OrderRequest
	ContainerPoints []models4logist.ContainerPoint `json:"containerPoints"`
}

// Validate returns error if request is invalid
func (v AddContainerPointsRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	for i, containerPoint := range v.ContainerPoints {
		if err := containerPoint.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("containerPoints[%v]: %v", i, containerPoint), err.Error())
		}
	}
	return nil
}

// UpdateContainerPointRequest is a request to update a container point in an order
type UpdateContainerPointRequest struct {
	ContainerPointRequest
	models4logist.FreightPoint
	ArrivesDate *string `json:"arrivesDate,omitempty"`
	DepartsDate *string `json:"departsDate,omitempty"`
}

// Validate returns error if request is invalid
func (v UpdateContainerPointRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if err := models4teamus.ValidateShippingPointID(v.ShippingPointID); err != nil {
		return validation.NewErrBadRequestFieldValue("shippingPointID", err.Error())
	}
	if err := validateContainerID("containerID", v.ContainerID); err != nil {
		return err
	}
	if err := v.FreightPoint.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("FreightPoint", err.Error())
	}
	return nil
}

// SetContainerPointTaskRequest is a request to set container point task
type SetContainerPointTaskRequest struct {
	ContainerPointRequest
	Task  models4logist.ShippingPointTask `json:"task"`
	Value bool                            `json:"value"`
}

// Validate returns error if SetContainerPointTaskRequest is invalid
func (v SetContainerPointTaskRequest) Validate() error {
	if err := v.ContainerPointRequest.Validate(); err != nil {
		return err
	}
	return models4logist.ValidateShippingPointTask(v.Task,
		models4logist.ValidatingRequest,
		func() string {
			return "task"
		},
	)
}

// SetContainerEndpointFieldsRequest is a request to set container point dates
type SetContainerEndpointFieldsRequest struct {
	ContainerPointRequest
	Side        models4logist.EndpointSide `json:"side"`
	Dates       map[string]string          `json:"dates"`
	Times       map[string]string          `json:"times"`
	ByContactID *string                    `json:"byContactID,omitempty"`
}

// Validate returns error if SetContainerEndpointFieldsRequest is invalid
func (v SetContainerEndpointFieldsRequest) Validate() error {
	if err := v.ContainerPointRequest.Validate(); err != nil {
		return err
	}
	switch v.Side {
	case models4logist.EndpointSideArrival, models4logist.EndpointSideDeparture:
	// OK
	case "":
		return validation.NewErrRequestIsMissingRequiredField("side")
	default:
		return validation.NewErrBadRequestFieldValue("side", "must be 'arrival' or 'departure'")
	}
	if v.ByContactID == nil && len(v.Dates) == 0 && len(v.Times) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("Either `byContactID` or `dates` or `times` field or all must be set")
	}
	for name, value := range v.Dates {
		if name == "" {
			return validation.NewErrRequestIsMissingRequiredField("dates")
		}
		if value != "" {
			if _, err := validate.DateString(value); err != nil {
				return validation.NewErrBadRequestFieldValue(fmt.Sprintf("dates.%s", name), err.Error())
			}
		}
	}
	for name, value := range v.Times {
		if name == "" {
			return validation.NewErrRequestIsMissingRequiredField("times")
		}
		if value != "" {
			if err := validate.IsValidateTime(value); err != nil {
				return validation.NewErrBadRequestFieldValue(fmt.Sprintf("times.%s", name), err.Error())
			}
		}
	}
	return nil
}

// SetContainerPointFieldsRequest is a request to set container point fields
type SetContainerPointFieldsRequest struct {
	ContainerPointRequest
	SetStrings map[string]string
}

// Validate returns error if SetContainerPointFieldsRequest is invalid
func (v SetContainerPointFieldsRequest) Validate() error {
	if err := v.ContainerPointRequest.Validate(); err != nil {
		return err
	}
	for name, value := range v.SetStrings {
		switch name {
		case "notes":
			if len(value) > 10000 {
				return validation.NewErrBadRequestFieldValue("setStrings.specialInstructions", "must be less than 10,000 characters")
			}
			continue
		case "refNumber":
			if len(value) > 50 {
				return validation.NewErrBadRequestFieldValue("setStrings.specialInstructions", "must be less than 50 characters")
			}
			continue
		case "":
			return validation.NewErrBadRequestFieldValue("setStrings", "name must not be empty")
		default:
			return validation.NewErrBadRequestFieldValue("setStrings", "unknown field name: "+name)
		}
	}
	return nil
}

// SetContainerPointFreightFieldsRequest is a request to set container point freight fields
type SetContainerPointFreightFieldsRequest struct {
	ContainerPointRequest
	Task     models4logist.ShippingPointTask `json:"task"`
	Integers map[string]*int                 `json:"integers"`
}

// Validate returns error if SetContainerPointFreightFieldsRequest is invalid
func (v SetContainerPointFreightFieldsRequest) Validate() error {
	if err := v.ContainerPointRequest.Validate(); err != nil {
		return err
	}
	for name, value := range v.Integers {
		if value != nil && *value < 0 {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("integers.%s", name), fmt.Sprintf("must be >= 0, got: %d", *value))
		}
	}
	if err := models4logist.ValidateShippingPointTask(v.Task,
		models4logist.ValidatingRequest,
		func() string {
			return "task"
		},
	); err != nil {
		return err
	}
	return nil
}
