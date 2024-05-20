package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"net/http"
)

func getListRequestParamsFromURL(r *http.Request) (request facade4listus.ListRequest) {
	query := r.URL.Query()
	request.TeamID = query.Get("teamID")
	request.ListID = query.Get("listID")
	return
}
