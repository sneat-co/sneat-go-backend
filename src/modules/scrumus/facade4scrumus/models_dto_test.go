package facade4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/strongo/validation"
	"testing"
	"time"
)

func TestAddTaskRequest_Validate(t *testing.T) {
	const validDateFormat = "2006-01-02"
	now := time.Now()
	validRequest := AddTaskRequest{
		TaskRequest: TaskRequest{
			Request: facade4meetingus.Request{
				SpaceRequest: dto4teamus.SpaceRequest{
					SpaceID: "space1",
				},
				MeetingID: now.Format(validDateFormat),
			},
			ContactID: "member1",
			Type:      "todo",
			Task:      "id1",
		},
		Title: "Test task",
	}
	t.Run("valid request", func(t *testing.T) {
		if err := validRequest.Validate(); err != nil {
			t.Errorf("A valid request failed validation: %v", err)
		}
	})
	t.Run("missing", func(t *testing.T) {
		t.Run("space", func(t *testing.T) {
			request := validRequest
			request.SpaceID = " "
			if err := request.Validate(); err == nil {
				t.Error("Should return error for empty team, got err == nil")
			} else {
				validation.MustBeFieldError(t, err, "space")
			}
		})
		t.Run("members", func(t *testing.T) {
			request := validRequest
			request.ContactID = " "
			if err := request.Validate(); err == nil {
				t.Error("Should return error for empty members, got")
			} else {
				validation.MustBeFieldError(t, err, "MemberDto")
			}
		})
		t.Run("date", func(t *testing.T) {
			request := validRequest
			request.MeetingID = " "
			if err := request.Validate(); err == nil {
				t.Error("Should return error for empty date")
			} else {
				validation.MustBeFieldError(t, err, "MeetingID")
			}
		})
	})
	t.Run("date", func(t *testing.T) {
		t.Run("format", func(t *testing.T) {
			request := validRequest
			request.MeetingID = now.Format("20060102")
			if err := request.Validate(); err == nil {
				t.Error("Should return error for date without dashes")
			} else {
				validation.MustBeFieldError(t, err, "MeetingID")
			}
		})
		t.Run("too_big", func(t *testing.T) {
			request := validRequest
			request.MeetingID = now.AddDate(0, 2, 1).Format(validDateFormat)
			if err := request.Validate(); err == nil {
				t.Error("Should return error for date too far in future")
			} else {
				validation.MustBeFieldError(t, err, "MeetingID")
			}
		})
		t.Run("to_small", func(t *testing.T) {
			request := validRequest
			request.MeetingID = now.AddDate(0, 02, 1).Format(validDateFormat)
			if err := request.Validate(); err == nil {
				t.Error("Should return error for date too far in past")
			} else {
				validation.MustBeFieldError(t, err, "MeetingID")
			}
		})
	})
}

func TestTaskRequest_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		request := TaskRequest{}
		if err := request.Validate(); err == nil {
			t.Fatalf("should return error for empty request")
		}
	})
	t.Run("valid", func(t *testing.T) {
		request := TaskRequest{
			Task:      "task1",
			Type:      "done",
			ContactID: "member1",
			Request: facade4meetingus.Request{
				SpaceRequest: dto4teamus.SpaceRequest{
					SpaceID: "space1",
				},
				MeetingID: "2020-12-13",
			},
		}
		if err := request.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
