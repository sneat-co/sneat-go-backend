package pages

import (
	"github.com/sneat-co/debtstracker-translations/trans"
	"google.golang.org/appengine/v2"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var adsPageTmpl *template.Template

func AdsPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	locale, err := getLocale(appengine.NewContext(r), w, r)
	if err != nil {
		return
	}
	translator, data := pageContext(r, locale)

	for _, key := range []string{
		trans.WS_ADS_TITLE,
		trans.WS_ADS_CONTENT,
	} {
		data[key] = template.HTML(translator.Translate(key))
	}

	if adsPageTmpl == nil {
		adsPageTmpl = template.Must(template.ParseFiles(
			BASE_TEMPLATE,
			TEMPLATES_PATH+"ads.html",
			TEMPLATES_PATH+"device-switcher.html",
			TEMPLATES_PATH+"device.js.html",
		))
	}
	RenderCachedPage(w, r, adsPageTmpl, locale, data, 0)
}
