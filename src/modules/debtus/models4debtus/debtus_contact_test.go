package models4debtus

import (
	"testing"
	"time"
)

func TestContactEntity_SetTransfersInfo(t *testing.T) {
	contact := DebtusSpaceContactDbo{}
	if err := contact.SetTransfersInfo(UserContactTransfersInfo{
		Count: 1,
		Last: LastTransfer{
			ID: "2",
			At: time.Now(),
		},
	}); err != nil {
		t.Fatal(err)
	}
	if contact.Transfers == nil {
		t.Fatal("contact.TransfersJson is not set")
	}
}
