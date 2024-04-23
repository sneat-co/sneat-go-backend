package models4assetus

import "github.com/strongo/validation"

type AssetDwellingExtra struct {
	BedRooms  int `json:"bedRooms,omitempty" firestore:"bedRooms,omitempty"`
	BathRooms int `json:"bathRooms,omitempty" firestore:"bathRooms,omitempty"`
}

func (v AssetDwellingExtra) Validate() error {
	if v.BedRooms < 0 {
		return validation.NewErrBadRecordFieldValue("bedRooms", "negative value")
	}
	if v.BathRooms < 0 {
		return validation.NewErrBadRecordFieldValue("bathRooms", "negative value")
	}
	return nil
}
