package dal4assetus

type Mileage struct {
	ID             string  `json:"id,omitempty" firestore:"id,omitempty"`
	FuelVolume     float32 `json:"fuelVolume,omitempty" firestore:"fuelVolume,omitempty"`
	FuelVolumeUnit string  `json:"fuelVolumeUnit,omitempty" firestore:"fuelVolumeUnit,omitempty"`
	FuelCost       float32 `json:"fuelCost,omitempty" firestore:"fuelCost,omitempty"`
	Currency       string  `json:"currency,omitempty" firestore:"currency,omitempty"`
	Mileage        float32 `json:"mileage,omitempty" firestore:"mileage,omitempty"`
	MileageUnit    string  `json:"mileageUnit,omitempty" firestore:"mileageUnit,omitempty"`
}
