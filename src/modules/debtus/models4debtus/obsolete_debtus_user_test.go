package models4debtus

import (
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/strongo/strongoapp/appuser"
	"testing"
	"time"
)

func TestLastLogin_SetLastLogin(t *testing.T) {
	user := NewUser(ClientInfo{})
	now := time.Now()
	user.Data.SetLastLogin(now)
	if user.Data.DtLastLogin != now {
		t.Errorf("user.DtLastLogin != now")
	}

	userGoogle := models4auth.UserAccountEntry{
		Data: &appuser.AccountDataBase{},
	}
	userGoogle.Data.SetLastLogin(now)
	if userGoogle.Data.DtLastLogin != now {
		t.Errorf("userGoogle.DtLastLogin != now")
	}

	type LastLoginSetter interface {
		SetLastLogin(v time.Time)
	}

	userGoogle = models4auth.UserAccountEntry{
		Data: &appuser.AccountDataBase{},
	}
	var lastLoginSetter LastLoginSetter = userGoogle.Data
	lastLoginSetter.SetLastLogin(now)
	if userGoogle.Data.DtLastLogin != now {
		t.Errorf("lastLoginSetter.DtLastLogin != now")
	}
}

func TestAppUserEntity_BalanceWithInterest(t *testing.T) {
	t.Skip("TODO: Fix test")
	//user := DebutsAppUserDataOBSOLETE{
	//	//TransfersWithInterestCount: 1,
	//	//Balanced: money.Balanced{
	//	//	BalanceCount: 1,
	//	//	BalanceJson:  `{"EUR":58}`,
	//	//},
	//	ContactsJsonActive: `[{"ContactID":"6296903092273152","Name":"Test1","Balance":{"EUR":58},"Transfers":{"Count":1,"Last":{"ContactID":"6156165603917824","At":"2017-11-04T23:05:30.847526702Z"},"OutstandingWithInterest":[{"TransferID":"6156165603917824","Starts":"2017-11-04T23:05:30.847526702Z","Currency":"EUR","Amount":14,"InterestType":"simple","InterestPeriod":3,"InterestPercent":3,"InterestMinimumPeriod":3}]}}]`,
	//}
	//balanceWithInterest, err := user.BalanceWithInterest(context.Background(), time.Now())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if balanceWithInterest.IsZero() {
	//	t.Fatal("balanceWithInterest.IsZero()")
	//}
	//t.Log(balanceWithInterest)
}
