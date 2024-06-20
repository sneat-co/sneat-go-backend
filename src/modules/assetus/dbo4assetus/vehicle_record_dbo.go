package dbo4assetus

import (
	"github.com/crediterra/money"
	"github.com/strongo/decimal"
	"github.com/strongo/strongoapp/with"
)

type VehicleFuelRecord struct {
	Volume decimal.Decimal64p2 `json:"volume,omitempty" firestore:"unit,omitempty"`
	Unit   string              `json:"unit,omitempty" firestore:"unit,omitempty"`
	Amount *money.Amount       `json:"amount,omitempty" firestore:"amount,omitempty"`
}

func (v *VehicleFuelRecord) Validate() error {
	// TODO: implement validation
	// if Volume !=0 Unit should be set
	// Unit should be one of known units (l, gal, etc)
	// Amount should be valid (call v.amount.Validate())
	return nil
}

type VehicleMileage struct {
	Value int    `json:"value" firestore:"value"`
	Unit  string `json:"unit" firestore:"unit"`
}

func (v VehicleMileage) Validate() error {
	// TODO: implement validation
	return nil
}

type VehicleRecordDbo struct {
	with.CreatedFields                    // Mandatory field
	Fuel               *VehicleFuelRecord `json:"fuel,omitempty" firestore:"fuel,omitempty"`
	Mileage            *VehicleMileage    `json:"mileage,omitempty" firestore:"mileage,omitempty"`
}

func (v VehicleRecordDbo) Validate() error {
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.Fuel.Validate(); err != nil {
		return err
	}
	if err := v.Mileage.Validate(); err != nil {
		return err
	}
	return nil
}
