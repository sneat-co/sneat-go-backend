package cmds4sneatbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/core4spaceus"
	"net/url"
)

var spacesCommand = botsfw.Command{
	Code:     "spaces",
	Commands: []string{"/spaces"},
	InputTypes: []botinput.WebhookInputType{
		botinput.WebhookInputText,
		botinput.WebhookInputCallbackQuery,
	},
	Action:         spacesAction,
	CallbackAction: spacesCallbackAction,
}

func spacesAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	return spacesCallbackAction(whc, nil)
}

func spacesCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	var spaceRef core4spaceus.SpaceRef

	if callbackUrl != nil {
		spaceRef = core4spaceus.SpaceRef(callbackUrl.Query().Get("s"))
		if m, err = whc.NewEditMessage("", botsfw.MessageFormatHTML); err != nil {
			return
		}
	}
	if spaceRef == "" {
		spaceRef = core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, "")
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "👪 Family",
				CallbackData: tghelpers.GetSpaceCallbackData(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, "")),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🔒 Private",
				CallbackData: tghelpers.GetSpaceCallbackData(core4spaceus.NewSpaceRef(core4spaceus.SpaceTypePrivate, "")),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "➕ Add new space (not implemented yet)",
				CallbackData: "add-space",
			},
		},
	)
	var currentSpaceTitle string
	var currentSpaceEmoji string

	switch spaceRef.SpaceType() {
	case core4spaceus.SpaceTypeFamily:
		currentSpaceTitle = "Family"
		currentSpaceEmoji = "👪"
		keyboard.InlineKeyboard[0][0].Text += " ✅"
		keyboard.InlineKeyboard[1][0].Text += " 🔲"
	case core4spaceus.SpaceTypePrivate:
		currentSpaceTitle = "Private"
		currentSpaceEmoji = "🔒"
		keyboard.InlineKeyboard[0][0].Text += " 🔲"
		keyboard.InlineKeyboard[1][0].Text += " ✅"
	}
	m.Format = botsfw.MessageFormatHTML

	m.Text = "<b>Your spaces</b>"
	m.Text += fmt.Sprintf("\nCurrent space: %s <b>%s</b>", currentSpaceEmoji, currentSpaceTitle)
	m.Text += "\nClick to switch current space."

	m.Keyboard = keyboard
	return
}
