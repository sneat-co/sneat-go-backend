package main

import (
	"github.com/pkg/profile"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"github.com/sneat-co/sneat-go-core/emails/email2writer"
	"github.com/strongo/logus"
	"github.com/strongo/slice"
	"io"
	"os"
)

func main() { // TODO: document why we need this wrapper

	logus.AddLogEntryHandler(logus.NewStandardGoLogger())

	if slice.Contains(os.Args, "pprof") {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	emailClient := email2writer.NewClient(func() (io.StringWriter, error) {
		return os.Stdout, nil
	})

	httpRouter := sneatgaeapp.CreateHttpRouter()

	serveStaticFiles(httpRouter)

	initBots(httpRouter)

	sneatgaeapp.Start(nil, nil, httpRouter, emailClient)
}
