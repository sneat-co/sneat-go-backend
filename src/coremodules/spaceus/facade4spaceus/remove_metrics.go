package facade4spaceus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// RemoveMetrics removes a metric
func RemoveMetrics(ctx context.Context, userCtx facade.UserContext, request dto4spaceus.SpaceMetricsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4spaceus.RunSpaceWorkerWithUserContext(ctx, userCtx, request.SpaceID,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceWorkerParams) (err error) {
			changed := false
			team := params.Space

			metrics := make([]*dbo4spaceus.SpaceMetric, 0, len(team.Data.Metrics))
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
				if err = dal4spaceus.TxUpdateSpace(ctx, tx, params.Started, params.Space, updates); err != nil {
					return err
				}
			}
			return nil
		})
	return
}
