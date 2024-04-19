package models4calendarium

import (
	"github.com/crediterra/money"
)

type HappeningAsset struct { // TODO: Should it be removed in favor of `related`?
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
