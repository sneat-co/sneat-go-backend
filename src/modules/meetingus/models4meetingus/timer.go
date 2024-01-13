package models4meetingus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
	"time"
)

var _ core.Validatable = (*Timer)(nil)

const (
	// TimerStatusActive represents "active" status
	TimerStatusActive = "active"

	// TimerStatusPaused represents "paused" status
	TimerStatusPaused = "paused"

	// TimerStatusStopped represents "stopped" status
	TimerStatusStopped = "stopped"
)

// Timer record
type Timer struct {
	Status          string          `json:"status" firestore:"status"`
	By              dbmodels.ByUser `json:"by" firestore:"by"`
	At              time.Time       `json:"at" firestore:"at"`
	ElapsedSeconds  int             `json:"elapsedSeconds,omitempty" firestore:"elapsedSeconds,omitempty"`
	ActiveMemberID  string          `json:"activeMemberId,omitempty" firestore:"activeMemberId,omitempty"`
	SecondsByMember map[string]int  `json:"secondsByMember,omitempty" firestore:"secondsByMember,omitempty"`
	SecondsByTopic  map[string]int  `json:"secondsByTopic,omitempty" firestore:"secondsByTopic,omitempty"`
}

// Validate validates Timer record
func (v *Timer) Validate() error {
	switch v.Status {
	case TimerStatusActive:
		break
	case TimerStatusStopped:
		if v.ElapsedSeconds == 0 {
			return validation.NewErrBadRecordFieldValue("Timer.ElapsedSeconds", "stopped timer should have elapsed seconds")
		}
		if v.ActiveMemberID != "" {
			return validation.NewErrBadRecordFieldValue("Timer.ActiveMemberID", "stopped timer should NOT have current members")
		}
	case TimerStatusPaused:
		if v.ElapsedSeconds == 0 {
			return validation.NewErrBadRecordFieldValue("Timer.ElapsedSeconds", "paused timer should have elapsed seconds")
		}
	case "":
		return validation.NewErrRecordIsMissingRequiredField("Timer.status")
	default:
		return validation.NewErrBadRecordFieldValue("Timer.status", fmt.Sprintf("unknown timer status: %v", v.Status))
	}
	if err := v.By.Validate(); err != nil {
		return fmt.Errorf("invalid 'by' field: %w", err)
	}
	if v.ActiveMemberID != strings.TrimSpace(v.ActiveMemberID) {
		return validation.NewErrBadRecordFieldValue("Timer.ActiveMemberID", fmt.Sprintf("unexpected spaces: [%v]", v.ActiveMemberID))
	}
	for k, v := range v.SecondsByMember {
		if v <= 0 {
			return validation.NewErrBadRecordFieldValue("Timer.SecondsByMember", fmt.Sprintf("members duration should be positive, got: %v", k))
		}
	}
	for k, v := range v.SecondsByTopic {
		if v <= 0 {
			return validation.NewErrBadRecordFieldValue("Timer.SecondsByTopic", fmt.Sprintf("topic duration should be positive, got: %v", k))
		}
	}
	return nil
}
