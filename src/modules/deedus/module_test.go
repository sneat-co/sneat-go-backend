package deedus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         "deedus",
		HandlersCount: 0,
		DelayersCount: 0,
	})
}
