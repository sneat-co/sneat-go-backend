package dbo4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
)

// SegmentDates is a dates of the ContainerSegment.
type SegmentDates struct {
	Departs string `json:"departs,omitempty" firestore:"departs,omitempty"`
	Arrives string `json:"arrives,omitempty" firestore:"arrives,omitempty"`
}

// Validate returns nil if the SegmentDates is valid, otherwise returns the first error.
func (v SegmentDates) Validate() error {
	if v.Departs != "" {
		if _, err := validate.DateString(v.Departs); err != nil {
			return validation.NewErrBadRecordFieldValue("departs", err.Error())
		}
	}
	if v.Arrives != "" {
		if _, err := validate.DateString(v.Arrives); err != nil {
			return validation.NewErrBadRecordFieldValue("arrives", err.Error())
		}
	}
	if v.Arrives != "" && v.Departs != "" && v.Arrives < v.Departs {
		return validation.NewErrBadRecordFieldValue("end", "must be greater than or equal to start")
	}
	return nil
}

// SegmentCounterparty is a counterparty of the ContainerSegment.
type SegmentCounterparty struct {
	ContactID string           `json:"contactID" firestore:"contactID"`
	Role      CounterpartyRole `json:"role" firestore:"role"`
}

// Validate returns nil if the SegmentCounterparty is valid, otherwise returns the first error.
func (v SegmentCounterparty) Validate() error {
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := ValidateSegmentCounterpartyRole(v.Role); err != nil {
		return err // No need to wrap the error as we pass the field name to the validation function
	}
	return nil
}

// SegmentEndpoint is a segment endpoint of the ContainerSegment.
type SegmentEndpoint struct {
	SegmentCounterparty

	// ShippingPointID can be empty for a port for example,
	// but should be present in database for CounterpartyRoleDispatchPoint contacts
	ShippingPointID string `json:"shippingPointID,omitempty" firestore:"shippingPointID,omitempty"`
}

// Validate returns nil if the SegmentEndpoint is valid, otherwise returns the first error.
func (v SegmentEndpoint) Validate() error {
	if err := v.SegmentCounterparty.Validate(); err != nil {
		return err
	}
	switch v.Role {
	case CounterpartyRoleDispatchPoint:
		if strings.TrimSpace(v.ShippingPointID) == "" {
			return validation.NewErrBadRecordFieldValue("shippingPointID",
				fmt.Sprintf("must be present for %s endpoint", v.Role))
		}
	}
	return nil
}

// ContainerSegmentKey contains the contacts of the ContainerSegment.
type ContainerSegmentKey struct {
	ContainerID string          `json:"containerID" firestore:"containerID"`
	From        SegmentEndpoint `json:"from" firestore:"from"`
	To          SegmentEndpoint `json:"to" firestore:"to"`
}

// SegmentsFilter is a filter for ContainerSegments.
type SegmentsFilter struct {
	ContainerIDs        []string `json:"containerIDs" firestore:"containerIDs"`
	FromShippingPointID string   `json:"fromShippingPointID,omitempty"`
	ToShippingPointID   string   `json:"toShippingPointID,omitempty"`
	ByContactID         string   `json:"byContactID,omitempty"`
}

// Validate returns nil if the SegmentsFilter is valid, otherwise returns the first error.
func (v SegmentsFilter) Validate() error {
	if err := dbo4spaceus.ValidateShippingPointID(v.FromShippingPointID); err != nil {
		return validation.NewErrBadRequestFieldValue("fromShippingPointID", err.Error())
	}
	if err := dbo4spaceus.ValidateShippingPointID(v.ToShippingPointID); err != nil {
		return validation.NewErrBadRequestFieldValue("toShippingPointID", err.Error())
	}
	if err := briefs4contactus.ValidateContactIDRecordField("byContactID", v.ByContactID, false); err != nil {
		return err
	}
	return nil
}

// String representation of the ContainerSegmentKey.
func (v ContainerSegmentKey) String() string {
	return fmt.Sprintf("container=%s&from=%s:%s&to=%s:%s",
		v.ContainerID, v.From.Role, v.From.ContactID, v.To.Role, v.To.ContactID)
}

//func (v ContainerSegmentKey) Equals(v2 ContainerSegmentKey) bool { // not needed for now
//	return v.ContainerID == v2.ContainerID &&
//		v.From.ContactID == v2.From.ContactID &&
//		v.To.ContactID == v2.To.ContactID
//}

