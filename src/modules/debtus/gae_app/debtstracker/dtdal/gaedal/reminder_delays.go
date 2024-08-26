package gaedal

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/core/queues"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	tgbots2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/platforms/debtustgbots"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"strings"
	"time"
)

func (ReminderDalGae) DelayCreateReminderForTransferUser(c context.Context, transferID string, userID string) (err error) {
	if transferID == "" {
		panic("transferID == 0")
	}
	if userID == "" {
		panic("userID == 0")
	}
	//if !dtdal.DB.IsInTransaction(c) {
	//	panic("This function should be called within transaction")
	//}
	if err = delayerCreateReminderForTransferUser.EnqueueWork(c, delaying.With(queues.QueueReminders, "create-reminder-4-transfer-user", 0), transferID, userID); err != nil {
		return fmt.Errorf("failed to create a task for reminder creation. transferID=%v, userID=%v: %w", transferID, userID, err)
	}
	logus.Debugf(c, "Added task to create reminder for transfer id=%s", transferID)
	return
}

func delayedCreateReminderForTransferUser(c context.Context, transferID string, userID string) (err error) {
	logus.Debugf(c, "delayedCreateReminderForTransferUser(transferID=%s, userID=%s)", transferID, userID)
	if transferID == "" {
		logus.Errorf(c, "transferID == 0")
		return nil
	}
	if userID == "" {
		logus.Errorf(c, "userID == 0")
		return nil
	}

	return facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var transfer models4debtus.TransferEntry
		transfer, err = facade4debtus.Transfers.GetTransferByID(c, tx, transferID)
		if err != nil {
			if dal.IsNotFound(err) {
				logus.Errorf(c, fmt.Errorf("not able to create reminder for specified transfer: %w", err).Error())
				return
			}
			return fmt.Errorf("failed to get transfer by id: %w", err)
		}
		transferUserInfo := transfer.Data.UserInfoByUserID(userID)
		if transferUserInfo.UserID != userID {
			panic("transferUserInfo.UserID != userID")
		}

		if transferUserInfo.ReminderID != "" {
			logus.Warningf(c, "TransferEntry user already has reminder # %v", transferUserInfo.ReminderID)
			return
		}

		if transferUserInfo.TgChatID == 0 { // TODO: Try to get TgChat from user record or check other channels?
			logus.Warningf(c, "TransferEntry user has no associated TgChatID: %+v", transferUserInfo)
			return
		}

		//reminderKey := NewReminderIncompleteKey(c)
		next := transfer.Data.DtDueOn
		isAutomatic := next.IsZero()
		if isAutomatic {
			if strings.Contains(strings.ToLower(transfer.Data.CreatedOnID), "dev") {
				next = time.Now().Add(2 * time.Minute)
			} else {
				next = time.Now().Add(7 * 24 * time.Hour)
			}
		}
		reminder := models4debtus.NewReminder("", models4debtus.NewReminderViaTelegram(transferUserInfo.TgBotID, transferUserInfo.TgChatID, userID, transferID, isAutomatic, next))
		if err = tx.Insert(c, reminder.Record); err != nil {
			return fmt.Errorf("failed to save reminder to db: %w", err)
		}
		reminderID := reminder.Key.ID.(string)
		logus.Infof(c, "Created reminder id=%v", reminderID)
		if err = QueueSendReminder(c, reminderID, time.Until(next)); err != nil {
			return fmt.Errorf("failed to queue reminder for sending: %w", err)
		}
		transferUserInfo.ReminderID = reminderID

		if err = facade4debtus.Transfers.SaveTransfer(c, tx, transfer); err != nil {
			return fmt.Errorf("failed to save transfer to db: %w", err)
		}

		return
	})
}

func (ReminderDalGae) DelayDiscardReminders(c context.Context, transferIDs []string, returnTransferID string) error {
	if len(transferIDs) > 0 {
		return delayerDiscardReminders.EnqueueWork(c, delaying.With(queues.QueueReminders, "discard-reminders", 0), transferIDs, returnTransferID)
	} else {
		logus.Warningf(c, "DelayDiscardReminders(): len(transferIDs)==0")
		return nil
	}
}

func delayedDiscardReminders(c context.Context, transferIDs []int, returnTransferID int) error {
	logus.Debugf(c, "delayedDiscardReminders(transferIDs=%+v, returnTransferID=%d)", transferIDs, returnTransferID)
	if len(transferIDs) == 0 {
		return errors.New("len(transferIDs) == 0")
	}
	const queueName = queues.QueueReminders
	args := make([][]interface{}, len(transferIDs))
	for i, transferID := range transferIDs {
		args[i] = []interface{}{transferID, returnTransferID}
	}
	return delayerDiscardRemindersForTransfer.EnqueueWorkMulti(c, delaying.With(queueName, "discard-reminders-for-transfer", 0), args...)
}

