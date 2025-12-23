package retrospectus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/const4retrospectus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4retrospectus.ExtensionID,
		HandlersCount: 9,
		DelayersCount: 0,
	})
}