// Validate returns nil if the ContainerSegmentKey is valid, otherwise returns the first error.
func (v ContainerSegmentKey) Validate() error {
	if err := validateContainerID(v.ContainerID, "containerID"); err != nil {
		return err
	}
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("from", err.Error())
	}
	if err := v.To.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("to", err.Error())
	}
	if v.From.ContactID == v.To.ContactID {
		return validation.NewErrBadRecordFieldValue("from.id", "must be different from `to.id`")
	}
	return nil
}

// ContainerSegment is a segment of the OrderDbo and is referenced by WithSegments.
type ContainerSegment struct {
	ContainerSegmentKey
	ByContactID string        `json:"byContactID,omitempty" firestore:"byContactID,omitempty"`
	Dates       *SegmentDates `json:"dates,omitempty" firestore:"dates,omitempty"`
}

// Validate returns nil if the segment is valid, otherwise returns the first error.
func (v ContainerSegment) Validate() error {
	if err := v.ContainerSegmentKey.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("endpoints", err.Error())
	}
	if v.Dates != nil {
		if err := v.Dates.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("dates", err.Error())
		}
	}
	return nil
}

// WithSegments is a struct that contains a list of segments.
type WithSegments struct {
	Segments []*ContainerSegment `json:"segments,omitempty" firestore:"segments,omitempty"`
}

// DeleteSegments deletes segments from an order.
func (v WithSegments) DeleteSegments(segmentsToDelete []*ContainerSegment) (segments []*ContainerSegment) {
	segments = make([]*ContainerSegment, 0, len(v.Segments)-len(segmentsToDelete))
	for _, segment := range v.Segments {
		for _, s2d := range segmentsToDelete {
			if s2d.ContainerSegmentKey == segment.ContainerSegmentKey {
				continue
			}
			segments = append(segments, segment)
		}
	}
	return segments
}

// GetSegmentsByFilter returns the segments that match the given filter.
func (v WithSegments) GetSegmentsByFilter(filter SegmentsFilter) (segment []*ContainerSegment) {
	segments := make([]*ContainerSegment, 0, len(v.Segments))
	for _, s := range v.Segments {
		if (len(filter.ContainerIDs) == 0 || slice.Index(filter.ContainerIDs, s.ContainerID) >= 0) &&
			(filter.FromShippingPointID == "" || s.From.ShippingPointID == filter.FromShippingPointID) &&
			(filter.ToShippingPointID == "" || s.To.ShippingPointID == filter.ToShippingPointID) &&
			(filter.ByContactID == "" || s.ByContactID == filter.ByContactID) {
			segments = append(segments, s)
		}
	}
	return segments
}

// GetSegmentByKey returns a segment by the given key.
func (v WithSegments) GetSegmentByKey(k ContainerSegmentKey) *ContainerSegment {
	if k.ContainerID == "" {
		panic("container ID is required")
	}
	for _, s := range v.Segments {
		if s.ContainerID != k.ContainerID {
			continue
		}
		if k.From.ShippingPointID != "" && s.From.ShippingPointID == k.From.ShippingPointID {
			return s
		}
		if k.To.ShippingPointID != "" && s.To.ShippingPointID == k.To.ShippingPointID {
			return s
		}
		if s.ContainerSegmentKey == k {
			return s
		}
	}
	return nil
}

// Updates returns updates for order segments.
func (v WithSegments) Updates() []dal.Update {
	if len(v.Segments) == 0 {
		return []dal.Update{
			{Field: "segments", Value: dal.DeleteField},
		}
	}
	return []dal.Update{
		{Field: "segments", Value: v.Segments},
	}
}

// Validate returns nil if all segments are valid, otherwise returns the first error.
func (v WithSegments) Validate() error {
	for i, s := range v.Segments {
		if err := s.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("segments[%v]", i), err.Error())
		}
		for j, s2 := range v.Segments {
			if j == i {
				continue
			}
			if s2.ContainerSegmentKey == s.ContainerSegmentKey {
				return validation.NewErrBadRecordFieldValue("segments", fmt.Sprintf("duplicate segment keys at indexes %d & %d: %+v", i, j, s.ContainerSegmentKey))
			}
		}
	}
	return nil
}