func delayedDiscardRemindersForTransfer(c context.Context, transferID, returnTransferID int) error {
	logus.Debugf(c, "delayedDiscardReminders(transferID=%v, returnTransferID=%v)", transferID, returnTransferID)
	if transferID == 0 {
		logus.Errorf(c, "transferID == 0")
		return nil
	}
	delayDuration := time.Millisecond * 10
	var _discard = func(
		getIDs func(context.Context, dal.ReadSession, int) ([]int, error),
		loadedFormat, notLoadedFormat string,
	) error {
		if reminderIDs, err := getIDs(c, nil, transferID); err != nil {
			return err
		} else if len(reminderIDs) > 0 {
			logus.Debugf(c, loadedFormat, len(reminderIDs), transferID)
			for _, reminderID := range reminderIDs {
				if err := delayerDiscardReminder.EnqueueWork(c, delaying.With(queues.QueueReminders, "discard-reminder", delayDuration), reminderID, transferID, returnTransferID); err != nil {
					return fmt.Errorf("failed to create a task for reminder ContactID=%v: %w", reminderID, err)
				}
				delayDuration += time.Millisecond * 10
			}
		} else {
			logus.Infof(c, notLoadedFormat, transferID)
		}
		return nil
	}
	if err := _discard(dtdal.Reminder.GetActiveReminderIDsByTransferID, "Loaded %v keys of active reminders for transfer id=%v", "The are no ative reminders for transfer id=%v"); err != nil {
		return err
	}
	if err := _discard(dtdal.Reminder.GetSentReminderIDsByTransferID, "Loaded %v keys of sent reminders for transfer id=%v", "The are no sent reminders for transfer id=%v"); err != nil {
		return err
	}
	return nil
}

func DiscardReminder(c context.Context, reminderID, transferID, returnTransferID string) (err error) {
	return facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		return discardReminder(c, tx, reminderID, transferID, returnTransferID)
	})
}

func delayedDiscardReminder(c context.Context, reminderID, transferID, returnTransferID string) (err error) {
	return facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = discardReminder(c, tx, reminderID, transferID, returnTransferID); err == ErrDuplicateAttemptToDiscardReminder {
			logus.Errorf(c, err.Error())
			return nil
		}
		return err
	})
}

func discardReminder(ctx context.Context, tx dal.ReadwriteTransaction, reminderID, transferID, returnTransferID string) (err error) {
	logus.Debugf(ctx, "discardReminder(reminderID=%v, transferID=%v, returnTransferID=%v)", reminderID, transferID, returnTransferID)

	var (
		transfer = models4debtus.NewTransfer(transferID, nil)
		reminder = models4debtus.NewReminder(reminderID, new(models4debtus.ReminderDbo))
	)

	if returnTransferID > "" {
		//returnTransferKey := models.NewTransferKey(returnTransferID)
		returnTransfer := models4debtus.NewTransfer(returnTransferID, nil)
		//keys := []*datastore.Key{reminderKey, transferKey, returnTransferKey}
		if err = tx.GetMulti(ctx, []dal.Record{reminder.Record, transfer.Record, returnTransfer.Record}); err != nil {
			return err
		}
	} else {
		if err = tx.GetMulti(ctx, []dal.Record{reminder.Record, transfer.Record}); err != nil {
			return err
		}
	}

	if reminder, err = dtdal.Reminder.SetReminderStatus(ctx, reminderID, returnTransferID, models4debtus.ReminderStatusDiscarded, time.Now()); err != nil {
		return err // DO NOT WRAP as there is check in delayedDiscardReminder() errors.Wrapf(err, "Failed to set reminder status to '%v'", models.ReminderStatusDiscarded)
	}

	switch reminder.Data.SentVia {
	case telegram.PlatformID: // We need to update a reminder message if it was already sent out
		if reminder.Data.BotID == "" {
			logus.Errorf(ctx, "reminder.BotID == ''")
			return nil
		}
		if reminder.Data.MessageIntID == 0 {
			//logus.Infof(ctx, "No need to update reminder message in Telegram as a reminder is not sent yet")
			return nil
		}
		logus.Infof(ctx, "Will try to update a reminder message as it was already sent to user, reminder.MessageIntID: %v", reminder.Data.MessageIntID)
		tgBotApi := tgbots2.GetTelegramBotApiByBotCode(ctx, reminder.Data.BotID)
		if tgBotApi == nil {
			return fmt.Errorf("not able to create API client as there no settings for telegram bot with id '%v'", reminder.Data.BotID)
		}

		if reminder.Data.Locale == "" {
			logus.Errorf(ctx, "reminder.Locale == ''")
			user := dbo4userus.NewUserEntry(reminder.Data.UserID)
			if err = dal4userus.GetUser(ctx, nil, user); err != nil { // Intentionally do not use transaction
				return err
			}
			if user.Data.PreferredLocale != "" {
				reminder.Data.Locale = user.Data.PreferredLocale
			} else if s, ok := tgbots2.Bots(dtdal.HttpAppHost.GetEnvironment(ctx, nil)).ByCode[reminder.Data.BotID]; ok {
				reminder.Data.Locale = s.Locale.Code5
			}
		}

		translator := GetTranslatorForReminder(ctx, reminder.Data)

		utmParams := common4debtus.UtmParams{
			Source:   "TODO", // TODO: Get bot ContactID
			Medium:   telegram.PlatformID,
			Campaign: common4debtus.UTM_CAMPAIGN_RECEIPT_DISCARD,
		}

		receiptMessageText := common4debtus.TextReceiptForTransfer(
			ctx,
			translator,
			transfer,
			reminder.Data.UserID,
			common4debtus.ShowReceiptToAutodetect,
			utmParams,
		)

		locale := i18n.GetLocaleByCode5(reminder.Data.Locale) // TODO: Check for supported locales

		transferUrlForUser := common4debtus.GetTransferUrlForUser(ctx, transferID, reminder.Data.UserID, locale, utmParams)

		receiptMessageText += "\n\n" + strings.Join([]string{
			translator.Translate(trans.MESSAGE_TEXT_DEBT_IS_RETURNED),
			fmt.Sprintf(`<a href="%v">%v</a>`, transferUrlForUser, translator.Translate(trans.MESSAGE_TEXT_DETAILS_ARE_HERE)),
		}, "\n")

		tgMessage := tgbotapi.NewEditMessageText(reminder.Data.ChatIntID, int(reminder.Data.MessageIntID), "", receiptMessageText)
		tgMessage.ParseMode = "HTML"
		if _, err = tgBotApi.Send(tgMessage); err != nil {
			return fmt.Errorf("failed to send message to Telegram: %w", err)
		}

	default:
		return errors.New("Unknown reminder channel: %v" + reminder.Data.SentVia)
	}

	return err
}

