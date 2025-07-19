package listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4listus.ExtensionID,
		HandlersCount: 6,
		DelayersCount: 0,
	})
}
