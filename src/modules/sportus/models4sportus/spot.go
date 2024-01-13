package models4sportus

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// SpotBrief DTO
type SpotBrief struct {
	Title    string `json:"title" firestore:"title"`
	GeoPoint GeoPoint
}

// Validate returns error if not valid
func (v SpotBrief) Validate() error {
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// Spot DTO
type Spot struct {
	SpotBrief
	dbmodels.WithUserIDs
	CheckedIn int `json:"checkedIn" firestore:"checkedIn"`
}

// Validate returns error if not valid
func (v Spot) Validate() error {
	if err := v.SpotBrief.Validate(); err != nil {
		return err
	}
	return nil
}

// SpotVisit DTO
type SpotVisit struct {
	UserID       string    `json:"userID" firestore:"userID"`
	SpotID       string    `json:"spotID" firestore:"spotID"`
	CheckedInAt  time.Time `json:"checkedInAt" firestore:"checkedInAt"`
	CheckedOutAt time.Time `json:"checkedOutAt" firestore:"checkedOutAt"`
	Comment      string
}

// Validate returns error if not valid
func (v SpotVisit) Validate() error {
	if strings.TrimSpace(v.UserID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("userID")
	}
	if strings.TrimSpace(v.SpotID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("spotID")
	}
	if !v.CheckedOutAt.IsZero() && v.CheckedOutAt.Before(v.CheckedInAt) {
		return validation.NewErrBadRecordFieldValue("checkedOutAt", "is before CheckedInAt")
	}
	return nil
}
