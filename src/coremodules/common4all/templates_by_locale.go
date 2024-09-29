package common4all

import (
	"bytes"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	html "html/template"
	"sync"
	text "text/template"

	"context"
)

type TextTemplateProvider struct {
	mutex    sync.Mutex
	compiled map[string]*text.Template
}

type HtmlTemplateProvider struct {
	mutex    sync.Mutex
	compiled map[string]*html.Template
}

func NewTextTemplates() TextTemplateProvider {
	return TextTemplateProvider{compiled: map[string]*text.Template{}}
}

func NewHtmlTemplates() HtmlTemplateProvider {
	return HtmlTemplateProvider{compiled: map[string]*html.Template{}}
}

var TextTemplates = NewTextTemplates()
var HtmlTemplates = NewHtmlTemplates()

func (templates *TextTemplateProvider) RenderTemplate(ctx context.Context, translator i18n.SingleLocaleTranslator, templateName string, params interface{}) (string, error) {
	var t *text.Template
	var err error
	var ok bool
	if t, ok = templates.compiled[templateName]; !ok {
		templateContent := translator.Translate(templateName)
		t, err = text.New(templateName).Parse(templateContent)
		if err != nil {
			return "", err
		}
		defer templates.mutex.Unlock()
		templates.mutex.Lock()
		templates.compiled[templateName] = t
		logus.Infof(ctx, "Compiled & cached template [%v], total cached: %v", templateName, len(templates.compiled))
	}

	var wr bytes.Buffer

	if err := t.Execute(&wr, params); err != nil {
		return "", err
	}
	return wr.String(), nil
}

func (templates *HtmlTemplateProvider) RenderTemplate(ctx context.Context, wr *bytes.Buffer, translator i18n.SingleLocaleTranslator, templateName string, params interface{}) (err error) {
	var t *html.Template
	var ok bool
	cacheCode := templateName + ":" + translator.Locale().Code5
	if t, ok = templates.compiled[cacheCode]; !ok {
		templateContent := translator.Translate(templateName)
		if t, err = html.New(templateName).Parse(templateContent); err != nil {
			return
		}
		defer templates.mutex.Unlock()
		templates.mutex.Lock()
		templates.compiled[cacheCode] = t

		logus.Infof(ctx, "Compiled & cached template [%v], total cached: %v", cacheCode, len(templates.compiled))
	}

	if err = t.Execute(wr, params); err != nil {
		return
	}
	return
}
