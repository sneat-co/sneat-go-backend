package api4calendarius

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func getHappeningRequestParamsFromURL(r *http.Request) (request dto4calendarius.HappeningRequest) {
	query := r.URL.Query()
	request.SpaceID = coretypes.SpaceID(query.Get("spaceID"))
	request.HappeningID = query.Get("happeningID")
	request.HappeningType = query.Get("happeningType")
	return
}
