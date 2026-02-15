package retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus/api4retrospectus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus/const4retrospectus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4retrospectus.ExtensionID, extension.RegisterRoutes(api4retrospectus.RegisterHttpRoutes))
}
