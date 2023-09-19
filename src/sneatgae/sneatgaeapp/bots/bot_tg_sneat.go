package bots

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp/bots/sneatbots"
	"github.com/sneat-co/sneat-go-modules/listus/bot4listus"
	"strconv"
)

func errorFooterText() string {
	return "Please report any issues to @trakhimenok"
}

var SneatBotProfile botsfw.BotProfile

func init() {
	var textAndContactCommands = []botsfw.Command{
		startCommand,
		pingCommand,
		counterCommand,
	}
	textAndContactCommands = append(textAndContactCommands,
		bot4listus.Commands...,
	)
	commandsByType := map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText: textAndContactCommands,
	}
	var router = botsfw.NewWebhookRouter(commandsByType, errorFooterText)
	SneatBotProfile = sneatbots.NewProfile("sneat", &router)
}

var startCommand = botsfw.Command{
	Code:       "start",
	Commands:   []string{"/start"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Hello, stranger!"
		return
	},
}

var pingCommand = botsfw.Command{
	Code:       "ping",
	Commands:   []string{"/ping"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Pong!"
		return
	},
}

var counterCommand = botsfw.Command{
	Code:       "count",
	Commands:   []string{"/count"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		chatData := whc.ChatData()
		v := 0
		s := chatData.GetWizardParam("v")
		if s != "" {
			v, err = strconv.Atoi(s)
			if err != nil {
				return m, err
			}
		}
		v += 1
		m.Text = fmt.Sprintf("Counter: %d", v)
		s = strconv.Itoa(v)
		chatData.SetAwaitingReplyTo("count")
		chatData.AddWizardParam("v", s)
		return
	},
}
