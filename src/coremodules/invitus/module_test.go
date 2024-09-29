package invitus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/const4invitus"
	"github.com/sneat-co/sneat-go-core/module"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	module.AssertModule(t, m, module.Expected{
		ModuleID:      const4invitus.ModuleID,
		HandlersCount: 6,
		DelayersCount: 0,
	})
}
