package dalmocks

import (
	"github.com/strongo/logus"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type ReminderDalMock struct {
}

func NewReminderDalMock() *ReminderDalMock {
	return &ReminderDalMock{}
}

func (mock *ReminderDalMock) DelayDiscardReminders(c context.Context, transferIDs []int, returntransferID int) error {
	logus.Warningf(c, "DelayDiscardReminders() is not implemented in mock")
	return nil
}

func (mock *ReminderDalMock) DelayCreateReminderForTransferCounterparty(c context.Context, transferID int) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) DelayCreateReminderForTransferUser(c context.Context, transferID, userID string) error {
	return nil
}

func (mock *ReminderDalMock) GetReminderByID(c context.Context, id int64) (models.Reminder, error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SaveReminder(c context.Context, reminder models.Reminder) (err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) RescheduleReminder(c context.Context, reminderID string, remindInDuration time.Duration) (oldReminder, newReminder models.Reminder, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SetReminderStatus(c context.Context, reminderID, returnTransferID string, status string, when time.Time) (reminder models.Reminder, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) DelaySetReminderIsSent(c context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SetReminderIsSent(c context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SetReminderIsSentInTransaction(c context.Context, reminder models.Reminder, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) GetActiveReminderIDsByTransferID(c context.Context, transferID int) ([]int64, error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) GetSentReminderIDsByTransferID(c context.Context, transferID int) ([]int64, error) {
	panic(NOT_IMPLEMENTED_YET)
}
