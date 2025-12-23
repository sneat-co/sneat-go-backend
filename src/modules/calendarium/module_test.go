package calendarium

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4calendarium.ExtensionID,
		HandlersCount: 14,
		DelayersCount: 1,
	})
}
