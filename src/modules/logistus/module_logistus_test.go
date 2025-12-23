package logistus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/const4logistus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4logistus.ExtensionID,
		HandlersCount: 21,
		DelayersCount: 0,
	})
}