func (v WithSegments) validateOrder(order OrderDbo) error {
	for i, segment := range v.Segments {
		if err := validateOrderSegment(order, segment); err != nil {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("segments[%v]", i),
				fmt.Errorf("invalid segment with key=%s: %w", segment, err).Error())
		}
	}
	return nil
}

func validateOrderSegment(order OrderDbo, segment *ContainerSegment) (err error) {
	if _, container := order.GetContainerByID(segment.ContainerID); container == nil {
		return validation.NewErrBadRecordFieldValue("containerID", fmt.Sprintf("container with ID=[%s] does not exist", segment.ContainerID))
	}
	//var fromPoint, toPoint *OrderShippingPoint
	if err = validateSegmentEndpoint(order, segment, "from", segment.ContainerSegmentKey.From); err != nil {
		return err
	}
	if err = validateSegmentEndpoint(order, segment, "to", segment.ContainerSegmentKey.To); err != nil {
		return err
	}
	if err := briefs4contactus.ValidateContactIDRecordField("byContactID", segment.ByContactID, false); err != nil {
		return err
	}
	return nil
}

func validateSegmentCounterparty(order OrderDbo, segmentCounterparty SegmentCounterparty) error {
	if _, c := order.GetCounterpartyByRoleAndContactID(segmentCounterparty.Role, segmentCounterparty.ContactID); c == nil {
		return validation.NewErrBadRecordFieldValue("contactID", fmt.Sprintf("referenced contact is not present in order: [%s:%s]", segmentCounterparty.ContactID, segmentCounterparty.Role))
	}
	return nil
}

func validateSegmentEndpoint(order OrderDbo, segment *ContainerSegment, field string, endpoint SegmentEndpoint) (err error) {
	if err = validateSegmentCounterparty(order, endpoint.SegmentCounterparty); err != nil {
		return validation.NewErrBadRecordFieldValue(field, err.Error())
	}
	switch endpoint.Role {
	case CounterpartyRolePortFrom, CounterpartyRolePortTo:
		if endpoint.ShippingPointID != "" {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("%s.shippingPointID", field),
				fmt.Sprintf("segment counterparty with role=[%s] should not reference an endpoint",
					endpoint.Role))
		}
	default:
		if endpoint.ShippingPointID == "" {
			return fmt.Errorf("segment endpoint with role=[%s], contactID=[%s] should reference a shipping point: %w",
				endpoint.Role, endpoint.ContactID,
				validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("%s.shippingPointID", field)))
		} else if _, shippingPoint := order.GetShippingPointByID(endpoint.ShippingPointID); shippingPoint == nil {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("%s.shippingPointID", field),
				fmt.Sprintf("referenced shippingPoint is not present in order: [%+v] ", endpoint))
		}
		containerPoint := order.GetContainerPoint(segment.ContainerID, endpoint.ShippingPointID)
		if containerPoint == nil {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("%s.shippingPointID", field),
				fmt.Sprintf("segment has no relevant container point: containerID=%s, [%+v]", segment.ContainerID, endpoint))
		}
		validateDates := func(dates *SegmentDates, field, containerPointDate string, getSegmentDate func() string) error {
			if dates == nil {
				if containerPointDate != "" {
					return validation.NewErrBadRecordFieldValue("dates."+field, "segment has no dates but container point has the date")
				}
			} else {
				segmentDate := getSegmentDate()
				if containerPointDate != segmentDate {
					return validation.NewErrBadRecordFieldValue(
						"dates."+field,
						fmt.Sprintf("segment departure date differs from container point '[%sDate]': segment.dates.%s=[%v] != containerPoint.%vDate=[%v] ", field, field, segmentDate, field, containerPointDate))
				}
			}
			return nil
		}
		switch field {
		case "from":
			if containerPoint.Departure != nil {
				if err := validateDates(segment.Dates, "departure", containerPoint.Departure.ScheduledDate, func() string {
					return segment.Dates.Departs
				}); err != nil {
					return err
				}
			}
		case "to":
			if containerPoint.Arrival != nil {
				if err := validateDates(segment.Dates, "arrival", containerPoint.Arrival.ScheduledDate, func() string {
					return segment.Dates.Arrives
				}); err != nil {
					return err
				}
			}
		default:
			panic("unknown value for parameter 'field': " + field)
		}
	}

	return nil
}
