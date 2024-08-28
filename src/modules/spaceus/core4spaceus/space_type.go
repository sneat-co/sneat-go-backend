package core4spaceus

type SpaceType string

const (
	// SpaceTypePrivate is a "private" space type
	SpaceTypePrivate SpaceType = "private"

	// SpaceTypeFamily is a "family" space type
	SpaceTypeFamily SpaceType = "family"

	// SpaceTypeCompany is a "company" space type
	SpaceTypeCompany SpaceType = "company"

	// SpaceTypeTeam is a "space" space type
	SpaceTypeTeam SpaceType = "team"

	// SpaceTypeClub is a "club" space type
	SpaceTypeClub SpaceType = "club"
)

// IsValidSpaceType checks if space has a valid/known type
func IsValidSpaceType(v SpaceType) bool {
	switch v {
	case SpaceTypeFamily, SpaceTypePrivate, SpaceTypeCompany, SpaceTypeTeam, SpaceTypeClub:
		return true
	default:
		return false
	}
}
