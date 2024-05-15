package facade4scrumus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/models4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------------------------------
// Request
var _ facade.Request = (*facade4meetingus.Request)(nil)

// ---------------------------------------------------------------------------------------------------
// TaskRequest
var _ facade.Request = (*TaskRequest)(nil)

// TaskRequest request
type TaskRequest struct {
	facade4meetingus.Request
	ContactID string `json:"contactID"`
	Type      string `json:"type"`
	Task      string `json:"task"`
}

// Validate validates request
func (v *TaskRequest) Validate() error {
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("contactID")
	}
	if strings.TrimSpace(v.Type) == "" {
		return validation.NewErrRecordIsMissingRequiredField("Role")
	}
	if strings.TrimSpace(v.Task) == "" {
		return validation.NewErrRecordIsMissingRequiredField("Task")
	}
	return v.Request.Validate()
}

// ---------------------------------------------------------------------------------------------------
// AddTaskRequest
var _ facade.Request = (*AddTaskRequest)(nil)

// AddTaskRequest request
type AddTaskRequest struct {
	TaskRequest
	Title string `json:"title"`
}

// Validate validates request
func (v *AddTaskRequest) Validate() error {
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("MemberDto")
	}
	if strings.TrimSpace(v.Type) == "" {
		return validation.NewErrRecordIsMissingRequiredField("Role")
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("Title")
	}
	date, err := time.Parse("2006-01-02", v.MeetingID)
	if err != nil {
		return validation.NewErrBadRequestFieldValue("MeetingID", err.Error())
	}
	now := time.Now()
	if date.After(now.AddDate(0, 1, 0)) {
		return validation.NewErrBadRequestFieldValue("MeetingID", "not allowed to create tasks more then 1 month in advance")
	}
	if date.Before(now.AddDate(0, 0, -7)) {
		return validation.NewErrBadRequestFieldValue("MeetingID", "not allowed to add tasks for api4meetingus more then 7 days old")
	}
	return v.Request.Validate()
}

var _ facade.Request = (*SetMetricRequest)(nil)

// SetMetricRequest request
type SetMetricRequest struct {
	facade4meetingus.Request
	Metric string `json:"metric"`
	Member string `json:"members,omitempty"`
	models4scrumus.MetricValue
}

// Validate validates request
func (v *SetMetricRequest) Validate() error {
	if strings.TrimSpace(v.Metric) == "" {
		return validation.NewErrRecordIsMissingRequiredField("metric")
	}
	if err := v.MetricValue.Validate(); err != nil {
		return err
	}
	return v.Request.Validate()
}

// AddTaskResponse response
type AddTaskResponse struct {
	// Task      string    `json:"id"`
	Created time.Time `json:"created"`
}

var _ facade.Request = (*DeleteTaskRequest)(nil)

// DeleteTaskRequest request
type DeleteTaskRequest = TaskRequest

var _ facade.Request = (*ThumbUpRequest)(nil)

// ThumbUpRequest request
type ThumbUpRequest struct {
	TaskRequest
	Value bool `json:"value"`
}

// Validate validates request
func (v *ThumbUpRequest) Validate() error {
	return v.TaskRequest.Validate()
}

var _ facade.Request = (*ReorderTaskRequest)(nil)

// ReorderTaskRequest request
type ReorderTaskRequest struct {
	TaskRequest
	Len    int    `json:"len"`
	From   int    `json:"from"`
	To     int    `json:"to"`
	After  string `json:"after"`
	Before string `json:"before"`
}

// Validate validates request
func (v *ReorderTaskRequest) Validate() error {
	if err := v.TaskRequest.Validate(); err != nil {
		return err
	}
	if v.Len == 0 {
		return validation.NewErrBadRequestFieldValue("Len", "can't perform reorder on empty list")
	}
	if v.From < 0 {
		return validation.NewErrBadRequestFieldValue("From", fmt.Sprintf("should be >= 0, got %d", v.From))
	}
	if v.To < 0 {
		return validation.NewErrBadRequestFieldValue("To", fmt.Sprintf("should be >= 0, got %d", v.To))
	}
	if v.From >= v.Len {
		return validation.NewErrBadRequestFieldValue("From", fmt.Sprintf("should be < len=%d, got %d", v.Len, v.From))
	}
	if v.To >= v.Len {
		return validation.NewErrBadRequestFieldValue("To", fmt.Sprintf("field 'to' should be >= len=%d, got %d", v.Len, v.To))
	}
	if v.From == v.To {
		return validation.NewErrBadRequestFieldValue("From == To", fmt.Sprintf("can't be equal, got %d", v.From))
	}
	if v.After == v.Before {
		return validation.NewErrBadRequestFieldValue("After == Before", fmt.Sprintf("can't be equal, got: %s", v.After))
	}
	if v.After == v.Task {
		return validation.NewErrBadRequestFieldValue("After == Task", fmt.Sprintf("can't be equal, got: %s", v.Task))
	}
	if v.Before == v.Task {
		return validation.NewErrBadRequestFieldValue("After == Task", fmt.Sprintf("can't be equal, got: %s", v.Task))
	}
	if v.To == 0 && v.After != "" {
		return validation.NewErrBadRequestFieldValue("After", fmt.Sprintf("can't have this field when moving to then beginning of list, got %s", v.After))
	}
	if v.To == v.Len-1 && v.Before != "" {
		return validation.NewErrBadRequestFieldValue("Before", fmt.Sprintf("can't have this field when moving to the end of list, got %s", v.Before))
	}
	return nil
}
