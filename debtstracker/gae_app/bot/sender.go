package bot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/log"
)

func SendRefreshOrNothingChanged(whc botsfw.WebhookContext, m botsfw.MessageFromBot) (m2 botsfw.MessageFromBot, err error) {
	c := whc.Context()
	if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		log.Debugf(c, "error type: %T", err)
		if apiResponse, ok := err.(tgbotapi.APIResponse); ok && apiResponse.ErrorCode == 400 {
			m.BotMessage = telegram.CallbackAnswer(tgbotapi.NewCallback("", whc.Translate(trans.ALERT_TEXT_NOTHING_CHANGED)))
			err = nil
		}
	}
	return m, err
}
