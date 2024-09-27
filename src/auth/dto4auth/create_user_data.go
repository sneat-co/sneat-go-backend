package dto4auth

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/validation"
)

// DataToCreateUser is NOT a DTO object - do not user at transport layer!
type DataToCreateUser struct {
	AuthAccount     appuser.AccountKey
	Email           string
	EmailIsVerified bool
	IanaTimezone    string
	LanguageCode    string
	Names           person.NameFields
	PhotoURL        string
	RemoteClient    dbmodels.RemoteClientInfo
}

func (v DataToCreateUser) Validate() error {
	if err := v.AuthAccount.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("authAccount", "mismatch with account provider")
	}
	if err := v.RemoteClient.Validate(); err != nil {
		return fmt.Errorf("invalid remote client: %w", err)
	}
	return nil
}
