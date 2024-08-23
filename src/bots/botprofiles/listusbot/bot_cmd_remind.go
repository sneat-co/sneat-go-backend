package listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/olebedev/when"
	"strings"
	"time"
)

var remindCommand = botsfw.Command{
	Code: "remind",
	Matcher: func(command botsfw.Command, context botsfw.WebhookContext) bool {
		switch input := context.Input().(type) {
		case botsfw.WebhookTextMessage:
			return strings.HasPrefix(input.Text(), "remind ")
		}
		return false
	},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		text := whc.Input().(botsfw.WebhookTextMessage).Text()
		now := time.Now()
		r, err := when.EN.Parse(text, now)
		if err != nil {
			return m, err
		}
		if r == nil {
			return m, nil
		}
		m = whc.NewMessage(fmt.Sprintf("I will remind you at %v: %s: %s", r.Time, r.Text, r.Source))
		return m, err
	},
}
