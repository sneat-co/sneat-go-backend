package facade4teamus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// RemoveMetrics removes a metric
func RemoveMetrics(ctx context.Context, user facade.User, request dto4teamus.TeamMetricsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4teamus.RunTeamWorker(ctx, user, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.TeamWorkerParams) (err error) {
			changed := false
			team := params.Team

			metrics := make([]*models4teamus.TeamMetric, 0, len(team.Data.Metrics))
		Metrics:
			for _, metric := range team.Data.Metrics {
				for i, metricID := range request.Metrics {
					if metric.ID == metricID {
						changed = true
						request.Metrics = append(request.Metrics[:i], request.Metrics[i+1:]...)
						continue Metrics
					}
				}
				metrics = append(metrics, metric)
			}
			if changed {
				var updates []dal.Update
				if len(metrics) == 0 {
					updates = []dal.Update{
						{Field: "metrics", Value: dal.DeleteField},
					}
				} else {
					updates = []dal.Update{
						{Field: "metrics", Value: metrics},
					}
				}
				if err = dal4teamus.TxUpdateTeam(ctx, tx, params.Started, params.Team, updates); err != nil {
					return err
				}
			}
			return nil
		})
	return
}
