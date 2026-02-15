package api4calendarium

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func getHappeningRequestParamsFromURL(r *http.Request) (request dto4calendarium.HappeningRequest) {
	query := r.URL.Query()
	request.SpaceID = coretypes.SpaceID(query.Get("spaceID"))
	request.HappeningID = query.Get("happeningID")
	request.HappeningType = query.Get("happeningType")
	return
}
