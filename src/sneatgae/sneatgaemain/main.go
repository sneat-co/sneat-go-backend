package main

import (
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
)

func main() { // TODO: document why we need this wrapper
	httpRouter := sneatgaeapp.CreateHttpRouter()
	sneatgaeapp.Start(httpRouter)
}
