package sportus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/const4sportus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4sportus.ExtensionID,
		HandlersCount: 8,
		DelayersCount: 0,
	})
}
