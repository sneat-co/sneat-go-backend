package webhooks

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func InitWebhooks(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/webhooks/twilio/", TwilioWebhook)
}
