package redirects

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func chooseCurrencyRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirectToWebApp(w, r, true, "/choose-currency/", map[string]string{}, []string{})
}
