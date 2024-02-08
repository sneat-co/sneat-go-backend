package dtb_transfer

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/log"
)

func reportReminderIsActed(whc botsfw.WebhookContext, action string) {
	ga := whc.GA()
	if err := ga.Queue(ga.GaEvent(
		"reminders",
		action,
	)); err != nil {
		log.Errorf(whc.Context(), err.Error())
		err = nil
	}
}
