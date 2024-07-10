package const4contactus

import "github.com/strongo/slice"

type SpaceMemberRole = string

const (
	SpaceMemberRoleMember = "member"

	SpaceMemberRoleAdult = "adult"

	SpaceMemberRoleChild = "child"

	// SpaceMemberRoleCreator role of a creator
	SpaceMemberRoleCreator SpaceMemberRole = "creator"

	// SpaceMemberRoleOwner role of an owner
	SpaceMemberRoleOwner SpaceMemberRole = "owner"

	SpaceMemberRoleAdmin SpaceMemberRole = "admin"

	// SpaceMemberRoleContributor role of a contributor
	SpaceMemberRoleContributor SpaceMemberRole = "contributor"

	// SpaceMemberRoleSpectator role of spectator
	SpaceMemberRoleSpectator SpaceMemberRole = "spectator"

	// SpaceMemberRoleExcluded if team members are excluded
	SpaceMemberRoleExcluded SpaceMemberRole = "excluded"
)

// SpaceMemberWellKnownRoles defines known roles
var SpaceMemberWellKnownRoles = []SpaceMemberRole{
	SpaceMemberRoleAdmin,
	SpaceMemberRoleContributor,
	SpaceMemberRoleCreator,
	SpaceMemberRoleMember,
	SpaceMemberRoleChild,
	SpaceMemberRoleAdult,
	SpaceMemberRoleSpectator,
	SpaceMemberRoleExcluded,
}

// IsKnownSpaceMemberRole checks if a role has a valid value
func IsKnownSpaceMemberRole(role SpaceMemberRole, teamRoles []SpaceMemberRole) bool {
	return teamRoles == nil || slice.Contains(SpaceMemberWellKnownRoles, role) || slice.Contains(teamRoles, role)
}
