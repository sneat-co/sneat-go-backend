package facade4calendarium

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// CancelHappening marks happening as canceled
func CancelHappening(ctx context.Context, user facade.User, request dto4calendarium.CancelHappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	happening := models4calendarium.NewHappeningContext(request.TeamID, request.HappeningID)
	err = dal4teamus.RunModuleTeamWorker(ctx, user, request.TeamRequest,
		const4calendarium.ModuleID,
		new(models4calendarium.CalendariumTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4calendarium.CalendariumTeamDto]) (err error) {
			if err = tx.Get(ctx, happening.Record); err != nil {
				return fmt.Errorf("failed to get happening: %w", err)
			}
			switch happening.Dbo.Type {
			case "":
				return fmt.Errorf("happening record has no type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
			case "single":
				return cancelSingleHappening(ctx, tx, params.UserID, happening)
			case "recurring":
				return cancelRecurringHappening(ctx, tx, params, params.UserID, happening, request)
			default:
				return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Dbo.Type)
			}
		})
	if err != nil {
		return fmt.Errorf("failed to cancel happening: %w", err)
	}
	return
}

func cancelSingleHappening(ctx context.Context, tx dal.ReadwriteTransaction, userID string, happening models4calendarium.HappeningContext) error {
	switch happening.Dbo.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case models4calendarium.HappeningStatusActive:
		happening.Dbo.Status = models4calendarium.HappeningStatusCanceled
		happening.Dbo.Canceled = &models4calendarium.Canceled{
			At: time.Now(),
			By: dbmodels.ByUser{UID: userID},
		}
		if err := happening.Dbo.Validate(); err != nil {
			return fmt.Errorf("happening record is not valid: %w", err)
		}
		happeningUpdates := []dal.Update{
			{Field: "status", Value: happening.Dbo.Status},
			{Field: "canceled", Value: happening.Dbo.Canceled},
		}
		if err := tx.Update(ctx, happening.Key, happeningUpdates); err != nil {
			return err
		}
	case models4calendarium.HappeningStatusDeleted:
		// Nothing to do
	default:
		return fmt.Errorf("only active happening can be canceled but happening is in status=[%v]", happening.Dbo.Status)
	}
	happening.Dbo.Status = "canceled"
	return nil
}

