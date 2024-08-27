package dtb_transfer

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	dtb_general2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/analytics"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/sms"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/logus"

	//"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/invites"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"context"
	"errors"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/debtstracker-translations/emoji"
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
		ctx := whc.Context()
		userCtx := facade.NewUserContext(whc.AppUserID())
		return m, dal4userus.RunUserWorker(ctx, userCtx, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
			logus.Debugf(ctx, "AskPhoneNumberForReceiptCommand.Action()")

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
				if params.User.Data.Names.FirstName == contact.FirstName() && params.User.Data.Names.LastName == contact.LastName() {
					phoneNumberStr = cleanPhoneNumber(contact.PhoneNumber())
					if phoneNumber, err = strconv.ParseInt(phoneNumberStr, 10, 64); err != nil {
						logus.Warningf(ctx, "Failed to parse contact's phone number: [%v]", phoneNumberStr)
						return err
					} else if len(params.User.Data.Phones) == 0 {
						params.User.Data.Phones = append(params.User.Data.Phones, dbmodels.PersonPhone{
							Number:   strconv.FormatInt(phoneNumber, 10),
							Verified: true,
						})
						params.UserUpdates = append(params.UserUpdates, dal.Update{
							Field: "phones",
							Value: params.User.Data.Phones,
						})
					}
					m = whc.NewMessage(trans.MESSAGE_TEXT_YOU_CAN_SEND_RECEIPT_TO_YOURSELF_BY_SMS)
					return nil
				}
				mt = contact.PhoneNumber()
			} else {
				mt = whc.Input().(botsfw.WebhookTextMessage).Text()
			}

			if twilioTestNumber, ok := common4debtus.TwilioTestNumbers[mt]; ok {
				logus.Debugf(ctx, "Using predefined test number [%v]: %v", mt, twilioTestNumber)
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

			transferID := awaitingUrl.Query().Get(WizardParamTransfer)
			if transferID == "" {
				return fmt.Errorf("empty transferID")
			}
			transfer, err := facade4debtus.Transfers.GetTransferByID(ctx, tx, transferID)
			if err != nil {
				return fmt.Errorf("failed to get transfer by ContactID: %w", err)
			}
			spaceID := params.User.Data.GetFamilySpaceID()
			counterparty, err := facade4debtus.GetDebtusSpaceContactByID(ctx, tx, spaceID, transfer.Data.Counterparty().ContactID)
			if err != nil {
				return err
			}
			phoneContact := dto4contactus.PhoneContact{PhoneNumber: phoneNumber, PhoneNumberConfirmed: false}

			m, err = sendReceiptBySms(whc, tx, spaceID, phoneContact, transfer, counterparty)
			return err
		})

	},
}

const SMS_STATUS_MESSAGE_ID_PARAM_NAME = "SmsStatusMessageId"
const SMS_STATUS_MESSAGE_UPDATES_COUNT_PARAM_NAME = "SmsStatusUpdatesCount"

