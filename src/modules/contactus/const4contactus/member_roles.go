package const4contactus

type TeamMemberRole = string

const (
	TeamMemberRoleMember = "member"

	TeamMemberRoleAdult = "adult"

	TeamMemberRoleChild = "child"

	// TeamMemberRoleCreator role of a creator
	TeamMemberRoleCreator TeamMemberRole = "creator"

	// TeamMemberRoleOwner role of an owner
	TeamMemberRoleOwner TeamMemberRole = "owner"

	TeamMemberRoleAdmin TeamMemberRole = "admin"

	// TeamMemberRoleContributor role of a contributor
	TeamMemberRoleContributor TeamMemberRole = "contributor"

	// TeamMemberRoleSpectator role of spectator
	TeamMemberRoleSpectator TeamMemberRole = "spectator"

	// TeamMemberRoleExcluded if team members is excluded
	TeamMemberRoleExcluded TeamMemberRole = "excluded"
)

// TeamMemberWellKnownRoles defines known roles
var TeamMemberWellKnownRoles = []TeamMemberRole{
	TeamMemberRoleAdmin,
	TeamMemberRoleContributor,
	TeamMemberRoleCreator,
	TeamMemberRoleMember,
	TeamMemberRoleSpectator,
	TeamMemberRoleExcluded,
}

// IsKnownTeamMemberRole checks if role has valid value
func IsKnownTeamMemberRole(role TeamMemberRole, teamRoles []TeamMemberRole) bool {
	for _, r := range TeamMemberWellKnownRoles {
		if r == role {
			return true
		}
	}
	if teamRoles == nil {
		return true
	}
	for _, r := range teamRoles {
		if r == role {
			return true
		}
	}
	return false
}
