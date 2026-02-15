package assetus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4assetus.ExtensionID,
		HandlersCount: 4,
		DelayersCount: 0,
	})
}
