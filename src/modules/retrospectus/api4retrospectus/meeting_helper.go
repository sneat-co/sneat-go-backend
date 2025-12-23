package api4retrospectus

import (
	"context"
	"errors"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/facade4retrospectus"
)

// meetingParams holds records settings for MeetingID entity
var meetingParams = facade4meetingus.Params{
	RecordFactory: facade4retrospectus.MeetingRecordFactory{},
	BeforeSafe:    beforeSafeRetrospective,
}

var beforeSafeRetrospective = func(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) error {
	retrospective := params.Meeting.Record.Data().(*dbo4retrospectus.Retrospective)
	if retrospective == nil {
		return errors.New("BeforeSafe: retrospective == nil")
	}
	if retrospective.Stage == "" {
		retrospective.Stage = dbo4retrospectus.StageFeedback
	}
	return tx.Set(ctx, params.Meeting.Record)
}
