package dto4contactus

import (
	"fmt"
	"github.com/strongo/strongoapp/person"
	"strconv"
	"strings"
)

type PhoneContact struct {
	// Part of ContactDetails => Part of User|DebtusSpaceContactEntry
	// DebtusSpaceContactEntry details
	PhoneNumber          int64 `firestore:"phoneNumber,omitempty"`
	PhoneNumberConfirmed bool  `firestore:"phoneNumberConfirmed,omitempty"`
	//+9223372036854775807
	//+353857403948
	//+79169743259
}

func (p PhoneContact) PhoneNumberAsString() string {
	return "+" + strconv.FormatInt(p.PhoneNumber, 10)
}

type EmailContact struct {
	EmailAddress         string `firestore:"emailAddress,omitempty"`
	EmailAddressOriginal string `firestore:"emailAddressOriginal,omitempty"`
	EmailConfirmed       bool   `firestore:"emailConfirmed,omitempty"`
}

func (ec *EmailContact) SetEmail(email string, confirmed bool) EmailContact {
	ec.EmailAddress = strings.ToLower(email)
	if ec.EmailAddress != email {
		ec.EmailAddressOriginal = email
	} else {
		ec.EmailAddressOriginal = ""
	}
	ec.EmailConfirmed = confirmed
	return *ec
}

type ContactDetails struct {
	// Helper struct, not stored as independent entity
	PhoneContact
	EmailContact
	person.NameFields
	//FirstName      string `firestore:",omitempty"`
	//LastName       string `firestore:",omitempty"`
	//ScreenName     string `firestore:",omitempty"`
	//Nickname       string `firestore:",omitempty"`
	//Username       string `firestore:",omitempty"` //TODO: Should it be "Name"?
	TelegramUserID int64 // When user ads Telegram contact we store Telegram user_id so we can link users later.
}

func (contact *ContactDetails) FullName() string {
	addUserNameIfNotSame := func(name string) string {
		if contact.UserName == "" || strings.EqualFold(contact.UserName, name) {
			return name
		} else {
			return fmt.Sprintf("%v (@%v)", name, contact.UserName)
		}
	}
	if contact.LastName != "" && contact.FirstName != "" {
		return addUserNameIfNotSame(contact.FirstName + " " + contact.LastName)
	} else if contact.FirstName != "" {
		return addUserNameIfNotSame(contact.FirstName)
	} else if contact.LastName != "" {
		return addUserNameIfNotSame(contact.LastName)
		//} else if contact.ScreenName != "" {
		//	return addUserNameIfNotSame(contact.ScreenName)
	} else if contact.UserName != "" {
		return contact.UserName
	} else if contact.NickName != "" {
		return addUserNameIfNotSame(contact.NickName)
	} else {
		return NoName
	}
}

const NoName = ">NO_NAME<"
