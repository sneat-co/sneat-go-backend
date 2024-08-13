package botcmds4splitus

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/shared_space"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/url"
)

const billsCommandCode = "bills"

var billsCommand = botsfw.Command{
	Code:     billsCommandCode,
	Commands: trans.Commands(trans.COMMAND_TEXT_BILLS, "/"+billsCommandCode),
	Icon:     emoji.CLIPBOARD_ICON,
	Title:    trans.COMMAND_TEXT_BILLS,
	Action:   billsAction,
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return billsAction(whc)
	},
}

func billsAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	if !whc.IsInGroup() {
		userID := whc.AppUserID()
		userDebtus := models4debtus.NewDebtusUserEntry(userID)

		var db dal.DB
		if db, err = facade.GetDatabase(c); err != nil {
			return
		}
		if err = db.Get(c, userDebtus.Record); err != nil {
			return
		}
		if len(userDebtus.Data.OutstandingBills) == 0 {
			m.Text = whc.Translate("You have no outstanding bills.")
			return
		}

		buf := new(bytes.Buffer)
		_, _ = fmt.Fprintf(buf, "<b>%v</b>\n", whc.Translate("Outstanding bills"))
		var i int
		for _, billJson := range userDebtus.Data.GetOutstandingBills() {
			i++
			_, _ = fmt.Fprintf(buf, "\n%v. %v: %v %v", i, billJson.Name, billJson.Total, billJson.Currency)
		}
		m.Text = buf.String()
		m.Format = botsfw.MessageFormatHTML
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		if !whc.IsInGroup() {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
				[]tgbotapi.InlineKeyboardButton{{
					Text:         whc.CommandText(trans.COMMAND_TEXT_SETTLE_BILLS, emoji.GREEN_CHECKBOX),
					CallbackData: settleBillsCommandCode,
				}},
			)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
					whc.CommandText(trans.COMMAND_TEXT_NEW_BILL, emoji.MEMO_ICON),
					"",
				),
			},
			[]tgbotapi.InlineKeyboardButton{
				shared_space.NewGroupTelegramInlineButton(whc, 0),
			},
		)
		m.Keyboard = keyboard
		return
	}
	m.Format = botsfw.MessageFormatHTML
	err = errors.New("not implemented yet")

	//var space dal4spaceus.SpaceEntry
	//if space, err = shared_space.GetSpaceEntryByCallbackUrl(whc, nil); err != nil {
	//	return
	//}

	//if space.Data.OutstandingBillsCount == 0 {
	//	mt := "This space has no outstanding bills"
	//	switch whc.InputType() {
	//	case botsfw.WebhookInputCallbackQuery:
	//		m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{Text: mt})
	//	case botsfw.WebhookInputText:
	//		m.Text = mt
	//	default:
	//		err = errors.New("unknown input type")
	//	}
	//	return
	//}
	//
	//buf := new(bytes.Buffer)
	//buf.WriteString("<b>Outstanding bills</b>\n\n")
	//
	//outstandingBills := space.Data.GetOutstandingBills()
	//
	//var i int
	//for billID, bill := range outstandingBills {
	//	i++
	//	_, _ = fmt.Fprintf(buf, `  %d. <a href="https://t.me/%v?start=bill-%v">%v</a>`+"\n", i, whc.GetBotCode(), billID, bill.Name)
	//}
	//
	//_, _ = fmt.Fprintf(buf, "\nSend /split@%v to close the bills.\nThe debts records will be available in @DebtsTrackerBot.", whc.GetBotCode())
	//
	//m.Text = buf.String()
	return
}
