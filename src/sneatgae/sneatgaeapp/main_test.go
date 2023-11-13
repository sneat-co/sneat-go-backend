package sneatgaeapp

import (
	"github.com/sneat-co/sneat-go-core/emails/email2console"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestInitFirebase(t *testing.T) {
	initFirebase()
}

func Test_start(t *testing.T) {
	serve = func(handler http.Handler) {
		assert.NotNil(t, handler)
	}
	httpRouter := CreateHttpRouter()
	Start(nil, nil, httpRouter, email2console.NewClient())
}
