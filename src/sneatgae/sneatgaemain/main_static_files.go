package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

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
