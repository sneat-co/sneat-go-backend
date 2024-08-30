package cmds4anybot

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"strconv"
)

var counterCommand = botsfw.Command{
	Code:       "count",
	Commands:   []string{"/count"},
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputText},
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
