package assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/api4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4assetus.ExtensionID,
		extension.RegisterRoutes(api4assetus.RegisterHttpRoutes),
	)
}
