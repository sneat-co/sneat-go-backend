package facade4auth

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
	"testing"
)

func TestNewUserEmailKey(t *testing.T) {
	const email = "test@example.come"
	testStringKey(t, email, models4auth.NewUserEmailKey(email))
}

func TestUserEmailGaeDal_GetUserEmailByID(t *testing.T) {
	//gaedb.Get = func(ctx context.Context, key *dal.Key, val interface{}) error {
	//	return nil
	//}
	t.Log("TestUserEmailGaeDal_GetUserEmailByID commented out")
	//userEmail, _ := NewUserEmailGaeDal().GetUserEmailByID(context.Background(), nil, " JackSmith@Example.com ")
	//
	//if userEmail.ContactID != "jacksmith@example.com" {
	//	t.Error("userEmail.ContactID expected to be lower case without spaces")
	//}
}
