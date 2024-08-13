package redirects

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	pages2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/website/pages"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2"
	"html/template"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine/v2/datastore"
)

var receiptOpenGraphPageTmpl *template.Template

func ReceiptRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)
	query := r.URL.Query()
	receiptCode := query.Get("id")
	if receiptCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	receiptID := receiptCode
	if receiptID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var err error
	logus.Debugf(c, "Receipt ContactID: %v", receiptID)
	_, err = dtdal.Receipt.GetReceiptByID(c, nil, receiptID)
	switch err {
	case nil: //pass
	case datastore.ErrNoSuchEntity:
		logus.Debugf(c, "Receipt not found by ContactID")
		http.NotFound(w, r)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		logus.Errorf(c, err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	//lang := query.Get("lang")
	//if lang == "" {
	//	lang = receipt.Data.Lang
	//}

	if strings.HasPrefix(r.UserAgent(), "facebookexternalhit/") || query.Get("for") == "fb" {
		if receiptOpenGraphPageTmpl == nil {
			receiptOpenGraphPageTmpl = template.Must(template.ParseFiles(pages2.TEMPLATES_PATH + "receipt-opengraph.html"))
		}
		locale := i18n.LocaleEnUS // strongoapp.GetLocaleByCode5(receipt.Lang) // TODO: Check for empty
		pages2.RenderCachedPage(w, r, receiptOpenGraphPageTmpl, locale, map[string]interface{}{
			"host":      r.Host,
			"ogUrl":     r.URL.String(),
			"ReceiptID": receiptID,
			//"ReceiptCode": shared.EncodeID(receiptID),
			"Title":       fmt.Sprintf("Receipt @ DebtsTracker.io #%v", receiptID),
			"Description": "Receipt description goes here",
		}, 9)
	} else {
		redirectToWebApp(w, r, false, common4debtus.Deeplink.AppHashPathToReceipt(receiptID), map[string]string{}, []string{})
	}
}
