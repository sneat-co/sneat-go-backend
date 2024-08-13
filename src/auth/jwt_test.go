package auth

import (
	"strings"
	"testing"
)

func TestIssueToken(t *testing.T) {
	token := IssueToken("123", "unit-test", false)
	if token == "" {
		t.Error("Token is empty")
	}
	parts := strings.Split(token, ".")
	const expected = 2
	if len(parts) != expected {
		t.Errorf("Unexpected number of token parts: expected %d, got %d, token: %v", expected, len(parts), token)
	}
}
