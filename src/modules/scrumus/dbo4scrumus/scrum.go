package dbo4scrumus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

// Comment record
type Comment struct {
	ID      string           `json:"id" firestore:"id"`
	Message string           `json:"message" firestore:"message"`
	By      *dbmodels.ByUser `json:"by,omitempty" firestore:"by,omitempty"`
}

// Validate validates Comment record
func (v *Comment) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if v.Message == "" {
		return validation.NewErrRecordIsMissingRequiredField("message")
	}
	if v.By == nil {
		return validation.NewErrRecordIsMissingRequiredField("by")
	}
	if err := v.By.Validate(); err != nil {
		return err
	}
	return nil
}

// Task record
type Task struct {
	ID       string     `json:"id" firestore:"id"`
	Title    string     `json:"title" firestore:"title"`
	ThumbUps []string   `json:"thumbUps,omitempty" firestore:"thumbUps,omitempty"`
	Comments []*Comment `json:"comments,omitempty" firestore:"comments,omitempty"`
}

// Validate validates Task record
func (v *Task) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	for i, comment := range v.Comments {
		if err := comment.Validate(); err != nil {
			return fmt.Errorf("invalid comment at index %d: %w", i, err)
		}
	}
	return nil
}

// ScrumMember record
type ScrumMember struct {
	ID    string `json:"id" firestore:"id"`
	Title string `json:"title" firestore:"title"`
}

// Validate validates ScrumMember record
func (v *ScrumMember) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// Tasks is a type alias for slice of tasks
type Tasks = []*Task

// ValidateTasks validates tasks
func ValidateTasks(v Tasks) error {
	for i, t1 := range v {
		if t1 == nil {
			return fmt.Errorf("nil task at index %d", i)
		}
		if err := t1.Validate(); err != nil {
			return err
		}
		for j, t2 := range v {
			if i != j && t1.ID == t2.ID {
				return fmt.Errorf("duplicate task taskID=%s at indexes %d, %d", t1.ID, i, j)
			}
		}
	}
	return nil
}

// TasksByType is a type alias for map of tasks by task ContactID
type TasksByType = map[string]Tasks

// MemberStatus record
type MemberStatus struct {
	Member  ScrumMember     `json:"members" firestore:"members"`
	ByType  TasksByType     `json:"byType,omitempty" firestore:"byType,omitempty"`
	Metrics []*MetricRecord `json:"metrics,omitempty" firestore:"metrics,omitempty"`
}

// Validate validate MemberStatus record
func (v *MemberStatus) Validate() error {
	if err := v.Member.Validate(); err != nil {
		return err
	}
	for k, tasks := range v.ByType {
		if strings.TrimSpace(k) == "" {
			return validation.NewErrRecordIsMissingRequiredField("MemberStatus.ByType.InviteID")
		}
		if err := ValidateTasks(tasks); err != nil {
			return err
		}
	}
	return nil
}

// GetTask returns task by type & ContactID
func (v MemberStatus) GetTask(taskType, id string) (*Task, int) {
	tasks := v.ByType[taskType]
	for i, t := range tasks {
		if t.ID == id {
			return t, i
		}
	}
	return &Task{}, -1
}

// ScrumStatusByMember and alias for map of members statuses by members ContactID
type ScrumStatusByMember = map[string]*MemberStatus

// ScrumIDs hold previous & next IDs
type ScrumIDs struct {
	Prev string `json:"prev,omitempty" firestore:"prev,omitempty"`
	Next string `json:"next,omitempty" firestore:"next,omitempty"`
}

// MetricValue record
type MetricValue struct {
	Bool *bool   `json:"bool,omitempty" firestore:"bool,omitempty"`
	Int  *int    `json:"int,omitempty" firestore:"int,omitempty"`
	Str  *string `json:"str,omitempty" firestore:"str,omitempty"`
}

// Validate validates MetricValue
func (v *MetricValue) Validate() error {
	if v.Bool == nil && v.Int == nil && v.Str == nil {
		return validation.NewErrRecordIsMissingRequiredField("value")
	}
	return nil
}

// MetricRecord db record
type MetricRecord struct {
	ID  string `json:"id" firestore:"id"`
	UID string `json:"uid,omitempty"  firestore:"uid"`
	MetricValue
}

// Validate validates MetricRecord
func (v *MetricRecord) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	return v.MetricValue.Validate()
}

//var _ core.Validatable = (*Scrum)(nil)

// Scrum record
type Scrum struct {
	dbo4meetingus.Meeting
	RisksCount     int                       `json:"risksCount,omitempty" firestore:"risksCount,omitempty"`
	QuestionsCount int                       `json:"questionsCount,omitempty" firestore:"questionsCount,omitempty"`
	Statuses       ScrumStatusByMember       `json:"statuses,omitempty" firestore:"statuses,omitempty"`
	ScrumIDs       *ScrumIDs                 `json:"scrumIDs,omitempty" firestore:"scrumIDs,omitempty"`
	Metrics        []*dbo4teamus.SpaceMetric `json:"metrics,omitempty" firestore:"metrics,omitempty"`
	SpaceMetrics   []*MetricRecord           `json:"spaceMetrics,omitempty" firestore:"spaceMetrics,omitempty"`
}

var _ dbo4meetingus.MeetingInstance = (*Scrum)(nil)

// BaseMeeting returns base information on a api4meetingus
func (v *Scrum) BaseMeeting() *dbo4meetingus.Meeting {
	return &v.Meeting
}

// Validate validates Scrum record
func (v *Scrum) Validate() error {
	if err := v.Meeting.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("api4meetingus", err.Error())
	}
	if v.RisksCount < 0 {
		return validation.NewErrBadRecordFieldValue("RisksCount", "should be 0 or positive")
	}
	if v.QuestionsCount < 0 {
		return validation.NewErrBadRecordFieldValue("QuestionsCount", "should be 0 or positive")
	}
	for k, status := range v.Statuses {
		if strings.TrimSpace(k) == "" {
			return validation.NewErrBadRecordFieldValue("Scrum.Statuses", "empty key")
		}
		if err := status.Validate(); err != nil {
			return err
		}
	}
	if v.Timer != nil {
		if err := v.Timer.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// GetOrCreateStatus returns status of a members
func (v *Scrum) GetOrCreateStatus(memberID string) *MemberStatus {
	if v.Statuses == nil {
		v.Statuses = make(map[string]*MemberStatus, 1)
	}
	for id, s := range v.Statuses {
		if id == memberID {
			return s
		}
	}
	status := &MemberStatus{Member: ScrumMember{ID: memberID}}
	v.Statuses[memberID] = status
	return status
}
