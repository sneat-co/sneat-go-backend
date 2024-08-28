package tghelpers

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"net/url"
)

func GetSpaceParams(callbackUrl *url.URL) (spaceType core4spaceus.SpaceType, spaceID string) {
	q := callbackUrl.Query()
	spaceType = core4spaceus.SpaceType(q.Get("spaceType"))
	spaceID = q.Get("spaceID")
	return spaceType, spaceID
}

func BackToSpaceMenuButton(callbackUrl *url.URL) tgbotapi.InlineKeyboardButton {
	spaceType, spaceID := GetSpaceParams(callbackUrl)
	return tgbotapi.InlineKeyboardButton{
		Text:         "⬅️ Back to space menu",
		CallbackData: fmt.Sprintf("space?spaceType=%s&spaceID=%s", spaceType, spaceID),
	}
}
