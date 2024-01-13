package dto4contactus

import "github.com/strongo/validation"

// SetContactRolesRequest request to set contact address
type SetContactRolesRequest struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// Validate returns error if request is invalid
func (v SetContactRolesRequest) Validate() error {
	if len(v.Add) == 0 && len(v.Remove) == 0 {
		return validation.NewErrBadRequestFieldValue("add", "either add or remove must be provided")
	}

	for _, add := range v.Add {
		for _, remove := range v.Remove {
			if add == remove {
				return validation.NewErrBadRequestFieldValue("add", "cannot add and remove the same role")
			}
		}
	}

	return nil
}
