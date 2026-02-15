package dbo4logist

import (
	"fmt"
	"strings"

	"github.com/strongo/validation"
)

// FreightLoad is a model for freight load
type FreightLoad struct {
	NumberOfPallets int    `json:"numberOfPallets,omitempty" firestore:"numberOfPallets,omitempty"`
	GrossWeightKg   int    `json:"grossWeightKg,omitempty" firestore:"grossWeightKg,omitempty"`
	VolumeM3        int    `json:"volumeM3,omitempty" firestore:"volumeM3,omitempty"` // 1m3 = 1000L
	Note            string `json:"note,omitempty" firestore:"note,omitempty"`
}

// IsEmpty returns true if freight load is empty
func (v *FreightLoad) IsEmpty() bool {
	return v == nil || *v == FreightLoad{}
}

// Add adds another freight load to this one
func (v *FreightLoad) Add(v2 *FreightLoad) {
	if v2 == nil {
		return
	}
	v.NumberOfPallets += v2.NumberOfPallets
	v.GrossWeightKg += v2.GrossWeightKg
	v.VolumeM3 += v2.VolumeM3
}

// Validate validates freight load
func (v *FreightLoad) Validate() error {
	if v.NumberOfPallets < 0 {
		return validation.NewErrBadRecordFieldValue("numberOfPallets", fmt.Sprintf("negative value: [%v]", v.NumberOfPallets))
	}
	if v.GrossWeightKg < 0 {
		return validation.NewErrBadRecordFieldValue("grossWeightKg", fmt.Sprintf("negative value: [%v]", v.GrossWeightKg))
	}
	if v.VolumeM3 < 0 {
		return validation.NewErrBadRecordFieldValue("volumeM3", fmt.Sprintf("negative value: [%v]", v.VolumeM3))
	}
	if v.Note != strings.TrimSpace(v.Note) {
		return validation.NewErrBadRecordFieldValue("note", "should be trimmed")
	}
	return nil
}
