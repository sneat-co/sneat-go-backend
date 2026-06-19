package calendarius

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/const4calendarius"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4calendarius.ExtensionID,
		HandlersCount: 14,
		DelayersCount: 1,
	})
}
