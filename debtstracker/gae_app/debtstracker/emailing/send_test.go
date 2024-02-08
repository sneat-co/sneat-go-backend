package emailing

import (
	"context"
	"testing"
)

func TestGetEmailTextWithoutTranslator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should fail")
		}
	}()
	c := context.Background()
	_, _ = GetEmailText(c, nil, "some-template", nil)
}

func TestGetEmailHtmlWithoutTranslator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should fail")
		}
	}()
	c := context.Background()
	_, _ = GetEmailHtml(c, nil, "some-template", nil)
}
