package dbo4logist

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"regexp"
	"strings"
)

const OrdersCollection = "orders"

// OrderDbo is a DTO for an order
type OrderDbo struct {
	dbmodels.WithModified
	dbmodels.WithSpaceID  // Owner of an order
	dbmodels.WithSpaceIDs // Spaces that can access an order
	dbmodels.WithUserIDs  // Users that can access an order
	with.CountryIDsField  // Countries related to an order
	with.DatesFields      // Dates related to an order
	with.KeysField        // Keys related to an order
	OrderBase             // Order fields that are stored in a record and also passed in a ValidatingRequest
	WithOrderContacts
	WithOrderContainers // Shipping containers included to an order
	WithShippingPoints  // Shipping points included to an order
	WithContainerPoints // Shipping points related to shipping containers
	WithSegments        //Order segments

	// Planned steps for an order (e.g. shipping dates)
	Steps []*OrderStep `json:"steps,omitempty" firestore:"steps,omitempty"`
}

// Validate validates OrderDbo
func (v *OrderDbo) Validate() error {
	//{ // This block should be run just in case before validation
	//	v.UpdateKeys()
	//	if err := v.WithKeys.Validate(); err != nil {
	//		return err
	//	}
	//	v.UpdateDates()
	//	if err := v.WithDates.Validate(); err != nil {
	//		return err
	//	}
	//}
	if err := v.OrderBase.Validate(); err != nil {
		return err
	}
	//
	if err := v.WithModified.Validate(); err != nil {
		return err
	}
	if err := v.CountryIDsField.Validate(); err != nil {
		return err
	}
	if err := v.WithSpaceID.Validate(); err != nil {
		return err
	}
	if err := v.WithSpaceIDs.Validate(); err != nil {
		return err
	}
	if slice.Index(v.SpaceIDs, v.SpaceID) < 0 {
		return validation.NewErrBadRecordFieldValue("field", " should contain value of teamID field")
	}
	if err := v.WithUserIDs.Validate(); err != nil {
		return err
	}
	if err := v.DatesFields.Validate(); err != nil {
		return err
	}
	if err := v.WithSegments.Validate(); err != nil {
		return err
	}

	if err := v.WithOrderContacts.Validate(); err != nil {
		return err
	}
	//
	if err := v.WithOrderContainers.validateOrder(*v); err != nil {
		return err
	}
	if err := v.WithShippingPoints.validateOrder(*v); err != nil {
		return err
	}
	if err := v.WithContainerPoints.validateOrder(*v); err != nil {
		return err
	}
	if err := v.validateDtoContacts(); err != nil {
		return err
	}
	if err := v.validateDtoShippingPoints(); err != nil {
		return err
	}
	if err := v.validateDtoCounterparties(); err != nil {
		return err
	}
	if err := v.validateDtoKeys(); err != nil {
		return err
	}
	if err := v.WithSegments.validateOrder(*v); err != nil {
		return err
	}

	if v.Route != nil {
		if slice.Index(v.CountryIDs, v.Route.Origin.CountryID) < 0 {
			return validation.NewErrBadRecordFieldValue("countryIDs", "missing value for origin countryID="+v.Route.Origin.CountryID)
		}
		if slice.Index(v.CountryIDs, v.Route.Destination.CountryID) < 0 {
			return validation.NewErrBadRecordFieldValue("countryIDs", "missing value for destination countryID="+v.Route.Destination.CountryID)
		}
	}
	if v.Route != nil {
		for _, t := range v.Route.TransitPoints {
			if slice.Index(v.CountryIDs, t.CountryID) < 0 {
				return validation.NewErrBadRecordFieldValue("countryIDs", "missing value for transit countryID="+t.CountryID)
			}
		}
	}
	for i, step := range v.Steps {
		if err := step.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("steps[%d]", i), err.Error())
		}
	}

	for _, cp := range v.Counterparties {
		_, contact := v.GetContactByID(cp.ContactID)
		if contact == nil {
			return validation.NewErrBadRecordFieldValue("contactID", "not found in `contacts` by ContactID="+cp.ContactID)
		}
		if cp.Title != contact.Title {
			return validation.NewErrBadRecordFieldValue("title", fmt.Sprintf(`"%s" does not match contact title "%s"`, cp.Title, contact.Title))
		}
		//if cp.CountryID != contact.Address.CountryID {
		//	return validation.NewErrBadRecordFieldValue("title", fmt.Sprintf(`"%s" does not match contact.Address.CountryID="%s"`, cp.CountryID, contact.Address.CountryID))
		//}
		if contact.ParentID != "" {
			if cp.Parent == nil {
				return validation.NewErrRecordIsMissingRequiredField("parent")
			}
			if cp.Parent.ContactID != contact.ParentID {
				return validation.NewErrBadRecordFieldValue("parent.ContactID", fmt.Sprintf(`"%s" does not match contact.ParentID="%s"`, cp.Parent.ContactID, contact.ParentID))
			}
		}
		switch cp.Role {
		case
			CounterpartyRoleDispatchPoint,
			CounterpartyRoleReceivePoint,
			CounterpartyRolePickPoint,
			CounterpartyRoleDropPoint,
			CounterpartyRoleShip:
			if strings.TrimSpace(contact.ParentID) == "" {
				return validation.NewErrRecordIsMissingRequiredField("parentContactID")
			}
		}
	}
	return nil
}

