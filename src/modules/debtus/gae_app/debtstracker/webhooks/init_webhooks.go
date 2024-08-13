package webhooks

import (
	"github.com/julienschmidt/httprouter"
)

func InitWebhooks(router *httprouter.Router) {
	router.HandlerFunc("POST", "/webhooks/twilio/", TwilioWebhook)
}
