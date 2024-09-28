package core4spaceus

import (
	"fmt"
	"strings"
)

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

type SpaceRef string

func (v SpaceRef) SpaceType() SpaceType {
	if i := strings.Index(string(v), SpaceRefSeparator); i > 0 {
		return SpaceType(v[:i])
	}
	return ""
}

func (v SpaceRef) SpaceID() string {
	if i := strings.Index(string(v), SpaceRefSeparator); i > 0 {
		return string(v[i+1:])
	}
	return ""
}

func (v SpaceRef) UrlPath() string {
	return fmt.Sprintf("%s/%s", v.SpaceType(), v.SpaceID())
}

const SpaceRefSeparator = "!"

func NewSpaceRef(spaceType SpaceType, spaceID string) SpaceRef {
	return SpaceRef(string(spaceType) + SpaceRefSeparator + spaceID)
}

// IsValidSpaceType checks if space has a valid/known type
func IsValidSpaceType(v SpaceType) bool {
	switch v {
	case SpaceTypeFamily, SpaceTypePrivate, SpaceTypeCompany, SpaceTypeTeam, SpaceTypeClub:
		return true
	default:
		return false
	}
}