func cancelRecurringHappening(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4teamus.ModuleTeamWorkerParams[*models4calendarium.CalendariumTeamDto],
	uid string,
	happening models4calendarium.HappeningContext,
	request dto4calendarium.CancelHappeningRequest,
) error {
	happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(happening.ID)
	if happeningBrief == nil {
		return errors.New("happening brief is not found in team record")
	}

	if request.Date == "" {
		if err := markRecurringHappeningRecordAsCanceled(ctx, tx, uid, happening, request); err != nil {
			return err
		}
		happeningBrief.Status = models4calendarium.HappeningStatusCanceled
		happeningBrief.Canceled = createCanceled(uid, request.Reason)
		if err := happeningBrief.Validate(); err != nil {
			return fmt.Errorf("happening brief in team record is not valid: %w", err)
		}
		params.TeamUpdates = append(params.TeamUpdates, dal.Update{
			Field: "recurringHappenings",
			Value: params.TeamModuleEntry.Data.RecurringHappenings,
		})
	} else {
		calendarDay := models4calendarium.NewCalendarDayContext(params.Team.ID, request.Date)

		var isNewRecord bool
		if err := tx.Get(ctx, calendarDay.Record); err != nil {
			if dal.IsNotFound(err) {
				isNewRecord = true
			} else {
				return fmt.Errorf("failed to get calendar day record by ContactID: %w", err)
			}
		}

		var dayUpdates []dal.Update
		_, adjustment := calendarDay.Dto.GetAdjustment(happening.ID, request.SlotID)
		if adjustment == nil {
			_, slot := happening.Dbo.GetSlot(request.SlotID)
			if slot == nil {
				return fmt.Errorf("%w: slot not found by ContactID=%v", facade.ErrBadRequest, request.SlotID)
			}
			adjustment = &models4calendarium.HappeningAdjustment{
				HappeningID: happening.ID,
				Slot:        *slot,
				Canceled: &models4calendarium.Canceled{
					At:     time.Now(),
					By:     dbmodels.ByUser{UID: uid},
					Reason: request.Reason,
				},
			}
			calendarDay.Dto.HappeningAdjustments = append(calendarDay.Dto.HappeningAdjustments, adjustment)
		}
		if i := slice.Index(calendarDay.Dto.HappeningIDs, happening.ID); i < 0 {
			calendarDay.Dto.HappeningIDs = append(calendarDay.Dto.HappeningIDs, happening.ID)
			dayUpdates = append(dayUpdates, dal.Update{
				Field: "happeningIDs", Value: calendarDay.Dto.HappeningIDs,
			})
		}
		var modified bool
		if adjustment.Slot.ID == request.SlotID {
			if strings.TrimSpace(request.Reason) != "" && (adjustment.Canceled == nil || request.Reason != adjustment.Canceled.Reason) {
				if adjustment.Canceled == nil {
					adjustment.Canceled = &models4calendarium.Canceled{
						At:     time.Now(),
						By:     dbmodels.ByUser{UID: uid},
						Reason: request.Reason,
					}
				}
				adjustment.Canceled.Reason = request.Reason
				modified = true
			}
		} else {
			_, slot := happening.Dbo.GetSlot(request.SlotID)
			if slot == nil {
				return fmt.Errorf("%w: unknown slot ContactID=%v", facade.ErrBadRequest, request.SlotID)
			}
			adjustment.Slot = *slot
			adjustment.Canceled = &models4calendarium.Canceled{
				At:     time.Now(),
				By:     dbmodels.ByUser{UID: uid},
				Reason: request.Reason,
			}
			modified = true
		}

		if err := calendarDay.Dto.Validate(); err != nil {
			return fmt.Errorf("calendar day record is not valid: %w", err)
		}

		if isNewRecord {
			if err := tx.Insert(ctx, calendarDay.Record); err != nil {
				return fmt.Errorf("failed to create calendar day record: %w", err)
			}
		} else if modified {
			dayUpdates = append(dayUpdates, dal.Update{
				Field: "cancellations", Value: calendarDay.Dto.HappeningAdjustments,
			})
			if err := tx.Update(ctx, calendarDay.Key, dayUpdates); err != nil {
				return fmt.Errorf("failed to update calendar day record: %w", err)
			}
		}
	}

	return nil
}

func createCanceled(uid, reason string) *models4calendarium.Canceled {
	return &models4calendarium.Canceled{
		At:     time.Now(),
		By:     dbmodels.ByUser{UID: uid},
		Reason: strings.TrimSpace(reason),
	}
}
func markRecurringHappeningRecordAsCanceled(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	happening models4calendarium.HappeningContext,
	request dto4calendarium.CancelHappeningRequest,
) error {
	var happeningUpdates []dal.Update
	happening.Dbo.Status = models4calendarium.HappeningStatusCanceled
	if happening.Dbo.Canceled == nil {
		happening.Dbo.Canceled = createCanceled(uid, request.Reason)
	} else if reason := strings.TrimSpace(request.Reason); reason != "" {
		happening.Dbo.Canceled.Reason = reason
	}
	happeningUpdates = append(happeningUpdates,
		dal.Update{
			Field: "status",
			Value: happening.Dbo.Status,
		},
		dal.Update{
			Field: "canceled",
			Value: happening.Dbo.Canceled,
		},
	)
	if err := happening.Dbo.Validate(); err != nil {
		return fmt.Errorf("happening record is not valid: %w", err)
	}
	if err := tx.Update(ctx, happening.Key, happeningUpdates); err != nil {
		return fmt.Errorf("faield to update happening record: %w", err)
	}
	return nil
}