func sendReceiptBySms(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, spaceID string, phoneContact dto4contactus.PhoneContact, transfer models4debtus.TransferEntry, counterparty models4debtus.DebtusSpaceContactEntry) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	if transfer.Data == nil {
		if transfer, err = facade4debtus.Transfers.GetTransferByID(ctx, tx, transfer.ID); err != nil {
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

	receiptData := models4debtus.NewReceiptEntity(whc.AppUserID(), transfer.ID, transfer.Data.Counterparty().UserID, whc.Locale().Code5, "sms", strconv.FormatInt(phoneContact.PhoneNumber, 10), general.CreatedOn{
		CreatedOnPlatform: whc.BotPlatform().ID(),
		CreatedOnID:       whc.GetBotCode(),
	})
	var receipt models4debtus.ReceiptEntry
	if receipt, err = dtdal.Receipt.CreateReceipt(ctx, receiptData); err != nil {
		return m, err
	}

	receiptUrl := common4debtus.GetReceiptUrl(receipt.ID, common4debtus.GetWebsiteHost(receiptData.CreatedOnID))

	if counterparty.Data.CounterpartyUserID == "" {
		//related := fmt.Sprintf("%v=%v", models.TransfersCollection, transferID)
		//inviteKey, invite, err := invites.CreatePersonalInvite(whc, whc.AppUserID(), invites.InviteBySms, strconv.FormatInt(phoneContact.PhoneNumber, 10), whc.BotPlatform().ContactID(), whc.GetBotCode(), related)
		//if err != nil {
		//	logus.Errorf(ctx, "Failed to create invite: %v", err)
		//	return m, err
		//}
		//inviteCode = inviteKey.StringID()
	} else {
		panic("Not implemented, need to call anybot.GetReceiptUrlForUser(...)")
	}

	// You've got $10 from Jack
	// You've given $10 to Jack

	switch transfer.Data.Direction() {
	case models4debtus.TransferDirectionUser2Counterparty:
		smsText = fmt.Sprintf(whc.Translate(trans.SMS_RECEIPT_YOU_GOT), transfer.Data.GetAmount(), user.FullName())
	case models4debtus.TransferDirectionCounterparty2User:
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
		//logus.Debugf(ctx, "whc.InputTypes(): %v, botsfw.WebhookInputCallbackQuery: %v, MessageID: %v", whc.InputTypes(), botsfw.WebhookInputCallbackQuery, whc.InputCallbackQuery().GetMessage().IntID())
		if whc.InputType() == botsfw.WebhookInputCallbackQuery {
			//logus.Debugf(ctx, "editMessage.MessageID: %v", editMessage.MessageID)
			if msgSmsStatus, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
				return err
			}
		} else {
			msgSmsStatus = whc.NewMessage(mt)
		}
		smsStatusMsg, err := whc.Responder().SendMessage(ctx, msgSmsStatus, botsfw.BotAPISendMessageOverHTTPS)
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

	if lastTwilioSmsese, err := dtdal.Twilio.GetLastTwilioSmsesForUser(ctx, tx, whc.AppUserID(), phoneContact.PhoneNumberAsString(), 1); err != nil {
		err = fmt.Errorf("failed to check latest SMS records: %w", err)
		return m, err
	} else if len(lastTwilioSmsese) > 0 {
		smsRecord := lastTwilioSmsese[0]
		if smsRecord.Data.To == phoneContact.PhoneNumberAsString() && (smsRecord.Data.Status == "delivered" || smsRecord.Data.Status == "queued") {
			// TODO: Do smarter check for limit
			m.Text = emoji.ERROR_ICON + " " + fmt.Sprintf("Exceeded limit for sending SMS to same number: %v", phoneContact.PhoneNumberAsString())
			logus.Warningf(ctx, m.Text)
			return m, err
		}
	}
	// TODO: Create SMS record before sending to ensure we don't spam user in case of bug after the API call.

	isTestSender, smsResponse, twilioException, err := sms.SendSms(whc.Context(), whc.GetBotSettings().Env == "prod", phoneContact.PhoneNumberAsString(), smsText)
	if err != nil {
		return m, fmt.Errorf("failed to send SMS: %w", err)
	}
	//sms := anybot.Sms{
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
		logus.Errorf(ctx, "Failed to send SMS via Twilio: %v", string(twilioExceptionStr))
		mt, tryAnotherNumber := sms.TwilioExceptionToMessage(whc, whc, twilioException)
		if tryAnotherNumber {
			logus.Infof(ctx, "Twilio identified invalid phone number, need to try another one.")
			if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatText); err != nil {
				return
			}
			m.EditMessageUID = telegram.NewChatMessageUID(tgChatID, smsStatusMessageID)
			return
		}
		if counterparty.Data.PhoneNumber == phoneContact.PhoneNumber {
			var counterparty models4debtus.DebtusSpaceContactEntry
			counterparty, err = facade4debtus.GetDebtusSpaceContactByID(ctx, tx, spaceID, transfer.Data.Counterparty().ContactID)
			if err != nil {
				return
			}
			if counterparty.Data.PhoneNumber != phoneContact.PhoneNumber {
				counterparty.Data.PhoneNumber = phoneContact.PhoneNumber
				err = facade4debtus.SaveContact(ctx, counterparty)
			}
		}
		if m, err = whc.NewEditMessage(fmt.Sprintf("<b>Exception</b>\n%v\n\n<b>SMS text</b>\n%v", twilioException, smsText), botsfw.MessageFormatHTML); err != nil {
			return
		}
		m.EditMessageUID = telegram.NewChatMessageUID(tgChatID, smsStatusMessageID)
		m.DisableWebPagePreview = true
		dtb_general2.SetMainMenuKeyboard(whc, &m)
		return
	}

	smsResponseStr, _ := json.Marshal(smsResponse)
	logus.Debugf(ctx, "Twilio response: %v", string(smsResponseStr))

	if err = analytics.ReceiptSentFromBot(whc, "sms"); err != nil {
		logus.Errorf(ctx, "Failed to send to analytics receipt sent event: %v", err)
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

	if _, err := whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		err = fmt.Errorf("failed to send bot response message over HTTPS: %w", err)
		return m, err
	}

	return dtb_general2.MainMenuCommand.Action(whc)
}
