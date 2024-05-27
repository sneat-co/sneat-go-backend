package facade4teamus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/gosimple/slug"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
)

// AddTeamMetricRequest request
type AddTeamMetricRequest struct {
	dto4teamus.TeamRequest
	Metric models4teamus.TeamMetric `json:"metric"`
}

// Validate validates request
func (v *AddTeamMetricRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if err := v.Metric.Validate(); err != nil {
		return err
	}
	return nil
}

// AddMetric adds metric
func AddMetric(ctx context.Context, user facade.User, request AddTeamMetricRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4teamus.RunTeamWorker(ctx, user, request.TeamID, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.TeamWorkerParams) (err error) {
		request.Metric.ID = strings.Replace(slug.Make(request.Metric.Title), "-", "_", -1)
		for _, m := range params.Team.Data.Metrics {
			if m.ID == request.Metric.ID {
				err = validation.NewErrBadRequestFieldValue("title", "duplicate slug(title)")
				return
			}
		}
		params.Team.Data.Metrics = append(params.Team.Data.Metrics, &request.Metric)
		if err = dal4teamus.TxUpdateTeam(ctx, tx, params.Started, params.Team, []dal.Update{
			{Field: "metrics", Value: params.Team.Data.Metrics},
		}); err != nil {
			return err
		}
		return nil
	})
	return
}
