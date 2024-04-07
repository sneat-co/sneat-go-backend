package api4logist

import (
	"github.com/sneat-co/sneat-go-core/apicoretest"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_httpAddContainers(t *testing.T) {
	var r = httptest.NewRequest(http.MethodPost, "/logistus/containers", strings.NewReader(`{}`))
	apicoretest.TestEndpoint(t, httpAddContainers, apicoretest.AssertOptions{
		AuthRequired: true,
	}, r)
}
