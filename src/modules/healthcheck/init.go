package healthcheck

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/extension"
)

// InitHealthCheck registers health check HTTP handlers
func InitHealthCheck(handle extension.HTTPHandleFunc) {
	handle(http.MethodGet, "/health-check", httpGetPage)
	handle(http.MethodGet, "/health-check/test-email", httpGetTestEmail)
}
