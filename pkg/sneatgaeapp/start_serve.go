package sneatgaeapp

import (
	"log"
	"net/http"
	"os"
)

var serve = func(handler http.Handler) {
	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "4300"
		//logus.Printf("Defaulting to port %s", port)
	}
	// [END setting_port]

	log.Printf("Listening on port %s, http://localhost:%s", port, port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
