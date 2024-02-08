package dtb_transfer

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/sms"
	//"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/invites"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"context"
	"errors"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/analytics"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/general"
	"github.com/strongo/log"
)

const ASK_PHONE_NUMBER_FOR_RECEIPT_COMMAND = "ask-phone-number-for-receipt"

func cleanPhoneNumber(phoneNumebr string) string {
	phoneNumebr = strings.Replace(phoneNumebr, " ", "", -1)
	phoneNumebr = strings.Replace(phoneNumebr, "(", "", -1)
	phoneNumebr = strings.Replace(phoneNumebr, ")", "", -1)
	return phoneNumebr
}

var AskPhoneNumberForReceiptCommand = botsfw.Command{
	Code: ASK_PHONE_NUMBER_FOR_RECEIPT_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		var db dal.DB
		if db, err = facade.GetDatabase(c); err != nil {
			return m, err
		}
		return m, db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
			log.Debugf(c, "AskPhoneNumberForReceiptCommand.Action()")

			input := whc.Input()

			var (
				mt             string
				phoneNumberStr string
				phoneNumber    int64
			)

			contact, isContactMessage := input.(botsfw.WebhookContactMessage)

			if isContactMessage {
				if contact == nil {
					m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVALID_PHONE_NUMBER)
					return nil
				}
				user, err := facade.User.GetUserByID(c, tx, whc.AppUserID())
				if err != nil {
					return err
				}
				if user.Data.FirstName == contact.FirstName() && user.Data.LastName == contact.LastName() {
					phoneNumberStr = cleanPhoneNumber(contact.PhoneNumber())
					if phoneNumber, err = strconv.ParseInt(phoneNumberStr, 10, 64); err != nil {
						log.Warningf(c, "Failed to parse contact's phone number: [%v]", phoneNumberStr)
						return err
					} else if user.Data.PhoneNumber == 0 {
						user, err := facade.User.GetUserByID(c, tx, whc.AppUserID())
						if err != nil {
							return err
						}
						if user.Data.PhoneNumber == 0 {
							user.Data.PhoneNumber = phoneNumber
							user.Data.PhoneNumberConfirmed = true
							if err = facade.User.SaveUser(c, tx, user); err != nil {
								log.Errorf(c, fmt.Errorf("failed to update user with phone number: %w", err).Error())
								return err
							}

						}
					}
					m = whc.NewMessage(trans.MESSAGE_TEXT_YOU_CAN_SEND_RECEIPT_TO_YOURSELF_BY_SMS)
					return nil
				}
				mt = contact.PhoneNumber()
			} else {
				mt = whc.Input().(botsfw.WebhookTextMessage).Text()
			}

			if twilioTestNumber, ok := common.TwilioTestNumbers[mt]; ok {
				log.Debugf(c, "Using predefined test number [%v]: %v", mt, twilioTestNumber)
				phoneNumberStr = twilioTestNumber
			} else {
				phoneNumberStr = cleanPhoneNumber(mt)
			}

			if phoneNumber, err = strconv.ParseInt(phoneNumberStr, 10, 64); err != nil {
				m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVALID_PHONE_NUMBER)
				return nil
			}

			chatEntity := whc.ChatData()

			awaitingUrl, err := url.Parse(chatEntity.GetAwaitingReplyTo())
			if err != nil {
				return fmt.Errorf("failed to parse chat state as URL: %w", err)
			}

			transferID := awaitingUrl.Query().Get(WIZARD_PARAM_TRANSFER)
			if transferID == "" {
				return fmt.Errorf("empty transferID")
			}
			transfer, err := facade.Transfers.GetTransferByID(c, tx, transferID)
			if err != nil {
				return fmt.Errorf("failed to get transfer by ID: %w", err)
			}
			counterparty, err := facade.GetContactByID(c, tx, transfer.Data.Counterparty().ContactID)
			if err != nil {
				return err
			}
			phoneContact := models.PhoneContact{PhoneNumber: phoneNumber, PhoneNumberConfirmed: false}

			m, err = sendReceiptBySms(whc, tx, phoneContact, transfer, counterparty)
			return err
		})

	},
}

const SMS_STATUS_MESSAGE_ID_PARAM_NAME = "SmsStatusMessageId"
const SMS_STATUS_MESSAGE_UPDATES_COUNT_PARAM_NAME = "SmsStatusUpdatesCount"

