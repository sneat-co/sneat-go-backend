package splitus

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
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
		var user models.AppUser
		if user, err = facade.User.GetUserByID(c, nil, whc.AppUserID()); err != nil {
			return
		}
		if user.Data.OutstandingBillsCount == 0 {
			m.Text = whc.Translate("You have no outstanding bills.")
			return
		}

		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "<b>%v</b>\n", whc.Translate("Outstanding bills"))
		for i, billJson := range user.Data.GetOutstandingBills() {
			fmt.Fprintf(buf, "\n%v. %v: %v %v", i+1, billJson.Name, billJson.Total, billJson.Currency)
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
				shared_group.NewGroupTelegramInlineButton(whc, 0),
			},
		)
		m.Keyboard = keyboard
		return
	}

	var group models.GroupEntry
	if group, err = shared_group.GetGroup(whc, nil); err != nil {
		return
	}

	m.Format = botsfw.MessageFormatHTML

	if group.Data.OutstandingBillsCount == 0 {
		mt := "This group has no outstanding bills"
		switch whc.InputType() {
		case botsfw.WebhookInputCallbackQuery:
			m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{Text: mt})
		case botsfw.WebhookInputText:
			m.Text = mt
		default:
			err = errors.New("Unknown input type")
		}
		return
	}

	buf := new(bytes.Buffer)
	buf.WriteString("<b>Outstanding bills</b>\n\n")

	outstandingBills := group.Data.GetOutstandingBills()

	for i, bill := range outstandingBills {
		fmt.Fprintf(buf, `  %d. <a href="https://t.me/%v?start=bill-%v">%v</a>`+"\n", i+1, whc.GetBotCode(), bill.ID, bill.Name)
	}

	fmt.Fprintf(buf, "\nSend /split@%v to close the bills.\nThe debts records will be available in @DebtsTrackerBot.", whc.GetBotCode())

	m.Text = buf.String()
	return
}
