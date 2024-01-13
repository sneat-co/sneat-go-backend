package briefs4contactus

import (
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

type ContactGroupBrief struct {
	Title string `json:"title,omitempty" firestore:"title,omitempty"`
}

func (v *ContactGroupBrief) Validate() error {
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if err := validate.RecordTitle(v.Title, "title"); err != nil {
		return err
	}
	return nil
}
