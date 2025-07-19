package meetingus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/const4meetingus"
	"github.com/sneat-co/sneat-go-core/extension"
	"net/http"
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
