package dbo4scrumus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

func TestScrum_Validate(t *testing.T) {
	t.Run("empty_record", func(t *testing.T) {
		record := Scrum{}
		if err := record.Validate(); err != nil {
			t.Fatalf("no error expected for empty value, got: %v", err)
		}
	})

	t.Run("invalid_meeting", func(t *testing.T) {
		record := Scrum{
			Meeting: dbo4meetingus.Meeting{
				WithUserIDs: dbmodels.WithUserIDs{
					UserIDs: []string{"u1"},
				},
			},
		}
		if err := record.Validate(); err == nil {
			t.Fatal("expected error for meeting with no contacts but has userIDs")
		}
	})

	t.Run("invalid_risks_count", func(t *testing.T) {
		record := Scrum{RisksCount: -1}
		if err := record.Validate(); err == nil {
			t.Fatal("expected error for negative RisksCount")
		}
	})

	t.Run("invalid_questions_count", func(t *testing.T) {
		record := Scrum{QuestionsCount: -1}
		if err := record.Validate(); err == nil {
			t.Fatal("expected error for negative QuestionsCount")
		}
	})

	t.Run("invalid_status_key", func(t *testing.T) {
		record := Scrum{
			Statuses: ScrumStatusByMember{
				"": &MemberStatus{},
			},
		}
		if err := record.Validate(); err == nil {
			t.Fatal("expected error for empty status key")
		}
	})

	t.Run("invalid_status_value", func(t *testing.T) {
		record := Scrum{
			Statuses: ScrumStatusByMember{
				"m1": &MemberStatus{}, // Empty Member status is invalid
			},
		}
		if err := record.Validate(); err == nil {
			t.Fatal("expected error for invalid status value")
		}
	})
}

func TestScrum_GetOrCreateStatus(t *testing.T) {
	record := Scrum{}
	status := record.GetOrCreateStatus("m1")
	if status == nil {
		t.Fatal("expected status to be created")
	}
	if status.Member.ID != "m1" {
		t.Errorf("expected member ID m1, got %s", status.Member.ID)
	}
	if len(record.Statuses) != 1 {
		t.Errorf("expected 1 status, got %d", len(record.Statuses))
	}

	status2 := record.GetOrCreateStatus("m1")
	if status != status2 {
		t.Fatal("expected same status to be returned")
	}
	if len(record.Statuses) != 1 {
		t.Errorf("expected 1 status, got %d", len(record.Statuses))
	}
}

func TestMemberStatus_Validate(t *testing.T) {
	t.Run("invalid_member", func(t *testing.T) {
		v := MemberStatus{}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty member")
		}
	})

	t.Run("invalid_type_key", func(t *testing.T) {
		v := MemberStatus{
			Member: ScrumMember{ID: "m1", Title: "Member 1"},
			ByType: TasksByType{"": nil},
		}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty type key")
		}
	})

	t.Run("invalid_tasks", func(t *testing.T) {
		v := MemberStatus{
			Member: ScrumMember{ID: "m1", Title: "Member 1"},
			ByType: TasksByType{"type1": {nil}},
		}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for nil task")
		}
	})

	t.Run("valid", func(t *testing.T) {
		v := MemberStatus{
			Member: ScrumMember{ID: "m1", Title: "Member 1"},
			ByType: TasksByType{"type1": {{ID: "t1", Title: "Task 1"}}},
		}
		if err := v.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestMemberStatus_GetTask(t *testing.T) {
	v := MemberStatus{
		ByType: TasksByType{
			"type1": {
				{ID: "t1", Title: "Task 1"},
				{ID: "t2", Title: "Task 2"},
			},
		},
	}

	t.Run("found", func(t *testing.T) {
		task, index := v.GetTask("type1", "t2")
		if index != 1 {
			t.Errorf("expected index 1, got %d", index)
		}
		if task.ID != "t2" {
			t.Errorf("expected task ID t2, got %s", task.ID)
		}
	})

	t.Run("not_found_type", func(t *testing.T) {
		_, index := v.GetTask("type2", "t1")
		if index != -1 {
			t.Errorf("expected index -1, got %d", index)
		}
	})

	t.Run("not_found_id", func(t *testing.T) {
		_, index := v.GetTask("type1", "t3")
		if index != -1 {
			t.Errorf("expected index -1, got %d", index)
		}
	})
}

func TestTask_Validate(t *testing.T) {
	t.Run("empty_id", func(t *testing.T) {
		v := Task{Title: "Title"}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty ID")
		}
	})

	t.Run("empty_title", func(t *testing.T) {
		v := Task{ID: "ID"}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty Title")
		}
	})

	t.Run("invalid_comment", func(t *testing.T) {
		v := Task{
			ID:    "t1",
			Title: "Task 1",
			Comments: []*Comment{
				{},
			},
		}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for invalid comment")
		}
	})
}

func TestComment_Validate(t *testing.T) {
	t.Run("empty_id", func(t *testing.T) {
		v := Comment{Message: "Msg", By: &dbmodels.ByUser{}}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty ID")
		}
	})

	t.Run("empty_message", func(t *testing.T) {
		v := Comment{ID: "c1", By: &dbmodels.ByUser{}}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty message")
		}
	})

	t.Run("nil_by", func(t *testing.T) {
		v := Comment{ID: "c1", Message: "Msg"}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for nil By")
		}
	})

	t.Run("invalid_by", func(t *testing.T) {
		// dbmodels.ByUser.Validate() usually checks UserID
		v := Comment{ID: "c1", Message: "Msg", By: &dbmodels.ByUser{}}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for invalid By (empty UserID)")
		}
	})
}

func TestValidateTasks(t *testing.T) {
	t.Run("nil_task", func(t *testing.T) {
		if err := ValidateTasks(Tasks{nil}); err == nil {
			t.Fatal("expected error for nil task")
		}
	})

	t.Run("duplicate_id", func(t *testing.T) {
		tasks := Tasks{
			{ID: "t1", Title: "Task 1"},
			{ID: "t1", Title: "Task 2"},
		}
		if err := ValidateTasks(tasks); err == nil {
			t.Fatal("expected error for duplicate ID")
		}
	})
}

func TestMetricValue_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := MetricValue{}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty metric value")
		}
	})

	t.Run("valid_bool", func(t *testing.T) {
		b := true
		v := MetricValue{Bool: &b}
		if err := v.Validate(); err != nil {
			t.Fatal("unexpected error")
		}
	})
}

func TestMetricRecord_Validate(t *testing.T) {
	t.Run("empty_id", func(t *testing.T) {
		b := true
		v := MetricRecord{MetricValue: MetricValue{Bool: &b}}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty ID")
		}
	})
}

func TestScrumMember_Validate(t *testing.T) {
	t.Run("empty_id", func(t *testing.T) {
		v := ScrumMember{Title: "Title"}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty ID")
		}
	})
	t.Run("empty_title", func(t *testing.T) {
		v := ScrumMember{ID: "ID"}
		if err := v.Validate(); err == nil {
			t.Fatal("expected error for empty Title")
		}
	})
}
