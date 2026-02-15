package api4logist

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-core/apicoretest"
)

func Test_httpDeleteContainer(t *testing.T) {
	var r = httptest.NewRequest(http.MethodDelete, RoutePathOrderDeleteContainer, strings.NewReader(`{}`))
	apicoretest.TestEndpoint(t, httpDeleteContainer, apicoretest.AssertOptions{AuthRequired: true}, r)
}
