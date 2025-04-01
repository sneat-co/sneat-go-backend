package dbo4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strconv"
)

// HappeningSlot DBO
type HappeningSlot struct {
	HappeningSlotTiming
	dbo4linkage.WithRelated
	Locations []Location `json:"locations,omitempty" firestore:"locations,omitempty"`
}

func (v *HappeningSlot) IsEmpty() bool {
	return v == nil || v.HappeningSlotTiming.IsEmpty() && len(v.Locations) == 0 && len(v.Related) == 0
}

// Validate returns error if not valid
func (v *HappeningSlot) Validate() error {
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
