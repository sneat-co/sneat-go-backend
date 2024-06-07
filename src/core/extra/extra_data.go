package extra

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/strongo/validation"
)

type Type string

type Data interface {
	//GetType() Type
	RequiredFields() []string
	IndexedFields() []string
	GetBrief() Data
	Validate() error
}

type BaseData struct {
	ExtraType string `json:"extraType" firestore:"extraType"` // This is mostly a workaround for TypeScript code-completion
}

//func (v *BaseData) GetType() Type {
//	return Type(v.ExtraType)
//}
//
//func (v *BaseData) Validate() error {
//	if v.ExtraType == "" {
//		return validation.NewErrRecordIsMissingRequiredField("extraType")
//	}
//	return nil
//
//}

// WithExtraField defines and `Extra` field to store extension data
type WithExtraField struct {
	ExtraType Type           `json:"extraType" firestore:"extraType"`
	Extra     map[string]any `json:"extra,omitempty" firestore:"extra,omitempty"`
	extraData Data
}

//func (v *WithExtraField) MarshalJSON() ([]byte, error) {
//	b, err := json.Marshal(v)
//	return b, err
//}
//
//func (v *WithExtraField) UnmarshalJSON(input []byte) error {
//	data := make(map[string]any)
//	err := json.Unmarshal(input, data)
//	v.ExtraType = data["extraType"].(Type)
//	v.Extra = data["extra"].(map[string]any)
//	return err
//}

// SetExtra sets extraData data
func (v *WithExtraField) SetExtra(extraType Type, extraData Data) (err error) {
	if extraType == "" {
		return errors.New("extraType is a required argument to set extraData")
	}
	v.ExtraType = extraType
	v.extraData = extraData
	if extraData == nil {
		v.Extra = make(map[string]any)
		return nil
	}
	var b []byte
	if b, err = json.Marshal(extraData); err != nil {
		return fmt.Errorf("failed to marshal extraData extraData to JSON: %w", err)
	}
	if err = json.Unmarshal(b, &v.Extra); err != nil {
		return fmt.Errorf("failed to unmarshal JSON extraData to extraData type %t: %w", extraData, err)
	}
	if len(v.Extra) == 0 {
		v.Extra = nil
	}
	return nil
}

// GetExtraData returns extraData data as module specific strongly typed Data
func (v *WithExtraField) GetExtraData() (extra Data, err error) {
	if v.extraData == nil {
		v.extraData = NewExtraData(v.ExtraType)
	}
	if len(v.Extra) == 0 {
		return v.extraData, nil
	}

	var b []byte
	if b, err = json.Marshal(v.Extra); err != nil {
		return nil, fmt.Errorf("failed to marshal extraData data to JSON: %w", err)
	}

	if err = json.Unmarshal(b, &v.extraData); err != nil {
		return nil, err
	}
	return v.extraData, nil
}

// Validate returns error if not valid
func (v *WithExtraField) Validate() error {
	if v.ExtraType == "" {
		return validation.NewErrRecordIsMissingRequiredField("extraType")
	}
	if v.Extra == nil {
		return validation.NewErrRecordIsMissingRequiredField("extraData")
	}
	if extra, err := v.GetExtraData(); err != nil {
		return validation.NewErrBadRecordFieldValue("extraData", fmt.Errorf("failed to get extraData data: %w", err).Error())
	} else if err = extra.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("extraData", err.Error())
	}
	return nil
}
