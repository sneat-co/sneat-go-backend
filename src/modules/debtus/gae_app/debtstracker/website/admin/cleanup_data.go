package admin

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CleanupPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch r.Method {
	case "GET":
		_, _ = w.Write([]byte("<form method=post><button type=submit></form>"))
	case "POST":
		_, _ = w.Write([]byte("Not implemented yet"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("unexpected method: " + r.Method))
	}
}
