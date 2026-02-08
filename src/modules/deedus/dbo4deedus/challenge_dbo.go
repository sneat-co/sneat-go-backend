package dbo4deedus

import (
	"strings"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

type ChallengeDbo struct {
	Title  string          `json:"title" firestore:"title"`
	Status dbmodels.Status `json:"status"  firestore:"status"`
	Stars  int             `json:"stars"  firestore:"stars"`
}

func (v ChallengeDbo) Validate() error {
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}
