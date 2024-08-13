package const4contactus

import (
	"slices"
)

type SpaceMemberRole = string

const (
	SpaceMemberRoleMember   = "member"
	SpaceMemberRoleExMember = "ex-member"

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

	// SpaceMemberRoleExcluded if space members are excluded
	SpaceMemberRoleExcluded SpaceMemberRole = "excluded"
)

// SpaceMemberWellKnownRoles defines known roles
var SpaceMemberWellKnownRoles = []SpaceMemberRole{
	SpaceMemberRoleAdmin,
	SpaceMemberRoleContributor,
	SpaceMemberRoleCreator,
	SpaceMemberRoleMember,
	SpaceMemberRoleExMember,
	SpaceMemberRoleChild,
	SpaceMemberRoleAdult,
	SpaceMemberRoleSpectator,
	SpaceMemberRoleExcluded,
}

// IsKnownSpaceMemberRole checks if a role has a valid value
func IsKnownSpaceMemberRole(role SpaceMemberRole, teamRoles []SpaceMemberRole) bool {
	return teamRoles == nil || slices.Contains(SpaceMemberWellKnownRoles, role) || slices.Contains(teamRoles, role)
}
