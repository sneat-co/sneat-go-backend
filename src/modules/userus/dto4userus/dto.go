package dto4userus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/security"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
)

var _ facade.Request = (*InitUserRecordRequest)(nil)

type InitSpaceInfo struct {
	Type  core4teamus.SpaceType `json:"type"`
	Title string                `json:"title"`
}

func (v InitSpaceInfo) Validate() error {
	if strings.TrimSpace(string(v.Type)) == "" {
		return validation.NewErrRequestIsMissingRequiredField("type")
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRequestIsMissingRequiredField("title")
	}
	return nil
}

// InitUserRecordRequest request
type InitUserRecordRequest struct {
	AuthProvider    string                         `json:"authProvider,omitempty"`
	Email           string                         `json:"email,omitempty"`
	EmailIsVerified bool                           `json:"emailIsVerified,omitempty"`
	IanaTimezone    string                         `json:"ianaTimezone,omitempty"`
	Names           *person.NameFields             `json:"names"`
	Space           *dto4teamus.CreateSpaceRequest `json:"space,omitempty"`
	//
	RemoteClient dbmodels.RemoteClientInfo `json:"remoteClient"`
}

// Validate validates request
func (v *InitUserRecordRequest) Validate() error {
	if v.AuthProvider != "" && !dbmodels.IsKnownAuthProviderID(v.AuthProvider) {
		return validation.NewErrBadRequestFieldValue("authProvider", "unknown value: "+v.AuthProvider)
	}
	if v.Names != nil {
		if err := v.Names.Validate(); err != nil {
			return fmt.Errorf("%w: %v", facade.ErrBadRequest, err)
		}
	}
	if v.Email != "" {
		if _, err := mail.ParseAddress(v.Email); err != nil {
			return validation.NewErrBadRequestFieldValue("email", err.Error())
		}
	}
	if v.Space != nil {
		if err := v.Space.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("space", err.Error())
		}
	}
	return nil
}

// CreateUserRequest DTO
type CreateUserRequest struct {
	Creator string `json:"creator"`
	Title   string `json:"title,omitempty"`
	Email   string `json:"email"`
}

// Validate validates request
func (v *CreateUserRequest) Validate() error {
	if err := validate.OptionalEmail(v.Email, "email"); err != nil {
		return err
	}
	if v.Creator == "" {
		return validation.NewErrRecordIsMissingRequiredField("creator")
	} else if !security.IsKnownHost(v.Creator) {
		return validation.NewErrBadRequestFieldValue("creator", "unknown creator: "+v.Creator)
	}
	return nil
}

// CreateUserRequestWithRemoteClientInfo a request DTO
type CreateUserRequestWithRemoteClientInfo struct {
	CreateUserRequest
	RemoteClient dbmodels.RemoteClientInfo
}
