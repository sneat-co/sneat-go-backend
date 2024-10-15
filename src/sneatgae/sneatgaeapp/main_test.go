package sneatgaeapp

import (
	"github.com/sneat-co/sneat-go-core/emails/email2writer"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/delaying"
	"io"
	"net/http"
	"os"
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
	emailClient := email2writer.NewClient(func() (io.StringWriter, error) {
		return os.Stdout, nil
	})

	delaying.InitNoopLogging()

	Start(nil, nil, httpRouter, emailClient)
}
