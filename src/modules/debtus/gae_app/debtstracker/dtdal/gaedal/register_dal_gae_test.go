package gaedal

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"testing"
)

func TestRegisterDal(t *testing.T) {
	// Pre-clean
	dtdal.Admin = nil
	dtdal.Contact = nil
	//dtdal.Group = nil
	dtdal.Twilio = nil
	dtdal.HttpClient = nil
	dtdal.Invite = nil
	dtdal.LoginCode = nil
	dtdal.LoginPin = nil
	//dtdal.Bill = nil
	dtdal.Receipt = nil
	dtdal.Reminder = nil
	facade4auth.TgUser = nil
	dtdal.Transfer = nil
	facade4auth.User = nil
	facade4auth.UserGooglePlus = nil
	facade4auth.UserFacebook = nil

	// Execute
	RegisterDal()
	// Assert
	if dtdal.Admin == nil {
		t.Error("dtdal.Admin == nil")
	}
	//if dtdal.Bill == nil {
	//	t.Error("dtdal.Bill == nil")
	//}
	if dtdal.Contact == nil {
		t.Error("dtdal.DebtusSpaceContactEntry == nil")
	}
	if dtdal.Receipt == nil {
		t.Error("dtdal.Receipt == nil")
	}
	if dtdal.Reminder == nil {
		t.Error("dtdal.Reminder == nil")
	}
	//if facade4auth.UserBrowser == nil {
	//	t.Error("dtdal.UserBrowser == nil")
	//}
	//if dtdal.Bill == nil {
	//	t.Error("dtdal.Bill == nil")
	//}
	if dtdal.HttpClient == nil {
		t.Error("dtdal.HttpClient == nil")
	}
	if dtdal.Invite == nil {
		t.Error("dtdal.Invite == nil")
	}
	//if dtdal.Group == nil {
	//	t.Error("dtdal.Invite == nil")
	//}
	if facade4auth.TgUser == nil {
		t.Error("dtdal.TgUser == nil")
	}
	if dtdal.Transfer == nil {
		t.Error("dtdal.Transfer == nil")
	}
	if dtdal.Twilio == nil {
		t.Error("dtdal.Twilio == nil")
	}
	if facade4auth.User == nil {
		t.Error("dtdal.User == nil")
	}
	//if facade4auth.UserBrowser == nil {
	//	t.Error("dtdal.UserBrowser == nil")
	//}
	//if facade4auth.UserGaClient == nil {
	//	t.Error("dtdal.UserGaClient == nil")
	//}
	if facade4auth.UserGooglePlus == nil {
		t.Error("dtdal.UserGooglePlus == nil")
	}
	if facade4auth.UserFacebook == nil {
		t.Error("dtdal.UserFacebook == nil")
	}
	//if facade4auth.UserOneSignal == nil {
	//	t.Error("dtdal.UserOneSignal == nil")
	//}
}
