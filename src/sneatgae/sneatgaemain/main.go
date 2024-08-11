package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/profile"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
	"github.com/sneat-co/sneat-go-core/emails/email2writer"
	"github.com/strongo/logus"
	"github.com/strongo/slice"
	"io"
	"net/http"
	"os"
)

func main() { // TODO: document why we need this wrapper

	logus.AddLogEntryHandler(logus.StandardGoLogger())

	if slice.Contains(os.Args, "pprof") {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	httpRouter := sneatgaeapp.CreateHttpRouter()

	staticDir := http.Dir("static")
	httpRouter.ServeFiles("/static/*filepath", staticDir)

	serveRootStaticFiles(httpRouter, staticDir)

	emailClient := email2writer.NewClient(func() (io.StringWriter, error) {
		return os.Stdout, nil
	})
	sneatgaeapp.Start(nil, nil, httpRouter, emailClient)
}

func serveRootStaticFiles(httpRouter *httprouter.Router, root http.FileSystem) {
	fileServer := http.FileServer(root)

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
