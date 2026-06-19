package facade4calendarius

import (
	"context"
	"errors"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dal4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

func CancelSlotAdjustment(ctx facade.ContextWithUser, request dto4calendarius.HappeningDateSlotIDRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) (err error) {
		switch params.Happening.Data.Type {
		case dbo4calendarius.HappeningTypeSingle:
			return errors.New("only recurring happenings can be adjusted, single happenings should be updated")
		case dbo4calendarius.HappeningTypeRecurring:
			if err = cancelAdjustment(ctx, tx, params, request); err != nil {
				return fmt.Errorf("failed to cancel adjustment of a recurring happening: %w", err)
			}
			return err
		}
		return
	}

	if err = dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func cancelAdjustment(
	ctx context.Context, tx dal.ReadwriteTransaction,
	_ *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.HappeningDateSlotIDRequest,
) (err error) {
	calendarDay := dbo4calendarius.NewCalendarDayEntry(request.SpaceID, request.Date)
	if err = tx.Get(ctx, calendarDay.Record); err != nil {
		if dal.IsNotFound(err) {
			return nil // nothing to cancel
		}
		return fmt.Errorf("failed to get calendar day record: %w", err)
	}
	happeningAdjustment, slotAdjustment := calendarDay.Data.GetAdjustment(request.HappeningID, request.SlotID)
	if slotAdjustment == nil {
		return nil // nothing to cancel
	}
	delete(happeningAdjustment.Slots, request.SlotID)
	if len(happeningAdjustment.Slots) == 0 {
		delete(calendarDay.Data.HappeningAdjustments, request.HappeningID)
		calendarDay.Data.HappeningIDs = slice.RemoveInPlaceByValue(calendarDay.Data.HappeningIDs, request.HappeningID)
		if len(calendarDay.Data.HappeningAdjustments) == 0 {
			if err = tx.Delete(ctx, calendarDay.Key); err != nil {
				return fmt.Errorf("failed to delete calendar day record: %w", err)
			}
			return nil
		} else if err = tx.Update(ctx, calendarDay.Key, []update.Update{
			update.ByFieldPath([]string{"happeningAdjustments", request.HappeningID}, update.DeleteField),
			update.ByFieldName("happeningIDs", calendarDay.Data.HappeningIDs),
		}); err != nil {
			return fmt.Errorf("failed to update calendar day record with adjustment cancellation: %w", err)
		}
	} else {
		if err = tx.Update(ctx, calendarDay.Key, []update.Update{
			update.ByFieldPath([]string{"happeningAdjustments", request.HappeningID, "slots", request.SlotID}, update.DeleteField),
		}); err != nil {
			return fmt.Errorf("failed to update calendar day record with adjustment cancellation: %w", err)
		}
	}
	return nil
}

// AdjustSlot temporary changes slot (for example, time changed for a specific date, or first class has been canceled)
func AdjustSlot(ctx facade.ContextWithUser, request dto4calendarius.HappeningSlotDateRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) (err error) {
		switch params.Happening.Data.Type {
		case dbo4calendarius.HappeningTypeSingle:
			return errors.New("only recurring happenings can be adjusted, single happenings should be updated")
		case dbo4calendarius.HappeningTypeRecurring:
			if err = adjustSlotInCalendarDay(ctx, tx, params, request); err != nil {
				return fmt.Errorf("failed to adjust recurring happening: %w", err)
			}
			return err
		}
		return
	}

	if err = dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func adjustSlotInCalendarDay(
	ctx context.Context, tx dal.ReadwriteTransaction,
	_ *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.HappeningSlotDateRequest,
) error {
	calendarDay := dbo4calendarius.NewCalendarDayEntry(request.SpaceID, request.Date)
	if err := tx.Get(ctx, calendarDay.Record); err != nil {
		if !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get calendar day record: %w", err)
		}
	}

	// TODO: We should cancel/delete adjustment if it matches original timing

	happeningAdjustment, slotAdjustment := calendarDay.Data.GetAdjustment(request.HappeningID, request.Slot.ID)
	if slotAdjustment == nil {
		if happeningAdjustment == nil {
			happeningAdjustment = &dbo4calendarius.HappeningAdjustment{}
			if calendarDay.Data.HappeningAdjustments == nil {
				calendarDay.Data.HappeningAdjustments = make(map[string]*dbo4calendarius.HappeningAdjustment, 1)
			}
			calendarDay.Data.HappeningAdjustments[request.HappeningID] = happeningAdjustment
		}
		slotAdjustment = new(dbo4calendarius.SlotAdjustment)
		if happeningAdjustment.Slots == nil {
			happeningAdjustment.Slots = make(map[string]*dbo4calendarius.SlotAdjustment, 1)
		}
		happeningAdjustment.Slots[request.Slot.ID] = slotAdjustment
	}
	slotAdjustment.Adjustment = &request.Slot.HappeningSlot
	var happeningIDsChanged bool
	if happeningIDsChanged = slice.Index(calendarDay.Data.HappeningIDs, request.HappeningID) < 0; happeningIDsChanged {
		calendarDay.Data.HappeningIDs = append(calendarDay.Data.HappeningIDs, request.HappeningID)
	}

	if err := calendarDay.Data.Validate(); err != nil {
		return fmt.Errorf("calednar day record is not valid: %w", err)
	}

	if calendarDay.Record.Exists() {
		updates := []update.Update{
			update.ByFieldName("happeningAdjustments", calendarDay.Data.HappeningAdjustments),
		}
		if happeningIDsChanged {
			updates = append(updates,
				update.ByFieldName("happeningIDs", calendarDay.Data.HappeningIDs))
		}
		if err := tx.Update(ctx, calendarDay.Key, updates); err != nil {
			return fmt.Errorf("failed to update calendar day record with happening slotAdjustment: %w", err)
		}
	} else {
		if err := tx.Insert(ctx, calendarDay.Record); err != nil {
			return fmt.Errorf("failed to insert calendar day record with happening slotAdjustment: %w", err)
		}
	}
	return nil
}
