package pages

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/debtusbotconst"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/platforms/debtustgbots"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"google.golang.org/appengine/v2"
	"html/template"
	"net/http"
	"strings"

	"context"
)

func pageContext(r *http.Request, locale i18n.Locale) (translator i18n.SingleLocaleTranslator, data map[string]interface{}) {
	userVoiceID := "6ed87444-76e3-43ee-8b6e-fd28d345e79c" // English
	c := appengine.NewContext(r)

	switch locale.Code5 {
	case i18n.LocalCodeRuRu:
		userVoiceID = "47c67b85-d064-4727-b149-bda58cfe6c2d"
	}

	//appTranslator := anybot.TheAppContext.GetTranslator(c)
	var appTranslator i18n.Translator
	translator = i18n.NewSingleMapTranslator(locale, appTranslator)

	if locale.Code5 != i18n.LocaleCodeEnUS {
		translator = i18n.NewSingleLocaleTranslatorWithBackup(translator, i18n.NewSingleMapTranslator(i18n.LocaleEnUS, appTranslator))
	}

	env := dtdal.HttpAppHost.GetEnvironment(c, r)
	if env == strongoapp.UnknownEnv {
		panic("Unknown host: " + r.Host)
	}
	botSettings, err := debtustgbots.GetBotSettingsByLang(env, debtusbotconst.DebtusBotProfileID, locale.Code5)
	if err != nil {
		panic(err)
	}

	data = map[string]interface{}{
		"lang":          locale.SiteCode(),
		"userVoiceID":   userVoiceID,
		"TgBotID":       botSettings.Code,
		"SubLocalePath": strings.Replace(r.URL.EscapedPath(), fmt.Sprintf("/%v/", locale.SiteCode()), "/", 1),
		trans.WS_ALEX_T: translator.TranslateNoWarning(trans.WS_ALEX_T),
		trans.WS_MOTTO:  translator.Translate(trans.WS_MOTTO),
	}
	return translator, data
}

func getLocale(ctx context.Context, w http.ResponseWriter, r *http.Request) (locale i18n.Locale, err error) {
	getLocaleBySiteCode := func(localeCode string) {
		for _, supportedLocale := range i18n.LocalesByCode5 {
			if supportedLocale.SiteCode() == localeCode {
				locale = supportedLocale
				break
			}
		}
	}

	path := r.URL.Path
	if path == "/" {
		if localeCode, ok := ctx.Value("locale").(string); !ok {
			locale = i18n.LocaleEnUS
		} else {
			getLocaleBySiteCode(localeCode)
			if locale.Code5 == "" {
				locale = i18n.LocaleEnUS
			}
		}
		return
	} else {
		if strings.HasPrefix(path, "/ru/") {
			locale = i18n.LocaleRuRu
		} else if strings.HasPrefix(path, "/zh/") {
			locale = i18n.LocaleZhCn
		} else if strings.HasPrefix(path, "/ja/") {
			locale = i18n.LocaleJaJp
		} else if strings.HasPrefix(path, "/fa/") {
			locale = i18n.LocaleFaIr
		} else {
			nextSlashIndex := strings.Index(path[1:], "/")
			if nextSlashIndex == -1 {
				err = fmt.Errorf("Unsupported path: %v", path)
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "text/plain")
				_, _ = w.Write(([]byte)(err.Error()))
				return
			} else {
				localeCode := path[1 : nextSlashIndex+1]
				getLocaleBySiteCode(localeCode)
				if locale.Code5 == "" {
					w.WriteHeader(http.StatusNotFound)
					w.Header().Set("Content-Type", "text/plain")
					if _, err := w.Write(([]byte)(fmt.Sprintf("Unsupported locale: %v", localeCode))); err != nil {
						logus.Errorf(ctx, err.Error())
					}
					return
				}
			}
		}
	}
	return
}

func RenderCachedPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template, locale i18n.Locale, data map[string]interface{}, maxAge int) {
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(buffer.Bytes())
		_, _ = w.Write([]byte("<hr><div style=color:red;position:absolute;padding:10px;background-color:white>" + err.Error() + "</div>"))
		return
	}
	eTag := fmt.Sprintf("%x", md5.Sum(buffer.Bytes()))
	if match := r.Header.Get("If-None-Match"); match == eTag {
		w.WriteHeader(http.StatusNotModified)
	} else {
		header := w.Header()
		header.Set("Content-Language", locale.Code5)
		if maxAge >= 0 {
			if maxAge == 0 {
				maxAge = 600
			}
			header.Set("Cache-Control", fmt.Sprintf("public, max-age=%v", maxAge))
		}
		header.Set("ETag", eTag)
		_, _ = w.Write(buffer.Bytes())
	}
}
