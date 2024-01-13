package models4companius

import (
	"fmt"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/validation"
)

// CompanyBase holds info shared by CompanyDto and CompanyBrief structs.
type CompanyBase struct {
	Kind  string `json:"kind" firestore:"kind"` // either "private" or "work"
	Type  string `json:"type" firestore:"type"` // e.g. personal, family, work
	Title string `json:"title,omitempty" firestore:"title,omitempty"`
}

// Validate returns error if not valid
func (v CompanyBase) Validate() error {
	if v.Kind == "" {
		return validation.NewErrRecordIsMissingRequiredField("kind")
	}
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if v.Kind == "private" {
		if v.Title != "" {
			return validation.NewErrBadRecordFieldValue("title", "should be empty for private kind of facade4meetingus")
		}
	} else if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// CompanyDto holds info about family or work or personal stuff.
type CompanyDto struct {
	CompanyBase
	// NumberOf keeps numbers of: members, documents, etc. It is used to show #s in company's menu.
	NumberOf map[string]int `json:"numberOf" firestore:"numberOf"`
}

// Validate returns error if not valid
func (v CompanyDto) Validate() error {
	if err := v.CompanyBase.Validate(); err != nil {
		return err
	}
	for k, v := range v.NumberOf {
		if v < 0 {
			return validation.NewErrBadRecordFieldValue("numberOf."+k, fmt.Sprintf("value expected to be positive, got: %v", v))
		}
	}
	return nil
}

type Company = record.DataWithID[string, *CompanyDto]
