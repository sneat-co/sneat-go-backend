package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"log"
)

// RevokeHappeningCancellation marks happening as canceled
func RevokeHappeningCancellation(ctx context.Context, user facade.User, request dto4calendarium.CancelHappeningRequest) (err error) {
	log.Printf("RevokeHappeningCancellation() %+v", request)
	if err = request.Validate(); err != nil {
		return err
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
				return revokeSingleHappeningCancellation(ctx, tx, happening)
			case dbo4calendarium.HappeningTypeRecurring:
				return revokeRecurringHappeningCancellation(ctx, tx, params, happening, request.Date, request.SlotID)
			default:
				return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Data.Type)
			}
		})
	if err != nil {
		return fmt.Errorf("failed to cancel happening: %w", err)
	}
	return
}

func revokeSingleHappeningCancellation(ctx context.Context, tx dal.ReadwriteTransaction, happening dbo4calendarium.HappeningEntry) error {
	return removeCancellationFromHappeningRecord(ctx, tx, happening)
}

func revokeRecurringHappeningCancellation(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo],
	happening dbo4calendarium.HappeningEntry,
	dateID string,
	slotID string,
) error {
	log.Printf("revokeRecurringHappeningCancellation(): teamID=%v, dateID=%v, happeningID=%v, slotID=%+v", params.Team.ID, dateID, happening.ID, slotID)
	if happening.Data.Status == dbo4calendarium.HappeningStatusCanceled {
		if err := removeCancellationFromHappeningRecord(ctx, tx, happening); err != nil {
			return fmt.Errorf("failed to remove cancellation from happening record: %w", err)
		}
	}
	if dateID == "" {
		if err := removeCancellationFromHappeningBrief(params, happening); err != nil {
			return fmt.Errorf("failed to remove cancellation from happening brief in team record: %w", err)
		}
	} else if err := removeCancellationFromCalendarDay(ctx, tx, params.Team.ID, dateID, happening.ID, slotID); err != nil {
		return fmt.Errorf("failed to remove cancellation from calendar day record: %w", err)
	}
	return nil
}

func removeCancellationFromHappeningBrief(params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo], happening dbo4calendarium.HappeningEntry) error {
	happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(happening.ID)
	if happeningBrief == nil {
		return nil
	}
	if happeningBrief.Status == dbo4calendarium.HappeningStatusCanceled {
		happeningBrief.Status = dbo4calendarium.HappeningStatusActive
		happeningBrief.Canceled = nil
		if err := happeningBrief.Validate(); err != nil {
			return err
		}
		params.TeamUpdates = append(params.TeamUpdates, dal.Update{
			Field: "recurringHappenings",
			Value: params.TeamModuleEntry.Data.RecurringHappenings,
		})
	}
	return nil
}

func removeCancellationFromHappeningRecord(ctx context.Context, tx dal.ReadwriteTransaction, happening dbo4calendarium.HappeningEntry) error {
	if happening.Data.Status != dbo4calendarium.HappeningStatusCanceled {
		return fmt.Errorf("not allowed to revoke cancelation for happening in status=" + happening.Data.Status)
	}
	happening.Data.Status = dbo4calendarium.HappeningStatusCanceled
	happening.Data.Canceled = nil
	if err := happening.Data.Validate(); err != nil {
		return err
	}
	updates := []dal.Update{
		{Field: "status", Value: dbo4calendarium.HappeningStatusActive},
		{Field: "canceled", Value: dal.DeleteField},
	}
	if err := tx.Update(ctx, happening.Key, updates); err != nil {
		return fmt.Errorf("failed to update happening record: %w", err)
	}
	return nil

}

func removeCancellationFromCalendarDay(ctx context.Context, tx dal.ReadwriteTransaction, teamID, dateID, happeningID string, slotID string) error {
	log.Printf("removeCancellationFromCalendarDay(): teamID=%v, dateID=%v, happeningID=%v, slotID=%+v", teamID, dateID, happeningID, slotID)
	if len(slotID) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("slotID")
	}
	calendarDay := dbo4calendarium.NewCalendarDayContext(teamID, dateID)
	if err := tx.Get(ctx, calendarDay.Record); err != nil {
		if dal.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to get calendar day record by ContactID")
	}
	if i, adjustment := calendarDay.Dto.GetAdjustment(happeningID, slotID); adjustment != nil && adjustment.Canceled != nil {
		a := calendarDay.Dto.HappeningAdjustments
		calendarDay.Dto.HappeningAdjustments = append(a[:i], a[i+1:]...)
		if len(calendarDay.Dto.HappeningAdjustments) == 0 {
			if err := tx.Delete(ctx, calendarDay.Key); err != nil {
				return fmt.Errorf("failed to delete calendar day record: %w", err)
			}
		}
	}
	return nil
}
