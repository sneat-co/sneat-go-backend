package facade

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

func SaveFeedback(c context.Context, tx dal.ReadwriteTransaction, feedbackID int64, feedbackEntity *models.FeedbackData) (feedback models.Feedback, user models.AppUser, err error) {
	if c == nil {
		panic("c == nil")
	}
	log.Debugf(c, "FeedbackDalGae.SaveFeedback(feedbackEntity:%v)", feedbackEntity)
	if feedbackEntity == nil {
		panic("feedbackEntity == nil")
	}
	if feedbackEntity.UserStrID == "" {
		panic("feedbackEntity.UserStrID is empty string")
	}
	if feedbackEntity.Rate == "" {
		panic("feedbackEntity.Rate is empty string")
	}
	feedback = models.Feedback{FeedbackData: feedbackEntity}
	if user, err = User.GetUserByID(c, tx, feedbackEntity.UserStrID); err != nil {
		return
	}
	user.Data.LastFeedbackRate = feedbackEntity.Rate
	if feedbackEntity.Created.IsZero() {
		now := time.Now()
		user.Data.LastFeedbackAt = now
		feedbackEntity.Created = now
	} else {
		user.Data.LastTransferAt = feedbackEntity.Created
	}
	if err = tx.SetMulti(c, []dal.Record{feedback.Record, user.Record}); err != nil {
		err = fmt.Errorf("failed to put feedback & user entities to datastore: %w", err)
	}
	return
}
