package healthcheck

import (
	"net/http"
	"testing"
)

func TestInitHealthCheck(t *testing.T) {
	var handle = func(method string, path string, handler http.HandlerFunc) {
	}
	InitHealthCheck(handle)
}
