package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"github.com/sneat-co/sneat-go-core/emails/email2writer"
	"github.com/sneat-co/sneat-go-core/security"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"io"
	"net/http"
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

func serveStaticFiles(httpRouter *httprouter.Router) {

	staticDir := http.Dir("static")
	httpRouter.ServeFiles("/static/*filepath", staticDir)

	fileServer := http.FileServer(staticDir)

	for _, file := range []struct {
		path        string
		contentType string
	}{
		{path: "robots.txt", contentType: "text/plain"},
		{path: "no-robots.txt", contentType: "text/plain"},
		{path: "favicon.ico", contentType: "image/x-icon"},
	} {
		httpRouter.GET("/"+file.path, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
			w.Header().Set("Content-Type", file.contentType)
			w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour (3600 seconds)
			fileServer.ServeHTTP(w, req)
		})
	}
}
