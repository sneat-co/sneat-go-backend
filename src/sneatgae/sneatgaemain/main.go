package main

import (
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"github.com/sneat-co/sneat-go-core/emails/email2writer"
	"github.com/sneat-co/sneat-go-core/security"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"io"
	"os"
	"strings"
)

func main() { // TODO: document why we need this wrapper

	logus.AddLogEntryHandler(logus.NewStandardGoLogger())

	knownHosts := os.Getenv("KNOWN_HOSTS")
	if knownHosts != "" {
		security.AddKnownHosts(strings.Split(knownHosts, ",")...)
	}

	emailClient := email2writer.NewClient(func() (io.StringWriter, error) {
		return os.Stdout, nil
	})

	httpRouter := sneatgaeapp.CreateHttpRouter()

	serveStaticFiles(httpRouter)

	delaying.Init(delaying.VoidWithLog)

	sneatgaeapp.Start(nil, nil, httpRouter, emailClient)
}
