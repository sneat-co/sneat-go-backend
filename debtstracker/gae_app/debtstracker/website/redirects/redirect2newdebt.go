package redirects

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func newDebtRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirectToWebApp(w, r, true, "/main/debts/new-debt", map[string]string{}, []string{})
}
