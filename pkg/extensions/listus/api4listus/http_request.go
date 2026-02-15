package api4listus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func getListRequestParamsFromURL(r *http.Request) (request dto4listus.ListRequest) {
	query := r.URL.Query()
	request.SpaceID = coretypes.SpaceID(query.Get("spaceID"))
	request.ListID = dbo4listus.ListKey(query.Get("listID"))
	return
}
