package facade4auth

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/validation"
)

// DataToCreateUser is NOT a DTO object - do not user at transport layer!
type DataToCreateUser struct {
	Account         appuser.AccountKey
	AuthProvider    string // TODO: Seems to be a duplicate of Account.Provider - document why needed or remove
	Email           string
	EmailIsVerified bool
	IanaTimezone    string
	LanguageCode    string
	Names           person.NameFields
	PhotoURL        string
	RemoteClient    dbmodels.RemoteClientInfo
}

func (v DataToCreateUser) Validate() error {
	if v.AuthProvider == "" {
		return validation.NewErrRequestIsMissingRequiredField("authProvider")
	}
	if v.Account.Provider != "" && v.Account.Provider != v.AuthProvider {
		return validation.NewErrBadRequestFieldValue("authProvider", "mismatch with account provider")
	}
	if err := v.RemoteClient.Validate(); err != nil {
		return fmt.Errorf("invalid remote client: %w", err)
	}
	return nil
}
