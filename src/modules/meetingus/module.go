package meetingus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/const4meetingus"
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

func Module() modules.Module {
	return modules.NewModule(const4meetingus.ModuleID, func(handle modules.HTTPHandleFunc) {
		handle("POST", "/api4meetingus/about", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("meetingus"))
		})
	})
}
