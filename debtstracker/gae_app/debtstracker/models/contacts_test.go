package models

import "testing"

func TestContactDetailsFullname(t *testing.T) {
	lastName := "Smith"
	contactDetails := &ContactDetails{
		LastName: lastName,
	}
	if fullName := contactDetails.FullName(); fullName != lastName {
		t.Errorf("Expected %v, got %v", lastName, fullName)
	}
}
