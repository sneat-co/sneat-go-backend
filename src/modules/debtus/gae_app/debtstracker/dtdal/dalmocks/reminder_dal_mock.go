package dalmocks

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/logus"
	"time"

	"context"
)

type ReminderDalMock struct {
}

func NewReminderDalMock() *ReminderDalMock {
	return &ReminderDalMock{}
}

func (mock *ReminderDalMock) DelayDiscardReminders(ctx context.Context, transferIDs []int, returntransferID int) error {
	logus.Warningf(ctx, "DelayDiscardReminders() is not implemented in mock")
	return nil
}

func (mock *ReminderDalMock) DelayCreateReminderForTransferCounterparty(ctx context.Context, transferID int) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) DelayCreateReminderForTransferUser(_ context.Context, transferID, userID string) error {
	return nil
}

func (mock *ReminderDalMock) GetReminderByID(_ context.Context, id int64) (models4debtus.Reminder, error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SaveReminder(_ context.Context, reminder models4debtus.Reminder) (err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) RescheduleReminder(_ context.Context, reminderID string, remindInDuration time.Duration) (oldReminder, newReminder models4debtus.Reminder, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SetReminderStatus(_ context.Context, reminderID, returnTransferID string, status string, when time.Time) (reminder models4debtus.Reminder, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) DelaySetReminderIsSent(_ context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SetReminderIsSent(_ context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) SetReminderIsSentInTransaction(_ context.Context, reminder models4debtus.Reminder, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) GetActiveReminderIDsByTransferID(_ context.Context, transferID int) ([]int64, error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *ReminderDalMock) GetSentReminderIDsByTransferID(_ context.Context, transferID int) ([]int64, error) {
	panic(NOT_IMPLEMENTED_YET)
}
