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

type PutMode int

const (
	AddSlot PutMode = iota
	UpdateSlot
)

func PutSlot(ctx context.Context, putMode PutMode, user facade.User, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}

	worker := func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		if err = params.GetRecords(ctx, tx); err != nil {
			return err
		}

		if existingSlot := params.Happening.Data.GetSlot(request.Slot.ID); existingSlot == nil && putMode == UpdateSlot {
			return validation.NewErrBadRequestFieldValue("slot.id", "slot not found by ID="+request.Slot.ID)
		} else {
			slot := &request.Slot.HappeningSlot
			params.Happening.Data.Slots[request.Slot.ID] = slot
			params.HappeningUpdates = []dal.Update{
				{
					Field: "slots." + request.Slot.ID,
					Value: slot,
				},
			}
		}

		if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
			if happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(params.Happening.ID); happeningBrief != nil {
				if err = happeningBrief.Validate(); err != nil {
					return fmt.Errorf("happening brief is not valid before update: %w", err)
				}
				if slot := happeningBrief.GetSlot(request.Slot.ID); slot == nil {
					return validation.NewErrBadRequestFieldValue("slotID", "exostingSlot not found by ID="+request.Slot.ID)
				}
				happeningBrief.Slots[request.Slot.ID] = &request.Slot.HappeningSlot
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
