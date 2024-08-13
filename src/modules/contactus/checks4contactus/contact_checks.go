package checks4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"slices"
)

func IsSpaceMember(roles []string) bool {
	return slices.Contains(roles, const4contactus.SpaceMemberRoleMember)
}
