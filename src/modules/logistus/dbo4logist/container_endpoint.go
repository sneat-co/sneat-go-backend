package dbo4logist

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

type EndpointSide = string

const (
	EndpointSideArrival   EndpointSide = "arrival"
	EndpointSideDeparture EndpointSide = "departure"
)

// ContainerEndpoints holds dates of a container.
type ContainerEndpoints struct {
	Arrival   *ContainerEndpoint `json:"arrival,omitempty" firestore:"arrival,omitempty"`
	Departure *ContainerEndpoint `json:"departure,omitempty" firestore:"departure,omitempty"`
}

// ContainerEndpoint represents a shipping arrival or departure point of a container.
type ContainerEndpoint struct {
	ByContactID   string `json:"byContactID,omitempty" firestore:"byContactID,omitempty"`
	ScheduledDate string `json:"scheduledDate,omitempty" firestore:"scheduledDate,omitempty"`
	ScheduledTime string `json:"scheduledTime,omitempty" firestore:"scheduledTime,omitempty"`
	ActualDate    string `json:"actualDate,omitempty" firestore:"actualDate,omitempty"`
	ActualTime    string `json:"actualTime,omitempty" firestore:"actualTime,omitempty"`
}

func (v ContainerEndpoint) IsEmpty() bool {
	return v.ScheduledDate == "" && v.ActualDate == "" && v.ByContactID == ""
}

// String returns a string representation of the ContainerEndpoint.
func (v ContainerEndpoint) String() string {
	return fmt.Sprintf("ContainerEndpoint{ByContactID=%s,ScheduledDate=%s,ActualDate=%v}", v.ByContactID, v.ScheduledDate, v.ActualDate)
}

// Validate returns an error if the ContainerEndpoint is invalid.
func (v ContainerEndpoint) Validate() error {
	if v.ScheduledDate != "" {
		if _, err := validate.DateString(v.ScheduledDate); err != nil {
			return validation.NewErrBadRecordFieldValue("scheduledDate", err.Error())
		}
	}
	if v.ActualDate != "" {
		if _, err := validate.DateString(v.ActualDate); err != nil {
			return validation.NewErrBadRecordFieldValue("actualDate", err.Error())
		}
	}
	if strings.TrimSpace(v.ByContactID) != v.ByContactID {
		return validation.NewErrBadRecordFieldValue("byContactID", "should not start or end with space")
	}
	return nil
}

// Validate returns an error if the ContainerEndpoints is invalid.
func (v ContainerEndpoints) Validate() error {
	if v.Arrival != nil {
		if err := v.Arrival.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("arrival", err.Error())
		}
	}
	if v.Departure != nil {
		if err := v.Departure.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("departure", err.Error())
		}
	}
	return nil
}

func (v ContainerEndpoints) Strings() string {
	return fmt.Sprintf("ContainerEndpoints{Arrival=%s,Departure=%s}", v.Arrival, v.Departure)
}
