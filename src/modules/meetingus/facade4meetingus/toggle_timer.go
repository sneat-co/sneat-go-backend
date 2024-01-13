package facade4meetingus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/models4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"math"
	"strings"
	"time"
)

const (
	// TimerOpStart "start"
	TimerOpStart = "start"

	// TimerOpStop "stop"
	TimerOpStop = "stop"

	// TimerOpPause "pause"
	TimerOpPause = "pause"
)

const (
	// TimerStatusActive "active"
	TimerStatusActive = "active"

	// TimerStatusStopped "stopped"
	TimerStatusStopped = "stopped"

	// TimerStatusPaused "paused"
	TimerStatusPaused = "paused"
)

//var timerActivationOps = [2]string{TimerOpStart, TimerOpResume}
//var timerSuspendingOps = [2]string{TimerOpStop, TimerOpPause}

// Params record
type Params struct {
	RecordFactory RecordFactory
	BeforeSafe    func(ctx context.Context, tx dal.ReadwriteTransaction, workerParams WorkerParams) error
}

// ToggleParams record
type ToggleParams struct {
	Params
	Request ToggleTimerRequest
}

// ToggleTimer toggles timer
func ToggleTimer(ctx context.Context, userContext facade.User, params ToggleParams) (response ToggleTimerResponse, err error) {
	if userContext == nil {
		err = errors.New("required parameter userContext == nil")
		return
	}
	if err = params.Request.Validate(); err != nil {
		err = fmt.Errorf("validation of MemberTimerRequest request failed: %w", err)
		return
	}

	uid := userContext.GetID()

	request := params.Request

	err = RunMeetingWorker(ctx, uid, request.Request, params.RecordFactory,
		func(ctx context.Context, tx dal.ReadwriteTransaction, workerParams WorkerParams) (err error) {
			now := time.Now()

			meeting := workerParams.Meeting.Data()
			timer := meeting.Timer
			if timer == nil {
				timer = &models4meetingus.Timer{}
				meeting.Timer = timer
			}

			isActive := meeting.Timer.Status == TimerStatusActive

			response.Timer = meeting.Timer

			setMemberDuration := func(seconds int) {
				if timer.SecondsByMember == nil {
					timer.SecondsByMember = make(map[string]int, 1)
				}
				timer.SecondsByMember[timer.ActiveMemberID] += seconds
			}

			getElapsedSeconds := func() int {
				elapsedSeconds := int(math.Round(now.Sub(timer.At).Seconds()))
				if elapsedSeconds == 0 {
					elapsedSeconds = 1
				}
				return elapsedSeconds
			}

			switch request.Operation {
			case TimerOpStart:
				if isActive && (request.Member == "" || timer.ActiveMemberID == request.Member) {
					return nil
				}
				timer.Status = TimerStatusActive
				if timer.ActiveMemberID != "" && request.Member != "" && timer.ActiveMemberID != request.Member {
					elapsedSeconds := getElapsedSeconds()
					setMemberDuration(elapsedSeconds)
				}
				timer.ActiveMemberID = strings.TrimSpace(request.Member)
			case TimerOpStop, TimerOpPause:
				if !isActive {
					return nil
				}
				switch request.Operation {
				case TimerOpStop:
					timer.Status = TimerStatusStopped
				case TimerOpPause:
					timer.Status = TimerStatusPaused
				default:
					return errors.New("coding error for timer operation stop or pause: " + request.Operation)
				}
				elapsedSeconds := getElapsedSeconds()
				timer.ElapsedSeconds += elapsedSeconds
				if timer.ActiveMemberID != "" {
					setMemberDuration(elapsedSeconds)
				}
				timer.ActiveMemberID = ""
			default:
				return validation.NewErrBadRequestFieldValue("operation",
					fmt.Sprintf("unknown timer operation: %v", request.Operation))
			}

			timer.At = now
			timer.By = dbmodels.ByUser{UID: uid}

			if params.BeforeSafe != nil {
				if err = params.BeforeSafe(ctx, tx, workerParams); err != nil {
					return err
				}
			}

			if record, ok := workerParams.Meeting.Record.(interface{ Validate() error }); ok {
				if err = record.Validate(); err != nil {
					return fmt.Errorf("api4meetingus record validation failed: %w", err)
				}
			}

			// This should be before updating or creating scrum record as it fetches other recs and will fail otherwise

			if workerParams.Meeting.Record.Exists() {
				err = tx.Update(ctx, workerParams.Meeting.Key, []dal.Update{
					{Field: "timer", Value: timer},
				})
				if err != nil {
					return fmt.Errorf("failed to update api4meetingus record: %w", err)
				}
			} else if err = tx.Insert(ctx, workerParams.Meeting.Record); err != nil {
				return fmt.Errorf("failed to create scrum record: %w", err)
			}
			return
		})
	return
}
