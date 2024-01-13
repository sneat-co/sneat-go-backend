package dto4contactus

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

type CreateLocationRequest struct {
	Title   string           `json:"title"`
	Address dbmodels.Address `json:"address"`
}

func (v CreateLocationRequest) Validate() error {
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRequestIsMissingRequiredField("title")
	}
	if err := v.Address.Validate(); err != nil {
		return err
	}
	return nil
}
