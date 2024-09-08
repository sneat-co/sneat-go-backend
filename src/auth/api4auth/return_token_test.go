package api4auth

import (
	"testing"
)

func TestIsAdminClaim(t *testing.T) {
	option := IsAdminClaim()
	var claims TokenClaims
	if claims.IsAdmin() {
		t.Error("Expected to be not admin")
	}
	option(&claims)
	if !claims.isAdmin {
		t.Error("Expected to be admin")
	}
	if !claims.IsAdmin() {
		t.Error("Expected to be admin")
	}
}
