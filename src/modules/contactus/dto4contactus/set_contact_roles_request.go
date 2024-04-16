package dto4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/strongo/validation"
)

// SetRolesRequest request to set contact address
type SetRolesRequest struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// Validate returns error if request is invalid
func (v SetRolesRequest) Validate() error {
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

	for _, remove := range v.Remove {
		if remove == const4contactus.TeamMemberRoleMember {
			return validation.NewErrBadRequestFieldValue("remove", "use remove_member endpoint to remove members from team")
		}
	}

	return nil
}
