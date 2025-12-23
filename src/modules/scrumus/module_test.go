package scrumus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/const4srumus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4srumus.ExtensionID,
		HandlersCount: 9,
		DelayersCount: 0,
	})
}
