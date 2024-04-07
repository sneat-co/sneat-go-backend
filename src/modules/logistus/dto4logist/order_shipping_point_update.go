package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/strongo/validation"
)

type shippingPointTask = models4logist.ShippingPointTask

// UpdateShippingPointRequest represents a request to update a shipping point of an order.
type UpdateShippingPointRequest struct {
	OrderShippingPointRequest
	Tasks []shippingPointTask `json:"tasks,omitempty"`
	SetFieldsRequest
}

// Validate returns an error if the request is invalid.
func (v UpdateShippingPointRequest) Validate() error {
	if err := v.OrderShippingPointRequest.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("OrderRequest", err.Error())
	}
	if err := models4logist.ValidateShippingPointTasksRequest(v.Tasks, false); err != nil {
		return err
	}
	if err := v.SetFieldsRequest.Validate(); err != nil {
		return err
	}
	if len(v.SetDates) == 0 && len(v.SetStrings) == 0 {
		return validation.NewErrBadRequestFieldValue("SetDates|SetStrings", "at least 1 must be set")
	}
	for name := range v.SetDates {
		switch name {
		case "scheduledStartDate", "scheduledEndDate":
			break // OK
		case "":
			return validation.NewErrBadRequestFieldValue("SetDates", "field name must be set")
		default:
			return validation.NewErrBadRequestFieldValue("SetDates."+name, "unknown field name: "+name)
		}
	}
	for name := range v.SetStrings {
		switch name {
		case "notes", "departsDate", "arrivedDate", "departedDate":
			break // OK
		case "":
			return validation.NewErrBadRequestFieldValue("SetDates", "field name must be set")
		default:
			return validation.NewErrBadRequestFieldValue("SetDates."+name, "unknown field name: "+name)
		}
	}
	return nil
}
