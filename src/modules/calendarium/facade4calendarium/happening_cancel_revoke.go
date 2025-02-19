package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
)

// RevokeHappeningCancellation marks happening as canceled
func RevokeHappeningCancellation(ctx context.Context, userCtx facade.UserContext, request dto4calendarium.CancelHappeningRequest) (err error) {
	logus.Debugf(ctx, "RevokeHappeningCancellation() %+v", request)
	if err = request.Validate(); err != nil {
		return err
	}

	err = dal4calendarium.RunHappeningSpaceWorker(ctx, userCtx, request.HappeningRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {

			var calendarDay dbo4calendarium.CalendarDayEntry
			{ // Load records that might need to be updated
				var recordsToGet []dal.Record

				if request.Date != "" {
					calendarDay = dbo4calendarium.NewCalendarDayEntry(request.SpaceID, request.Date)
					recordsToGet = append(recordsToGet, calendarDay.Record)
				}
				if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
					recordsToGet = append(recordsToGet, params.SpaceModuleEntry.Record)
				}
				if err = tx.GetMulti(ctx, recordsToGet); err != nil {
					return
				}
			}

			params.HappeningUpdates = params.Happening.Data.RemoveCancellation()
			if len(params.HappeningUpdates) > 0 {
				params.Happening.Record.MarkAsChanged()
			}

			if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
				if err = removeCancellationFromHappeningBriefInSpaceModuleEntry(params); err != nil {
					return fmt.Errorf("failed to remove cancellation from happening brief in team record: %w", err)
				}
			}
			if request.Date != "" {
				if err = removeCancellationFromCalendarDay(ctx, tx, params.Happening.ID, request.SlotID, calendarDay); err != nil {
					return fmt.Errorf("failed to remove cancellation from calendar day record: %w", err)
				}
			}
			return
		})
	if err != nil {
		return fmt.Errorf("failed to revoke happening cancellation: %w", err)
	}
	return
}

func removeCancellationFromHappeningBriefInSpaceModuleEntry(params *dal4calendarium.HappeningWorkerParams) (err error) {
	happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID)
	if happeningBrief == nil {
		return nil
	}

	if updates := happeningBrief.RemoveCancellation(); len(updates) > 0 {
		for _, u := range updates {
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
				update.ByFieldName("recurringHappenings."+params.Happening.ID+"."+u.FieldName(), u.Value()),
			)
		}
		params.SpaceModuleEntry.Record.MarkAsChanged()
	}
	return nil
}

func removeCancellationFromCalendarDay(ctx context.Context, tx dal.ReadwriteTransaction, happeningID string, slotID string, calendarDay dbo4calendarium.CalendarDayEntry) error {
	if happeningAdjustment, slotAdjustment := calendarDay.Data.GetAdjustment(happeningID, slotID); slotAdjustment != nil && slotAdjustment.Cancellation != nil {
		slotAdjustment.Cancellation = nil
		if slotAdjustment.IsEmpty() {
			delete(happeningAdjustment.Slots, slotID)
		}
		if happeningAdjustment.IsEmpty() {
			delete(calendarDay.Data.HappeningAdjustments, happeningID)
		}
		if len(calendarDay.Data.HappeningAdjustments) == 0 {
			if err := tx.Delete(ctx, calendarDay.Key); err != nil {
				return fmt.Errorf("failed to delete calendar day record: %w", err)
			}
		}
	}
	return nil
}
