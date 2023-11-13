package main

import (
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"github.com/sneat-co/sneat-go-core/emails/email2console"
)

func main() { // TODO: document why we need this wrapper
	httpRouter := sneatgaeapp.CreateHttpRouter()
	emailClient := email2console.NewClient()
	sneatgaeapp.Start(nil, nil, httpRouter, emailClient)
}