// UpdateDates updates dates related to an order
func (v *OrderDbo) UpdateDates() {
	for _, cp := range v.ShippingPoints {
		v.updateDatesFromShippingPoint(cp)
	}
}

func (v *OrderDbo) updateDatesFromShippingPoint(p *OrderShippingPoint) {
	if p.ScheduledStartDate != "" && slice.Index(v.Dates, p.ScheduledStartDate) < 0 {
		v.DatesFields.AddDate(p.ScheduledStartDate)
	}
	if p.ScheduledEndDate != "" && slice.Index(v.Dates, p.ScheduledEndDate) < 0 {
		v.DatesFields.AddDate(p.ScheduledEndDate)
	}
}

// UpdateCalculatedFields updates calculated fields in OrderDbo - calls UpdateKeys and UpdateDates
func (v *OrderDbo) UpdateCalculatedFields() {
	v.UpdateKeys()
	v.UpdateDates()
}

// UpdateKeys updates keys field in OrderDbo
func (v *OrderDbo) UpdateKeys() {
	keys := make([]string, 0, len(v.CountryIDs)+len(v.Counterparties)*2)
	for _, counterparty := range v.Counterparties {
		contactKey := getContactKey(counterparty.ContactID)
		countryKey := getCountryKey(counterparty.CountryID)
		refNumberKey := getRefNumberKey(counterparty.RefNumber)
		{ // Add individual keys
			if slice.Index(keys, contactKey) < 0 {
				keys = append(keys, contactKey)
			}
			if slice.Index(keys, countryKey) < 0 {
				keys = append(keys, countryKey)
			}
			if refNumberKey != "" && slice.Index(keys, refNumberKey) < 0 {
				keys = append(keys, refNumberKey)
			}

		}

		countryPlusCounterpartyKey := fmt.Sprintf("%s&%s", countryKey, contactKey)
		if slice.Index(keys, countryPlusCounterpartyKey) < 0 {
			keys = append(keys, countryPlusCounterpartyKey)
		}

		if refNumberKey != "" { // Add refNumber to country, counterparty and combined keys
			withRefNumberKey := fmt.Sprintf("%s&%s", countryKey, refNumberKey)
			if slice.Index(keys, withRefNumberKey) < 0 {
				keys = append(keys, withRefNumberKey)
			}
			withRefNumberKey = fmt.Sprintf("%s&%s", contactKey, refNumberKey)
			if slice.Index(keys, withRefNumberKey) < 0 {
				keys = append(keys, withRefNumberKey)
			}
			withRefNumberKey = fmt.Sprintf("%s&%s", countryPlusCounterpartyKey, refNumberKey)
			if slice.Index(keys, withRefNumberKey) < 0 {
				keys = append(keys, withRefNumberKey)
			}
		}
	}
	v.Keys = keys
}

const contactKeyPrefix = "contact="
const countryKeyPrefix = "country="
const refNumberKeyPrefix = "refNumber="

