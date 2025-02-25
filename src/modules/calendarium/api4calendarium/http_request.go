package api4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"net/http"
)

func getHappeningRequestParamsFromURL(r *http.Request) (request dto4calendarium.HappeningRequest) {
	query := r.URL.Query()
	request.SpaceID = coretypes.SpaceID(query.Get("spaceID"))
	request.HappeningID = query.Get("happeningID")
	request.HappeningType = query.Get("happeningType")
	return
}
