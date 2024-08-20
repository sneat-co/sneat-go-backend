package facade4debtus

//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/dtdal"
//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
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
//	user, err := facade4userus.GetUserByID(c, 2)
//	if err != nil {
//		t.Error("Faled to get user by ContactID: ", err)
//		return
//	}
//
//	Contact, err := contactDal.CreateContactWithinTransaction(c, user, 1, 3, models.ContactDetails{
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
//	if Contact.ContactID == 0 {
//		t.Error("Contact.ContactID == 0")
//		return
//	}
//	if Contact.FirstName != "Jack" {
//		t.Errorf("Contact.FirstName != 'Jack': %v", Contact.FirstName)
//		return
//	}
//	if Contact.LastName != "Brown" {
//		t.Errorf("Contact.FirstName != 'Jack': %v", Contact.FirstName)
//		return
//	}
//	if Contact.CounterpartyUserID != 1 {
//		t.Errorf("Contact.CounterpartyUserID != 1: %d", Contact.CounterpartyUserID)
//		return
//	}
//	if Contact.CounterpartyContactID != 1 {
//		t.Errorf("Contact.CounterpartyContactID != 3: %d", Contact.CounterpartyContactID)
//		return
//	}
//	var isUserHasTheCounterparty bool
//	for _, c := range user.Contacts() {
//		if c.ContactID == Contact.ContactID {
//			isUserHasTheCounterparty = true
//			break
//		}
//	}
//	if !isUserHasTheCounterparty {
//		t.Errorf("User.Contacts() does not have Contact with ContactID==%d", Contact.ContactID)
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
//	user, err := userDal.GetUserByIdOBSOLETE(c, 2)
//	if err != nil {
//		t.Error("Faled to get user by ContactID: ", err)
//		return
//	}
//	Contact, _, err := counterpartyDal.CreateContact(c, user.ContactID, models.ContactDetails{
//		FirstName: "Jack",
//		LastName:  "Brown",
//	})
//	if err != nil {
//		t.Error("Unexpected error:", err)
//		return
//	}
//	if Contact.ContactID == 0 {
//		t.Error("Contact.ContactID == 0")
//		return
//	}
//}
