package facade4retrospectus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
)

var runRetroWorker = func(ctx context.Context, userCtx facade.UserContext, request facade4meetingus.Request, worker facade4meetingus.Worker) error {
	return facade4meetingus.RunMeetingWorker(ctx, userCtx, request, MeetingRecordFactory{}, worker)
}
