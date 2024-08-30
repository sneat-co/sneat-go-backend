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
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"slices"
	"strings"
	"time"
)

// CancelHappening marks happening as canceled
func CancelHappening(ctx context.Context, userCtx facade.UserContext, request dto4calendarium.CancelHappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	happening := dbo4calendarium.NewHappeningEntry(request.SpaceID, request.HappeningID)
	err = dal4calendarium.RunCalendariumSpaceWorker(ctx, userCtx, request.SpaceRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.CalendariumSpaceWorkerParams) (err error) {
			if err = tx.Get(ctx, happening.Record); err != nil {
				return fmt.Errorf("failed to get happening: %w", err)
			}
			switch happening.Data.Type {
			case "":
				return fmt.Errorf("happening record has no type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
			case dbo4calendarium.HappeningTypeSingle:
				return cancelSingleHappening(ctx, tx, params.UserID, happening)
			case dbo4calendarium.HappeningTypeRecurring:
				return cancelRecurringHappeningInstance(ctx, tx, params, params.UserID, happening, request)
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
		happening.Data.Cancellation = &dbo4calendarium.Cancellation{
			At: time.Now(),
			By: dbmodels.ByUser{UID: userID},
		}
		if err := happening.Data.Validate(); err != nil {
			return fmt.Errorf("happening record is not valid: %w", err)
		}
		happeningUpdates := []dal.Update{
			{Field: "status", Value: happening.Data.Status},
			{Field: "canceled", Value: happening.Data.Cancellation},
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

func cancelRecurringHappeningInstance(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4calendarium.CalendariumSpaceWorkerParams,
	uid string,
	happening dbo4calendarium.HappeningEntry,
	request dto4calendarium.CancelHappeningRequest,
) (err error) {
	if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
		return fmt.Errorf("failed to get team module record: %w", err)
	}
	happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(happening.ID)
	if happeningBrief == nil {
		return errors.New("happening brief is not found in team record")
	}

	cancellation := CreateCancellation(uid, request.Reason)
	if request.Date == "" {
		if err := markRecurringHappeningRecordAsCanceled(ctx, tx, uid, happening, request); err != nil {
			return err
		}
		happeningBrief.Status = dbo4calendarium.HappeningStatusCanceled
		happeningBrief.Cancellation = cancellation
		if err := happeningBrief.Validate(); err != nil {
			return fmt.Errorf("happening brief in team record is not valid: %w", err)
		}
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{
			Field: "recurringHappenings",
			Value: params.SpaceModuleEntry.Data.RecurringHappenings,
		})
	} else {
		if err = addCancellationToHappeningDbo(happening.Data, cancellation); err != nil {
			return err
		}
		if err = addCancellationToCalendarDayAdjustments(ctx, tx, request, params, happening, cancellation); err != nil {
			return err
		}
	}

	return nil
}

func addCancellationToHappeningDbo(
	happeningDbo *dbo4calendarium.HappeningDbo,
	cancellation *dbo4calendarium.Cancellation,
) (err error) {
	if happeningDbo.Cancellation == nil {
		happeningDbo.Cancellation = cancellation
	} else if *happeningDbo.Cancellation != *cancellation {
		happeningDbo.Cancellation = cancellation
	}
	return err
}

func addCancellationToCalendarDayAdjustments(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4calendarium.CancelHappeningRequest,
	params *dal4calendarium.CalendariumSpaceWorkerParams,
	happening dbo4calendarium.HappeningEntry,
	cancellation *dbo4calendarium.Cancellation,
) (err error) {
	calendarDay := dbo4calendarium.NewCalendarDayEntry(params.Space.ID, request.Date)

	var isNewRecord bool
	if err := tx.Get(ctx, calendarDay.Record); err != nil {
		if dal.IsNotFound(err) {
			isNewRecord = true
		} else {
			return fmt.Errorf("failed to get calendar day record by ContactID: %w", err)
		}
	}

	var dayUpdates []dal.Update
	happeningAdjustment, slotAdjustment := calendarDay.Data.GetAdjustment(happening.ID, request.SlotID)
	if slotAdjustment == nil {
		slot := happening.Data.GetSlot(request.SlotID)
		if slot == nil {
			return fmt.Errorf("%w: slot not found by SlotID=%v", facade.ErrBadRequest, request.SlotID)
		}
		slotAdjustment = &dbo4calendarium.SlotAdjustment{
			Cancellation: cancellation,
		}
		if happeningAdjustment == nil {
			happeningAdjustment = new(dbo4calendarium.HappeningAdjustment)
			if calendarDay.Data.HappeningAdjustments == nil {
				calendarDay.Data.HappeningAdjustments = make(map[string]*dbo4calendarium.HappeningAdjustment, 1)
			}
			calendarDay.Data.HappeningAdjustments[happening.ID] = happeningAdjustment
		}

		if happeningAdjustment.Slots == nil {
			happeningAdjustment.Slots = make(map[string]*dbo4calendarium.SlotAdjustment, 1)
		}
		happeningAdjustment.Slots[request.SlotID] = slotAdjustment
	}
	if !slices.Contains(calendarDay.Data.HappeningIDs, happening.ID) {
		calendarDay.Data.HappeningIDs = append(calendarDay.Data.HappeningIDs, happening.ID)
		dayUpdates = append(dayUpdates, dal.Update{
			Field: "happeningIDs", Value: calendarDay.Data.HappeningIDs,
		})
	}
	var modified bool
	if slotAdjustment.Cancellation == nil || *slotAdjustment.Cancellation != *cancellation {
		modified = true
		slotAdjustment.Cancellation = cancellation
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
	return err
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
	if happening.Data.Cancellation == nil {
		happening.Data.Cancellation = CreateCancellation(uid, request.Reason)
	} else if reason := strings.TrimSpace(request.Reason); reason != "" {
		happening.Data.Cancellation.Reason = reason
	}
	happeningUpdates = append(happeningUpdates,
		dal.Update{
			Field: "status",
			Value: happening.Data.Status,
		},
		dal.Update{
			Field: "canceled",
			Value: happening.Data.Cancellation,
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
