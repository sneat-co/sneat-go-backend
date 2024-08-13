package redirects

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type router interface {
	GET(path string, handle httprouter.Handle)
}

func InitRedirects(router router) {
	router.GET("/receipt", ReceiptRedirect)

	router.GET("/transfer",
		RedirectHandlerToEntityPageWithIntID("transfer=%d", "send"))

	router.GET("/contact",
		RedirectHandlerToEntityPageWithIntID("contact=%d"))

	router.GET("/open/new-debt", newDebtRedirect)

	router.GET("/choose-currency", chooseCurrencyRedirect)

	router.GET("/confirm", confirmEmailRedirect)

}

func RedirectHandlerToEntityPageWithIntID(path string, optionalParams ...string) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Value of 'id' parameter is not an integer"))
			return
		} else {
			redirectToWebApp(w, r, true, fmt.Sprintf(path, id), nil, optionalParams)
		}
	}
}
