package extras4assetus

import (
	"encoding/json"
	"fmt"
	"github.com/strongo/validation"
)

type AssetExtra interface {
	GetType() AssetExtraType
	RequiredFields() []string
	IndexedFields() []string
	GetBrief() AssetExtra
	Validate() error
}

type AssetExtraBase struct {
	Type AssetExtraType `json:"type" firestore:"type"`
}

func (v *AssetExtraBase) GetType() AssetExtraType {
	return v.Type
}

func (v *AssetExtraBase) Validate() error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	return nil

}

func NewAssetNoExtra() AssetExtra {
	return &assetNoExtra{AssetExtraBase{Type: "empty"}}
}

var _ AssetExtra = (*assetNoExtra)(nil)

// assetNoExtra is used if no extension data is required by an asset type
type assetNoExtra struct {
	AssetExtraBase
}

func (e assetNoExtra) RequiredFields() []string {
	return nil
}

func (e assetNoExtra) IndexedFields() []string {
	return nil
}

func (e assetNoExtra) GetBrief() AssetExtra {
	return nil
}

// Validate always returns nil
func (assetNoExtra) Validate() error {
	return nil
}

// WithAssetExtraField defines and `Extra` field to store extension data
type WithAssetExtraField struct {
	ExtraType AssetExtraType `json:"extraType" firestore:"extraType"`
	Extra     map[string]any `json:"extra" firestore:"extra"`
	extra     AssetExtra
}

func (v *WithAssetExtraField) SetExtra(extra AssetExtra) (err error) {
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

func (v *WithAssetExtraField) GetExtra() (extra AssetExtra, err error) {
	var b []byte
	if v.extra == nil {
		switch v.ExtraType {
		case AssetExtraTypeVehicle:
			v.extra = new(AssetVehicleExtra)
		case AssetExtraTypeDwelling:
			v.extra = new(AssetDwellingExtra)
		case AssetExtraTypeDocument:
			v.extra = new(AssetDocumentExtra)
		default:
			return nil, fmt.Errorf("unsupported extra type: %s", v.ExtraType)
		}
	}
	if len(v.Extra) == 0 {
		return v.extra, nil
	}
	if b, err = json.Marshal(v.Extra); err != nil {
		return nil, fmt.Errorf("failed to marshal extra data to JSON: %w", err)
	}

	if err = json.Unmarshal(b, &v.extra); err != nil {
		return nil, err
	}
	return v.extra, nil
}

func (v *WithAssetExtraField) Validate() error {
	if v.ExtraType == "" {
		return validation.NewErrRecordIsMissingRequiredField("extraType")
	}
	if v.Extra == nil {
		return validation.NewErrRecordIsMissingRequiredField("extra")
	}
	if extra, err := v.GetExtra(); err != nil {
		return validation.NewErrBadRecordFieldValue("extra", fmt.Errorf("failed to get extra data: %w", err).Error())
	} else if err = extra.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("extra", err.Error())
	}
	return nil
}
