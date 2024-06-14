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
	"github.com/strongo/slice"
)

// AdjustSlot temporary changes slot (for example, time changed for a specific date, or first class has been canceled)
func AdjustSlot(ctx context.Context, user facade.User, request dto4calendarium.HappeningSlotDateRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		switch params.Happening.Data.Type {
		case dbo4calendarium.HappeningTypeSingle:
			return errors.New("only recurring happenings can be adjusted, single happenings should be updated")
		case dbo4calendarium.HappeningTypeRecurring:
			if err = adjustRecurringSlot(ctx, tx, params.Happening, request); err != nil {
				return fmt.Errorf("failed to adjust recurring happening: %w", err)
			}
			return err
		}
		return
	}

	if err = dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func adjustRecurringSlot(ctx context.Context, tx dal.ReadwriteTransaction, happening dbo4calendarium.HappeningEntry, request dto4calendarium.HappeningSlotDateRequest) (err error) {
	//for _, teamID := range happening.Data.TeamIDs { // TODO: run in parallel in go routine if > 1
	if err := adjustSlotInCalendarDay(ctx, tx, request.TeamID, happening.ID, request); err != nil {
		return fmt.Errorf("failed to adjust slot in calendar day record for teamID=%v: %w", request.TeamID, err)
	}
	//}
	return nil
}

func adjustSlotInCalendarDay(ctx context.Context, tx dal.ReadwriteTransaction, teamID, happeningID string, request dto4calendarium.HappeningSlotDateRequest) error {
	calendarDay := dbo4calendarium.NewCalendarDayContext(teamID, request.Date)
	if err := tx.Get(ctx, calendarDay.Record); err != nil {
		if !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get calendar day record: %w", err)
		}
	}
	_, adjustment := calendarDay.Data.GetAdjustment(happeningID, request.Slot.ID)
	if adjustment == nil {
		adjustment = &dbo4calendarium.HappeningAdjustment{
			HappeningID: happeningID,
		}
		calendarDay.Data.HappeningAdjustments = append(calendarDay.Data.HappeningAdjustments, adjustment)
	}
	adjustment.Slot = request.Slot
	var happeningIDsChanged bool
	if happeningIDsChanged = slice.Index(calendarDay.Data.HappeningIDs, happeningID) < 0; happeningIDsChanged {
		calendarDay.Data.HappeningIDs = append(calendarDay.Data.HappeningIDs, happeningID)
	}

	if err := calendarDay.Data.Validate(); err != nil {
		return fmt.Errorf("calednar day record is not valid: %w", err)
	}

	if calendarDay.Record.Exists() {
		updates := []dal.Update{
			{Field: "happeningAdjustments", Value: calendarDay.Data.HappeningAdjustments},
		}
		if happeningIDsChanged {
			updates = append(updates, dal.Update{
				Field: "happeningIDs", Value: calendarDay.Data.HappeningIDs,
			})
		}
		if err := tx.Update(ctx, calendarDay.Key, updates); err != nil {
			return fmt.Errorf("failed to update calendar day record with happening adjustment: %w", err)
		}
	} else {
		if err := tx.Insert(ctx, calendarDay.Record); err != nil {
			return fmt.Errorf("failed to insert calendar day record with happening adjustment: %w", err)
		}
	}
	return nil
}
