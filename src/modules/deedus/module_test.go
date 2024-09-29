package deedus

import (
	"github.com/sneat-co/sneat-go-core/module"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	module.AssertModule(t, m, module.Expected{
		ModuleID:      "deedus",
		HandlersCount: 0,
		DelayersCount: 0,
	})
}
