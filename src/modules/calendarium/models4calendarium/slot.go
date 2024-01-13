package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strconv"
	"strings"
)

// HappeningSlot DTO
type HappeningSlot struct {
	ID string `json:"id" firestore:"id"`
	HappeningSlotTiming
	Locations []Location `json:"locations,omitempty" firestore:"locations,omitempty"`
}

// Validate returns error if not valid
func (v HappeningSlot) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := validate.RecordID(v.ID); err != nil {
		return validation.NewErrBadRecordFieldValue("id", err.Error())
	}
	numberOfPhysicalLocations := 0
	for i, l := range v.Locations {
		if err := l.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("locations[%v]", i),
				err.Error(),
			)
		}
		if l.Type == "physical" {
			numberOfPhysicalLocations++
		}
	}
	if numberOfPhysicalLocations > 1 {
		return validation.NewErrBadRecordFieldValue("locations", "only one physical location is allowed, got "+strconv.Itoa(numberOfPhysicalLocations))
	}
	if err := v.HappeningSlotTiming.Validate(); err != nil {
		return fmt.Errorf("failed validation of happening slot timing: %w", err)
	}
	return nil
}

// Location DTO
type Location struct {
	Type  string // e.g. physical or online
	Title string
}

// Validate returns error if not valid
func (v Location) Validate() error {
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case "physical", "online":
		break
	default:
		return validation.NewErrBadRecordFieldValue("type", "expected to be either 'physical' or 'online'")
	}
	if err := validate.RecordTitle(v.Title, "title"); err != nil {
		return err
	}
	return nil
}
