package tghelpers

import "github.com/bots-go-framework/bots-api-telegram/tgbotapi"

func BackToSpaceMenuButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text:         "⬅️ Back to space menu",
		CallbackData: "space",
	}
}
