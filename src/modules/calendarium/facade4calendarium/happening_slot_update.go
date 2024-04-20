package facade4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateSlot(ctx context.Context, user facade.User, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *happeningWorkerParams) (err error) {
		//teamKey := models4teamus.NewTeamKey(request.TeamID)
		//teamDto := new(models4teamus.TeamDto)
		//teamRecord := dal.NewRecordWithData(teamKey, teamDto)
		//
		//if err = tx.Get(ctx, teamRecord); err != nil {
		//	return nil, fmt.Errorf("failed to get team record: %w", err)
		//}

		if params.Happening.Dbo.Type == "single" {
			params.Happening.Dbo.Slots[0] = &request.Slot
			params.HappeningUpdates = []dal.Update{
				{
					Field: "slots",
					Value: params.Happening.Dbo.Slots,
				},
			}
		}
		return
	}

	if err = modifyHappening(ctx, user, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}
