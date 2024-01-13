package healthcheck

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// InitHealthCheck registers health check HTTP handlers
func InitHealthCheck(handle modules.HTTPHandleFunc) {
	handle(http.MethodGet, "/health-check", httpGetPage)
	handle(http.MethodGet, "/health-check/test-email", httpGetTestEmail)
}
