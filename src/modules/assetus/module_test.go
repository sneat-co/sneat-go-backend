package assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/module"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	module.AssertModule(t, m, module.Expected{
		ModuleID:      const4assetus.ModuleID,
		HandlersCount: 4,
		DelayersCount: 0,
	})
}
