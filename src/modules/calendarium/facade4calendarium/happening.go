package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"log"
)

type happeningWorkerParams struct {
	*dal4calendarium.CalendariumTeamWorkerParams
	Happening        models4calendarium.HappeningContext
	HappeningUpdates []dal.Update
}

type happeningWorker = func(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	param *happeningWorkerParams,
) (err error)

func modifyHappening(ctx context.Context, user facade.User, request dto4calendarium.HappeningRequest, worker happeningWorker) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4calendarium.RunCalendariumTeamWorker(ctx, user, request.TeamRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.CalendariumTeamWorkerParams) error {
		happeningParams := happeningWorkerParams{
			CalendariumTeamWorkerParams: params,
			Happening:                   models4calendarium.NewHappeningContext(request.TeamID, request.HappeningID),
		}
		if err = worker(ctx, tx, &happeningParams); err != nil {
			return fmt.Errorf("failed in happening worker: %w", err)
		}
		if len(happeningParams.HappeningUpdates) > 0 {
			if err = happeningParams.Happening.Dto.Validate(); err != nil {
				return fmt.Errorf("happening record is not valid after running worker: %w", err)
			}
			log.Printf("updating happening: %s", happeningParams.Happening.Key)
			if err = tx.Update(ctx, happeningParams.Happening.Key, happeningParams.HappeningUpdates); err != nil {
				return fmt.Errorf("failed to update happening record: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update happening in transaction: %w", err)
	}
	return err
}
