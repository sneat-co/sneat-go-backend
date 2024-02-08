package gaedal

import (
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"testing"
)

func TestNewUserEmailKey(t *testing.T) {
	const email = "test@example.come"
	testDatastoreStringKey(t, email, models.NewUserEmailKey(email))
}

func TestUserEmailGaeDal_GetUserEmailByID(t *testing.T) {
	//gaedb.Get = func(c context.Context, key *dal.Key, val interface{}) error {
	//	return nil
	//}
	t.Log("TestUserEmailGaeDal_GetUserEmailByID commented out")
	//userEmail, _ := NewUserEmailGaeDal().GetUserEmailByID(context.Background(), nil, " JackSmith@Example.com ")
	//
	//if userEmail.ID != "jacksmith@example.com" {
	//	t.Error("userEmail.ID expected to be lower case without spaces")
	//}
}
