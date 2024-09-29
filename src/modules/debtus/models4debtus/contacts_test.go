package models4debtus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/strongo/strongoapp/person"
	"testing"
)

func TestContactDetailsFullname(t *testing.T) {
	lastName := "Smith"
	contactDetails := &dto4contactus.ContactDetails{
		NameFields: person.NameFields{
			LastName: lastName,
		},
	}
	if fullName := contactDetails.FullName(); fullName != lastName {
		t.Errorf("Expected %v, got %v", lastName, fullName)
	}
}
