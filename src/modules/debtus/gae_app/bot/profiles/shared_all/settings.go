package shared_all

import (
	"errors"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
)

const SettingsCommandCode = "settings"

var SettingsCommandTemplate = botsfw.Command{
	Code:     SettingsCommandCode,
	Commands: trans.Commands(trans.COMMAND_TEXT_SETTING, trans.COMMAND_SETTINGS, emoji.SETTINGS_ICON),
	Icon:     emoji.SETTINGS_ICON,
}

func SettingsMainAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	switch whc.BotPlatform().ID() {
	case telegram.PlatformID:
		m, _, err = SettingsMainTelegram(whc)
	default:
		err = errors.New("Unsupported platform")
	}
	return
}

func SettingsMainTelegram(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, keyboard *tgbotapi.InlineKeyboardMarkup, err error) {
	m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_SETTINGS))
	m.IsEdit = whc.InputType() == botsfw.WebhookInputCallbackQuery
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.CommandText(trans.COMMAND_TEXT_LANGUAGE, emoji.EARTH_ICON),
				CallbackData: SettingsLocaleListCallbackPath,
			},
		},
	)
	m.Keyboard = keyboard
	return
}
