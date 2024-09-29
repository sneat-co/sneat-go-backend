package models4debtus

import (
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"testing"
)

func TestOnSaveSerializeJson(t *testing.T) {
	transferEntity := TransferData{
		from: &TransferCounterpartyInfo{
			UserID: "11",
		},
		to: &TransferCounterpartyInfo{
			ContactID: "22",
		},
	}

	if err := transferEntity.onSaveSerializeJson(); err != nil {
		t.Fatal("unexpected error", err)
	}

	if transferEntity.FromJson == "" {
		t.Error("transferEntity.FromJson is empty")
	}
	if transferEntity.ToJson == "" {
		t.Error("transferEntity.ToJson is empty")
	}
}

func TestTransferFromToUpdate(t *testing.T) {
	transferEntity := TransferData{
		CreatorUserID: "11",
		from: &TransferCounterpartyInfo{
			UserID: "11",
		},
		to: &TransferCounterpartyInfo{
			ContactID: "22",
		},
	}

	from := transferEntity.From()
	if v := from.UserID; v != "11" {
		t.Errorf("from.UserID != 11: %s", v)
		return
	}

	to := transferEntity.To()
	if v := to.ContactID; v != "22" {
		t.Errorf("to.ContactID != 22: %s", v)
		return
	}

	from.ContactID = "33"
	if v := transferEntity.From().ContactID; v != "33" {
		t.Errorf("transferEntity.From().ContactID != 33: %s", v)
		return
	}

	to.UserID = "44"
	if v := transferEntity.To().UserID; v != "44" {
		t.Errorf("transferEntity.To().UserID != 44: %s", v)
		return
	}

	transfer := NewTransfer("55", &transferEntity)

	from = transfer.Data.From()

	from.ContactID = "77"
	if v := transfer.Data.From().ContactID; v != "77" {
		t.Errorf("transferEntity.From().ContactID != 77: %s", v)
		return
	}

	creator := transfer.Data.Creator()
	creator.ContactID = "88"
	if v := transfer.Data.Creator().ContactID; v != "88" {
		t.Errorf("transfer.Creator().ContactID != 88: %s", v)
	}
	if v := transfer.Data.From().ContactID; v != "88" {
		t.Errorf("transfer.From().ContactID != 88: %s", v)
	}
}

func TestTransferCounterpartyInfo_Name(t *testing.T) {
	var contact TransferCounterpartyInfo
	if contact.ContactName = "Alex (Alex)"; contact.Name() != "Alex" {
		t.Errorf("Exected contact.ContactName() == 'Alex', got: %v", contact.Name())
	}
	if contact.ContactName = "Alex1 (Alex2)"; contact.Name() != "Alex1 (Alex2)" {
		t.Errorf("Exected contact.ContactName() == 'Alex1 (Alex2)', got: %v", contact.Name())
	}
	if contact.ContactName = "John Smith (John Smith)"; contact.Name() != "John Smith" {
		t.Errorf("Exected contact.ContactName() == 'John Smith', got: %v", contact.Name())
	}
}

func TestFixContactName(t *testing.T) {
	if isFixed, _ := models4auth.FixContactName(""); isFixed {
		t.Error("Should not fix empty string")
	}
	if _, s := models4auth.FixContactName(""); s != "" {
		t.Errorf("Expected empty string, got: %v", s)
	}
	if isFixed, _ := models4auth.FixContactName("Alex (Alex)"); !isFixed {
		t.Error("Exected 'Alex (Alex)' to be fixed")
	}
	if _, s := models4auth.FixContactName("Alex (Alex)"); s != "Alex" {
		t.Errorf("Exected contact.ContactName() == 'Alex', got: %v", s)
	}
	if isFixed, s := models4auth.FixContactName("Alex1 (Alex2)"); isFixed || s != "Alex1 (Alex2)" {
		t.Errorf("Exected isFiexed=false, s='Alex1 (Alex2)'. Got: isFiexed=%v, s=%v", isFixed, s)
	}
}
