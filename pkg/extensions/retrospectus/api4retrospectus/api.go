package api4retrospectus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

var verifyRequest = apicore.VerifyRequestAndCreateUserContext /* func(
	w http.ResponseWriter, r *http.Request,
	options verify.RequestOptions,
) (ctx context.Context, userContext facade4debtus.User, err error) {
	return apicore.VerifyRequestAndCreateUserContext(w, r, options)
}
*/

func verifyAuthorizedJSONRequest(
	w http.ResponseWriter, r *http.Request,
	minSize, maxSize int64,
) (ctx facade.ContextWithUser, err error) {
	o := verify.Request(
		verify.AuthenticationRequired(true),
		verify.MinimumContentLength(minSize),
		verify.MaximumContentLength(maxSize),
	)
	return verifyRequest(w, r, o)
}
