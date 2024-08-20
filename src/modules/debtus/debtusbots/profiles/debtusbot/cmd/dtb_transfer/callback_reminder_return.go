package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/dtb_common"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/logus"
	"net/url"
	"time"

	"github.com/sneat-co/debtstracker-translations/emoji"
)

var ReturnCallbackCommand = botsfw.NewCallbackCommand(dtb_common.CALLBACK_DEBT_RETURNED_PATH, ProcessReturnAnswer)

func ProcessReturnAnswer(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	//
	c := whc.Context()
	logus.Debugf(c, "ProcessReturnAnswer()")
	q := callbackUrl.Query()
	reminderID := q.Get("reminder")
	if reminderID == "" {
		return m, fmt.Errorf("missing reminder ContactID")
	}
	var transferID string
	if err != nil {
		if q.Get("reminder") == "" { // TODO: Remove this obsolete branch
			if transferID = q.Get("id"); transferID == "" {
				return m, fmt.Errorf("missing reminder ContactID and transfer ContactID")
			}
		}
	} else {
		if reminder, err := dtdal.Reminder.SetReminderStatus(c, reminderID, "", models4debtus.ReminderStatusUsed, time.Now()); err != nil {
			return m, err
		} else {
			transferID = reminder.Data.TransferID
		}
	}

	howMuch := q.Get("how-much")
	transfer, err := facade4debtus.Transfers.GetTransferByID(c, nil, transferID)
	if err != nil {
		return m, err
	}
	switch howMuch {
	case "":
		panic("Missing how-much parameter")
	case dtb_common.RETURNED_FULLY:
		return ProcessFullReturn(whc, transfer)
	case dtb_common.RETURNED_PARTIALLY:
		return ProcessPartialReturn(whc, transfer)
	case dtb_common.RETURNED_NOTHING:
		return ProcessNoReturn(whc, reminderID, transfer)
	default:
		panic(fmt.Sprintf("Unknown how-much: %v", howMuch))
	}
}

const commandCodeEnableReminderAgain = "enable-reminder-again"

var EnableReminderAgainCallbackCommand = botsfw.NewCallbackCommand(commandCodeEnableReminderAgain, func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	logus.Debugf(c, "EnableReminderAgainCallbackCommand()")
	q := callbackUrl.Query()
	var (
		reminderID string
		transfer   models4debtus.TransferEntry
	)
	if reminderID = q.Get("reminder"); reminderID == "" {
		err = fmt.Errorf("parameter 'reminder' is empty")
		return
	}
	if transfer.ID = q.Get("transfer"); transfer.ID == "" {
		err = fmt.Errorf("parameter 'transfer' is empty")
		return
	}

	if transfer, err = facade4debtus.Transfers.GetTransferByID(c, nil, transfer.ID); err != nil {
		return
	}

	return askWhenToRemindAgain(whc, reminderID, transfer)
})

func ProcessFullReturn(whc botsfw.WebhookContext, transfer models4debtus.TransferEntry) (m botsfw.MessageFromBot, err error) {
	amountValue := transfer.Data.GetOutstandingValue(time.Now())
	if amountValue == 0 {
		return dtb_general.EditReminderMessage(whc, transfer, whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ALREADY_FULLY_RETURNED))
	} else if amountValue < 0 {
		err = fmt.Errorf("data integrity error -> transfer.GetOutstandingValue():%v < 0", amountValue)
		return
	}

	amount := money.NewAmount(transfer.Data.GetAmount().Currency, amountValue)

	var (
		counterpartyID string
		direction      models4debtus.TransferDirection
	)
	userID := whc.AppUserID()
	if transfer.Data.CreatorUserID == userID {
		counterpartyID = transfer.Data.Counterparty().ContactID
		switch transfer.Data.Direction() {
		case models4debtus.TransferDirectionCounterparty2User:
			direction = models4debtus.TransferDirectionUser2Counterparty
		case models4debtus.TransferDirectionUser2Counterparty:
			direction = models4debtus.TransferDirectionCounterparty2User
		default:
			return m, fmt.Errorf("transfer %v has unknown direction '%v'", transfer.ID, transfer.Data.Direction())
		}
	} else if transfer.Data.Counterparty().UserID == userID {
		switch transfer.Data.Direction() {
		case models4debtus.TransferDirectionCounterparty2User:
		case models4debtus.TransferDirectionUser2Counterparty:
		default:
			return m, fmt.Errorf("transfer %v has unknown direction '%v'.", transfer.ID, transfer.Data.Direction())
		}
		counterpartyID = transfer.Data.Creator().ContactID
		direction = transfer.Data.Direction()
	}

	if m, err = dtb_general.EditReminderMessage(whc, transfer, whc.Translate(trans.MESSAGE_TEXT_REPLIED_DEBT_RETURNED_FULLY)); err != nil {
		return
	}

	if _, err = whc.Responder().SendMessage(whc.Context(), m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		return m, err
	}

	if m, err = CreateReturnAndShowReceipt(whc, transfer.ID, counterpartyID, direction, amount); err != nil {
		return m, err
	}

	reportReminderIsActed(whc, "reminder-acted-returned-fully")

	//TODO: edit message
	return m, err
}

