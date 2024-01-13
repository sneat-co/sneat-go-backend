package models4brands

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"net/url"
	"strings"
)

// Maker defines a DB object model
type Maker struct {
	Title      string   `json:"title" firestore:"title"`
	AssetTypes []string `json:"assetTypes" firestore:"assetTypes"` // e.g. "vehicle:car", "vehicle:bicycle", "real-estate:house", "real-estate:apartment"
	WebsiteURL string   `json:"websiteURL,omitempty" firestore:"websiteURL,omitempty"`
	Models     []string `json:"models,omitempty" firestore:"models,omitempty"`
}

// Validate returns error if not valid
func (v *Maker) Validate() error {
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if err := validate.RecordTitle(v.Title, "title"); err != nil {
		return err
	}
	if strings.TrimSpace(v.WebsiteURL) != "" {
		if _, err := url.Parse(v.WebsiteURL); err != nil {
			return validation.NewErrBadRecordFieldValue("websiteURL", err.Error())
		}
	}
	if len(v.AssetTypes) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("assetTypes")
	}
	for i, model := range v.Models {
		if strings.TrimSpace(model) == "" {
			return validation.NewErrRequestIsMissingRequiredField(fmt.Sprintf("models[%v]", i))
		}
	}
	return nil
}
