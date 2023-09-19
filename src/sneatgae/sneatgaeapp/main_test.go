package sneatgaeapp

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestInitEmail(t *testing.T) {
	initEmail()
}

func TestInitSentry(t *testing.T) {
	initSentry()
}

func TestInitFirebase(t *testing.T) {
	initFirebase()
}

func Test_start(t *testing.T) {
	serve = func(handler http.Handler) {
		assert.NotNil(t, handler)
	}
	Start()
}
