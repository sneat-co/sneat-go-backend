package redirects

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func confirmEmailRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()
	email, pin := query.Get("email"), query.Get("pin")
	if email == "" || pin == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	redirectToWebApp(w, r, false,
		fmt.Sprintf("confirm-email=%v:%v", url.QueryEscape(email), url.QueryEscape(pin)),
		map[string]string{}, []string{},
	)
}
