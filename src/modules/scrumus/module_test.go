package scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/const4srumus"
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4srumus.ExtensionID,
		HandlersCount: 9,
		DelayersCount: 0,
	})
}
