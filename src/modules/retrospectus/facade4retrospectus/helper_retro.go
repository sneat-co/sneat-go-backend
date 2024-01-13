package facade4retrospectus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
)

var runRetroWorker = func(ctx context.Context, userID string, request facade4meetingus.Request, worker facade4meetingus.Worker) error {
	return facade4meetingus.RunMeetingWorker(ctx, userID, request, MeetingRecordFactory{}, worker)
}
