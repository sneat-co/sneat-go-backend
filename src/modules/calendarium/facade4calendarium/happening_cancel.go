package facade4calendarium

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
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

	happening := dbo4calendarium.NewHappeningEntry(request.TeamID, request.HappeningID)
	err = dal4teamus.RunModuleTeamWorker(ctx, user, request.TeamRequest,
		const4calendarium.ModuleID,
		new(dbo4calendarium.CalendariumTeamDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo]) (err error) {
			if err = tx.Get(ctx, happening.Record); err != nil {
				return fmt.Errorf("failed to get happening: %w", err)
			}
			switch happening.Data.Type {
			case "":
				return fmt.Errorf("happening record has no type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
			case dbo4calendarium.HappeningTypeSingle:
				return cancelSingleHappening(ctx, tx, params.UserID, happening)
			case dbo4calendarium.HappeningTypeRecurring:
				return cancelRecurringHappening(ctx, tx, params, params.UserID, happening, request)
			default:
				return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Data.Type)
			}
		})
	if err != nil {
		return fmt.Errorf("failed to cancel happening: %w", err)
	}
	return
}

func cancelSingleHappening(ctx context.Context, tx dal.ReadwriteTransaction, userID string, happening dbo4calendarium.HappeningEntry) error {
	switch happening.Data.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case dbo4calendarium.HappeningStatusActive:
		happening.Data.Status = dbo4calendarium.HappeningStatusCanceled
		happening.Data.Canceled = &dbo4calendarium.Canceled{
			At: time.Now(),
			By: dbmodels.ByUser{UID: userID},
		}
		if err := happening.Data.Validate(); err != nil {
			return fmt.Errorf("happening record is not valid: %w", err)
		}
		happeningUpdates := []dal.Update{
			{Field: "status", Value: happening.Data.Status},
			{Field: "canceled", Value: happening.Data.Canceled},
		}
		if err := tx.Update(ctx, happening.Key, happeningUpdates); err != nil {
			return err
		}
	case dbo4calendarium.HappeningStatusDeleted:
		// Nothing to do
	default:
		return fmt.Errorf("only active happening can be canceled but happening is in status=[%v]", happening.Data.Status)
	}
	happening.Data.Status = "canceled"
	return nil
}

func cancelRecurringHappening(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo],
	uid string,
	happening dbo4calendarium.HappeningEntry,
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
		happeningBrief.Status = dbo4calendarium.HappeningStatusCanceled
		happeningBrief.Canceled = createCanceled(uid, request.Reason)
		if err := happeningBrief.Validate(); err != nil {
			return fmt.Errorf("happening brief in team record is not valid: %w", err)
		}
		params.TeamUpdates = append(params.TeamUpdates, dal.Update{
			Field: "recurringHappenings",
			Value: params.TeamModuleEntry.Data.RecurringHappenings,
		})
	} else {
		calendarDay := dbo4calendarium.NewCalendarDayContext(params.Team.ID, request.Date)

		var isNewRecord bool
		if err := tx.Get(ctx, calendarDay.Record); err != nil {
			if dal.IsNotFound(err) {
				isNewRecord = true
			} else {
				return fmt.Errorf("failed to get calendar day record by ContactID: %w", err)
			}
		}

		var dayUpdates []dal.Update
		_, adjustment := calendarDay.Data.GetAdjustment(happening.ID, request.SlotID)
		if adjustment == nil {
			_, slot := happening.Data.GetSlot(request.SlotID)
			if slot == nil {
				return fmt.Errorf("%w: slot not found by ContactID=%v", facade.ErrBadRequest, request.SlotID)
			}
			adjustment = &dbo4calendarium.HappeningAdjustment{
				HappeningID: happening.ID,
				Slot:        *slot,
				Canceled: &dbo4calendarium.Canceled{
					At:     time.Now(),
					By:     dbmodels.ByUser{UID: uid},
					Reason: request.Reason,
				},
			}
			calendarDay.Data.HappeningAdjustments = append(calendarDay.Data.HappeningAdjustments, adjustment)
		}
		if i := slice.Index(calendarDay.Data.HappeningIDs, happening.ID); i < 0 {
			calendarDay.Data.HappeningIDs = append(calendarDay.Data.HappeningIDs, happening.ID)
			dayUpdates = append(dayUpdates, dal.Update{
				Field: "happeningIDs", Value: calendarDay.Data.HappeningIDs,
			})
		}
		var modified bool
		if adjustment.Slot.ID == request.SlotID {
			if strings.TrimSpace(request.Reason) != "" && (adjustment.Canceled == nil || request.Reason != adjustment.Canceled.Reason) {
				if adjustment.Canceled == nil {
					adjustment.Canceled = &dbo4calendarium.Canceled{
						At:     time.Now(),
						By:     dbmodels.ByUser{UID: uid},
						Reason: request.Reason,
					}
				}
				adjustment.Canceled.Reason = request.Reason
				modified = true
			}
		} else {
			_, slot := happening.Data.GetSlot(request.SlotID)
			if slot == nil {
				return fmt.Errorf("%w: unknown slot ContactID=%v", facade.ErrBadRequest, request.SlotID)
			}
			adjustment.Slot = *slot
			adjustment.Canceled = &dbo4calendarium.Canceled{
				At:     time.Now(),
				By:     dbmodels.ByUser{UID: uid},
				Reason: request.Reason,
			}
			modified = true
		}

		if err := calendarDay.Data.Validate(); err != nil {
			return fmt.Errorf("calendar day record is not valid: %w", err)
		}

		if isNewRecord {
			if err := tx.Insert(ctx, calendarDay.Record); err != nil {
				return fmt.Errorf("failed to create calendar day record: %w", err)
			}
		} else if modified {
			dayUpdates = append(dayUpdates, dal.Update{
				Field: "cancellations", Value: calendarDay.Data.HappeningAdjustments,
			})
			if err := tx.Update(ctx, calendarDay.Key, dayUpdates); err != nil {
				return fmt.Errorf("failed to update calendar day record: %w", err)
			}
		}
	}

	return nil
}

func createCanceled(uid, reason string) *dbo4calendarium.Canceled {
	return &dbo4calendarium.Canceled{
		At:     time.Now(),
		By:     dbmodels.ByUser{UID: uid},
		Reason: strings.TrimSpace(reason),
	}
}
func markRecurringHappeningRecordAsCanceled(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	happening dbo4calendarium.HappeningEntry,
	request dto4calendarium.CancelHappeningRequest,
) error {
	var happeningUpdates []dal.Update
	happening.Data.Status = dbo4calendarium.HappeningStatusCanceled
	if happening.Data.Canceled == nil {
		happening.Data.Canceled = createCanceled(uid, request.Reason)
	} else if reason := strings.TrimSpace(request.Reason); reason != "" {
		happening.Data.Canceled.Reason = reason
	}
	happeningUpdates = append(happeningUpdates,
		dal.Update{
			Field: "status",
			Value: happening.Data.Status,
		},
		dal.Update{
			Field: "canceled",
			Value: happening.Data.Canceled,
		},
	)
	if err := happening.Data.Validate(); err != nil {
		return fmt.Errorf("happening record is not valid: %w", err)
	}
	if err := tx.Update(ctx, happening.Key, happeningUpdates); err != nil {
		return fmt.Errorf("faield to update happening record: %w", err)
	}
	return nil
}
