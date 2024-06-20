package extras4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
)

// MileageExtra is an extension of asset data for vehicles
type WithMileageExtra struct {
	Mileages []dal4assetus.Mileage `json:"mileage,omitempty" firestore:"mileage,omitempty"`
}

func (v *WithMileageExtra) Validate() error {
	return nil
}
