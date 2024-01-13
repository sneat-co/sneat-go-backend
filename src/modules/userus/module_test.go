package userus

import (
	"github.com/sneat-co/sneat-go-core/tests"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	tests.VerifyModule(t, m)
}
