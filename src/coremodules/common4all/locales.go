package common4all

import (
	"fmt"
	"github.com/strongo/i18n"
	"strings"
)

func Locale2to5(locale2 string) string {
	if len(locale2) != 2 {
		panic("len(locale2) != 2")
	}
	if strings.ToLower(locale2) == "en" {
		return i18n.LocaleCodeEnUS
	} else {
		return fmt.Sprintf("%v-%v", strings.ToLower(locale2), strings.ToUpper(locale2))
	}
}
