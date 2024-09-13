package cmds4anybot

import (
	"errors"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
)

const SettingsCommandCode = "settings"

var SettingsCommandTemplate = botsfw.Command{
	Code: SettingsCommandCode,
	Commands: trans.Commands(
		trans.COMMAND_TEXT_SETTING,
		trans.COMMAND_SETTINGS,
		emoji.SETTINGS_ICON,
	),
	Icon: emoji.SETTINGS_ICON,
}

func SettingsMainAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	switch botPlatformID := whc.BotPlatform().ID(); botPlatformID {
	case "":
		panic("whc.BotPlatform().ID() is empty string")
	case telegram.PlatformID:
		m, _, err = SettingsMainTelegram(whc)
	default:
		err = errors.New("unsupported platform: " + botPlatformID)
	}
	return
}

func SettingsMainTelegram(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, keyboard *tgbotapi.InlineKeyboardMarkup, err error) {
	m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_SETTINGS))
	m.IsEdit = whc.Input().InputType() == botinput.WebhookInputCallbackQuery
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
