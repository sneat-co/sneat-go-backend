package dto4logist

import (
	"fmt"
	"strings"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/strongo/validation"
)

// PointOfNewContainer is used in NewContainer
type PointOfNewContainer struct {
	ShippingPointID string `json:"shippingPointID"`
	Tasks           []dbo4logist.ShippingPointTask
}

func (v PointOfNewContainer) Validate() error {
	if strings.TrimSpace(v.ShippingPointID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("shippingPointID")
	}
	return nil
}

// NewContainer defines a new container
type NewContainer struct {
	dbo4logist.OrderContainerBase
	Points []PointOfNewContainer `json:"points"`
}

// Validate validates container request
func (v NewContainer) Validate() error {
	if err := v.OrderContainerBase.Validate(); err != nil {
		return err
	}
	for i, point := range v.Points {
		if err := point.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("points[%d]", i), err.Error())
		}
	}
	return nil
}

// AddContainersRequest defines a request to add containers to an order
type AddContainersRequest struct {
	OrderRequest
	Containers []NewContainer `json:"containers"`
}

// Validate validates the request
func (v AddContainersRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if len(v.Containers) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("containers")
	}
	for i, c := range v.Containers {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("containers[%v]: %v", i, err)
		}
	}
	return nil
}

// ContainerRequest defines a request related to a container
type ContainerRequest struct {
	OrderRequest
	ContainerID string `json:"containerID"`
}

// Validate validates the request
func (v ContainerRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if err := validateContainerID("containerID", v.ContainerID); err != nil {
		return validation.NewBadRequestError(err)
	}
	return nil
}

// SetContainerFieldsRequest defines a request to set fields of a container
type SetContainerFieldsRequest struct {
	ContainerRequest
	SetFieldsRequest
}

// Validate returns an error if the SetContainerFieldsRequest is invalid
func (v SetContainerFieldsRequest) Validate() error {
	if err := v.ContainerRequest.Validate(); err != nil {
		return err
	}
	if err := v.SetFieldsRequest.Validate(); err != nil {
		return err
	}
	for name := range v.SetStrings {
		switch name {
		case "instructions", "number": // OK
		case "":
			return validation.NewErrRequestIsMissingRequiredField("setStrings")
		default:
			return validation.NewErrBadRecordFieldValue("setStrings", "unknown field name: "+name)
		}
	}
	for name := range v.SetDates {
		return validation.NewErrBadRecordFieldValue("setDates", "unknown field name: "+name)
	}
	return nil
}
