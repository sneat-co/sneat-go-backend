package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"strconv"
	"strings"

	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_general"
	dtb_common "github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/dtb_common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

const NEW_COUNTERPARTY_COMMAND = "new-counterparty"

func NewCounterpartyCommand(nextCommand botsfw.Command) botsfw.Command {
	return botsfw.Command{
		Code:    NEW_COUNTERPARTY_COMMAND,
		Title:   trans.COMMAND_TEXT_NEW_COUNTERPARTY,
		Replies: []botsfw.Command{nextCommand},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {

			chatEntity := whc.ChatData()
			if chatEntity.IsAwaitingReplyTo(NEW_COUNTERPARTY_COMMAND) {
				var user models.AppUser
				if user, err = dtb_common.GetUser(whc); err != nil {
					return
				}

				input := whc.Input()
				input.LogRequest()

				var contact models.ContactEntry

				var (
					contactDetails  models.ContactDetails
					existingContact bool
				)

				switch input2 := input.(type) {
				case botsfw.WebhookTextMessage:
					mt := strings.TrimSpace(input2.Text())
					if mt == "." {
						return dtb_general.MainMenuAction(whc, "", false)
					}
					if mt == "" {
						return m, errors.New("failed to get userContactJson details: mt is empty && inputMessage == nil")
					}
					if _, err = strconv.ParseFloat(mt, 64); err == nil {
						// User entered a number
						return whc.NewMessageByCode(trans.MESSAGE_TEXT_CONTACT_NAME_IS_NUMBER), nil
					}
					contactDetails = models.ContactDetails{
						Username: mt,
					}
				case botsfw.WebhookContactMessage:
					if input == nil {
						return m, errors.New("failed to get WebhookContactMessage: contactMessage == nil")
					}

					contactDetails = models.ContactDetails{
						FirstName: input2.FirstName(),
						LastName:  input2.LastName(),
						//Username: username,
					}
					phoneStr := input2.PhoneNumber()
					if phoneNum, err := strconv.ParseInt(phoneStr, 10, 64); err != nil {
						log.Warningf(whc.Context(), "Failed to parse phone string to int (%v)", phoneStr)
					} else {
						contactDetails.PhoneContact = models.PhoneContact{
							PhoneNumber:          phoneNum,
							PhoneNumberConfirmed: true,
						}
					}

					switch input.InputType() {
					case botsfw.WebhookInputContact:
						contactDetails.TelegramUserID = input2.UserID().(int64) // TODO: check we are on Telegram
						if contactDetails.TelegramUserID != 0 {
							for _, userContactJson := range user.Data.Contacts() {
								if userContactJson.TgUserID == contactDetails.TelegramUserID {
									log.Debugf(whc.Context(), "Matched contact my TelegramUserID=%d", contactDetails.TelegramUserID)
									existingContact = true
									contact.ID = userContactJson.ID
								}
							}
						}
					}
				default:
					err = fmt.Errorf("unknown input, expected text or contact message, got: %T", input)
					return
				}

				if !existingContact {
					var user models.AppUser
					if user, err = facade.User.GetUserByID(whc.Context(), nil, whc.AppUserID()); err != nil {
						return
					}

					contactFullName := contactDetails.FullName()

					for _, userContact := range user.Data.Contacts() {
						if userContact.Name == contactFullName {
							m.Text = whc.Translate(trans.MESSAGE_TEXT_ALREADY_HAS_CONTACT_WITH_SUCH_NAME)
							return
						}
					}
				}

				if !existingContact {
					if contact, user, err = facade.CreateContact(whc.Context(), nil, whc.AppUserID(), contactDetails); err != nil {
						return m, err
					}
					ga := whc.GA()
					if err = ga.Queue(ga.GaEventWithLabel(
						"contacts",
						"contact-created",
						fmt.Sprintf("user-%v", whc.AppUserID()),
					)); err != nil {
						return m, err
					}
					if contact.Data.PhoneNumber != 0 && contact.Data.PhoneNumberConfirmed {
						if err = ga.Queue(ga.GaEventWithLabel(
							"contacts",
							"contact-details-added",
							"phone-number",
						)); err != nil {
							return m, err
						}
					}
				}
				if contact.ID == "" {
					panic("contact.ID == 0")
				}
				chatEntity.AddWizardParam(WIZARD_PARAM_COUNTERPARTY, contact.ID)
				return nextCommand.Action(whc)
				//m = whc.NewMessageByCode(fmt.Sprintf("ContactEntry Created: %v", counterpartyKey))
			} else {
				m = whc.NewMessageByCode(trans.MESSAGE_TEXT_ASK_NEW_COUNTERPARTY_NAME)
				m.Format = botsfw.MessageFormatHTML
				m.Keyboard = tgbotapi.NewHideKeyboard(true)
				chatEntity.PushStepToAwaitingReplyTo(NEW_COUNTERPARTY_COMMAND)
			}
			return m, err
		},
	}
}
