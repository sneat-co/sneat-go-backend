package facade4calendarius

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dal4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
)

// CancelHappening cancel a happening or it's slot or a single occurrence at specific date
func CancelHappening(ctx facade.ContextWithUser, request dto4calendarius.CancelHappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	err = dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) (err error) {
			switch params.Happening.Data.Type {
			case "":
				return fmt.Errorf("happening record has no type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
			case dbo4calendarius.HappeningTypeSingle:
				return cancelSingleHappening(ctx, params, request)
			case dbo4calendarius.HappeningTypeRecurring:
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

func cancelSingleHappening(ctx context.Context, params *dal4calendarius.HappeningWorkerParams, request dto4calendarius.CancelHappeningRequest) error {
	switch params.Happening.Data.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case dbo4calendarius.HappeningStatusActive:
		cancellation := CreateCancellation(params.UserID(), request.Reason)
		happeningUpdates := params.Happening.Data.MarkAsCanceled(cancellation)
		if err := params.Happening.Data.Validate(); err != nil {
			return fmt.Errorf("happening record is not valid: %w", err)
		}
		params.HappeningUpdates = append(params.HappeningUpdates, happeningUpdates...)
		params.Happening.Record.MarkAsChanged()
	case // Nothing to do
		dbo4calendarius.HappeningStatusCanceled,
		dbo4calendarius.HappeningStatusDeleted:
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
	params *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.CancelHappeningRequest,
) (err error) {

	var calendarDay dbo4calendarius.CalendarDayEntry
	{ // Load records that might need to be updated
		var recordsToGet []dal.Record

		if request.Date != "" {
			calendarDay = dbo4calendarius.NewCalendarDayEntry(request.SpaceID, request.Date)
			recordsToGet = append(recordsToGet, calendarDay.Record)
		} else if params.Happening.Data.Type == dbo4calendarius.HappeningTypeRecurring {
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
		for _, u := range updates {
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
				update.ByFieldPath([]string{dbo4calendarius.RecurringHappeningsField, params.Happening.ID, u.FieldName()}, u.Value()),
			)
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
	request dto4calendarius.CancelHappeningRequest,
	params *dal4calendarius.HappeningWorkerParams,
	cancellation dbo4calendarius.Cancellation,
	calendarDay dbo4calendarius.CalendarDayEntry,
) (err error) {
	var dayUpdates []update.Update
	happeningAdjustment, slotAdjustment := calendarDay.Data.GetAdjustment(params.Happening.ID, request.SlotID)
	if slotAdjustment == nil {
		slot := params.Happening.Data.GetSlot(request.SlotID)
		if slot == nil {
			return fmt.Errorf("%w: slot not found by SlotID=%v", facade.ErrBadRequest, request.SlotID)
		}
		slotAdjustment = &dbo4calendarius.SlotAdjustment{
			Cancellation: &cancellation,
		}
		if happeningAdjustment == nil {
			happeningAdjustment = new(dbo4calendarius.HappeningAdjustment)
			if calendarDay.Data.HappeningAdjustments == nil {
				calendarDay.Data.HappeningAdjustments = make(map[string]*dbo4calendarius.HappeningAdjustment, 1)
			}
			calendarDay.Data.HappeningAdjustments[params.Happening.ID] = happeningAdjustment
		}

		if happeningAdjustment.Slots == nil {
			happeningAdjustment.Slots = make(map[string]*dbo4calendarius.SlotAdjustment, 1)
		}
		happeningAdjustment.Slots[request.SlotID] = slotAdjustment
	}
	if !slices.Contains(calendarDay.Data.HappeningIDs, params.Happening.ID) {
		calendarDay.Data.HappeningIDs = append(calendarDay.Data.HappeningIDs, params.Happening.ID)
		dayUpdates = append(dayUpdates, update.ByFieldName("happeningIDs", calendarDay.Data.HappeningIDs))
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
		dayUpdates = append(dayUpdates, update.ByFieldName("cancellations", calendarDay.Data.HappeningAdjustments))
		if err = tx.Update(ctx, calendarDay.Key, dayUpdates); err != nil {
			return fmt.Errorf("failed to update calendar day record: %w", err)
		}
	}
	return err
}
