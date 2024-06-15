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

		if slotIndex, _ := params.Happening.Data.GetSlot(request.Slot.ID); slotIndex < 0 {
			return validation.NewErrBadRequestFieldValue("slot.id", "slot not found by ID="+request.Slot.ID)
		} else {
			params.Happening.Data.Slots[slotIndex] = &request.Slot
			params.HappeningUpdates = []dal.Update{
				{
					Field: "slots",
					Value: params.Happening.Data.Slots,
				},
			}
		}

		if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
			if happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID); happeningBrief != nil {
				if slotIndex, _ := happeningBrief.GetSlot(request.Slot.ID); slotIndex >= 0 {
					if err = happeningBrief.Validate(); err != nil {
						return fmt.Errorf("happening brief is not valid before update: %w", err)
					}
					happeningBrief.Slots[slotIndex] = &request.Slot
				} else {
					happeningBrief.Slots = append(happeningBrief.Slots, &request.Slot)
				}
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
