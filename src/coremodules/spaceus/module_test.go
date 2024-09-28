package spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/const4spaceus"
	"github.com/sneat-co/sneat-go-core/module"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	module.AssertModule(t, m, module.Expected{
		ModuleID:      const4spaceus.ModuleID,
		HandlersCount: 7,
		DelayersCount: 0,
	})
}
