package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"net/http"
)

func getListRequestParamsFromURL(r *http.Request) (request dto4listus.ListRequest) {
	query := r.URL.Query()
	request.SpaceID = coretypes.SpaceID(query.Get("spaceID"))
	request.ListID = dbo4listus.ListKey(query.Get("listID"))
	return
}
