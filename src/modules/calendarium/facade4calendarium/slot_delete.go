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
	"github.com/strongo/logus"
	"github.com/strongo/validation"
)

// DeleteSlot deletes happening
func DeleteSlot(ctx context.Context, user facade.User, request dto4calendarium.DeleteHappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	err = dal4teamus.RunModuleTeamWorker(ctx, user, request.TeamRequest,
		const4calendarium.ModuleID,
		new(dbo4calendarium.CalendariumTeamDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo]) (err error) {
			return deleteSlotTxWorker(ctx, tx, params, request)
		})
	if err != nil {
		return fmt.Errorf("failed to delete happening: %w", err)
	}
	return
}

func deleteSlotTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo], request dto4calendarium.DeleteHappeningSlotRequest) (err error) {
	happening := dbo4calendarium.NewHappeningEntry(request.TeamID, request.HappeningID)
	hasHappeningRecord := true
	if err = tx.GetMulti(ctx, []dal.Record{happening.Record, params.TeamModuleEntry.Record}); err != nil {
		return err
	}
	if err = happening.Record.Error(); err != nil {
		if dal.IsNotFound(err) {
			hasHappeningRecord = false
			happening.Data.Type = request.HappeningType
		} else {
			return fmt.Errorf("failed to get happening: %w", err)
		}
	}
	switch happening.Data.Type {
	case "":
		return fmt.Errorf("unknown happening type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
	case dbo4calendarium.HappeningTypeSingle:
		if err := removeSlotFromSingleHappening(ctx, tx, happening, request); err != nil {
			return fmt.Errorf("failed to delete slot from single happening: %w", err)
		}
	case dbo4calendarium.HappeningTypeRecurring:
		if err := removeSlotFromRecurringHappening(ctx, tx, params, happening, request); err != nil {
			return fmt.Errorf("failed to delete slot from recurrign happening: %w", err)
		}
	default:
		return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Data.Type)
	}
	if request.SlotID == "" && hasHappeningRecord || len(happening.Data.Slots) == 0 {
		if err = tx.Delete(ctx, happening.Key); err != nil {
			return fmt.Errorf("faield to delete happening record: %w", err)
		}
	}
	return nil
}

func removeSlotFromSingleHappening(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	happening dbo4calendarium.HappeningEntry,
	request dto4calendarium.DeleteHappeningSlotRequest,
) error {
	if err := removeSlotFromHappeningDbo(ctx, tx, happening, request); err != nil {
		return fmt.Errorf("faile to remove slot from happening record: %w", err)
	}
	return nil
}

func removeSlotFromRecurringHappening(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo],
	happening dbo4calendarium.HappeningEntry,
	request dto4calendarium.DeleteHappeningSlotRequest,
) error {
	if err := removeSlotFromHappeningDbo(ctx, tx, happening, request); err != nil {
		return fmt.Errorf("failed to remove slot from happening record: %w", err)
	}
	if params.TeamModuleEntry.Record.Exists() {
		if err := removeSlotFromHappeningBriefInTeamRecord(params, happening, request); err != nil {
			return fmt.Errorf("failed to remove slot from happening brief in team record: %w", err)
		}
	}
	return nil
}

func removeSlotFromHappeningDbo(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	happening dbo4calendarium.HappeningEntry,
	request dto4calendarium.DeleteHappeningSlotRequest,
) error {
	logus.Debugf(ctx, "removeSlotFromHappeningDbo() %+v", request)
	if len(happening.Data.Slots) == 0 {
		return nil
	}
	var updates []dal.Update
	if request.Weekday == "" {
		removedSlotIDs := removeSlots(happening.Data.Slots, []string{request.SlotID})
		if len(removedSlotIDs) > 0 {
			updates = append(updates, dal.Update{
				Field: "slots." + request.SlotID,
				Value: dal.DeleteField,
			})
			if len(happening.Data.Slots) == 0 {
				happening.Data.Status = dbo4calendarium.HappeningStatusDeleted
				updates = append(updates, dal.Update{
					Field: "status", Value: happening.Data.Status,
				})
			}
		}
	} else {
		slot := happening.Data.GetSlot(request.SlotID)
		if changed := removeWeekday(slot, request.Weekday); changed {
			updates = append(updates, dal.Update{
				Field: "slots",
				Value: happening.Data.Slots,
			})
		}
	}
	if err := happening.Data.Validate(); err != nil {
		return fmt.Errorf("happening record is not valid after removal of slots: %w", err)
	}
	if len(updates) > 0 {
		if err := tx.Update(ctx, happening.Key, updates); err != nil {
			return fmt.Errorf("faile to update happening record: %w", err)
		}
	}
	return nil
}

func removeSlotFromHappeningBriefInTeamRecord(
	params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo],
	happening dbo4calendarium.HappeningEntry,
	request dto4calendarium.DeleteHappeningSlotRequest,
) error {
	brief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(happening.ID)
	if brief == nil {
		return nil
	}
	if request.Weekday == "" {
		removedSlotIDs := removeSlots(brief.Slots, []string{request.SlotID})
		if len(removedSlotIDs) > 0 {
			params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{
				Field: fmt.Sprintf("recurringHappenings.%s.slots.%s", happening.ID, request.SlotID),
				Value: dal.DeleteField,
			})
			if len(happening.Data.Slots) == 0 {
				happening.Data.Status = dbo4calendarium.HappeningStatusDeleted
				params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{
					Field: fmt.Sprintf("recurringHappenings.%s.status", happening.ID),
					Value: happening.Data.Status,
				})
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

func removeSlots(slots map[string]*dbo4calendarium.HappeningSlot, slotIDs []string) (removedSlotIDs []string) {
	for _, slotID := range slotIDs {
		if _, ok := slots[slotID]; ok {
			removedSlotIDs = append(removedSlotIDs, slotID)
			delete(slots, slotID)
		}
	}
	return
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
