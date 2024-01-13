package models4calendarium

import (
	"fmt"
	"github.com/crediterra/money"
	"github.com/strongo/validation"
	"strings"
)

type HappeningParticipant struct {
	Roles []string `json:"roles,omitempty" firestore:"roles,omitempty"`
}

func (v HappeningParticipant) Validate() error {
	for i, role := range v.Roles {
		if strings.TrimSpace(role) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("roles[%d]", i), "role is empty")
		}
	}
	return nil
}

type HappeningAsset struct {
	RentAmount *money.Amount `json:"rentAmount,omitempty" firestore:"rentAmount,omitempty"`
}

func (v HappeningAsset) Validate() error {
	if v.RentAmount != nil {
		if err := v.RentAmount.Validate(); err != nil {
			return err
		}
	}
	return nil
}
