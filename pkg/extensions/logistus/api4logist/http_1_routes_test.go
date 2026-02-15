package api4logist

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicoretest"
	"github.com/sneat-co/sneat-go-core/sneatauth"
)

func TestRegisterLogistRoutes(t *testing.T) {
	var handle = func(method, path string, handler http.HandlerFunc) {
		const prefix = "/v0/logistus/"
		if !strings.HasPrefix(path, prefix) {
			t.Errorf("Unexpected path: %s - should start with %s", path, prefix)
		}
		switch method {
		case http.MethodGet, http.MethodPost, http.MethodDelete:
			break // OK
		default:
			t.Errorf("Unexpected method [%s] for path [%s]", method, path)
		}
		var r = httptest.NewRequest(method, path, nil)
		apicore.GetAuthTokenFromHttpRequest = func(r *http.Request, authRequired bool) (token *sneatauth.Token, err error) {
			return nil, nil
		}
		//apicore.NewContextWithToken = func(r *http.Request, authRequired bool) (ctx context.Context, err error) {
		//	return sneatfb.NewContextWithFirebaseToken(r.Context(), &auth.Token{UID: "user1"}), nil
		//}
		apicoretest.TestEndpoint(t, handler, apicoretest.AssertOptions{
			AuthRequired: true,
		}, r)
	}
	RegisterHttpRoutes(handle)
}
