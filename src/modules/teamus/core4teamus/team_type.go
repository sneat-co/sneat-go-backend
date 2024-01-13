package core4teamus

type TeamType string

const (
	// TeamTypeFamily is a "family" team type
	TeamTypeFamily TeamType = "family"

	// TeamTypeCompany is a "company" team type
	TeamTypeCompany TeamType = "company"

	// TeamTypeTeam is a "team" team type
	TeamTypeTeam TeamType = "team"

	// TeamTypeClub is a "club" team type
	TeamTypeClub TeamType = "club"
)

// IsValidTeamType checks if team has a valid/known type
func IsValidTeamType(v TeamType) bool {
	switch v {
	case TeamTypeFamily, TeamTypeCompany, TeamTypeTeam, TeamTypeClub:
		return true
	default:
		return false
	}
}
