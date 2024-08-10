package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/logus"
	"net/url"
	"strings"
	"time"

	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

const BALANCE_COMMAND = "balance"

var BalanceCallbackCommand = botsfw.NewCallbackCommand(BALANCE_COMMAND, balanceCallbackAction)

var BalanceCommand = botsfw.Command{ //TODO: Write unit tests!
	Code:     BALANCE_COMMAND,
	Title:    trans.COMMAND_TEXT_BALANCE,
	Icon:     emoji.BALANCE_ICON,
	Commands: trans.Commands(trans.COMMAND_BALANCE),
	Action:   balanceAction,
}

func balanceCallbackAction(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
	return balanceAction(whc)
}

func balanceAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	logus.Debugf(c, "BalanceCommand.Action()")

	var user models.AppUser

	if user, err = facade2debtus.User.GetUserByID(c, nil, whc.AppUserID()); err != nil {
		return
	}

	var buffer bytes.Buffer
	if user.Data.BalanceCount == 0 {
		if _, err = buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_BALANCE_IS_ZERO)); err != nil {
			return
		}
	} else {
		balanceMessageBuilder := NewBalanceMessageBuilder(whc)
		contacts := user.Data.Contacts()
		if len(contacts) == 0 {
			return m, fmt.Errorf("Integrity issue: User{ID=%v} has non zero balance and no contacts.", whc.AppUserID())
		}
		buffer.WriteString(fmt.Sprintf("<b>%v</b>", whc.Translate(trans.MESSAGE_TEXT_BALANCE_HEADER)) + common.HORIZONTAL_LINE)
		linker := common.NewLinkerFromWhc(whc)
		buffer.WriteString(balanceMessageBuilder.ByContact(c, linker, contacts))

		var thereAreFewDebtsForSingleCurrency = func() bool {
			//TODO: Duplicate call to Balance() - consider move inside BalanceMessageBuilder
			//logus.Debugf(c, "thereAreFewDebtsForSingleCurrency()")
			var currencies []money.CurrencyCode
			for _, counterparty := range contacts {
				//logus.Debugf(c, "counterparty: %v", counterparty)
				for currency := range counterparty.Balance() {
					//logus.Debugf(c, "currency: %v", currency)
					for _, curr := range currencies {
						//logus.Debugf(c, "curr: %v; curr == currency: %v", curr, curr == currency)
						if curr == currency {
							return true
						}
					}
					currencies = append(currencies, currency)
				}
			}
			//logus.Debugf(c, "thereAreFewDebtsForSingleCurrency: %v", currencies)
			return false
		}

		if len(contacts) > 1 && thereAreFewDebtsForSingleCurrency() {
			userBalanceWithInterest, err := user.Data.BalanceWithInterest(c, time.Now())
			if err != nil {
				m := fmt.Sprintf("Failed to get balance with interest for user %v: %v", user.ID, err)
				logus.Errorf(c, m)
				buffer.WriteString(m)
			} else {
				buffer.WriteString("\n" + strings.Repeat("─", 16) + "\n" + balanceMessageBuilder.ByCurrency(true, userBalanceWithInterest))
			}
		}

		//if len(contacts) > 0 {
		//	//for i, counterparty := range contacts {
		//	//	telegramKeyboard = append(telegramKeyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(counterparty.GetFullName(), fmt.Sprintf("transfer-history?counterparty=%v", counterpartyKeys[i].IntID()))})
		//	//}
		//	telegramKeyboard = append(telegramKeyboard, []tgbotapi.InlineKeyboardButton{
		//		tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("balance?counterparty=%v", counterpartyKeys[len(counterpartyKeys)-1].IntID())),
		//		tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("balance?counterparty=%v", counterpartyKeys[0].IntID())),
		//	})
		//}
	}
	buffer.WriteString(common.HORIZONTAL_LINE)
	//buffer.WriteString(dtb_general.AdSlot(whc, "balance"))
	const THUMB_UP = "👍"
	buffer.WriteString(THUMB_UP + " " + whc.Translate(trans.MESSAGE_TEXT_PLEASE_HELP_MAKE_IT_BETTER))
	if whc.InputType() == botsfw.WebhookInputCallbackQuery {
		if m, err = whc.NewEditMessage(buffer.String(), botsfw.MessageFormatHTML); err != nil {
			return
		}
	} else {
		m = whc.NewMessage(buffer.String())
		m.Format = botsfw.MessageFormatHTML
	}

	m.DisableWebPagePreview = true

	if user.Data.HasDueTransfers {
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.COMMAND_TEXT_DUE_RETURNS),
					CallbackData: DUE_RETURNS_COMMAND,
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.COMMAND_TEXT_INVITE_FIREND),
					CallbackData: "invite",
				},
			},
		)
	}

	//err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS)
	return m, err
	//SetMainMenuKeyboard(whc, &m) - Bad idea! Need to cleanup AwaitingReplyTo
}
