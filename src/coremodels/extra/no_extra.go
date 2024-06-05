package extra

import (
	"fmt"
	"github.com/strongo/validation"
)

func NewDataNoExtra() Data {
	return &noExtra{BaseData{Type: "empty"}}
}

var _ Data = (*noExtra)(nil)

// noExtra is used if no extension data is required by an asset type
type noExtra struct {
	BaseData
}

func (noExtra) RequiredFields() []string {
	return nil
}

func (noExtra) IndexedFields() []string {
	return nil
}

func (noExtra) GetBrief() Data {
	return nil
}

// Validate always returns nil
func (v noExtra) Validate() error {
	if v.Type != "" {
		return validation.NewErrBadRecordFieldValue("type", fmt.Sprintf("unexpected value: %s", v.Type))
	}
	return nil
}
