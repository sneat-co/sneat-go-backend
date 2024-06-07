package gaedal

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"reflect"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

func NewReminderIncompleteKey(_ context.Context) *dal.Key {
	return dal.NewIncompleteKey(models.ReminderKind, reflect.Int, nil)
}

func NewReminderKey(reminderID string) *dal.Key {
	if reminderID == "" {
		panic("reminderID == 0")
	}
	return dal.NewKeyWithID(models.ReminderKind, reminderID)
}

type ReminderDalGae struct {
}

func NewReminderDalGae() ReminderDalGae {
	return ReminderDalGae{}
}

var _ dtdal.ReminderDal = (*ReminderDalGae)(nil)

func (reminderDalGae ReminderDalGae) GetReminderByID(c context.Context, tx dal.ReadSession, id string) (reminder models.Reminder, err error) {
	reminder = models.NewReminder(id, nil)
	return reminder, tx.Get(c, reminder.Record)
}

func (reminderDalGae ReminderDalGae) SaveReminder(c context.Context, tx dal.ReadwriteTransaction, reminder models.Reminder) (err error) {
	return tx.Set(c, reminder.Record)
}

func (reminderDalGae ReminderDalGae) GetSentReminderIDsByTransferID(c context.Context, tx dal.ReadSession, transferID int) ([]int, error) {
	q := dal.From(models.ReminderKind).Where(
		dal.WhereField("TransferID", dal.Equal, transferID),
		dal.WhereField("Status", dal.Equal, models.ReminderStatusSent),
	).SelectKeysOnly(reflect.Int)

	records, err := tx.QueryAllRecords(c, q)
	if err != nil {
		return nil, err
	}
	reminderIDs := make([]int, len(records))
	for i, record := range records {
		reminderIDs[i] = record.Key().ID.(int)
	}
	return reminderIDs, nil
}

func (reminderDalGae ReminderDalGae) GetActiveReminderIDsByTransferID(c context.Context, tx dal.ReadSession, transferID int) ([]int, error) {
	q := dal.From(models.ReminderKind).Where(
		dal.WhereField("TransferID", dal.Equal, transferID),
		dal.WhereField("DtNext", dal.GreaterThen, time.Time{}),
	).SelectKeysOnly(reflect.Int)
	records, err := tx.QueryAllRecords(c, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reminders by transfer id=%v: %w", transferID, err)
	}
	reminderIDs := make([]int, len(records))
	for i, record := range records {
		reminderIDs[i] = record.Key().ID.(int)
	}
	return reminderIDs, nil
}

func (reminderDalGae ReminderDalGae) SetReminderIsSent(c context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error) {
	//gaehost.GaeLogger.Debugf(c, "setReminderIsSent(reminderID=%v, sentAt=%v, messageIntID=%v, messageStrID=%v)", reminderID, sentAt, messageIntID, messageStrID)
	if err := _validateSetReminderIsSentMessageIDs(messageIntID, messageStrID, sentAt); err != nil {
		return err
	}
	reminder := models.NewReminder(reminderID, nil)
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return reminderDalGae.SetReminderIsSentInTransaction(c, tx, reminder, sentAt, messageIntID, messageStrID, locale, errDetails)
	})
}

func (reminderDalGae ReminderDalGae) SetReminderIsSentInTransaction(c context.Context, tx dal.ReadwriteTransaction, reminder models.Reminder, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error) {
	if reminder.Data == nil {
		reminder, err = reminderDalGae.GetReminderByID(c, tx, reminder.ID)
		if err != nil {
			if dal.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("failed to get reminder by ID: %w", err)
		}
	}
	if reminder.Data.Status != models.ReminderStatusSending {
		log.Errorf(c, "reminder.Status:%v != models.ReminderStatusSending:%v", reminder.Data.Status, models.ReminderStatusSending)
		return nil
	} else {
		reminder.Data.Status = models.ReminderStatusSent
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
		if err = tx.Set(c, reminder.Record); err != nil {
			err = fmt.Errorf("failed to save reminder to datastore: %w", err)
		}
		return err
	}
}

func (reminderDalGae ReminderDalGae) RescheduleReminder(_ context.Context, reminderID string, remindInDuration time.Duration) (oldReminder, newReminder models.Reminder, err error) {
	return models.Reminder{}, models.Reminder{}, errors.New("not implemented - needs to be refactored")
	//var (
	//	newReminderKey    *datastore.Key
	//	newReminderEntity *models.ReminderDbo
	//)
	//var db dal.DB
	//if db, err = facade.GetDatabase(c); err != nil {
	//	return
	//}
	//err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
	//	oldReminder, err = reminderDalGae.GetReminderByID(c, reminderID)
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
