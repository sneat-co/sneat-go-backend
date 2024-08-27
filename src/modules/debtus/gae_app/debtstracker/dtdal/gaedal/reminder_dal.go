package gaedal

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"reflect"
	"time"

	"context"
)

func NewReminderIncompleteKey(_ context.Context) *dal.Key {
	return dal.NewIncompleteKey(models4debtus.ReminderKind, reflect.Int, nil)
}

func NewReminderKey(reminderID string) *dal.Key {
	if reminderID == "" {
		panic("reminderID == 0")
	}
	return dal.NewKeyWithID(models4debtus.ReminderKind, reminderID)
}

type ReminderDalGae struct {
}

func NewReminderDalGae() ReminderDalGae {
	return ReminderDalGae{}
}

var _ dtdal.ReminderDal = (*ReminderDalGae)(nil)

func (reminderDalGae ReminderDalGae) GetReminderByID(ctx context.Context, tx dal.ReadSession, id string) (reminder models4debtus.Reminder, err error) {
	reminder = models4debtus.NewReminder(id, nil)
	return reminder, tx.Get(ctx, reminder.Record)
}

func (reminderDalGae ReminderDalGae) SaveReminder(ctx context.Context, tx dal.ReadwriteTransaction, reminder models4debtus.Reminder) (err error) {
	return tx.Set(ctx, reminder.Record)
}

func (reminderDalGae ReminderDalGae) GetSentReminderIDsByTransferID(ctx context.Context, tx dal.ReadSession, transferID int) ([]int, error) {
	q := dal.From(models4debtus.ReminderKind).Where(
		dal.WhereField("TransferID", dal.Equal, transferID),
		dal.WhereField("Status", dal.Equal, models4debtus.ReminderStatusSent),
	).SelectKeysOnly(reflect.Int)

	records, err := tx.QueryAllRecords(ctx, q)
	if err != nil {
		return nil, err
	}
	reminderIDs := make([]int, len(records))
	for i, record := range records {
		reminderIDs[i] = record.Key().ID.(int)
	}
	return reminderIDs, nil
}

func (reminderDalGae ReminderDalGae) GetActiveReminderIDsByTransferID(ctx context.Context, tx dal.ReadSession, transferID int) ([]int, error) {
	q := dal.From(models4debtus.ReminderKind).Where(
		dal.WhereField("TransferID", dal.Equal, transferID),
		dal.WhereField("DtNext", dal.GreaterThen, time.Time{}),
	).SelectKeysOnly(reflect.Int)
	records, err := tx.QueryAllRecords(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reminders by transfer id=%v: %w", transferID, err)
	}
	reminderIDs := make([]int, len(records))
	for i, record := range records {
		reminderIDs[i] = record.Key().ID.(int)
	}
	return reminderIDs, nil
}

func (reminderDalGae ReminderDalGae) SetReminderIsSent(ctx context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error) {
	//gaehost.GaeLogger.Debugf(ctx, "delayedSetReminderIsSent(reminderID=%v, sentAt=%v, messageIntID=%v, messageStrID=%v)", reminderID, sentAt, messageIntID, messageStrID)
	if err := _validateSetReminderIsSentMessageIDs(messageIntID, messageStrID, sentAt); err != nil {
		return err
	}
	reminder := models4debtus.NewReminder(reminderID, nil)
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return reminderDalGae.SetReminderIsSentInTransaction(ctx, tx, reminder, sentAt, messageIntID, messageStrID, locale, errDetails)
	})
}

func (reminderDalGae ReminderDalGae) SetReminderIsSentInTransaction(ctx context.Context, tx dal.ReadwriteTransaction, reminder models4debtus.Reminder, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error) {
	if reminder.Data == nil {
		reminder, err = reminderDalGae.GetReminderByID(ctx, tx, reminder.ID)
		if err != nil {
			if dal.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("failed to get reminder by ContactID: %w", err)
		}
	}
	if reminder.Data.Status != models4debtus.ReminderStatusSending {
		logus.Errorf(ctx, "reminder.Status:%v != models.ReminderStatusSending:%v", reminder.Data.Status, models4debtus.ReminderStatusSending)
		return nil
	} else {
		reminder.Data.Status = models4debtus.ReminderStatusSent
		reminder.Data.DtSent = sentAt
		reminder.Data.DtScheduled = reminder.Data.DtNext
		reminder.Data.DtNext = time.Time{}
		reminder.Data.ErrDetails = errDetails
		reminder.Data.Locale = locale
		if messageIntID != 0 {
			reminder.Data.MessageIntID = messageIntID
		}
		if messageStrID != "" {
			reminder.Data.MessageStrID = messageStrID
		}
		if err = tx.Set(ctx, reminder.Record); err != nil {
			err = fmt.Errorf("failed to save reminder to datastore: %w", err)
		}
		return err
	}
}

func (reminderDalGae ReminderDalGae) RescheduleReminder(ctx context.Context, reminderID string, remindInDuration time.Duration) (oldReminder, newReminder models4debtus.Reminder, err error) {
	return models4debtus.Reminder{}, models4debtus.Reminder{}, errors.New("not implemented - needs to be refactored")
	//var (
	//	newReminderKey    *datastore.Key
	//	newReminderEntity *models.ReminderDbo
	//)
	//err = facade.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {
	//	oldReminder, err = reminderDalGae.GetReminderByID(tctx, reminderID)
	//	if err != nil {
	//		return fmt.Errorf("failed to get oldReminder by id: %w", err)
	//	}
	//	if oldReminder.IsRescheduled {
	//		err = dtdal.ErrReminderAlreadyRescheduled
	//		return err
	//	}
	//	reminder := models.NewReminder(reminderID)
	//	if remindInDuration == time.Duration(0) {
	//		if _, err = tx.Set(tc, reminderKey, oldReminder.ReminderDbo); err != nil {
	//			return err
	//		}
	//	} else {
	//		nextReminderOn := time.Now().Add(remindInDuration)
	//		newReminderEntity = oldReminder.ScheduleNextReminder(reminderID, nextReminderOn)
	//		newReminderKey = NewReminderIncompleteKey(tc)
	//		keys, err := gaedb.PutMulti(tc, []*datastore.Key{reminderKey, newReminderKey}, []interface{}{oldReminder.ReminderDbo, newReminderEntity})
	//		if err != nil {
	//			err = fmt.Errorf("failed to reschedule oldReminder: %w", err)
	//		}
	//		newReminderKey = keys[1]
	//		if err = QueueSendReminder(tc, newReminderKey.IntID(), remindInDuration); err != nil { // TODO: Should be outside of DAL?
	//			return err
	//		}
	//	}
	//	return err
	//})
	//if err != nil {
	//	return
	//}
	//if newReminderKey != nil && newReminderEntity != nil {
	//	newReminder = models.Reminder{
	//		IntegerID:      db.NewIntID(newReminderKey.IntID()),
	//		ReminderDbo: newReminderEntity,
	//	}
	//}
	//return
}
