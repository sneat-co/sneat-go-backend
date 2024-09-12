package cmds4anybot

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/facade4anybot"
	"net/url"

	"github.com/sneat-co/debtstracker-translations/emoji"
)

const REFERRERS_COMMAND = "referrers"

var ReferrersCommand = botsfw.Command{
	Code:       REFERRERS_COMMAND,
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputText},
	Commands:   trans.Commands(trans.COMMAND_TEXT_REFERRERS, emoji.PUBLIC_LOUDSPEAKER),
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		if m, err = topReferrersMessageText(whc); err != nil {
			return
		}
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(whc.Translate(trans.COMMAND_TEXT_ADD_MY_TG_CHANNEL), ADD_REFERRER_COMMAND),
			),
		)
		return
	},
}

func topReferrersMessageText(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	var topTelegramReferrers []string
	if topTelegramReferrers, err = facade4anybot.Referer.TopTelegramReferrers(ctx, whc.GetBotCode(), 5); err != nil {
		return
	} else if len(topTelegramReferrers) == 0 {
		topTelegramReferrers = []string{"meduzalive", "varlamov"}
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "<b>%v</b>\n", whc.Translate(trans.MESSAGE_TEXT_REFERRERS_TITLE))
	buf.WriteString("\n")
	for i, channel := range topTelegramReferrers {
		fmt.Fprintf(buf, "  %v. @%v\n", i+1, channel)
	}
	m.Text = buf.String()
	m.Format = botsfw.MessageFormatHTML
	return
}

const ADD_REFERRER_COMMAND = "add-referrer"

var AddReferrerCommand = botsfw.Command{
	Code: ADD_REFERRER_COMMAND,
	CallbackAction: func(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
		if m, err = topReferrersMessageText(whc); err != nil {
			return
		}
		url := fmt.Sprintf("https://t.me/%v?start=refbytguser-YOUR_CHANNEL", whc.GetBotCode())
		botID := whc.GetBotCode()
		m.Text += "\n" + whc.Translate(trans.MESSAGE_TEXT_HOW_TO_ADD_TG_CHANNEL,
			url,
			url, botID,
			url, botID,
		)
		m.IsEdit = true
		return
	},
}
