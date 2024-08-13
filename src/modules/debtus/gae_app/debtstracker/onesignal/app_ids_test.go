package onesignal

import "testing"

//goland:noinspection ALL
func TestConstants(t *testing.T) {
	if APP_ID_LOCAL == "" {
		t.Error("APP_ID_LOCAL is not set")
	}
	if APP_ID_DEV1 == "" {
		t.Error("APP_ID_DEV1 is not set")
	}
	if APP_ID_PROD == "" {
		t.Error("APP_ID_PROD is not set")
	}
}
