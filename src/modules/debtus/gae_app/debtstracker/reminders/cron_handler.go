package reminders

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal/gaedal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/http"
	"reflect"
	"time"

	"context"
)

func CronSendReminders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	query := dal.From(models4debtus.ReminderKind).
		WhereField("Status", dal.Equal, models4debtus.ReminderStatusCreated).
		WhereField("DtNext", dal.GreaterThen, time.Time{}).
		WhereField("DtNext", dal.LessThen, time.Now()).
		OrderBy(dal.AscendingField("DtNext")).
		Limit(100).
		SelectKeysOnly(reflect.Int)

	db, err := facade.GetSneatDB(ctx)
	if err != nil {
		logus.Errorf(ctx, "Failed to get database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reader dal.Reader
	if reader, err = db.QueryReader(ctx, query); err != nil {
		logus.Errorf(ctx, "Failed to load due api4transfers: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reminderIDs []string
	if reminderIDs, err = dal.SelectAllIDs[string](reader, dal.WithLimit(query.Limit())); err != nil {
		logus.Errorf(ctx, "Failed to load due api4transfers: %v", err)
		return
	}

	if len(reminderIDs) == 0 {
		logus.Debugf(ctx, "No reminders to send")
		return
	}

	logus.Debugf(ctx, "Loaded %d reminder(s)", len(reminderIDs))

	for _, reminderID := range reminderIDs {
		/*task,*/ err = gaedal.CreateSendReminderTask(ctx, reminderID)
		panic(fmt.Errorf("TODO: implement CreateSendReminderTask: %w", err))
		//task.Name = fmt.Sprintf("r_%s_%s", reminderID, time.Now().Format("200601021504"))
		//if _, err := apphostgae.AddTaskToQueue(ctx, task, queues.QueueReminders); err != nil {
		//	logus.Errorf(ctx, "Failed to add delayed task for reminder %s", reminderID)
		//	return
		//}
	}

	w.WriteHeader(http.StatusOK)
}
