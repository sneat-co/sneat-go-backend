package facade4calendarius

import (
	"context"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dal4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
)

// DeleteSlot deletes happening
func DeleteSlot(ctx facade.ContextWithUser, request dto4calendarius.DeleteHappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) (err error) {
			return deleteSlotTxWorker(ctx, tx, params, request)
		})
	if err != nil {
		return fmt.Errorf("failed to delete happening: %w", err)
	}
	return
}

func deleteSlotTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams, request dto4calendarius.DeleteHappeningSlotRequest) (err error) {
	switch params.Happening.Data.Type {
	case "":
		return fmt.Errorf("unknown happening type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
	case dbo4calendarius.HappeningTypeSingle:
		removeSlotFromHappeningDbo(ctx, params, request)
	case dbo4calendarius.HappeningTypeRecurring:
		if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
			return
		}
		if err = removeSlotFromRecurringHappening(ctx, params, request); err != nil {
			return fmt.Errorf("failed to delete slot from recurrign happening: %w", err)
		}
	default:
		return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+params.Happening.Data.Type)
	}
	//if request.SlotID == "" && len(params.Happening.Data.Slots) == 0 {
	//	if err = tx.Delete(ctx, params.Happening.Key); err != nil {
	//		return fmt.Errorf("faield to delete happening record: %w", err)
	//	}
	//}
	return nil
}

func removeSlotFromRecurringHappening(
	ctx context.Context,
	params *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.DeleteHappeningSlotRequest,
) (err error) {
	removeSlotFromHappeningDbo(ctx, params, request)
	if params.SpaceModuleEntry.Record.Exists() {
		if err = removeSlotFromHappeningBriefInSpaceRecord(params, request); err != nil {
			return fmt.Errorf("failed to remove slot from happening brief in team record: %w", err)
		}
	}
	return nil
}

func removeSlotFromHappeningDbo(
	ctx context.Context,
	params *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.DeleteHappeningSlotRequest,
) {
	logus.Debugf(ctx, "removeSlotFromHappeningDbo() %+v", request)
	if len(params.Happening.Data.Slots) == 0 {
		return
	}
	if request.Weekday == "" {
		if params.Happening.Data.Slots[request.SlotID] != nil {
			if len(params.Happening.Data.Slots) == 1 {
				params.Happening.Data.Status = dbo4calendarius.HappeningStatusDeleted

				params.HappeningUpdates = append(params.HappeningUpdates,
					update.ByFieldName("status", params.Happening.Data.Status))
			} else {
				params.HappeningUpdates = append(params.HappeningUpdates,
					update.ByFieldPath([]string{dbo4calendarius.SlotsField, request.SlotID}, update.DeleteField))
			}
			params.Happening.Record.MarkAsChanged()
		}
	} else {
		slot := params.Happening.Data.GetSlot(request.SlotID)
		if changed := removeWeekday(slot, request.Weekday); changed {
			params.HappeningUpdates = append(params.HappeningUpdates,
				update.ByFieldName(dbo4calendarius.SlotsField, params.Happening.Data.Slots))
			params.Happening.Record.MarkAsChanged()
		}
	}
}

func removeSlotFromHappeningBriefInSpaceRecord(
	params *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.DeleteHappeningSlotRequest,
) error {
	brief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID)
	if brief == nil {
		return nil
	}
	if request.Weekday == "" {
		if brief.Slots[request.SlotID] != nil {
			if len(brief.Slots) == 1 {
				if brief.Status != dbo4calendarius.HappeningStatusDeleted {
					brief.Status = dbo4calendarius.HappeningStatusDeleted

					params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath(
						[]string{dbo4calendarius.RecurringHappeningsField, params.Happening.ID, "status"},
						params.Happening.Data.Status,
					))
					params.SpaceModuleEntry.Record.MarkAsChanged()
				}
			} else {
				delete(brief.Slots, request.SlotID)

				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath(
					[]string{dbo4calendarius.RecurringHappeningsField, params.Happening.ID, dbo4calendarius.SlotsField, request.SlotID},
					update.DeleteField,
				))
				params.SpaceModuleEntry.Record.MarkAsChanged()
			}
		}
	} else {
		slot := brief.GetSlot(request.SlotID)
		if slot == nil {
			return nil
		}
		if changed := removeWeekday(slot, request.Weekday); !changed {
			return nil
		}
	}
	return nil
}

func removeWeekday(slot *dbo4calendarius.HappeningSlot, weekday dbo4calendarius.WeekdayCode) (changed bool) {
	weekdays := make([]dbo4calendarius.WeekdayCode, 0, len(slot.Weekdays))
	for _, wd := range slot.Weekdays {
		if wd != weekday {
			weekdays = append(weekdays, wd)
		}
	}
	if changed = len(weekdays) != len(slot.Weekdays); changed {
		slot.Weekdays = weekdays
	}
	return
}
