package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type FeedbackDalGae struct {
}

func NewFeedbackDalGae() FeedbackDalGae {
	return FeedbackDalGae{}
}

func (FeedbackDalGae) GetFeedbackByID(ctx context.Context, tx dal.ReadSession, feedbackID int64) (feedback models4debtus.Feedback, err error) {
	if tx == nil {
		if tx, err = facade.GetDatabase(ctx); err != nil {
			return
		}
	}
	feedback = models4debtus.NewFeedback(feedbackID, nil)
	return feedback, tx.Get(ctx, feedback.Record)
}
