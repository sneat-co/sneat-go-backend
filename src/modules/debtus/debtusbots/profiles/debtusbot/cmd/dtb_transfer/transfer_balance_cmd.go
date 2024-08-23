package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"net/url"
	"strings"
	"time"

	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/emoji"
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
	ctx := whc.Context()

	logus.Debugf(ctx, "BalanceCommand.Action()")

	var buffer bytes.Buffer

	appUserID := whc.AppUserID()

	if appUserID == "" {
		if _, err = buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_BALANCE_IS_ZERO)); err != nil {
			return
		}
	} else {
		user := dbo4userus.NewUserEntry(appUserID)

		if err = dal4userus.GetUser(ctx, nil, user); err != nil {
			return
		}

		spaceID := user.Data.GetFamilySpaceID()

		debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)
		if err = models4debtus.GetDebtusSpace(ctx, nil, debtusSpace); err != nil {
			return
		}

		contactusSpace := dal4contactus.NewContactusSpaceEntry(spaceID)
		if err = dal4contactus.GetContactusSpace(ctx, nil, contactusSpace); err != nil {
			return
		}

		if len(debtusSpace.Data.Balance) == 0 {
			if _, err = buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_BALANCE_IS_ZERO)); err != nil {
				return
			}
		} else {
			balanceMessageBuilder := NewBalanceMessageBuilder(whc)
			if len(debtusSpace.Data.Contacts) == 0 {
				return m, fmt.Errorf("integrity issue: UserEntry{ContactID=%s} has non zero balance and no contacts", whc.AppUserID())
			}
			buffer.WriteString(fmt.Sprintf("<b>%v</b>", whc.Translate(trans.MESSAGE_TEXT_BALANCE_HEADER)) + common4debtus.HORIZONTAL_LINE)
			linker := common4debtus.NewLinkerFromWhc(whc)
			buffer.WriteString(balanceMessageBuilder.ByContact(ctx, linker, contactusSpace.Data.Contacts, debtusSpace.Data.Contacts))

			var thereAreFewDebtsForSingleCurrency = func() bool {
				//TODO: Duplicate call to Balance() - consider move inside BalanceMessageBuilder
				//logus.Debugf(ctx, "thereAreFewDebtsForSingleCurrency()")
				var currencies []money.CurrencyCode
				for _, counterparty := range debtusSpace.Data.Contacts {
					//logus.Debugf(ctx, "counterparty: %v", counterparty)
					for currency := range counterparty.Balance {
						//logus.Debugf(ctx, "currency: %v", currency)
						for _, curr := range currencies {
							//logus.Debugf(ctx, "curr: %v; curr == currency: %v", curr, curr == currency)
							if curr == currency {
								return true
							}
						}
						currencies = append(currencies, currency)
					}
				}
				//logus.Debugf(ctx, "thereAreFewDebtsForSingleCurrency: %v", currencies)
				return false
			}

			if len(debtusSpace.Data.Contacts) > 1 && thereAreFewDebtsForSingleCurrency() {
				userBalanceWithInterest, err := debtusSpace.Data.BalanceWithInterest(ctx, time.Now())
				if err != nil {
					m := fmt.Sprintf("Failed to get balance with interest for user %v: %v", user.ID, err)
					logus.Errorf(ctx, m)
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
			//}	}
			if debtusSpace.Data.HasDueTransfers {
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
		}
	}
	buffer.WriteString(common4debtus.HORIZONTAL_LINE)
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

	//err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS)
	return m, err
	//SetMainMenuKeyboard(whc, &m) - Bad idea! Need to cleanup AwaitingReplyTo
}