func getRefNumberKey(refNumber string) string {
	if refNumber == "" {
		return ""
	}
	return refNumberKeyPrefix + refNumber
}

func getCountryKey(countryID string) string {
	return countryKeyPrefix + countryID
}

func getContactKey(contactID string) string {
	return contactKeyPrefix + contactID
}

func isContactKey(key string) bool {
	return strings.Contains(key, contactKeyPrefix)
}

var contactIDRegexp = regexp.MustCompile(contactKeyPrefix + `(.+?)(&|$)`)

func getContactIdFromOrderKey(key string) string {
	match := contactIDRegexp.FindStringSubmatch(key)
	return match[1]
}

// NewOrderShippingPointID generates a new ContactID for a shipping point
func (v *WithShippingPoints) NewOrderShippingPointID() string {
	var i int
	for {
		id := random.ID(2)
		for _, p := range v.ShippingPoints {
			if p.ID == id {
				if i++; i > 100 {
					panic("Too many attempts to generate new ShippingPointID")
				}
				continue
			}
		}
		return id
	}
}

// DeleteShippingPoint deletes a shipping point from an order
func (v *WithShippingPoints) DeleteShippingPoint(pointType, contactID string) (deletedShippingPointID string, shippingPoints []*OrderShippingPoint) {
	shippingPoints = make([]*OrderShippingPoint, 0, len(v.ShippingPoints))
	for _, shippingPoint := range v.ShippingPoints {
		if shippingPoint.Location.ContactID == contactID {
			deletedShippingPointID = shippingPoint.ID
			continue
		}
		shippingPoints = append(shippingPoints, shippingPoint)
	}
	v.ShippingPoints = shippingPoints
	return deletedShippingPointID, shippingPoints
}

func (v *OrderDbo) validateDtoCounterparties() error {
	for _, counterparty := range v.Counterparties {
		key := getContactKey(counterparty.ContactID)
		if slice.Index(v.Keys, key) < 0 {
			return fmt.Errorf("key [%s] is not in `keys`", key)
		}
	}
	return nil
}

func (v *OrderDbo) validateDtoContacts() error {
	for i, contact := range v.Contacts {
		_, counterparty := v.GetCounterpartyByContactID(contact.ID)
		if counterparty == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contacts[%d].id", i), fmt.Sprintf("no entry in `counterparties` for orphaned contact{ContactID=%s, Title=%s}", contact.ID, contact.Title))
		}
		key := getContactKey(contact.ID)
		if slice.Index(v.Keys, key) < 0 {
			return fmt.Errorf("key [%s] is not in `keys`", key)
		}
	}
	return nil
}

func (v *OrderDbo) validateDtoShippingPoints() error {
	for i, sp := range v.ShippingPoints {
		for _, task := range sp.Tasks {
			var role CounterpartyRole
			switch task {
			case ShippingPointTaskLoad:
				role = CounterpartyRoleDispatchPoint
			case ShippingPointTaskUnload:
				role = CounterpartyRoleReceivePoint
			case ShippingPointTaskDrop:
				role = CounterpartyRoleDropPoint
			case ShippingPointTaskPick:
				role = CounterpartyRolePickPoint
			}
			if role != "" && sp.Location != nil {
				if _, counterparty := v.GetCounterpartyByRoleAndContactID(role, sp.Location.ContactID); counterparty == nil {
					return validation.NewErrBadRecordFieldValue(fmt.Sprintf("shippingPoints[%d]", i),
						fmt.Sprintf("no entry in `counterparties` with role=%s and contactID=%s for shipping point with ContactID=%s",
							role, sp.Location.ContactID, sp.ID))
				}
			}
		}

	}
	return nil
}

func (v *OrderDbo) validateDtoKeys() error {
	for _, key := range v.Keys {
		if isContactKey(key) {
			contactID := getContactIdFromOrderKey(key)
			if i, _ := v.GetContactByID(contactID); i < 0 {
				return fmt.Errorf("key [%s] is referecing entry that is not in `contacts` field", key)
			}
			if i, _ := v.GetCounterpartyByContactID(contactID); i < 0 {
				return fmt.Errorf("key [%s] is referecing entry that is not in `counterparties` field", key)
			}
		}
	}
	return nil
}
