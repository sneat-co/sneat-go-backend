package api4retrospectus

import (
	"context"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var verifyRequest = func(
	w http.ResponseWriter, r *http.Request,
	options verify.RequestOptions,
) (ctx context.Context, userContext facade.User, err error) {
	return apicore.VerifyRequestAndCreateUserContext(w, r, options)
}

func verifyAuthorizedJSONRequest(
	w http.ResponseWriter, r *http.Request,
	minSize, maxSize int64,
) (ctx context.Context, userContext facade.User, err error) {
	return verifyRequest(w, r, verify.DefaultJsonWithAuthRequired)
}
