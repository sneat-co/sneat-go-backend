package gaedal

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
)

func TestRegisterDal(t *testing.T) {
	// Pre-clean
	dtdal.Admin = nil
	dtdal.Contact = nil
	dtdal.Group = nil
	dtdal.Twilio = nil
	dtdal.HttpClient = nil
	dtdal.Invite = nil
	dtdal.LoginCode = nil
	dtdal.LoginPin = nil
	dtdal.Bill = nil
	dtdal.Receipt = nil
	dtdal.Reminder = nil
	dtdal.TgUser = nil
	dtdal.Transfer = nil
	dtdal.User = nil
	dtdal.UserBrowser = nil
	dtdal.UserGaClient = nil
	dtdal.UserGooglePlus = nil
	dtdal.UserFacebook = nil
	dtdal.UserOneSignal = nil

	// Execute
	RegisterDal()
	// Assert
	if dtdal.Admin == nil {
		t.Error("dtdal.Admin == nil")
	}
	if dtdal.Bill == nil {
		t.Error("dtdal.Bill == nil")
	}
	if dtdal.Contact == nil {
		t.Error("dtdal.ContactEntry == nil")
	}
	if dtdal.Receipt == nil {
		t.Error("dtdal.Receipt == nil")
	}
	if dtdal.Reminder == nil {
		t.Error("dtdal.Reminder == nil")
	}
	if dtdal.UserBrowser == nil {
		t.Error("dtdal.UserBrowser == nil")
	}
	if dtdal.Bill == nil {
		t.Error("dtdal.Bill == nil")
	}
	if dtdal.HttpClient == nil {
		t.Error("dtdal.HttpClient == nil")
	}
	if dtdal.Invite == nil {
		t.Error("dtdal.Invite == nil")
	}
	if dtdal.Group == nil {
		t.Error("dtdal.Invite == nil")
	}
	if dtdal.TgUser == nil {
		t.Error("dtdal.TgUser == nil")
	}
	if dtdal.Transfer == nil {
		t.Error("dtdal.Transfer == nil")
	}
	if dtdal.Twilio == nil {
		t.Error("dtdal.Twilio == nil")
	}
	if dtdal.User == nil {
		t.Error("dtdal.User == nil")
	}
	if dtdal.UserBrowser == nil {
		t.Error("dtdal.UserBrowser == nil")
	}
	if dtdal.UserGaClient == nil {
		t.Error("dtdal.UserGaClient == nil")
	}
	if dtdal.UserGooglePlus == nil {
		t.Error("dtdal.UserGooglePlus == nil")
	}
	if dtdal.UserFacebook == nil {
		t.Error("dtdal.UserFacebook == nil")
	}
	if dtdal.UserOneSignal == nil {
		t.Error("dtdal.UserOneSignal == nil")
	}
}
