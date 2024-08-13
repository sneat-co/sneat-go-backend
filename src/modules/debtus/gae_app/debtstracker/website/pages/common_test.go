package pages

import (
	"github.com/strongo/i18n"
	"testing"
)

func TestRenderCachedPageWithoutArguemnts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	RenderCachedPage(nil, nil, nil, i18n.LocaleEnUS, nil, 0)
}
