package tghelpers

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"net/url"
)

func GetSpaceRef(callbackUrl *url.URL) (spaceRef core4spaceus.SpaceRef) {
	q := callbackUrl.Query()
	return core4spaceus.SpaceRef(q.Get("s"))
}

func BackToSpaceMenuButton(spaceRef core4spaceus.SpaceRef) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text:         "⬅️ Back to space",
		CallbackData: GetSpaceCallbackData(spaceRef),
	}
}

func GetSpaceCallbackData(spaceRef core4spaceus.SpaceRef) string {
	return fmt.Sprintf("space?s=%s", spaceRef)
}
