package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
)

// RevokeHappeningCancellation marks happening as canceled
func RevokeHappeningCancellation(ctx context.Context, userCtx facade.UserContext, request dto4calendarium.CancelHappeningRequest) (err error) {
	logus.Debugf(ctx, "RevokeHappeningCancellation() %+v", request)
	if err = request.Validate(); err != nil {
		return err
	}

	err = dal4calendarium.RunHappeningSpaceWorker(ctx, userCtx, request.HappeningRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
			params.HappeningUpdates = params.Happening.Data.RemoveCancellation()
			if len(params.HappeningUpdates) > 0 {
				params.Happening.Record.MarkAsChanged()
			}
			if request.Date != "" {
				if err = removeCancellationFromCalendarDay(ctx, tx, params.Space.ID, request.Date, params.Happening.ID, request.SlotID); err != nil {
					return fmt.Errorf("failed to remove cancellation from calendar day record: %w", err)
				}
			} else if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
				if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
					return
				}
				if err = removeCancellationFromHappeningBriefInSpaceModuleEntry(params); err != nil {
					return fmt.Errorf("failed to remove cancellation from happening brief in team record: %w", err)
				}
			}
			return
		})
	if err != nil {
		return fmt.Errorf("failed to revoke happening cancellation: %w", err)
	}
	return
}

func removeCancellationFromHappeningBriefInSpaceModuleEntry(params *dal4calendarium.HappeningWorkerParams) error {
	happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID)
	if happeningBrief == nil {
		return nil
	}

	if updates := happeningBrief.RemoveCancellation(); len(updates) > 0 {
		for _, update := range updates {
			update.Field = "recurringHappenings." + params.Happening.ID + "." + update.Field
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update)
		}
		params.SpaceModuleEntry.Record.MarkAsChanged()
	}
	return nil
}

func removeCancellationFromCalendarDay(ctx context.Context, tx dal.ReadwriteTransaction, teamID, dateID, happeningID string, slotID string) error {
	logus.Debugf(ctx, "removeCancellationFromCalendarDay(): teamID=%v, dateID=%v, happeningID=%v, slotID=%+v", teamID, dateID, happeningID, slotID)
	if len(slotID) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("slotID")
	}
	calendarDay := dbo4calendarium.NewCalendarDayEntry(teamID, dateID)
	if err := tx.Get(ctx, calendarDay.Record); err != nil {
		if dal.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to get calendar day record by ContactID")
	}
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
