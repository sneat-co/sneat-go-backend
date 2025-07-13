package meetingus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/const4meetingus"
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

func Module() module.Module {
	return module.NewExtension(const4meetingus.ModuleID, module.RegisterRoutes(func(handle module.HTTPHandleFunc) {
		handle(http.MethodPost, "/api4meetingus/about", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("meetingus"))
		})
	}))
}
