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
	"github.com/strongo/validation"
)

type PutMode string

const (
	AddSlot    PutMode = "AddSlot"
	UpdateSlot PutMode = "UpdateSlot"
)

func PutSlot(ctx context.Context, userCtx facade.UserContext, putMode PutMode, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}

	worker := func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		return putSlotTxWorker(ctx, tx, params, putMode, request)
	}

	return dal4calendarium.RunHappeningSpaceWorker(ctx, userCtx, request.HappeningRequest, worker)
}

func putSlotTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, putMode PutMode, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}

	if existingSlot := params.Happening.Data.GetSlot(request.Slot.ID); existingSlot == nil && putMode == UpdateSlot {
		return validation.NewErrBadRequestFieldValue("slot.id", "slot not found by SlotID="+request.Slot.ID)
	} else {
		slot := &request.Slot.HappeningSlot
		params.Happening.Record.MarkAsChanged()
		params.Happening.Data.Slots[request.Slot.ID] = slot
		params.HappeningUpdates = []update.Update{update.ByFieldName("slots."+request.Slot.ID, slot)}
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
			happeningBrief.Slots[request.Slot.ID] = &request.Slot.HappeningSlot
			if err = happeningBrief.Validate(); err != nil {
				return fmt.Errorf("happening brief is not valid after update: %w", err)
			}
			params.SpaceModuleEntry.Record.MarkAsChanged()
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldName(
				"recurringHappenings."+params.Happening.ID+".slots",
				happeningBrief.Slots,
			))
		}
	}
	return
}
