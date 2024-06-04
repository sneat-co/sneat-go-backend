package facade4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
)

func updateTeamDtoWithNumberOfContact(numberOfContacts int) (update dal.Update) {
	var value interface{}
	if numberOfContacts == 0 {
		value = dal.DeleteField
	} else {
		value = numberOfContacts
	}
	return dal.Update{
		Field: dbo4teamus.NumberOfUpdateField(const4contactus.ContactsField),
		Value: value,
	}
}
