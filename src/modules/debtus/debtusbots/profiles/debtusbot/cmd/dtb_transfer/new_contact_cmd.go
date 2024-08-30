package dtb_transfer

import (
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"strconv"
	"strings"
)

const newCounterpartyCommandCode = "new-counterparty"

func NewCounterpartyCommand(nextCommand botsfw.Command) botsfw.Command {
	return botsfw.Command{
		Code:    newCounterpartyCommandCode,
		Title:   trans.COMMAND_TEXT_NEW_COUNTERPARTY,
		Replies: []botsfw.Command{nextCommand},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {

			ctx := whc.Context()
			userCtx := facade.NewUserContext(whc.AppUserID())

			chatEntity := whc.ChatData()
			spaceID := ""
			if chatEntity.IsAwaitingReplyTo(newCounterpartyCommandCode) {
				contactusSpace := dal4contactus.NewContactusSpaceEntry(spaceID)

				input := whc.Input()
				input.LogRequest()

				var debtusContact models4debtus.DebtusSpaceContactEntry

				var (
					contactDetails  dto4contactus.ContactDetails
					existingContact bool
				)

				switch input := input.(type) {
				case botinput.WebhookTextMessage:
					mt := strings.TrimSpace(input.Text())
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
					contactDetails = dto4contactus.ContactDetails{
						NameFields: person.NameFields{
							UserName: mt,
						},
					}
				case botinput.WebhookContactMessage:
					if input == nil {
						return m, errors.New("failed to get WebhookContactMessage: contactMessage == nil")
					}

					contactDetails = dto4contactus.ContactDetails{
						NameFields: person.NameFields{
							FirstName: input.GetFirstName(),
							LastName:  input.GetLastName(),
						},
						//Username: username,
					}
					phoneStr := input.GetPhoneNumber()
					if phoneNum, err := strconv.ParseInt(phoneStr, 10, 64); err != nil {
						logus.Warningf(ctx, "Failed to parse phone string to int (%v)", phoneStr)
					} else {
						contactDetails.PhoneContact = dto4contactus.PhoneContact{
							PhoneNumber:          phoneNum,
							PhoneNumberConfirmed: true,
						}
					}

					contactBotUserID := input.GetBotUserID()
					if contactBotUserID != "" {
						contactDetails.TelegramUserID, err = strconv.ParseInt(input.GetBotUserID(), 10, 64) // TODO: check we are on Telegram
						if err != nil {
							err = fmt.Errorf("failed to parse contactBotUserID: %w", err)
							return
						}
					}
					var telegramUserID string
					if contactDetails.TelegramUserID != 0 {
						telegramUserID = strconv.FormatInt(contactDetails.TelegramUserID, 10)
					}

					if contactDetails.TelegramUserID != 0 {
						for contactID, contactBrief := range contactusSpace.Data.Contacts {
							var tgAccount *appuser.AccountKey
							if tgAccount, err = contactBrief.AccountsOfUser.GetAccount(const4userus.TelegramAuthProvider, ""); err != nil {
								return
							}
							if tgAccount != nil && tgAccount.ID == telegramUserID {
								logus.Debugf(ctx, "Matched debtusContact my TelegramUserID=%d", contactDetails.TelegramUserID)
								existingContact = true
								debtusContact.ID = contactID
							}
						}
					}
				default:
					err = fmt.Errorf("unknown input, expected text or debtusContact message, got: %T", input)
					return
				}

				if !existingContact {
					if err = dal4contactus.GetContactusSpace(ctx, nil, contactusSpace); err != nil {
						return
					}

					contactFullName := contactDetails.FullName()

					for _, userContact := range contactusSpace.Data.Contacts {
						if userContact.Names.FullName == contactFullName {
							m.Text = whc.Translate(trans.MESSAGE_TEXT_ALREADY_HAS_CONTACT_WITH_SUCH_NAME)
							return
						}
					}
				}

				if !existingContact {
					if _, contactusSpace, _, err = facade4debtus.CreateContact(ctx, nil, userCtx.GetUserID(), spaceID, contactDetails); err != nil {
						return m, err
					}
					ga := whc.GA()
					if err = ga.Queue(ga.GaEventWithLabel(
						"contacts",
						"debtusContact-created",
						fmt.Sprintf("user-%v", whc.AppUserID()),
					)); err != nil {
						return m, err
					}
					if debtusContact.Data.PhoneNumber != 0 && debtusContact.Data.PhoneNumberConfirmed {
						if err = ga.Queue(ga.GaEventWithLabel(
							"contacts",
							"debtusContact-details-added",
							"phone-number",
						)); err != nil {
							return m, err
						}
					}
				}
				if debtusContact.ID == "" {
					panic("debtusContact.ContactID == 0")
				}
				chatEntity.AddWizardParam(WizardParamCounterparty, debtusContact.ID)
				return nextCommand.Action(whc)
				//m = whc.NewMessageByCode(fmt.Sprintf("DebtusSpaceContactEntry Created: %v", counterpartyKey))
			} else {
				m = whc.NewMessageByCode(trans.MESSAGE_TEXT_ASK_NEW_COUNTERPARTY_NAME)
				m.Format = botsfw.MessageFormatHTML
				m.Keyboard = tgbotapi.NewHideKeyboard(true)
				chatEntity.PushStepToAwaitingReplyTo(newCounterpartyCommandCode)
			}
			return m, err
		},
	}
}
