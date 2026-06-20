package facade4scrumus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/dbo4scrumus"
	"github.com/stretchr/testify/assert"
)

func ptrBool(b bool) *bool { return &b }

func newScrumRequest() facade4meetingus.Request {
	return facade4meetingus.Request{
		SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: "space1"},
		MeetingID:    "meeting1",
	}
}

func TestSetMetricRequest_Validate(t *testing.T) {
	newRequest := func() SetMetricRequest {
		return SetMetricRequest{
			Request:     newScrumRequest(),
			Metric:      "velocity",
			MetricValue: dbo4scrumus.MetricValue{Bool: ptrBool(true)},
		}
	}
	tests := []struct {
		name    string
		mutate  func(r *SetMetricRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*SetMetricRequest) {}, wantErr: false},
		{name: "missing_metric", mutate: func(r *SetMetricRequest) { r.Metric = " " }, wantErr: true},
		{name: "missing_value", mutate: func(r *SetMetricRequest) { r.MetricValue = dbo4scrumus.MetricValue{} }, wantErr: true},
		{name: "missing_space", mutate: func(r *SetMetricRequest) { r.SpaceID = "" }, wantErr: true},
		{name: "missing_meeting", mutate: func(r *SetMetricRequest) { r.MeetingID = "" }, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRequest()
			tt.mutate(&r)
			err := r.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReorderTaskRequest_Validate(t *testing.T) {
	newRequest := func() ReorderTaskRequest {
		return ReorderTaskRequest{
			TaskRequest: TaskRequest{
				Request:   newScrumRequest(),
				ContactID: "member1",
				Type:      "todo",
				Task:      "task1",
			},
			Len:    3,
			From:   0,
			To:     1,
			After:  "task2",
			Before: "task3",
		}
	}
	tests := []struct {
		name    string
		mutate  func(r *ReorderTaskRequest)
		wantErr bool
	}{
		{name: "valid", mutate: func(*ReorderTaskRequest) {}, wantErr: false},
		{name: "missing_task_request_field", mutate: func(r *ReorderTaskRequest) { r.ContactID = "" }, wantErr: true},
		{name: "empty_len", mutate: func(r *ReorderTaskRequest) { r.Len = 0 }, wantErr: true},
		{name: "negative_from", mutate: func(r *ReorderTaskRequest) { r.From = -1 }, wantErr: true},
		{name: "negative_to", mutate: func(r *ReorderTaskRequest) { r.To = -1 }, wantErr: true},
		{name: "from_ge_len", mutate: func(r *ReorderTaskRequest) { r.From = 3 }, wantErr: true},
		{name: "to_ge_len", mutate: func(r *ReorderTaskRequest) { r.To = 3 }, wantErr: true},
		{name: "from_eq_to", mutate: func(r *ReorderTaskRequest) { r.To = 0 }, wantErr: true},
		{name: "after_eq_before", mutate: func(r *ReorderTaskRequest) { r.After = "x"; r.Before = "x" }, wantErr: true},
		{name: "after_eq_task", mutate: func(r *ReorderTaskRequest) { r.After = "task1" }, wantErr: true},
		{name: "before_eq_task", mutate: func(r *ReorderTaskRequest) { r.Before = "task1" }, wantErr: true},
		{name: "after_set_when_moving_to_start", mutate: func(r *ReorderTaskRequest) { r.To = 0; r.From = 1; r.After = "task2" }, wantErr: true},
		{name: "before_set_when_moving_to_end", mutate: func(r *ReorderTaskRequest) { r.To = 2; r.Before = "task3" }, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRequest()
			tt.mutate(&r)
			err := r.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
