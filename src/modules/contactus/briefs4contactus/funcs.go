package briefs4contactus

// GetFullContactID returns full member ID
func GetFullContactID(teamID, memberID string) string {
	if teamID == "" {
		panic("teamID is required parameter")
	}
	if memberID == "" {
		panic("memberID is required parameter")
	}
	return teamID + ":" + memberID
}

// IsUniqueShortTitle checks if a given value is an unique member title
func IsUniqueShortTitle(v string, contacts map[string]*ContactBrief, role string) bool {
	for _, c := range contacts {
		if c.ShortTitle == v && (role == "" || c.HasRole(role)) {
			return false
		}
	}
	return true
}
