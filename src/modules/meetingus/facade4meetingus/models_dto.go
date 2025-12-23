package facade4meetingus

import (
	"strings"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// Request base for api4meetingus requests
type Request struct {
	dto4spaceus.SpaceRequest
	MeetingID string `json:"meetingID"`
}

// Validate validates api4meetingus requests
func (v *Request) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.MeetingID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("api4meetingus")
	}
	return nil
}

var _ facade.Request = (*ToggleTimerRequest)(nil)

// ToggleTimerRequest toggle timer parameters
type ToggleTimerRequest struct {
	Request
	Operation string `json:"operation"`
	Member    string `json:"members,omitempty"`
}

// Validate validate request
func (v ToggleTimerRequest) Validate() error {
	if strings.TrimSpace(v.Operation) == "" {
		return validation.NewErrRecordIsMissingRequiredField("operation")
	}
	return v.Request.Validate()
}

// ToggleTimerResponse response
type ToggleTimerResponse struct {
	Timer   *dbo4meetingus.Timer `json:"timer,omitempty"`
	Message string               `json:"message,omitempty"`
}

// Validate validates response
func (v *ToggleTimerResponse) Validate() error {
	if err := v.Timer.Validate(); err != nil {
		return err
	}
	return nil
}
