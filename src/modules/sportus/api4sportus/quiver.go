package api4sportus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/dbo4sportus"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/facade4sportus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

const (
	quiverPathPrefix = "/v0/quiver/"
	//quiverMyPathPrefix     = quiverPathPrefix + "my/"
	quiverWantedPathPrefix = quiverPathPrefix + "wanted/"
)

func registerQuiverHandlers(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, quiverWantedPathPrefix+"create_wanted", createWantedItem)
	handle(http.MethodPut, quiverWantedPathPrefix+"update_wanted", updateWantedItem)
	handle(http.MethodDelete, quiverWantedPathPrefix+"delete_wanted", deleteWantedItem)
}

func createWantedItem(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	wanted := dbo4sportus.Wanted{}
	request := facade4sportus.CreateWantedRequest{
		Wanted: wanted,
	}
	if err := apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	var id string
	if id, err = facade4sportus.CreateWanted(ctx, userContext, request); err != nil {
		httpserver.HandleError(ctx, err, "createWantedItem", w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(id))
}

func updateWantedItem(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func deleteWantedItem(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	request := facade4sportus.DeleteWantedRequest{ID: r.URL.Query().Get("id")}
	if err = facade4sportus.DeleteWanted(ctx, userContext, request); err != nil {
		httpserver.HandleError(ctx, err, "deleteWantedItem", w, r)
	}
	w.WriteHeader(http.StatusOK)
}
