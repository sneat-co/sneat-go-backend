package models4sportus

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// PriceRange DTO
type PriceRange struct {
	PriceMin int
	PriceMax int
}

// LengthRange DTO
type LengthRange struct {
	LengthMin int
	LengthMax int
}

// WidthRange DTO
type WidthRange struct {
	WidthMin int
	WidthMax int
}

// WeightRange DTO
type WeightRange struct {
	WeightMin int
	WeightMax int
}

// SizeRange DTO
type SizeRange struct {
	SizeMin int
	SizeMax int
}

// YearRange DTO
type YearRange struct {
	YearMin int
	YearMax int
}

// RepairsRange DTO
type RepairsRange struct {
	RepairsMin int
	RepairsMax int
}

// PinholesRange DTO
type PinholesRange struct {
	PinholesMin int
	PinholesMax int
}

// Item DTO
type Item struct {
	UserID    string    `json:"userID" firestore:"userID"`
	Status    string    `json:"status" firestore:"status"`
	DtCreated time.Time `json:"dtCreated" firestore:"dtCreated"`
	DtUpdated time.Time `json:"dtUpdated" firestore:"dtUpdated"`
}

// Validate returns error if not valid
func (v Item) Validate() error {
	if strings.TrimSpace(v.UserID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("userID")
	}
	if v.DtCreated.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("dtCreated")
	}
	if v.DtUpdated.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("dtUpdated")
	}
	if v.DtUpdated.Before(v.DtCreated) {
		return validation.NewErrBadRecordFieldValue("dtUpdated", "is before dtCreated")
	}
	if strings.TrimSpace(v.Status) != v.Status {
		return validation.NewErrBadRecordFieldValue("status", fmt.Sprintf("unknown value: %v", v.Status))
	}
	switch v.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case "active", "deleted":
		// Known statuses
	default:
		return validation.NewErrBadRecordFieldValue("status", "unknown value: "+v.Status)
	}
	return nil
}

// Location DTO
type Location struct {
	Location []string
	GeoPoint
}
