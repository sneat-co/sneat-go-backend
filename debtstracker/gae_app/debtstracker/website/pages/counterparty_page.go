package pages

import (
	"fmt"
	"google.golang.org/appengine/v2"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/strongo/log"
	"golang.org/x/net/html"
)

func CounterpartyPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)
	log.Infof(c, "CounterpartyPage: %v", r.Method)
	counterpartyID := r.URL.Query().Get("id")
	if counterpartyID == "" {
		w.WriteHeader(404)
		_, _ = w.Write([]byte("missing required parameter: id"))
		return
	}

	counterparty, err := facade.GetContactByID(c, nil, counterpartyID)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte(fmt.Sprintf(`<html>
	<head><title>Contact: %v</title>
	<meta name="description" content="Transfered amount: %v">
	<link rel="canonical" href="./counterparty?id=%v" />
	<style>
	body{padding: 50px; font-family: Verdana; font-size: small;}
	th{padding-right:10px;text-align:left;}
	</style>
	</head>
	<body>
	<header><a href="/">DebtsTracker.io</a></header>
	<hr>
	<h1>Contact: %v</h1>

	<footer style="margin-top:50px; border-top: 1px solid lightgrey; padding-top:10px">
	<small style="color:grey">2016 &copy; Powered by <a href="https://golang.org/" target="_blank">Go lang</a> & <a href="https://cloud.google.com/appengine/" target="_blank">AppEngine</a></small>
	</footer>
	%v
	</body></html>`, html.EscapeString(counterparty.Data.FullName()), counterpartyID, html.EscapeString(counterparty.Data.FullName()), html.EscapeString(counterparty.Data.FullName()), GA_CODE)))
}
