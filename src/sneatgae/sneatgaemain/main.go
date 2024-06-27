package main

import (
	"github.com/pkg/profile"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"github.com/sneat-co/sneat-go-core/emails/email2writer"
	"github.com/strongo/log"
	"github.com/strongo/slice"
	"io"
	golog "log"
	"os"
)

func main() { // TODO: document why we need this wrapper

	defaultLogger := golog.Default()
	log.AddLogger(log.NewPrinter("log.Default()", func(format string, a ...any) (n int, err error) {
		defaultLogger.Printf(format, a...)
		return 0, nil
	}))

	if slice.Contains(os.Args, "pprof") {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	httpRouter := sneatgaeapp.CreateHttpRouter()

	emailClient := email2writer.NewClient(func() (io.StringWriter, error) {
		return os.Stdout, nil
	})
	sneatgaeapp.Start(nil, nil, httpRouter, emailClient)
}
