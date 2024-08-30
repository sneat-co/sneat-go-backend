package shared_all

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"testing"
)

func TestAddSharedRoutes(t *testing.T) {
	router := botsfw.NewWebhookRouter(map[botinput.WebhookInputType][]botsfw.Command{}, nil)
	AddSharedRoutes(router, BotParams{})
}
