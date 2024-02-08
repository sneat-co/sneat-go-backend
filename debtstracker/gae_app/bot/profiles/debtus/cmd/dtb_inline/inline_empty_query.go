package dtb_inline

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/log"
)

func InlineEmptyQuery(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	log.Debugf(whc.Context(), "InlineEmptyQuery()")
	inlineQuery := whc.Input().(botsfw.WebhookInlineQuery)
	m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
		InlineQueryID:     inlineQuery.GetInlineQueryID(),
		CacheTime:         60,
		SwitchPMText:      "Help: How to use this bot?",
		SwitchPMParameter: "help_inline",
	})
	return m, err
}
