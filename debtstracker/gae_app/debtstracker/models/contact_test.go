package models

import (
	"testing"
	"time"
)

func TestContactEntity_SetTransfersInfo(t *testing.T) {
	contact := DebtusContactData{}
	if err := contact.SetTransfersInfo(UserContactTransfersInfo{
		Count: 1,
		Last: LastTransfer{
			ID: "2",
			At: time.Now(),
		},
	}); err != nil {
		t.Fatal(err)
	}
	if contact.TransfersJson == "" {
		t.Fatal("contact.TransfersJson is not set")
	}
}
