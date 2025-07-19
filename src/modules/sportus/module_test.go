package sportus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/const4sportus"
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4sportus.ExtensionID,
		HandlersCount: 8,
		DelayersCount: 0,
	})
}
