package dbo4logist

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// CounterpartyRole is a role of a counterparty in an order
type CounterpartyRole = string

const (
	CounterpartyRoleBuyer         CounterpartyRole = "buyer"
	CounterpartyRoleConsignee     CounterpartyRole = "consignee"
	CounterpartyRoleCustomBroker  CounterpartyRole = "custom_broker"
	CounterpartyRoleDispatcher    CounterpartyRole = "dispatcher"
	CounterpartyRoleDispatchPoint CounterpartyRole = "dispatch_point" // Point where container needs to be brought and loaded
	CounterpartyRoleDriver        CounterpartyRole = "driver"         // Final point to drop a container
	CounterpartyRoleDropPoint     CounterpartyRole = "drop_point"     // Final point to drop a container
	CounterpartyRoleDispatchAgent CounterpartyRole = "dispatch_agent"
	CounterpartyRoleReceiveAgent  CounterpartyRole = "receive_agent"
	CounterpartyRoleNotifyParty   CounterpartyRole = "notify_party"
	CounterpartyRoleReceiver      CounterpartyRole = "receiver"
	CounterpartyRolePortFrom      CounterpartyRole = "port_from"
	CounterpartyRolePortTo        CounterpartyRole = "port_to"
	CounterpartyRolePickPoint     CounterpartyRole = "pick_point"    // Original point to pick up a container
	CounterpartyRoleReceivePoint  CounterpartyRole = "receive_point" // Point to unload a container without leaving the container
	CounterpartyRoleShip          CounterpartyRole = "ship"
	CounterpartyRoleShippingLine  CounterpartyRole = "shipping_line"
	CounterpartyRoleTruck         CounterpartyRole = "truck"
	CounterpartyRoleTrucker       CounterpartyRole = "trucker"
)

var KnownCounterpartyRoles = []CounterpartyRole{
	CounterpartyRoleBuyer,
	CounterpartyRoleConsignee,
	CounterpartyRoleCustomBroker,
	CounterpartyRoleDispatcher,
	CounterpartyRoleDispatchAgent,
	CounterpartyRoleReceiveAgent,
	CounterpartyRoleReceiver,
	CounterpartyRoleDispatchPoint,
	CounterpartyRoleDriver,
	CounterpartyRoleReceivePoint,
	CounterpartyRoleNotifyParty,
	CounterpartyRolePortFrom,
	CounterpartyRolePortTo,
	CounterpartyRolePickPoint,
	CounterpartyRoleDropPoint,
	CounterpartyRoleShippingLine,
	CounterpartyRoleTruck,
	CounterpartyRoleTrucker,
	CounterpartyRoleShip,
}

// CounterpartyType is a type of a counterparty in an order
type CounterpartyType = int

const (
	CounterpartyTypeUnknown = iota
	CounterpartyTypeCompany
	CounterpartyTypeLocation
)

// GetCounterpartyTypeByRole returns counterparty type by role
func GetCounterpartyTypeByRole(role CounterpartyRole) CounterpartyType {
	switch role {
	case
		CounterpartyRoleReceivePoint,
		CounterpartyRoleDispatchPoint,
		CounterpartyRolePickPoint,
		CounterpartyRoleDropPoint:
		return CounterpartyTypeLocation
	case
		CounterpartyRoleDispatcher,
		CounterpartyRoleReceiver,
		CounterpartyRolePortFrom,
		CounterpartyRolePortTo:
		return CounterpartyTypeCompany
	default:
		return CounterpartyTypeUnknown
	}
}

type CounterpartyParent struct {
	ContactID string           `json:"contactID" firestore:"contactID"`
	Role      CounterpartyRole `json:"role" firestore:"role"`
}

// OrderCounterparty is a counterparty in an order
type OrderCounterparty struct {
	ContactID string `json:"contactID" firestore:"contactID"`
	// We keep single role for counterparty because same contact can have different ref numbers for different roles.
	Role      CounterpartyRole    `json:"role" firestore:"role"` // buyer, seller, buyer_agent, seller_agent, transporter
	Parent    *CounterpartyParent `json:"parent,omitempty" firestore:"parent,omitempty"`
	RefNumber string              `json:"refNumber,omitempty" firestore:"refNumber,omitempty"`
	CountryID string              `json:"countryID" firestore:"countryID"` // ISO 3166-1 alpha-2
	Title     string              `json:"title" firestore:"title"`
}

// String returns string representation of OrderCounterparty
func (v OrderCounterparty) String() string {
	return fmt.Sprintf(`OrderCounterparty{ContactID="%s", Role="%s", CountryID="%s", Title="%s"}`, v.ContactID, v.Role, v.CountryID, v.Title)
}

// ValidateSegmentCounterpartyRole returns nil if role is valid, otherwise returns bad record field error.
func ValidateSegmentCounterpartyRole(value string) error {
	const field = "role"
	switch value {
	case "":
		return validation.NewErrRecordIsMissingRequiredField(field)
	case
		CounterpartyRoleDispatchPoint,
		CounterpartyRoleReceivePoint,
		CounterpartyRolePickPoint,
		CounterpartyRoleDropPoint,
		CounterpartyRoleTrucker:
		break // OK, no special checks required
	default:
		return validation.NewErrBadRecordFieldValue(field, fmt.Sprintf(`unknown value: "%v"`, value))
	}
	return nil
}

// ValidateOrderCounterpartyRoles returns nil if valid, or error if not
func ValidateOrderCounterpartyRoles(field string, value CounterpartyRole) error {
	if strings.TrimSpace(value) == "" {
		return validation.NewErrRecordIsMissingRequiredField(field)
	}
	if slice.Index(KnownCounterpartyRoles, value) == -1 {
		return validation.NewErrBadRecordFieldValue(field, fmt.Sprintf("unknown value: [%v], supported values: %s", value, strings.Join(KnownCounterpartyRoles, ", ")))
	}
	return nil
}

// Validate validates OrderCounterparty
func (v OrderCounterparty) Validate() error {
	if err := briefs4contactus.ValidateContactIDRecordField("contactID", v.ContactID, true); err != nil {
		return err
	}
	isCountryRequired := v.Role != CounterpartyRoleShip
	if err := with.ValidateCountryID("countryID", v.CountryID, isCountryRequired); err != nil {
		return err
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if err := ValidateOrderCounterpartyRoles("role", v.Role); err != nil {
		return err
	}
	if v.Parent != nil && v.Parent.ContactID == v.ContactID {
		return validation.NewErrBadRecordFieldValue("parent.contactID", "cannot be the same as contactID")
	}
	return nil
}
