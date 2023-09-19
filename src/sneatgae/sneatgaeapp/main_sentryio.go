package sneatgaeapp

import (
	"context"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/sneat-co/sneat-go-core/capturer"
	"log"
	"os"
)

var sentryHandler = sentryhttp.New(sentryhttp.Options{Repanic: true})

var sentryInitialized = false

func initSentry() {
	if sentryInitialized { // to account for `go test -race`
		return
	}
	sentryInitialized = true
	options := sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: "",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug:            true,
		AttachStacktrace: true,
		DebugWriter:      os.Stderr,
	}
	if err := sentry.Init(options); err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	capturer.AddErrorLogger(sentryIoLogger{})
}

type sentryIoLogger struct {
}

func (v sentryIoLogger) LogError(ctx context.Context, err error) {
	sentry.CaptureException(err)
}
