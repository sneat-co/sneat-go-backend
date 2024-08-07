package api4spaceus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpPostCreateSpace is an API endpoint that creates a new team
func httpPostCreateSpace(w http.ResponseWriter, r *http.Request) {
	var request dto4spaceus.CreateSpaceRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			facadeResponse, err := facade4spaceus.CreateSpace(ctx, userCtx, request)
			if err != nil {
				return nil, err
			}
			var apiResponse dto4spaceus.CreateSpaceResponse
			apiResponse.Space.ID = facadeResponse.Space.ID
			apiResponse.Space.Dto = *facadeResponse.Space.Data
			return apiResponse, err
		})
}
