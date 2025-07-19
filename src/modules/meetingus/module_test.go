package meetingus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/const4meetingus"
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4meetingus.ExtensionID,
		HandlersCount: 1,
		DelayersCount: 0,
	})
}
