package dbo4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

const (
	OrderShippingPointStatusPending    = "pending"
	OrderShippingPointStatusProcessing = "processing"
	OrderShippingPointStatusCompleted  = "completed"
)

func validateOrderShippingPointStatus(field, v string) error {
	switch v {
	case OrderShippingPointStatusPending, OrderShippingPointStatusCompleted, OrderShippingPointStatusProcessing: // OK
		return nil
	case "":
		return validation.NewErrRecordIsMissingRequiredField(field)
	default:
		return validation.NewErrBadRecordFieldValue(field, fmt.Sprintf("unsupported value: [%v]", v))
	}
}

// ShippingPointCounterparty is used in OrderShippingPoint
type ShippingPointCounterparty struct {
	ContactID string `json:"contactID" firestore:"contactID"`
	Title     string `json:"title" firestore:"title"`
}

// Validate returns an error if the ShippingPointBase is invalid.
func (v ShippingPointCounterparty) Validate() error {
	if v.ContactID == "" {
		return validation.NewErrRecordIsMissingRequiredField("contactID")
	}
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// ShippingPointLocation is used in OrderShippingPoint
type ShippingPointLocation struct {
	ContactID string            `json:"contactID" firestore:"contactID"`
	Title     string            `json:"title" firestore:"title"`
	Address   *dbmodels.Address `json:"address" firestore:"address"`
}

// Validate returns an error if the ShippingPointLocation is invalid.
func (v ShippingPointLocation) Validate() error {
	if v.ContactID == "" {
		return validation.NewErrRecordIsMissingRequiredField("contactID")
	}
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if v.Address == nil {
		return validation.NewErrRequestIsMissingRequiredField("address")
	} else {
		if err := v.Address.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("address", err.Error())
		}
	}
	return nil
}

// OrderShippingPoint represents a pick-up or drop-off point for an order.
// It's the one referenced by ContainerPoint and ContainerSegment.
type OrderShippingPoint struct {
	ID string `json:"id" firestore:"id"`
	ShippingPointBase

	// ScheduledStartDate is a ISO8601 date of the 1st container to arrive
	ScheduledStartDate string `json:"scheduledStartDate,omitempty" firestore:"scheduledStartDate,omitempty"`
	// ScheduledEndDate is an ISO8601 date of the last container to depart
	ScheduledEndDate string `json:"scheduledEndDate,omitempty" firestore:"scheduledEndDate,omitempty"`

	Counterparty ShippingPointCounterparty `json:"counterparty" firestore:"counterparty"`

	// Location where the order processed. TODO: Consider if it should be mandatory (e.g. not the * pointer)
	Location *ShippingPointLocation `json:"location" firestore:"location"`
}

// Validate validates the OrderShippingPoint.
func (v OrderShippingPoint) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := ValidateShippingPointTasksRecord(v.Tasks); err != nil {
		return err // do not wrap error here
	}
	if err := validateOrderShippingPointStatus("status", v.Status); err != nil {
		return err // do not wrap error here
	}
	if err := v.Counterparty.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("counterparty", err.Error())
	}
	if v.Location != nil {
		if err := v.Location.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("location", err.Error())
		}
		if v.Location.ContactID == v.Counterparty.ContactID {
			return validation.NewErrBadRecordFieldValue("location.contactID", "must be different from counterparty.contactID")
		}
	}
	return nil
}

// WithShippingPoints holds a list of OrderShippingPoint's
// We do need the `ShippingPoints` field as it can not be fully deducted
// from `Contacts`, `Containers` & `Segments`.
type WithShippingPoints struct {
	ShippingPoints []*OrderShippingPoint `json:"shippingPoints,omitempty" firestore:"shippingPoints,omitempty"`
}

func (v *WithShippingPoints) Updates() []dal.Update {
	const field = "shippingPoints"
	if len(v.ShippingPoints) == 0 {
		return []dal.Update{
			{Field: field, Value: dal.DeleteField},
		}
	}
	return []dal.Update{
		{Field: field, Value: v.ShippingPoints},
	}
}

func (v *WithShippingPoints) validateOrder(orderDto OrderDbo) error {
	if err := v.Validate(); err != nil {
		return err
	}
	for i, shippingPoint := range v.ShippingPoints {
		if err := v.validateShippingPoint(i, shippingPoint, orderDto); err != nil {
			return err
		}
	}
	return nil
}