func GetTranslatorForReminder(c context.Context, reminder *models4debtus.ReminderDbo) i18n.SingleLocaleTranslator {
	return i18n.NewSingleMapTranslator(i18n.GetLocaleByCode5(reminder.Locale), i18n.NewMapTranslator(c, trans.TRANS))
}

var ErrDuplicateAttemptToDiscardReminder = errors.New("Duplicate attempt to close reminder by same return transfer")

func (ReminderDalGae) SetReminderStatus(c context.Context, reminderID, returnTransferID string, status string, when time.Time) (reminder models4debtus.Reminder, err error) {
	var (
		changed        bool
		previousStatus string
	)
	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if reminder, err = dtdal.Reminder.GetReminderByID(c, tx, reminderID); err != nil {
			return
		} else {
			switch status {
			case string(models4debtus.ReminderStatusDiscarded):
				reminder.Data.DtDiscarded = when
			case string(models4debtus.ReminderStatusSent):
				reminder.Data.DtSent = when
			case string(models4debtus.ReminderStatusSending):
				// pass
			case string(models4debtus.ReminderStatusViewed):
				reminder.Data.DtViewed = when
			case string(models4debtus.ReminderStatusUsed):
				reminder.Data.DtUsed = when
			default:
				return errors.New("unsupported status: " + status)
			}
			previousStatus = reminder.Data.Status
			changed = previousStatus != status
			if returnTransferID != "" && status == string(models4debtus.ReminderStatusDiscarded) {
				for _, id := range reminder.Data.ClosedByTransferIDs { // TODO: WTF are we doing here?
					if id == returnTransferID {
						logus.Infof(c, "new status: '%v', Reminder{Status: '%v', ClosedByTransferIDs: %v}", status, reminder.Data.Status, reminder.Data.ClosedByTransferIDs)
						return ErrDuplicateAttemptToDiscardReminder
					}
				}
				reminder.Data.ClosedByTransferIDs = append(reminder.Data.ClosedByTransferIDs, returnTransferID)
				changed = true
			}
			if changed {
				reminder.Data.Status = status
				if err = tx.Set(c, reminder.Record); err != nil {
					err = fmt.Errorf("failed to save reminder to db (id=%v): %w", reminderID, err)
				}
			}
			return
		}
	}, nil)
	if err == nil {
		if changed {
			logus.Debugf(c, "Reminder(id=%v) status changed from '%v' to '%v'", reminderID, previousStatus, status)
		} else {
			logus.Debugf(c, "Reminder(id=%v) status not changed as already '%v'", reminderID, status)
		}
	}
	return
}
