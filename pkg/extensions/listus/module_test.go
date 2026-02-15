package listus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/const4listus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4listus.ExtensionID,
		HandlersCount: 6,
		DelayersCount: 0,
	})
}
