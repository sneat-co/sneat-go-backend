package deedus

import (
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         "deedus",
		HandlersCount: 0,
		DelayersCount: 0,
	})
}
