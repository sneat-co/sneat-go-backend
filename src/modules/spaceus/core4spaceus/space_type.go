package core4spaceus

type SpaceType string

const (
	// SpaceTypeFamily is a "family" team type
	SpaceTypeFamily SpaceType = "family"

	// SpaceTypeCompany is a "company" team type
	SpaceTypeCompany SpaceType = "company"

	// SpaceTypeTeam is a "space" team type
	SpaceTypeTeam SpaceType = "team"

	// SpaceTypeClub is a "club" team type
	SpaceTypeClub SpaceType = "club"
)

// IsValidSpaceType checks if team has a valid/known type
func IsValidSpaceType(v SpaceType) bool {
	switch v {
	case SpaceTypeFamily, SpaceTypeCompany, SpaceTypeTeam, SpaceTypeClub:
		return true
	default:
		return false
	}
}
