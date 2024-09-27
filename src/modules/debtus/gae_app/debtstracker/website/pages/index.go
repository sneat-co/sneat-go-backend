package pages

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/logus"
	"html/template"
	"net/http"
	"strings"
)

const TEMPLATES_PATH = "templates/" //"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/pages/templates/"
const BASE_TEMPLATE = TEMPLATES_PATH + "_base.html"

var countryToLocale = map[string]string{
	"RU": "ru",
	"IR": "fa",
	"IT": "it",
	"FR": "fr",
	"DE": "de",
	"PL": "pl",
	"PT": "pt",
	"KP": "ko",
	"KR": "ko",
	"JP": "jp",
	"CN": "zh",
	// TODO: Spanish speaking - add all from https://en.wikipedia.org/wiki/List_of_countries_where_Spanish_is_an_official_language
	"ES": "es",
	"MX": "es",
	"CO": "es",
	"AR": "es",
	"PE": "es",
	"VE": "es",
	"CL": "es",
	"EC": "es",
	"GT": "es",
	"CU": "es",
	"BO": "es",
	"DO": "es",
}

var supportedLocales = []string{"ru", "es", "it", "fr", "fa", "de", "pl", "pt", "ko", "jp", "zh"}

func IndexRootPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := r.Context()
	logus.Debugf(c, "IndexRootPage")

	if r.URL.Path != "/" { // This handler should work just for a root path.
		w.WriteHeader(http.StatusNotFound)
		return
	}

	acceptLanguages := r.Header.Get("Accept-Language")
	for _, acceptLanguage := range strings.Split(acceptLanguages, ";") {
		for _, al := range strings.Split(acceptLanguage, ",") {
			al = strings.TrimSpace(al)
			if len(al) == 2 {
				for _, l := range supportedLocales {
					if l == al {
						c = context.WithValue(c, ContextLocale, al)
						indexPage(c, w, r)
						//w.Header().Add("Location", fmt.Sprintf("/%v/", al))
						//w.WriteHeader(http.StatusTemporaryRedirect)
						return
					}
				}
			}
		}
	}

	{ // Try to detect language by country
		country := r.Header.Get("CF-IPCountry")
		if country != "" {
			logus.Debugf(c, "CF-IPCountry: %s", country)
		} else {
			country = r.Header.Get("X-AppEngine-Country")
			logus.Debugf(c, "X-AppEngine-Country: %s", country)
		}

		if country != "" {
			if localeCode, ok := countryToLocale[strings.ToUpper(country)]; ok {
				c = context.WithValue(c, ContextLocale, localeCode)
				indexPage(c, w, r)
				//w.Header().Add("Location", "/"+localeCode+"/")
				//w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
	}

	indexPage(c, w, r)
}

type contextLocale struct {
	id string
}

var ContextLocale = contextLocale{id: "locale"}

var indexTmpl *template.Template

func IndexPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !strings.HasSuffix(r.URL.Path, "/") {
		w.WriteHeader(http.StatusPermanentRedirect)
		path := r.URL.Path
		w.Header().Add("Location", strings.Replace(r.URL.RequestURI(), path, path+"/", 1))
		return
	}
	indexPage(r.Context(), w, r)
}

func indexPage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if referrer := r.Referer(); referrer != "" && !strings.Contains(referrer, "://"+r.Host) {
		logus.Debugf(ctx, "Referer: %v", referrer)
	}
	locale, err := getLocale(ctx, w, r)
	if err != nil {
		return
	}
	translator, data := pageContext(r, locale)
	for _, key := range []string{
		trans.WS_SHORT_DESC,
		trans.WS_INDEX_TITLE,
		trans.WS_INDEX_TG_BOT_H2,
		trans.WS_INDEX_TG_BOT_P,
		trans.WS_INDEX_TG_BOT_OPEN,
		trans.WS_LIVE_DEMO,
	} {
		data[key] = template.HTML(translator.Translate(key))
	}

	if indexTmpl == nil {
		indexTmpl = template.Must(template.ParseFiles(
			BASE_TEMPLATE,
			TEMPLATES_PATH+"index.html",
			TEMPLATES_PATH+"device-switcher.html",
			TEMPLATES_PATH+"device.js.html",
		))
	}
	RenderCachedPage(w, r, indexTmpl, locale, data, 0)
}
