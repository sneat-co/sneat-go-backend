package dto4contactus

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

type CreateCompanyRequest struct {
	Title     string            `json:"title"`
	CountryID string            `json:"countryID"`
	VATNumber string            `json:"vatNumber,omitempty"`
	Address   *dbmodels.Address `json:"address"`
}

func (v CreateCompanyRequest) Validate() error {
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRequestIsMissingRequiredField("title")
	}
	if v.Address != nil {
		if err := v.Address.Validate(); err != nil {
			return err
		}
		if v.CountryID != "" && v.Address.CountryID != v.CountryID {
			return validation.NewErrBadRequestFieldValue("address", "countryID should be equal to address.countryID")
		}

	} else if v.CountryID == "" {
		return validation.NewErrRequestIsMissingRequiredField("countryID|address")
	}
	if strings.TrimSpace(v.VATNumber) != v.VATNumber {
		return validation.NewErrBadRecordFieldValue("vatNumber", "should be trimmed")
	}
	return nil
}
