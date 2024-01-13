package api4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var thumbUp = facade4scrumus.ThumbUp

// httpPostThumbUp add a thumb up
func httpPostThumbUp(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request facade4scrumus.ThumbUpRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = thumbUp(ctx, userContext, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
