package extra

import (
	"encoding/json"
	"fmt"
	"github.com/strongo/validation"
)

type Type string

type Data interface {
	GetType() Type
	RequiredFields() []string
	IndexedFields() []string
	GetBrief() Data
	Validate() error
}

type BaseData struct {
	Type Type `json:"type" firestore:"type"`
}

func (v *BaseData) GetType() Type {
	return v.Type
}

func (v *BaseData) Validate() error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	return nil

}

// WithExtraField defines and `Extra` field to store extension data
type WithExtraField struct {
	ExtraType Type           `json:"extraType" firestore:"extraType"`
	Extra     map[string]any `json:"extra,omitempty" firestore:"extra,omitempty"`
	extra     Data
}

// SetExtra sets extra data
func (v *WithExtraField) SetExtra(extra Data) (err error) {
	v.extra = extra
	if extra == nil {
		v.Extra = make(map[string]any)
	} else {
		var b []byte
		if b, err = json.Marshal(extra); err != nil {
			return fmt.Errorf("failed to marshal extra data to JSON: %w", err)
		}
		if err = json.Unmarshal(b, &v.Extra); err != nil {
			return fmt.Errorf("failed to unmarshal JSON data to extra type %t: %w", extra, err)
		}
	}
	return nil
}

// GetExtraData returns extra data as module specific strongly typed Data
func (v *WithExtraField) GetExtraData() (extra Data, err error) {
	if v.extra == nil {
		v.extra = newExtra(v.ExtraType)
	}
	if len(v.Extra) == 0 {
		return v.extra, nil
	}

	var b []byte
	if b, err = json.Marshal(v.Extra); err != nil {
		return nil, fmt.Errorf("failed to marshal extra data to JSON: %w", err)
	}

	if err = json.Unmarshal(b, &v.extra); err != nil {
		return nil, err
	}
	return v.extra, nil
}

// Validate returns error if not valid
func (v *WithExtraField) Validate() error {
	if v.ExtraType == "" {
		return validation.NewErrRecordIsMissingRequiredField("extraType")
	}
	if v.Extra == nil {
		return validation.NewErrRecordIsMissingRequiredField("extra")
	}
	if extra, err := v.GetExtraData(); err != nil {
		return validation.NewErrBadRecordFieldValue("extra", fmt.Errorf("failed to get extra data: %w", err).Error())
	} else if err = extra.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("extra", err.Error())
	}
	return nil
}
