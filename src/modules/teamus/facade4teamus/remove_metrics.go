package facade4teamus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// RemoveMetrics removes a metric
func RemoveMetrics(ctx context.Context, user facade.User, request dto4teamus.SpaceMetricsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4teamus.RunSpaceWorker(ctx, user, request.SpaceID,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.SpaceWorkerParams) (err error) {
			changed := false
			team := params.Space

			metrics := make([]*dbo4teamus.SpaceMetric, 0, len(team.Data.Metrics))
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
				if err = dal4teamus.TxUpdateSpace(ctx, tx, params.Started, params.Space, updates); err != nil {
					return err
				}
			}
			return nil
		})
	return
}
