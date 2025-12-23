package facade4retrospectus

import (
	"fmt"

	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// RetroRequest is a type alias
type RetroRequest = facade4meetingus.Request

var _ facade.Request = (*MoveRetroItemRequest)(nil)

// MoveRetroItemRequest parameters
type MoveRetroItemRequest struct {
	facade4meetingus.Request
	Item string                        `json:"item"`
	From dbo4retrospectus.TreePosition `json:"from"`
	To   dbo4retrospectus.TreePosition `json:"to"`
}

// Validate validates request
func (v *MoveRetroItemRequest) Validate() error {
	if err := v.Request.Validate(); err != nil {
		return err
	}
	if err := v.From.Validate(); err != nil {
		return err
	}
	if err := v.To.Validate(); err != nil {
		return err
	}
	if v.From == v.To {
		return fmt.Errorf("an attempt to move to the same position: %+v", v.To)
	}
	if v.To.Index > 100 {
		return validation.NewErrBadRequestFieldValue("to.index", fmt.Sprintf("too large value (>100): %v", v.To.Index))
	}
	return nil
}

// RetroDurations record
type RetroDurations struct {
	Feedback int `json:"feedback"`
	Review   int `json:"review"`
}

// StartRetrospectiveRequest request
type StartRetrospectiveRequest struct {
	RetroRequest
	DurationsInMinutes RetroDurations `json:"durationsInMinutes"`
}

// Validate validates
func (v *StartRetrospectiveRequest) Validate() error {
	if err := v.RetroRequest.Validate(); err != nil {
		return err
	}
	if v.DurationsInMinutes.Feedback < 0 {
		return validation.NewErrBadRequestFieldValue("durationsInMinutes.feedback", "should be positive")
	}
	if v.DurationsInMinutes.Review < 0 {
		return validation.NewErrBadRequestFieldValue("durationsInMinutes.review", "should be positive")
	}
	return nil
}

// RetrospectiveResponse response
type RetrospectiveResponse struct {
	ID   string                          `json:"id"`
	Data *dbo4retrospectus.Retrospective `json:"data"`
}

// FixCountsRequest request
type FixCountsRequest struct {
	RetroRequest
}

// Validate validates
func (v *FixCountsRequest) Validate() error {
	return v.RetroRequest.Validate()
}
