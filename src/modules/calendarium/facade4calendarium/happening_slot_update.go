package facade4calendarium

import (
	"context"
	"errors"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

type PutMode string

const (
	AddSlot    PutMode = "AddSlot"
	UpdateSlot PutMode = "UpdateSlot"
)

func UpdateHappeningSlot(ctx facade.ContextWithUser, putMode PutMode, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}

	worker := func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		return updateHappeningSlotTxWorker(ctx, tx, params, putMode, request)
	}

	return dal4calendarium.RunHappeningSpaceWorker(ctx, request.HappeningRequest, worker)
}

func updateHappeningSlotTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, putMode PutMode, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}

	slot := &request.Slot.HappeningSlot

	if existingSlot := params.Happening.Data.GetSlot(request.Slot.ID); existingSlot == nil && putMode == UpdateSlot {
		return validation.NewErrBadRequestFieldValue("slot.id", "slot not found by SlotID="+request.Slot.ID)
	} else {
		params.Happening.Data.Slots[request.Slot.ID] = slot
		params.Happening.Record.MarkAsChanged()
		params.HappeningUpdates = []update.Update{
			update.ByFieldPath([]string{dbo4calendarium.SlotsField, request.Slot.ID}, slot),
		}
		if params.Happening.Data.Type == dbo4calendarium.HappeningTypeSingle {
			//TODO: For single happening we need to update the dates fields
			return errors.New("needs implementation for updating dates fields for 'single' type happening")
			//params.HappeningUpdates = append(params.HappeningUpdates,
			//	params.Happening.Data.UpdatesWhenDatesChanged()...)
		}
	}

	if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
		if happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID); happeningBrief != nil {
			if err = happeningBrief.Validate(); err != nil {
				return fmt.Errorf("happening brief is not valid before update: %w", err)
			}
			switch putMode {
			case AddSlot:
				if happeningBrief.HasSlot(request.Slot.ID) {
					return validation.NewErrBadRequestFieldValue("slotID",
						"happening already have slot with ContactID="+request.Slot.ID)
				}
			case UpdateSlot:
				if !happeningBrief.HasSlot(request.Slot.ID) {
					return validation.NewErrBadRequestFieldValue("slotID", "slot not found by SlotID="+request.Slot.ID)
				}
			default:
				return fmt.Errorf("unsupported put mode: %v", putMode)
			}
			happeningBrief.Slots[request.Slot.ID] = slot
			if err = happeningBrief.Validate(); err != nil {
				return fmt.Errorf("happening brief is not valid after update: %w", err)
			}
			params.SpaceModuleEntry.Record.MarkAsChanged()
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath(
				[]string{dbo4calendarium.RecurringHappeningsField, params.Happening.ID, dbo4calendarium.SlotsField, request.Slot.ID},
				slot,
			))
		}
	}
	return
}
