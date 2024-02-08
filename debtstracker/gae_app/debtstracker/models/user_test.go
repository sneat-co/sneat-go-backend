package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/appuser"
	"testing"
	"time"
)

func TestAppUserEntity_Contacts(t *testing.T) {
	var userEntity DebutsAppUserDataOBSOLETE

	userEntity.ContactsJsonActive = `[{"ID":"1","Name":"Alex (Alex)"}]`

	contacts := userEntity.Contacts()

	contact := contacts[0]
	assert.Equal(t, contact.Name, "Alex")
	assert.Equal(t, contact.Status, "active")
}

func TestAppUserEntity_SetLastCurrency(t *testing.T) {
	userEntity := DebutsAppUserDataOBSOLETE{}
	userEntity.SetLastCurrency("EUR")
	if len(userEntity.LastCurrencies) != 1 {
		t.Errorf("Expected 1 value in LastCurrencies, got: %d", len(userEntity.LastCurrencies))
	}
	userEntity.SetLastCurrency("USD")
	if len(userEntity.LastCurrencies) != 2 {
		t.Errorf("Expected 2 values in LastCurrencies, got: %d", len(userEntity.LastCurrencies))
	}
	if userEntity.LastCurrencies[0] != "USD" {
		t.Errorf("First currency should be USD, got: %v", userEntity.LastCurrencies[0])
	}
	if userEntity.LastCurrencies[1] != "EUR" {
		t.Errorf("Second currency should be EUR, got: %v", userEntity.LastCurrencies[1])
	}

	userEntity.SetLastCurrency("EUR")
	if len(userEntity.LastCurrencies) != 2 {
		t.Errorf("Expected 2 values in LastCurrencies, got: %d", len(userEntity.LastCurrencies))
	}
	if userEntity.LastCurrencies[0] != "EUR" {
		t.Errorf("Second currency should be EUR, got: %v", userEntity.LastCurrencies[0])
	}
	if userEntity.LastCurrencies[1] != "USD" {
		t.Errorf("First currency should be USD, got: %v", userEntity.LastCurrencies[1])
	}
}

func TestLastLogin_SetLastLogin(t *testing.T) {
	user := NewUser(ClientInfo{})
	now := time.Now()
	user.Data.SetLastLogin(now)
	if user.Data.DtLastLogin != now {
		t.Errorf("user.DtLastLogin != now")
	}

	userGoogle := UserAccount{
		data: &appuser.AccountDataBase{},
	}
	userGoogle.Data().SetLastLogin(now)
	if userGoogle.data.DtLastLogin != now {
		t.Errorf("userGoogle.DtLastLogin != now")
	}

	type LastLoginSetter interface {
		SetLastLogin(v time.Time)
	}

	userGoogle = UserAccount{
		data: &appuser.AccountDataBase{},
	}
	var lastLoginSetter LastLoginSetter = userGoogle.data
	lastLoginSetter.SetLastLogin(now)
	if userGoogle.data.DtLastLogin != now {
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
	//	ContactsJsonActive: `[{"ID":"6296903092273152","Name":"Test1","Balance":{"EUR":58},"Transfers":{"Count":1,"Last":{"ID":"6156165603917824","At":"2017-11-04T23:05:30.847526702Z"},"OutstandingWithInterest":[{"TransferID":"6156165603917824","Starts":"2017-11-04T23:05:30.847526702Z","Currency":"EUR","Amount":14,"InterestType":"simple","InterestPeriod":3,"InterestPercent":3,"InterestMinimumPeriod":3}]}}]`,
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
