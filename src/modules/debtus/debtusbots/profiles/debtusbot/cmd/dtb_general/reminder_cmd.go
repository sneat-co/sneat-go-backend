package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
)

func EditReminderMessage(whc botsfw.WebhookContext, transfer models4debtus.TransferEntry, message string) (m botsfw.MessageFromBot, err error) {
	utm := common4all.NewUtmParams(whc, common4all.UTM_CAMPAIGN_REMINDER)
	appUserID := whc.AppUserID()
	mt := fmt.Sprintf(
		"<b>%v</b>\n%v\n\n%v",
		whc.Translate(trans.MESSAGE_TEXT_REMINDER),
		common4debtus.TextReceiptForTransfer(whc.Context(), whc, transfer, appUserID, common4debtus.ShowReceiptToAutodetect, utm),
		message,
	)
	if whc.Input().InputType() == botinput.WebhookInputCallbackQuery {
		if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
			return
		}
	} else {
		m = whc.NewMessage(mt)
		m.Format = botsfw.MessageFormatHTML
		SetMainMenuKeyboard(whc, &m)
	}

	return
}
