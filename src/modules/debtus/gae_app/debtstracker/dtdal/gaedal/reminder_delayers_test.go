package gaedal

import (
	"fmt"
	"github.com/strongo/i18n"
	"strings"
	"testing"
	"time"

	"context"
)

func Test__validateSetReminderIsSentMessageIDs(t *testing.T) {
	var err error
	now := time.Now()
	if err = _validateSetReminderIsSentMessageIDs(0, "", now); err == nil {
		t.Error("Should fail: _validateSetReminderIsSentMessageIDs(0, '')")
	}
	if err = _validateSetReminderIsSentMessageIDs(1, "not empty", now); err == nil {
		t.Error("Should fail: _validateSetReminderIsSentMessageIDs(1, 'not empty')")
	}
	if err = _validateSetReminderIsSentMessageIDs(1, "", time.Time{}); err == nil {
		t.Error("Should fail as sentAt is zero")
		if !strings.Contains(err.Error(), "sentAt.IsZero()") {
			t.Error("Error message does not contain 'sentAt.IsZero()'")
		}
	}
}

func TestDelaySetReminderIsSent(t *testing.T) {
	var err error

	reminderDal := NewReminderDalGae()

	if err = reminderDal.DelaySetReminderIsSent(context.TODO(), "", time.Now(), 1, "", i18n.LocaleCodeEnUS, ""); err == nil {
		t.Error("Should fail as reminder is 0")
	}
	if err = reminderDal.DelaySetReminderIsSent(context.TODO(), "1", time.Now(), 0, "", i18n.LocaleCodeEnUS, ""); err == nil {
		t.Error("Should fail as no message id supplied")
	}
	if err = reminderDal.DelaySetReminderIsSent(context.TODO(), "1", time.Now(), 1, "not empty", i18n.LocaleCodeEnUS, ""); err == nil {
		t.Error("Should fail as both int and string message ids supplied")
	}
	if err = reminderDal.DelaySetReminderIsSent(context.TODO(), "1", time.Time{}, 1, "not empty", i18n.LocaleCodeEnUS, ""); err == nil {
		t.Error("Should fail as both int and string message ids supplied")
	}
	if err = reminderDal.DelaySetReminderIsSent(context.TODO(), "1", time.Time{}, 1, "", i18n.LocaleCodeEnUS, ""); err == nil {
		t.Error("Should fail as both sentAt is zero")
	}

	//countOfCallsToDelay := 0
	//apphostgae.CallDelayFunc = func(ctx context.Context, queueName, subPath string, f *delay.Function, args ...interface{}) error {
	//	countOfCallsToDelay += 1
	//	return nil
	//}
	if err = reminderDal.DelaySetReminderIsSent(context.TODO(), "1", time.Now(), 1, "", i18n.LocaleCodeEnUS, ""); err != nil {
		t.Error(fmt.Errorf("should NOT fail: %w", err).Error())
	}
	//if countOfCallsToDelay != 1 {
	//	t.Errorf("Expeted to get 1 call to delay, got: %v", countOfCallsToDelay)
	//}
}
