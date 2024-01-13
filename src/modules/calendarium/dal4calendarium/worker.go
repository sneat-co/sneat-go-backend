package dal4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type CalendariumTeamWorkerParams = dal4teamus.ModuleTeamWorkerParams[*models4calendarium.CalendariumTeamDto]

func NewCalendariumTeamWorkerParams(userID, teamID string) *CalendariumTeamWorkerParams {
	teamWorkerParams := dal4teamus.NewTeamWorkerParams(userID, teamID)
	return dal4teamus.NewTeamModuleWorkerParams(const4calendarium.ModuleID, teamWorkerParams, new(models4calendarium.CalendariumTeamDto))
}

type HappeningWorkerParams struct {
	CalendariumTeamWorkerParams
	Happening models4calendarium.HappeningContext
}

func RunHappeningTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4calendarium.HappeningRequest,
	moduleID string,
	happeningWorker func(ctx context.Context, tx dal.ReadwriteTransaction, params *HappeningWorkerParams) (err error),
) (err error) {
	calendariumTeamDto := new(models4calendarium.CalendariumTeamDto)

	moduleTeamWorker := func(
		ctx context.Context,
		tx dal.ReadwriteTransaction,
		moduleTeamParams *dal4teamus.ModuleTeamWorkerParams[*models4calendarium.CalendariumTeamDto],
	) (err error) {
		params := &HappeningWorkerParams{
			CalendariumTeamWorkerParams: *moduleTeamParams,
			Happening:                   models4calendarium.NewHappeningContext(request.TeamID, request.HappeningID),
		}
		if err = tx.Get(ctx, params.Happening.Record); err != nil {
			if dal.IsNotFound(err) {
				params.Happening.Dto.Type = request.HappeningType
			} else {
				return fmt.Errorf("failed to get happening: %w", err)
			}
		}

		return happeningWorker(ctx, tx, params)
	}
	return dal4teamus.RunModuleTeamWorker(ctx, user, request.TeamRequest, moduleID, calendariumTeamDto, moduleTeamWorker)
}
