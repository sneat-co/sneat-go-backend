package pages

import (
	"github.com/julienschmidt/httprouter"
	"github.com/strongo/i18n"
	"html/template"
	"net/http"

	"github.com/sneat-co/debtstracker-translations/trans"
)

var iouADanceTmpl *template.Template

func AnnieIOUaDancePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if iouADanceTmpl == nil {
		iouADanceTmpl = template.Must(template.ParseFiles(
			BASE_TEMPLATE,
			TEMPLATES_PATH+"song-iou-a-dance.html",
			TEMPLATES_PATH+"device-switcher.html",
			TEMPLATES_PATH+"device.js.html",
		))
	}

	translator, data := pageContext(r, i18n.LocaleEnUS)
	for _, key := range []string{
		trans.WS_SHORT_DESC,
		trans.WS_LIVE_DEMO,
	} {
		data[key] = template.HTML(translator.Translate(key))
	}
	data["SubLocalePath"] = "/"
	RenderCachedPage(w, r, iouADanceTmpl, i18n.LocaleEnUS, data, 0)
}

var iouDappyTmpl *template.Template

func IOWDappyPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if iouDappyTmpl == nil {
		iouDappyTmpl = template.Must(template.ParseFiles(
			BASE_TEMPLATE,
			TEMPLATES_PATH+"song-iou-dappy.html",
			TEMPLATES_PATH+"device-switcher.html",
			TEMPLATES_PATH+"device.js.html",
		))
	}

	translator, data := pageContext(r, i18n.LocaleEnUS)
	data["SubLocalePath"] = "/"
	for _, key := range []string{
		trans.WS_SHORT_DESC,
		trans.WS_LIVE_DEMO,
	} {
		data[key] = template.HTML(translator.Translate(key))
	}
	RenderCachedPage(w, r, iouDappyTmpl, i18n.LocaleEnUS, data, 0)
}
