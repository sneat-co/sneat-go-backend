package redirects

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/common4all"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	pages2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/website/pages"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"html/template"
	"net/http"
	"strings"
)

var receiptOpenGraphPageTmpl *template.Template

func ReceiptRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := r.Context()
	query := r.URL.Query()
	receiptCode := query.Get("id")
	if receiptCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	receiptID := receiptCode
	var err error
	logus.Debugf(c, "Receipt ContactID: %v", receiptID)
	_, err = dtdal.Receipt.GetReceiptByID(c, nil, receiptID)
	if err != nil {
		if dal.IsNotFound(err) {
			logus.Debugf(c, "Receipt not found by ContactID")
			http.NotFound(w, r)
			return
		}
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
			//"ReceiptCode": anybot.EncodeID(receiptID),
			"Title":       fmt.Sprintf("Receipt @ DebtsTracker.io #%v", receiptID),
			"Description": "Receipt description goes here",
		}, 9)
	} else {
		redirectToWebApp(w, r, false, common4all.Deeplink.AppHashPathToReceipt(receiptID), map[string]string{}, []string{})
	}
}
