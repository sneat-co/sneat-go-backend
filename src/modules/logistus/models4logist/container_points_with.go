package models4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/validation"
)

// WithContainerPoints represents containers with their shipping points.
type WithContainerPoints struct {
	ContainerPoints []*ContainerPoint `json:"containerPoints,omitempty" firestore:"containerPoints,omitempty"`
}

// Updates populates update instructions for DAL
func (v *WithContainerPoints) Updates() []dal.Update {
	if len(v.ContainerPoints) == 0 {
		return []dal.Update{
			{Field: "containerPoints", Value: dal.DeleteField},
		}
	}
	return []dal.Update{
		{Field: "containerPoints", Value: v.ContainerPoints},
	}
}

// Validate returns an error if the WithContainerPoints is invalid.
func (v *WithContainerPoints) Validate() error {
	panic("not implemented, call `WithContainerPoints.validateOrder(order OrderDto)` instead")
}

func (v *WithContainerPoints) validateOrder(order OrderDto) error {
	for i, p := range v.ContainerPoints {
		field := fmt.Sprintf("containerPoints[%v]", i)
		if err := p.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(field, err.Error())
		}
		for j, p2 := range v.ContainerPoints {
			if j != i && p2.ContainerID == p.ContainerID && p2.ShippingPointID == p.ShippingPointID {
				msg := fmt.Sprintf("duplicate container point ID for items %d & %d: [containerID=%v & shippingPointID=%v]", i, j, p.ContainerID, p.ShippingPointID)
				return validation.NewErrBadRecordFieldValue(field, msg)
			}
			{ // Validate dates for overlapping
				overlappingError := func(side EndpointSide, dateField string) error {
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("%s.%s.%s", field, side, dateField),
						fmt.Sprintf("overlaps with dates of container point at index %d", j))
				}
				if p2.Arrival != nil && p2.Arrival.ScheduledDate != "" && p2.Departure != nil && p2.Departure.ScheduledDate != "" {
					if p.Arrival != nil && p.Arrival.ScheduledDate > p2.Arrival.ScheduledDate && p.Arrival.ScheduledDate < p2.Departure.ScheduledDate {
						return overlappingError(EndpointSideArrival, "scheduledDate")
					}
					if p.Departure != nil && p.Departure.ScheduledDate > p2.Arrival.ScheduledDate && p.Departure.ScheduledDate < p2.Departure.ScheduledDate {
						return overlappingError(EndpointSideDeparture, "scheduledDate")
					}
				}
				if p2.Arrival != nil && p2.Arrival.ActualDate != "" && p2.Departure != nil && p2.Departure.ActualDate != "" {
					if p.Arrival != nil && p.Arrival.ActualDate > p2.Arrival.ActualDate && p.Arrival.ActualDate < p2.Departure.ActualDate {
						return overlappingError(EndpointSideArrival, "actualDate")
					}
					if p.Departure != nil && p.Departure.ActualDate > p2.Arrival.ActualDate && p.Departure.ActualDate < p2.Departure.ActualDate {
						return overlappingError(EndpointSideDeparture, "actualDate")
					}
				}
			}
		}
		if _, c := order.GetContainerByID(p.ContainerID); c == nil {
			return validation.NewErrBadRecordFieldValue(field, fmt.Sprintf("container not found by id=[%v]", p.ContainerID))
		}
		if _, sp := order.GetShippingPointByID(p.ShippingPointID); sp == nil {
			return validation.NewErrBadRecordFieldValue(field, fmt.Sprintf("shipping point not found by id=[%v]", p.ShippingPointID))
		}
		if p.Arrival != nil && p.Arrival.ByContactID != "" {
			if _, c := order.GetCounterpartyByRoleAndContactID(CounterpartyRoleTrucker, p.Arrival.ByContactID); c == nil {
				return validation.NewErrBadRecordFieldValue(field+".arrival.byContactID", fmt.Sprintf("trucker counterparty not found in order by contact id=[%v]", p.Arrival.ByContactID))
			}
		}
		if p.Departure != nil && p.Departure.ByContactID != "" {
			if _, c := order.GetCounterpartyByRoleAndContactID(CounterpartyRoleTrucker, p.Departure.ByContactID); c == nil {
				return validation.NewErrBadRecordFieldValue(field+".arrival.byContactID", fmt.Sprintf("trucker counterparty not found in order by contact id=[%v]", p.Departure.ByContactID))
			}
		}
	}
	return nil
}

// RemoveContainerPointsByShippingPointID removes container points by shipping point ContactID.
func (v *WithContainerPoints) RemoveContainerPointsByShippingPointID(shippingPointID string) (containerPoints []*ContainerPoint) {
	return v.removeContainerPoints(func(p *ContainerPoint) bool { return p.ShippingPointID == shippingPointID })
}

func (v *WithContainerPoints) GetContainerPoint(containerID, shippingPointID string) *ContainerPoint {
	if containerID == "" {
		panic("containerID is a required parameter and is empty")
	}
	if shippingPointID == "" {
		panic("shippingPointID is a required parameter and is empty")
	}
	for i, p := range v.ContainerPoints {
		if p == nil {
			panic(fmt.Sprintf("nil container point at index %d out of %d items", i, len(v.ContainerPoints)))
		}
		if p.ContainerID == containerID && p.ShippingPointID == shippingPointID {
			return p
		}
	}
	return nil
}

// RemoveContainerPointsByContainerID removes container points by container ContactID.
func (v *WithContainerPoints) RemoveContainerPointsByContainerID(containerID string) (containerPoints []*ContainerPoint) {
	return v.removeContainerPoints(func(p *ContainerPoint) bool { return p.ContainerID == containerID })
}

func (v *WithContainerPoints) removeContainerPoints(match func(p *ContainerPoint) bool) (containerPoints []*ContainerPoint) {
	containerPoints = make([]*ContainerPoint, 0, len(v.ContainerPoints))
	for _, p := range v.ContainerPoints {
		if match(p) {
			continue
		}
		containerPoints = append(containerPoints, p)
	}
	if len(containerPoints) == 0 {
		v.ContainerPoints = containerPoints
	} else {
		v.ContainerPoints = nil
	}
	return containerPoints
}
