package assetus

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/api4assetus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/extension"
)

// Extension returns description of the assetus extension
func Extension() extension.Config {
	return extension.NewExtension(const4assetus.ExtensionID,
		extension.RegisterRoutes(api4assetus.RegisterHttpRoutes),
	)
}
