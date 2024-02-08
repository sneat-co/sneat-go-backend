package dtb_transfer

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"net/url"

	"errors"
)

const CHANGE_RECEIPT_LANG_COMMAND = "change-lang-receipt"

var ChangeReceiptAnnouncementLangCommand = botsfw.NewCallbackCommand(
	CHANGE_RECEIPT_LANG_COMMAND,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		query := callbackUrl.Query()
		code5 := query.Get("locale")
		if len(code5) != 5 {
			return m, errors.New("ChangeReceiptAnnouncementLangCommand: len(code5) != 5")
		}
		if err = whc.SetLocale(code5); err != nil {
			return m, err
		}
		receiptID := query.Get("id")
		if err != nil {
			return m, err
		}
		return showReceiptAnnouncement(whc, receiptID, "")
	},
)
