package models4debtus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/strongo/strongoapp/appuser"
	"testing"
	"time"
)

func TestLastLogin_SetLastLogin(t *testing.T) {
	user := NewUser(common4all.ClientInfo{})
	now := time.Now()
	user.Data.SetLastLoginAt(now)
	if user.Data.LastLoginAt != now {
		t.Errorf("user.DtLastLogin != now")
	}

	userGoogle := models4auth.UserAccountEntry{
		Data: &appuser.AccountDataBase{},
	}
	userGoogle.Data.SetLastLoginAt(now)
	if userGoogle.Data.LastLoginAt != now {
		t.Errorf("userGoogle.DtLastLogin != now")
	}

	type LastLoginSetter interface {
		SetLastLoginAt(v time.Time) dal.Update
	}

	userGoogle = models4auth.UserAccountEntry{
		Data: &appuser.AccountDataBase{},
	}
	var lastLoginSetter LastLoginSetter = userGoogle.Data
	lastLoginSetter.SetLastLoginAt(now)
	if userGoogle.Data.LastLoginAt != now {
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