func (v *WithShippingPoints) validateShippingPoint(i int, shippingPoint *OrderShippingPoint, orderDto OrderDbo) error {
	field := func() string {
		return fmt.Sprintf("shippingPoints[%v]", i)
	}

	if shippingPoint.Location != nil && shippingPoint.Location.Address != nil {
		_, contact := orderDto.GetContactByID(shippingPoint.Location.ContactID)
		if contact == nil {
			return validation.NewErrBadRecordFieldValue(field()+"location.contactID", "not found in order contacts by ContactID="+shippingPoint.Location.ContactID)
		}
		//if shippingPoint.Location.Address != nil && shippingPoint.Location.Address.CountryID != contact.CountryID {
		//	return validation.NewErrBadRecordFieldValue(
		//		field()+"location.countryID",
		//		fmt.Sprintf(`"%s" does not match value in order contact "%s"`,
		//			shippingPoint.Location.Address.CountryID, contact.CountryID))
		//}
	}

	if err := ValidateShippingPointTasksRequest(shippingPoint.Tasks, false); err != nil {
		return err // do not wrap error here
	}

	for _, task := range shippingPoint.Tasks {
		var counterpartyRoles []string
		switch task {
		case ShippingPointTaskPick, ShippingPointTaskDrop:
			counterpartyRoles = append(counterpartyRoles, CounterpartyRolePortFrom)
		case ShippingPointTaskLoad:
			counterpartyRoles = append(counterpartyRoles, CounterpartyRoleDispatchPoint)
		case ShippingPointTaskUnload:
			counterpartyRoles = append(counterpartyRoles, CounterpartyRoleReceivePoint)
		default:
			return validation.NewErrBadRecordFieldValue(field()+".type", fmt.Sprintf("unsupported value: [%v]", shippingPoint.Tasks))
		}
		for _, counterpartyRole := range counterpartyRoles {
			var contactID string
			switch GetCounterpartyTypeByRole(counterpartyRole) {
			case CounterpartyTypeLocation:
				if shippingPoint.Location != nil {
					contactID = shippingPoint.Location.ContactID
				}
			case CounterpartyTypeCompany:
				contactID = shippingPoint.Counterparty.ContactID // No need to check shippingPoint.Counterparty as it is NOT by reference
			default:
				panic("unexpected counterparty role: " + counterpartyRole)
			}
			if contactID == "" {
				continue
			}
			if _, c := orderDto.GetCounterpartyByRoleAndContactID(counterpartyRole, contactID); c == nil {
				return validation.NewErrBadRecordFieldValue(
					field(),
					fmt.Sprintf("contact of shipping point not found in order's field `counterparties` by {role=%v&id=%v}",
						counterpartyRole,
						contactID))
			}
		}
	}

	return nil
}

// Validate validates the field	ShippingPoints []*OrderShippingPoint.
func (v *WithShippingPoints) Validate() error {
	for i, p := range v.ShippingPoints {
		field := func() string {
			return fmt.Sprintf("shippingPoints[%v]", i)
		}
		if err := p.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(field(), err.Error())
		}
		for j, p2 := range v.ShippingPoints {
			if j != i {
				if p2.ID == p.ID {
					return validation.NewErrBadRecordFieldValue(field(), fmt.Sprintf("duplicate ContactID at indexes %v & %v: [%v]", i, j, p.ID))
				}
				if p2.Location != nil && p.Location != nil && p2.Location.ContactID == p.Location.ContactID {
					return validation.NewErrBadRecordFieldValue(field(), fmt.Sprintf("duplicate location ContactID at indexes %v & %v: [%v]", i, j, p.ID))
				}
			}
		}
	}
	return nil
}

// GetShippingPointByID returns an *OrderShippingPoint & index by ContactID.
func (v *WithShippingPoints) GetShippingPointByID(id string) (i int, shippingPoint *OrderShippingPoint) {
	for i, shippingPoint = range v.ShippingPoints {
		if shippingPoint.ID == id {
			return i, shippingPoint
		}
	}
	return -1, nil
}

//// GetShippingPointByID returns an *OrderShippingPoint & index by ContactID.
//func (v *WithShippingPoints) GetShippingPointstByCounterID(containerID string) (i int, shippingPoints []*OrderShippingPoint) {
//	for i, shippingPoint := range v.ShippingPoints {
//		if shippingPoint. == id {
//			return i, shippingPoint
//		}
//	}
//	return -1, nil
//}

// GetShippingPointByContactID returns an *OrderShippingPoint & index by ContactID.
func (v *WithShippingPoints) GetShippingPointByContactID(contactID string) (i int, shippingPoint *OrderShippingPoint) {
	for i, shippingPoint = range v.ShippingPoints {
		if shippingPoint.Counterparty.ContactID == contactID ||
			shippingPoint.Location != nil && shippingPoint.Location.ContactID == contactID {
			return i, shippingPoint
		}
	}
	return -1, nil
}
