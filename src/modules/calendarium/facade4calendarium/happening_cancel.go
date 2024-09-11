package facade4calendarium

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
	"slices"
)

// CancelHappening cancel a happening or it's slot or a single occurrence at specific date
func CancelHappening(ctx context.Context, userCtx facade.UserContext, request dto4calendarium.CancelHappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	err = dal4calendarium.RunHappeningSpaceWorker(ctx, userCtx, request.HappeningRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		switch params.Happening.Data.Type {
		case "":
			return fmt.Errorf("happening record has no type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
		case dbo4calendarium.HappeningTypeSingle:
			return cancelSingleHappening(ctx, params, request)
		case dbo4calendarium.HappeningTypeRecurring:
			return cancelRecurringHappeningInstance(ctx, tx, params, request)
		default:
			return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+params.Happening.Data.Type)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to cancel happening: %w", err)
	}
	return
}

func cancelSingleHappening(ctx context.Context, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.CancelHappeningRequest) error {
	switch params.Happening.Data.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case dbo4calendarium.HappeningStatusActive:
		cancellation := CreateCancellation(params.UserID(), request.Reason)
		happeningUpdates := params.Happening.Data.MarkAsCanceled(cancellation)
		if err := params.Happening.Data.Validate(); err != nil {
			return fmt.Errorf("happening record is not valid: %w", err)
		}
		params.HappeningUpdates = append(params.HappeningUpdates, happeningUpdates...)
		params.Happening.Record.MarkAsChanged()
	case // Nothing to do
		dbo4calendarium.HappeningStatusCanceled,
		dbo4calendarium.HappeningStatusDeleted:
		logus.Infof(ctx, "An attempt to cancel happening that is already canceled or deleted: happeningID=%s, status=%s",
			params.Happening.ID, params.Happening.Data.Status)
	default:
		return fmt.Errorf("only active happening can be canceled but happening is in status=[%s]", params.Happening.Data.Status)
	}
	return nil
}

func cancelRecurringHappeningInstance(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4calendarium.HappeningWorkerParams,
	request dto4calendarium.CancelHappeningRequest,
) (err error) {

	var calendarDay dbo4calendarium.CalendarDayEntry
	{ // Load records that might need to be updated
		var recordsToGet []dal.Record

		if request.Date != "" {
			calendarDay = dbo4calendarium.NewCalendarDayEntry(request.SpaceID, request.Date)
			recordsToGet = append(recordsToGet, calendarDay.Record)
		} else if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
			// We do not update space module entry if it's a single happening or if a specific date is provided
			recordsToGet = append(recordsToGet, params.SpaceModuleEntry.Record)
		}
		if len(recordsToGet) > 0 {
			if err = tx.GetMulti(ctx, recordsToGet); err != nil {
				return
			}
		}
	}

	uid := params.UserID()
	cancellation := CreateCancellation(uid, request.Reason)

	if request.Date == "" {
		happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID)
		if happeningBrief == nil {
			return errors.New("happening brief is not found in space record")
		}

		if params.HappeningUpdates = params.Happening.Data.MarkAsCanceled(cancellation); len(params.HappeningUpdates) > 0 {
			params.Happening.Record.MarkAsChanged()
		}
		updates := happeningBrief.MarkAsCanceled(cancellation)
		if err = happeningBrief.Validate(); err != nil {
			return fmt.Errorf("happening brief in team record is not valid: %w", err)
		}
		for _, update := range updates {
			update.Field = fmt.Sprintf("recurringHappenings.%s.%s", params.Happening.ID, update.Field)
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update)
		}
		params.SpaceModuleEntry.Record.MarkAsChanged()
	} else if err = addCancellationToCalendarDayAdjustments(ctx, tx, request, params, cancellation, calendarDay); err != nil {
		return fmt.Errorf("failed to add cancellation to calendar day adjustments: %w", err)
	}

	return nil
}

func addCancellationToCalendarDayAdjustments(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4calendarium.CancelHappeningRequest,
	params *dal4calendarium.HappeningWorkerParams,
	cancellation dbo4calendarium.Cancellation,
	calendarDay dbo4calendarium.CalendarDayEntry,
) (err error) {
	var dayUpdates []dal.Update
	happeningAdjustment, slotAdjustment := calendarDay.Data.GetAdjustment(params.Happening.ID, request.SlotID)
	if slotAdjustment == nil {
		slot := params.Happening.Data.GetSlot(request.SlotID)
		if slot == nil {
			return fmt.Errorf("%w: slot not found by SlotID=%v", facade.ErrBadRequest, request.SlotID)
		}
		slotAdjustment = &dbo4calendarium.SlotAdjustment{
			Cancellation: &cancellation,
		}
		if happeningAdjustment == nil {
			happeningAdjustment = new(dbo4calendarium.HappeningAdjustment)
			if calendarDay.Data.HappeningAdjustments == nil {
				calendarDay.Data.HappeningAdjustments = make(map[string]*dbo4calendarium.HappeningAdjustment, 1)
			}
			calendarDay.Data.HappeningAdjustments[params.Happening.ID] = happeningAdjustment
		}

		if happeningAdjustment.Slots == nil {
			happeningAdjustment.Slots = make(map[string]*dbo4calendarium.SlotAdjustment, 1)
		}
		happeningAdjustment.Slots[request.SlotID] = slotAdjustment
	}
	if !slices.Contains(calendarDay.Data.HappeningIDs, params.Happening.ID) {
		calendarDay.Data.HappeningIDs = append(calendarDay.Data.HappeningIDs, params.Happening.ID)
		dayUpdates = append(dayUpdates, dal.Update{
			Field: "happeningIDs", Value: calendarDay.Data.HappeningIDs,
		})
	}
	var modified bool
	if slotAdjustment.Cancellation == nil || *slotAdjustment.Cancellation != cancellation {
		modified = true
		slotAdjustment.Cancellation = &cancellation
	}

	if err = calendarDay.Data.Validate(); err != nil {
		return fmt.Errorf("calendar day record is not valid: %w", err)
	}

	if !calendarDay.Record.Exists() {
		if err = tx.Insert(ctx, calendarDay.Record); err != nil {
			return fmt.Errorf("failed to create calendar day record: %w", err)
		}
	} else if modified {
		dayUpdates = append(dayUpdates, dal.Update{
			Field: "cancellations", Value: calendarDay.Data.HappeningAdjustments,
		})
		if err = tx.Update(ctx, calendarDay.Key, dayUpdates); err != nil {
			return fmt.Errorf("failed to update calendar day record: %w", err)
		}
	}
	return err
}
