package main

import (
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"github.com/sneat-co/sneat-go-core/emails/email2writer"
	"io"
	"os"
)

func main() { // TODO: document why we need this wrapper
	httpRouter := sneatgaeapp.CreateHttpRouter()

	emailClient := email2writer.NewClient(func() (io.StringWriter, error) {
		return os.Stdout, nil
	})
	sneatgaeapp.Start(nil, nil, httpRouter, emailClient)
}
