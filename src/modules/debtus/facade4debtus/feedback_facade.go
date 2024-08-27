package facade4debtus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"time"
)

func SaveFeedback(ctx context.Context, tx dal.ReadwriteTransaction, feedbackID int64, feedbackEntity *models4debtus.FeedbackData) (feedback models4debtus.Feedback, user dbo4userus.UserEntry, err error) {
	if ctx == nil {
		panic("ctx == nil")
	}
	logus.Debugf(ctx, "FeedbackDalGae.SaveFeedback(feedbackEntity:%v)", feedbackEntity)
	if feedbackEntity == nil {
		panic("feedbackEntity == nil")
	}
	if feedbackEntity.UserStrID == "" {
		panic("feedbackEntity.UserStrID is empty string")
	}
	if feedbackEntity.Rate == "" {
		panic("feedbackEntity.Rate is empty string")
	}
	feedback = models4debtus.Feedback{FeedbackData: feedbackEntity}
	user = dbo4userus.NewUserEntry(feedbackEntity.UserStrID)
	if err = dal4userus.GetUser(ctx, tx, user); err != nil {
		return
	}
	user.Data.LastFeedbackRate = feedbackEntity.Rate
	if feedbackEntity.Created.IsZero() {
		now := time.Now()
		user.Data.LastFeedbackAt = now
		feedbackEntity.Created = now
	} else {
		user.Data.LastFeedbackAt = feedbackEntity.Created
	}
	if err = tx.SetMulti(ctx, []dal.Record{feedback.Record, user.Record}); err != nil {
		err = fmt.Errorf("failed to put feedback & user entities to datastore: %w", err)
	}
	return
}