func sendReceiptBySms(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, phoneContact models.PhoneContact, transfer models.Transfer, counterparty models.Contact) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	if transfer.Data == nil {
		if transfer, err = facade.Transfers.GetTransferByID(c, tx, transfer.ID); err != nil {
			return m, err
		}
	}

	whc.ChatData() //TODO: Workaround to make whc.GetAppUser() working
	appUser, err := whc.AppUserData()
	if err != nil {
		return m, err
	}
	user := appUser.(interface{ FullName() string })
	//if err != nil {
	//	return
	//}

	var (
		smsText string
		//inviteCode string
	)

	receiptData := models.NewReceiptEntity(whc.AppUserID(), transfer.ID, transfer.Data.Counterparty().UserID, whc.Locale().Code5, "sms", strconv.FormatInt(phoneContact.PhoneNumber, 10), general.CreatedOn{
		CreatedOnPlatform: whc.BotPlatform().ID(),
		CreatedOnID:       whc.GetBotCode(),
	})
	var receipt models.Receipt
	if receipt, err = dtdal.Receipt.CreateReceipt(c, receiptData); err != nil {
		return m, err
	}

	receiptUrl := common.GetReceiptUrl(receipt.ID, common.GetWebsiteHost(receiptData.CreatedOnID))

	if counterparty.Data.CounterpartyUserID == "" {
		//related := fmt.Sprintf("%v=%v", models.TransferKind, transferID)
		//inviteKey, invite, err := invites.CreatePersonalInvite(whc, whc.AppUserID(), invites.InviteBySms, strconv.FormatInt(phoneContact.PhoneNumber, 10), whc.BotPlatform().ID(), whc.GetBotCode(), related)
		//if err != nil {
		//	log.Errorf(c, "Failed to create invite: %v", err)
		//	return m, err
		//}
		//inviteCode = inviteKey.StringID()
	} else {
		panic("Not implemented, need to call common.GetReceiptUrlForUser(...)")
	}

	// You've got $10 from Jack
	// You've given $10 to Jack

	switch transfer.Data.Direction() {
	case models.TransferDirectionUser2Counterparty:
		smsText = fmt.Sprintf(whc.Translate(trans.SMS_RECEIPT_YOU_GOT), transfer.Data.GetAmount(), user.FullName())
	case models.TransferDirectionCounterparty2User:
		smsText = fmt.Sprintf(whc.Translate(trans.SMS_RECEIPT_YOU_GAVE), transfer.Data.GetAmount(), user.FullName())
	default:
		return m, errors.New("Unknown direction: " + string(transfer.Data.Direction()))
	}
	smsText += "\n\n" + whc.Translate(trans.SMS_CLICK_TO_CONFIRM_OR_DECLINE, receiptUrl)

	chatEntity := whc.ChatData()

	var (
		smsStatusMessageID int
		//smsStatusMessageUpdatesCount int
	)

	var createSmsStatusMessage = func() error {
		var msgSmsStatus botsfw.MessageFromBot
		mt := whc.Translate(trans.MESSAGE_TEXT_SMS_QUEUING_FOR_SENDING, phoneContact.PhoneNumberAsString())
		//log.Debugf(c, "whc.InputTypes(): %v, botsfw.WebhookInputCallbackQuery: %v, MessageID: %v", whc.InputTypes(), botsfw.WebhookInputCallbackQuery, whc.InputCallbackQuery().GetMessage().IntID())
		if whc.InputType() == botsfw.WebhookInputCallbackQuery {
			//log.Debugf(c, "editMessage.MessageID: %v", editMessage.MessageID)
			if msgSmsStatus, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
				return err
			}
		} else {
			msgSmsStatus = whc.NewMessage(mt)
		}
		smsStatusMsg, err := whc.Responder().SendMessage(c, msgSmsStatus, botsfw.BotAPISendMessageOverHTTPS)
		if err != nil {
			return err
		}
		smsStatusMessageID = smsStatusMsg.TelegramMessage.(tgbotapi.Message).MessageID
		chatEntity.AddWizardParam(SMS_STATUS_MESSAGE_ID_PARAM_NAME, strconv.Itoa(smsStatusMessageID))
		return nil
	}

	if err = createSmsStatusMessage(); err != nil {
		return m, err
	}
	//if smsStatusMessageID, err = strconv.Atoi(chatEntity.GetWizardParam(SMS_STATUS_MESSAGE_ID_PARAM_NAME)); err != nil {
	//	if err = createSmsStatusMessage(); err != nil {
	//		return m, err
	//	}
	//}
	//if smsStatusMessageUpdatesCount, err = strconv.Atoi(chatEntity.GetWizardParam(SMS_STATUS_MESSAGE_UPDATES_COUNT_PARAM_NAME)); err == nil {
	//	if smsStatusMessageUpdatesCount > 2 {
	//		if err = createSmsStatusMessage(); err != nil {
	//			return m, err
	//		}
	//		chatEntity.AddWizardParam(SMS_STATUS_MESSAGE_UPDATES_COUNT_PARAM_NAME, "1")
	//	} else {
	//		chatEntity.AddWizardParam(SMS_STATUS_MESSAGE_UPDATES_COUNT_PARAM_NAME, strconv.Itoa(smsStatusMessageUpdatesCount + 1))
	//	}
	//} else {
	//	chatEntity.AddWizardParam(SMS_STATUS_MESSAGE_UPDATES_COUNT_PARAM_NAME, "1")
	//}

	tgChatID, err := strconv.ParseInt(whc.MustBotChatID(), 10, 64)

	if err != nil {
		return m, fmt.Errorf("failed to parse whc.BotChatID() to int: %w", err)
	}

	if lastTwilioSmsese, err := dtdal.Twilio.GetLastTwilioSmsesForUser(c, tx, whc.AppUserID(), phoneContact.PhoneNumberAsString(), 1); err != nil {
		err = fmt.Errorf("failed to check latest SMS records: %w", err)
		return m, err
	} else if len(lastTwilioSmsese) > 0 {
		smsRecord := lastTwilioSmsese[0]
		if smsRecord.Data.To == phoneContact.PhoneNumberAsString() && (smsRecord.Data.Status == "delivered" || smsRecord.Data.Status == "queued") {
			// TODO: Do smarter check for limit
			m.Text = emoji.ERROR_ICON + " " + fmt.Sprintf("Exceeded limit for sending SMS to same number: %v", phoneContact.PhoneNumberAsString())
			log.Warningf(c, m.Text)
			return m, err
		}
	}
	// TODO: Create SMS record before sending to ensure we don't spam user in case of bug after the API call.

	isTestSender, smsResponse, twilioException, err := sms.SendSms(whc.Context(), whc.GetBotSettings().Env == "prod", phoneContact.PhoneNumberAsString(), smsText)
	if err != nil {
		return m, fmt.Errorf("failed to send SMS: %w", err)
	}
	//sms := common.Sms{
	//	DtCreated: smsResponse.DateCreated,
	//	DtUpdate: smsResponse.DateUpdate,
	//	DtSent: smsResponse.DateSent,
	//	InviteCode: inviteCode,
	//	To: smsResponse.To,
	//	From: smsResponse.From,
	//	Status: smsResponse.Status,
	//}
	//if smsResponse.Price != nil {
	//	sms.Price = *smsResponse.Price
	//}

	if twilioException != nil {
		twilioExceptionStr, _ := json.Marshal(twilioException)
		log.Errorf(c, "Failed to send SMS via Twilio: %v", string(twilioExceptionStr))
		mt, tryAnotherNumber := sms.TwilioExceptionToMessage(whc, whc, twilioException)
		if tryAnotherNumber {
			log.Infof(c, "Twilio identified invalid phone number, need to try another one.")
			if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatText); err != nil {
				return
			}
			m.EditMessageUID = telegram.NewChatMessageUID(tgChatID, smsStatusMessageID)
			return
		}
		if counterparty.Data.PhoneNumber == phoneContact.PhoneNumber {
			var counterparty models.Contact
			counterparty, err = facade.GetContactByID(c, tx, transfer.Data.Counterparty().ContactID)
			if err != nil {
				return
			}
			if counterparty.Data.PhoneNumber != phoneContact.PhoneNumber {
				counterparty.Data.PhoneNumber = phoneContact.PhoneNumber
				err = facade.SaveContact(c, counterparty)
			}
		}
		if m, err = whc.NewEditMessage(fmt.Sprintf("<b>Exception</b>\n%v\n\n<b>SMS text</b>\n%v", twilioException, smsText), botsfw.MessageFormatHTML); err != nil {
			return
		}
		m.EditMessageUID = telegram.NewChatMessageUID(tgChatID, smsStatusMessageID)
		m.DisableWebPagePreview = true
		dtb_general.SetMainMenuKeyboard(whc, &m)
		return
	}

	smsResponseStr, _ := json.Marshal(smsResponse)
	log.Debugf(c, "Twilio response: %v", string(smsResponseStr))

	if err = analytics.ReceiptSentFromBot(whc, "sms"); err != nil {
		log.Errorf(c, "Failed to send to analytics receipt sent event: %v", err)
	}

	if _, err = dtdal.Twilio.SaveTwilioSms(
		whc.Context(),
		smsResponse,
		transfer,
		phoneContact,
		whc.AppUserID(),
		tgChatID,
		smsStatusMessageID,
	); err != nil {
		return
	}

	mt := whc.Translate(trans.MESSAGE_TEXT_SMS_QUEUED_FOR_SENDING, phoneContact.PhoneNumberAsString())

	if isTestSender {
		mt += "\n\n<b>SMS text</b>\n" + smsText
	}

	if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
		return
	}
	m.EditMessageUID = telegram.NewChatMessageUID(tgChatID, smsStatusMessageID)
	m.DisableWebPagePreview = true

	if _, err := whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		err = fmt.Errorf("failed to send bot response message over HTTPS: %w", err)
		return m, err
	}

	return dtb_general.MainMenuCommand.Action(whc)
}
