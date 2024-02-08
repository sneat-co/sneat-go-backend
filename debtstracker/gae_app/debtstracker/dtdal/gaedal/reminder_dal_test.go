package gaedal

import (
	"testing"

	"context"
)

func TestNewReminderKey(t *testing.T) {
	const reminderID = "135"
	testStrKey(t, reminderID, NewReminderKey(reminderID))
}

func TestNewReminderIncompleteKey(t *testing.T) {
	testIncompleteKey(t, NewReminderIncompleteKey(context.Background()))
}
