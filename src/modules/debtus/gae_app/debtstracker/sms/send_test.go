package sms

import "testing"

func TestTwilioExceptionToMessageWithoutArguments(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should fail")
		}
	}()
	TwilioExceptionToMessage(nil, nil, nil)
}
