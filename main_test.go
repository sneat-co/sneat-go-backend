package main

import (
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"testing"
)

func TestServeStaticFiles(m *testing.T) {
	httpRouter := sneatgaeapp.CreateHttpRouter()
	serveStaticFiles(httpRouter)
}
