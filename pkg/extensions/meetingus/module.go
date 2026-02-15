package meetingus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/const4meetingus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Module() extension.Config {
	return extension.NewExtension(const4meetingus.ExtensionID,
		extension.RegisterRoutes(func(handle extension.HTTPHandleFunc) {
			handle(http.MethodPost, "/api4meetingus/about", func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("meetingus"))
			})
		}),
	)
}
