package vkapp

import (
	"net/http"
)

type router interface {
	HandlerFunc(method, path string, handler http.HandlerFunc)
}

func InitVkIFrameApp(router router) {
	router.HandlerFunc("GET", "/apps/vk/iframe", IFrameHandler)
}

//var vkIFrameTemplate *template.Template

func IFrameHandler(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
	//if vkIFrameTemplate == nil {
	//	vkIFrameTemplate = template.Must(
	//		template.ParseFiles(
	//			pages.TEMPLATES_PATH+"vk-iframe.html",
	//			pages.TEMPLATES_PATH+"device-switcher.html",
	//			pages.TEMPLATES_PATH+"device.js.html",
	//		),
	//	)
	//}
	//query := r.URL.Query()
	//apiID := query.Get("api_id")
	//_, ok := vkbots.BotsBy.ByCode[apiID]
	//if !ok {
	//	w.WriteHeader(http.StatusBadRequest)
	//	_, _ = w.Write([]byte("Unknown app id"))
	//	return
	//}
	//
	//lang := "ru"
	//if query.Get("language") == "3" {
	//	lang = "en"
	//}
	//
	//data := map[string]interface{}{
	//	"vkApiId": apiID,
	//	"lang":    lang,
	//	"hash":    query.Get("hash"),
	//}
	//
	//pages.RenderCachedPage(w, r, vkIFrameTemplate, i18n.LocaleRuRu, data, 0)
}
