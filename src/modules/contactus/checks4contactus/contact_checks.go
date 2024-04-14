package checks4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/strongo/slice"
)

func IsTeamMember(roles []string) bool {
	return slice.Contains(roles, const4contactus.TeamMemberRoleMember)
}
