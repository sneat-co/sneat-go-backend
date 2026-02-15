package api4logist

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-core/apicoretest"
)

func Test_httpAddContainers(t *testing.T) {
	var r = httptest.NewRequest(http.MethodPost, "/logistus/containers", strings.NewReader(`{}`))
	apicoretest.TestEndpoint(t, httpAddContainers, apicoretest.AssertOptions{
		AuthRequired: true,
	}, r)
}
