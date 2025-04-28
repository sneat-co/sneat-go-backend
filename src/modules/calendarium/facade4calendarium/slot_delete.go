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
	"github.com/strongo/validation"
)

// DeleteSlot deletes happening
func DeleteSlot(ctx facade.ContextWithUser, request dto4calendarium.DeleteHappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4calendarium.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
			return deleteSlotTxWorker(ctx, tx, params, request)
		})
	if err != nil {
		return fmt.Errorf("failed to delete happening: %w", err)
	}
	return
}

func deleteSlotTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.DeleteHappeningSlotRequest) (err error) {
	switch params.Happening.Data.Type {
	case "":
		return fmt.Errorf("unknown happening type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
	case dbo4calendarium.HappeningTypeSingle:
		removeSlotFromHappeningDbo(ctx, params, request)
	case dbo4calendarium.HappeningTypeRecurring:
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
	params *dal4calendarium.HappeningWorkerParams,
	request dto4calendarium.DeleteHappeningSlotRequest,
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
	params *dal4calendarium.HappeningWorkerParams,
	request dto4calendarium.DeleteHappeningSlotRequest,
) {
	logus.Debugf(ctx, "removeSlotFromHappeningDbo() %+v", request)
	if len(params.Happening.Data.Slots) == 0 {
		return
	}
	if request.Weekday == "" {
		if params.Happening.Data.Slots[request.SlotID] != nil {
			if len(params.Happening.Data.Slots) == 1 {
				params.Happening.Data.Status = dbo4calendarium.HappeningStatusDeleted

				params.HappeningUpdates = append(params.HappeningUpdates,
					update.ByFieldName("status", params.Happening.Data.Status))
			} else {
				params.HappeningUpdates = append(params.HappeningUpdates,
					update.ByFieldName("slots."+request.SlotID, update.DeleteField))
			}
			params.Happening.Record.MarkAsChanged()
		}
	} else {
		slot := params.Happening.Data.GetSlot(request.SlotID)
		if changed := removeWeekday(slot, request.Weekday); changed {
			params.HappeningUpdates = append(params.HappeningUpdates,
				update.ByFieldName("slots", params.Happening.Data.Slots))
			params.Happening.Record.MarkAsChanged()
		}
	}
}

func removeSlotFromHappeningBriefInSpaceRecord(
	params *dal4calendarium.HappeningWorkerParams,
	request dto4calendarium.DeleteHappeningSlotRequest,
) error {
	brief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID)
	if brief == nil {
		return nil
	}
	if request.Weekday == "" {
		if brief.Slots[request.SlotID] != nil {
			if len(brief.Slots) == 1 {
				if brief.Status != dbo4calendarium.HappeningStatusDeleted {
					brief.Status = dbo4calendarium.HappeningStatusDeleted

					params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldName(
						fmt.Sprintf("recurringHappenings.%s.status", params.Happening.ID),
						params.Happening.Data.Status,
					))
					params.SpaceModuleEntry.Record.MarkAsChanged()
				}
			} else {
				delete(brief.Slots, request.SlotID)

				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldName(
					fmt.Sprintf("recurringHappenings.%s.slots.%s", params.Happening.ID, request.SlotID),
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

func removeWeekday(slot *dbo4calendarium.HappeningSlot, weekday dbo4calendarium.WeekdayCode) (changed bool) {
	weekdays := make([]dbo4calendarium.WeekdayCode, 0, len(slot.Weekdays))
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
