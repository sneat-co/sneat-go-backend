package models4userus

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	user := NewUser("user1")
	if user.ID == "" {
		t.Error("user.ID is empty")
	}
	if user.ID != "user1" {
		t.Errorf("user.ID is %v, expected %v", user.ID, "user1")
	}
	if user.Key == nil {
		t.Error("user.Key is nil")
	}
	if user.Record.Key() == nil {
		t.Error("user.Record.Key() is nil")
	}
	if user.Key != user.Record.Key() {
		t.Error("user.Key != user.Record.Key()")
	}
	if user.Data == nil {
		t.Error("user.Data is nil")
	}
	user.Record.SetError(nil)
	recordData := user.Record.Data()
	if recordData == nil {
		t.Error("user.Record.Data() is nil")
	}
	if recordData != user.Data {
		t.Error("user.Data != user.Record.Data()")
	}
}
