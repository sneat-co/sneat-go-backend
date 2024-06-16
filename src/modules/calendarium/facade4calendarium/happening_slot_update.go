package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

func UpdateSlot(ctx context.Context, user facade.User, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}

	worker := func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		if err = params.GetRecords(ctx, tx); err != nil {
			return err
		}

		if slot := params.Happening.Data.GetSlot(request.SlotID); slot == nil {
			return validation.NewErrBadRequestFieldValue("slot.id", "slot not found by ID="+request.SlotID)
		} else {
			slot := &request.Slot
			params.Happening.Data.Slots[request.SlotID] = slot
			params.HappeningUpdates = []dal.Update{
				{
					Field: "slots." + request.SlotID,
					Value: slot,
				},
			}
		}

		if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
			if happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID); happeningBrief != nil {
				if err = happeningBrief.Validate(); err != nil {
					return fmt.Errorf("happening brief is not valid before update: %w", err)
				}
				if slot := happeningBrief.GetSlot(request.SlotID); slot == nil {
					return validation.NewErrBadRequestFieldValue("slotID", "slot not found by ID="+request.SlotID)
				}
				happeningBrief.Slots[request.SlotID] = &request.Slot
				if err = happeningBrief.Validate(); err != nil {
					return fmt.Errorf("happening brief is not valid after update: %w", err)
				}
				params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{
					Field: "recurringHappenings." + params.Happening.ID + ".slots",
					Value: happeningBrief.Slots,
				})
			}
		}
		return
	}

	return dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, worker)
}