func ProcessPartialReturn(whc botsfw.WebhookContext, transfer models4debtus.TransferEntry) (botsfw.MessageFromBot, error) {
	var counterpartyID string
	switch whc.AppUserID() {
	case transfer.Data.CreatorUserID:
		counterpartyID = transfer.Data.Counterparty().ContactID
	case transfer.Data.Counterparty().UserID:
		counterpartyID = transfer.Data.Creator().ContactID
	default:
		panic(fmt.Sprintf("whc.whc.AppUserID()=%v not in (transfer.Counterparty().ContactID=%v, transfer.Creator().ContactID=%v)",
			whc.AppUserID(), transfer.Data.Counterparty().ContactID, transfer.Data.Creator().ContactID))
	}
	chatEntity := whc.ChatData()
	chatEntity.SetAwaitingReplyTo("")
	chatEntity.AddWizardParam(WizardParamCounterparty, counterpartyID)
	chatEntity.AddWizardParam(WizardParamTransfer, transfer.ID)
	chatEntity.AddWizardParam("currency", string(transfer.Data.Currency))

	reportReminderIsActed(whc, "reminder-acted-returned-partially")

	return AskHowMuchHaveBeenReturnedCommand.Action(whc)
}

func askWhenToRemindAgain(whc botsfw.WebhookContext, reminderID string, transfer models4debtus.TransferEntry) (m botsfw.MessageFromBot, err error) {
	if m, err = dtb_general.EditReminderMessage(whc, transfer, whc.Translate(trans.MESSAGE_TEXT_ASK_WHEN_TO_REMIND_AGAIN)); err != nil {
		return
	}
	callbackData := fmt.Sprintf("%s?id=%s&in=%s", dtb_common.CALLBACK_REMIND_AGAIN, reminderID, "%v")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         emoji.CALENDAR_ICON + " " + whc.Translate(trans.COMMAND_TEXT_SET_DATE),
				CallbackData: fmt.Sprintf("%s?id=%s", SET_NEXT_REMINDER_DATE_COMMAND, reminderID),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: whc.Translate(trans.COMMAND_TEXT_TOMORROW), CallbackData: fmt.Sprintf(callbackData, "24h")},
			{Text: whc.Translate(trans.COMMAND_TEXT_DAY_AFTER_TOMORROW), CallbackData: fmt.Sprintf(callbackData, "48h")},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: whc.Translate(trans.COMMAND_TEXT_IN_1_WEEK), CallbackData: fmt.Sprintf(callbackData, "168h")},
			{Text: whc.Translate(trans.COMMAND_TEXT_IN_1_MONTH), CallbackData: fmt.Sprintf(callbackData, "720h")},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: whc.Translate(trans.COMMAND_TEXT_DISABLE_REMINDER), CallbackData: fmt.Sprintf(callbackData, dtb_common.C_REMIND_IN_DISABLE)},
		},
	)

	if whc.GetBotSettings().Env == "dev" {
		keyboard.InlineKeyboard = append(
			[][]tgbotapi.InlineKeyboardButton{
				{
					{
						Text:         whc.Translate(trans.COMMAND_TEXT_IN_FEW_MINUTES),
						CallbackData: fmt.Sprintf(callbackData, "1m"),
					},
				},
			},
			keyboard.InlineKeyboard...,
		)
	}
	m.IsEdit = true
	m.Keyboard = keyboard
	return
}

func ProcessNoReturn(whc botsfw.WebhookContext, reminderID string, transfer models4debtus.TransferEntry) (m botsfw.MessageFromBot, err error) {
	return askWhenToRemindAgain(whc, reminderID, transfer)
}

const (
	SET_NEXT_REMINDER_DATE_COMMAND = "set-next-reminder-date"
)

var SetNextReminderDateCallbackCommand = botsfw.Command{
	Code: SET_NEXT_REMINDER_DATE_COMMAND,
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()

		reminderID := callbackUrl.Query().Get("id")
		if reminderID == "" {
			return m, fmt.Errorf("missing reminder ContactID")
		}

		chatEntity := whc.ChatData()
		chatEntity.SetAwaitingReplyTo(SET_NEXT_REMINDER_DATE_COMMAND)
		chatEntity.AddWizardParam(WizardParamReminder, reminderID)

		reminder, err := dtdal.Reminder.GetReminderByID(c, nil, reminderID)
		if err != nil {
			return m, fmt.Errorf("failed to get reminder by id: %w", err)
		}
		transfer, err := facade4debtus.Transfers.GetTransferByID(c, nil, reminder.Data.TransferID)
		if err != nil {
			return m, fmt.Errorf("failed to get transfer by id: %w", err)
		}

		if m, err = dtb_general.EditReminderMessage(whc, transfer, whc.Translate(trans.MESSAGE_TEXT_ASK_WHEN_TO_REMIND_AGAIN)); err != nil {
			return
		}

		if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
			return m, err
		}

		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_ASK_DATE_TO_REMIND)

		return m, err
	},
	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
		m, date, err := processSetDate(whc)
		if !date.IsZero() {
			chatEntity := whc.ChatData()

			reminderID := chatEntity.GetWizardParam(WizardParamReminder)
			if err != nil {
				return m, fmt.Errorf("failed to decode reminder id: %w", err)
			}
			now := time.Now()
			sinceToday := now.Sub(now.Truncate(24 * time.Hour))

			date = date.Add(sinceToday)
			remindInDuration := date.Sub(now)
			return rescheduleReminder(whc, reminderID, remindInDuration)
		}
		return m, err
	},
}
