package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type FeedbackDalGae struct {
}

func NewFeedbackDalGae() FeedbackDalGae {
	return FeedbackDalGae{}
}

func (FeedbackDalGae) GetFeedbackByID(c context.Context, tx dal.ReadSession, feedbackID int64) (feedback models.Feedback, err error) {
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	feedback = models.NewFeedback(feedbackID, nil)
	return feedback, tx.Get(c, feedback.Record)
}
