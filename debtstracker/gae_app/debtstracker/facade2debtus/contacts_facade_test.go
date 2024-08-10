package facade2debtus

//"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
//"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
//"context"
//"testing"
//"time"
//"github.com/strongo/db/gaedb"

//func TestNewContactIncompleteKey(t *testing.T) {
//	testDatastoreIncompleteKey(t, NewContactIncompleteKey(context.Background()))
//}
//
//func TestNewContactKey(t *testing.T) {
//	const contactID = 112233
//	testDatastoreIntKey(t, contactID, NewContactKey(context.Background(), contactID))
//}
//
//func TestContactDalGae_CreateContactWithinTransaction(t *testing.T) {
//	gaedb.SetupNdsMock()
//
//	contactDal := NewContactDalGae()
//
//	c := context.Background()
//
//	user, err := User.GetUserByID(c, 2)
//	if err != nil {
//		t.Error("Faled to get user by ID: ", err)
//		return
//	}
//
//	contact, err := contactDal.CreateContactWithinTransaction(c, user, 1, 3, models.ContactDetails{
//		FirstName: "Jack",
//		LastName:  "Brown",
//	}, money.Balanced{
//		BalanceCount:   1,
//		BalanceJson:    `{"EUR":10.25}`,
//		LastTransferID: 15,
//		LastTransferAt: time.Now(),
//	})
//	if err != nil {
//		t.Error("Unexpected error:", err)
//	}
//	if contact.ID == 0 {
//		t.Error("contact.ID == 0")
//		return
//	}
//	if contact.FirstName != "Jack" {
//		t.Errorf("contact.FirstName != 'Jack': %v", contact.FirstName)
//		return
//	}
//	if contact.LastName != "Brown" {
//		t.Errorf("contact.FirstName != 'Jack': %v", contact.FirstName)
//		return
//	}
//	if contact.CounterpartyUserID != 1 {
//		t.Errorf("contact.CounterpartyUserID != 1: %d", contact.CounterpartyUserID)
//		return
//	}
//	if contact.CounterpartyCounterpartyID != 1 {
//		t.Errorf("contact.CounterpartyCounterpartyID != 3: %d", contact.CounterpartyCounterpartyID)
//		return
//	}
//	var isUserHasTheCounterparty bool
//	for _, c := range user.Contacts() {
//		if c.ID == contact.ID {
//			isUserHasTheCounterparty = true
//			break
//		}
//	}
//	if !isUserHasTheCounterparty {
//		t.Errorf("User.Contacts() does not have contact with ID==%d", contact.ID)
//	}
//}
//
//func TestCounterpartyDalGae_CreateContact(t *testing.T) {
//	gaedb.SetupNdsMock()
//
//	counterpartyDal := NewContactDalGae()
//	userDal := NewUserDalGae()
//
//	c := context.Background()
//
//	user, err := userDal.GetUserByID(c, 2)
//	if err != nil {
//		t.Error("Faled to get user by ID: ", err)
//		return
//	}
//	contact, _, err := counterpartyDal.CreateContact(c, user.ID, models.ContactDetails{
//		FirstName: "Jack",
//		LastName:  "Brown",
//	})
//	if err != nil {
//		t.Error("Unexpected error:", err)
//		return
//	}
//	if contact.ID == 0 {
//		t.Error("contact.ID == 0")
//		return
//	}
//}
