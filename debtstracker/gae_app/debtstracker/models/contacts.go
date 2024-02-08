package models

import (
	"fmt"
	"strconv"
	"strings"
)

type PhoneContact struct {
	// Part of ContactDetails => Part of User|Contact
	// Contact details
	PhoneNumber            int64 `datastore:",omitempty"`
	PhoneNumberConfirmed   bool
	PhoneNumberIsConfirmed bool `datastore:",noindex"` // Deprecated
	//+9223372036854775807
	//+353857403948
	//+79169743259
}

func (p PhoneContact) PhoneNumberAsString() string {
	return "+" + strconv.FormatInt(p.PhoneNumber, 10)
}

type EmailContact struct {
	EmailAddress         string `datastore:",omitempty"`
	EmailAddressOriginal string `datastore:",noindex,omitempty"`
	EmailConfirmed       bool   `datastore:",noindex,omitempty"`
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
	FirstName      string `datastore:",noindex,omitempty"`
	LastName       string `datastore:",noindex,omitempty"`
	ScreenName     string `datastore:",noindex,omitempty"`
	Nickname       string `datastore:",noindex,omitempty"`
	Username       string `datastore:",noindex,omitempty"` //TODO: Should it be "Name"?
	TelegramUserID int64  // When user ads Telegram contact we store Telegram user_id so we can link users later.
}

func (contact *ContactDetails) FullName() string {
	addUserNameIfNotSame := func(name string) string {
		if contact.Username == "" || strings.EqualFold(contact.Username, name) {
			return name
		} else {
			return fmt.Sprintf("%v (@%v)", name, contact.Username)
		}
	}
	if contact.LastName != "" && contact.FirstName != "" {
		return addUserNameIfNotSame(contact.FirstName + " " + contact.LastName)
	} else if contact.FirstName != "" {
		return addUserNameIfNotSame(contact.FirstName)
	} else if contact.LastName != "" {
		return addUserNameIfNotSame(contact.LastName)
	} else if contact.ScreenName != "" {
		return addUserNameIfNotSame(contact.ScreenName)
	} else if contact.Username != "" {
		return contact.Username
	} else if contact.Nickname != "" {
		return addUserNameIfNotSame(contact.Nickname)
	} else {
		return NoName
	}
}

const NoName = ">NO_NAME<"
