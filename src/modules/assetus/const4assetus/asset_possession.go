package const4assetus

import (
	"fmt"
	"github.com/strongo/validation"
	"slices"
)

// AssetPossession is a type of asset possession
type AssetPossession = string

const (
	AssetPossessionUnknown     AssetPossession = "unknown"
	AssetPossessionUndisclosed AssetPossession = "undisclosed"
	AssetPossessionOwning      AssetPossession = "owning"
	AssetPossessionLeasing     AssetPossession = "leasing"
	AssetPossessionRenting     AssetPossession = "renting"
)

// AssetPossessions is a list of all possible possession values
var AssetPossessions = []AssetPossession{
	AssetPossessionUnknown,
	AssetPossessionUndisclosed,
	AssetPossessionOwning,
	AssetPossessionLeasing,
	AssetPossessionRenting,
}

// ValidateAssetPossession validates possession
func ValidateAssetPossession(v AssetPossession, required bool) error {
	if required && v == "" {
		return validation.NewErrRecordIsMissingRequiredField("possession")
	}
	if !slices.Contains(AssetPossessions, v) {
		return validation.NewErrBadRecordFieldValue("possession", fmt.Sprintf("unknown possession '%s', expected values: %s", v, AssetPossessions))
	}
	return nil
}
