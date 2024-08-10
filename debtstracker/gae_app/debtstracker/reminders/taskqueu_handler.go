package reminders

import (
	"context"
	"errors"
	"fmt"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal/gaedal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/http"
	"time"
)

func SendReminderHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
	//func sendNotificationForDueTransfer(c context.Context, key *datastore.Key) {
	err := r.ParseForm()
	if err != nil {
		logus.Errorf(c, "Failed to parse form")
		return
	}
	reminderID := r.FormValue("id")
	if reminderID == "" {
		logus.Errorf(c, "Failed to convert reminder ID to int")
		return
	}
	if err = sendReminder(c, reminderID); err != nil {
		logus.Errorf(c, err.Error())
		if !dal.IsNotFound(err) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func sendReminder(c context.Context, reminderID string) (err error) {
	logus.Debugf(c, "sendReminder(reminderID=%v)", reminderID)
	if reminderID == "" {
		return errors.New("reminderID == 0")
	}

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	reminder, err := dtdal.Reminder.GetReminderByID(c, db, reminderID)
	if err != nil {
		return err
	}
	if reminder.Data.Status != models.ReminderStatusCreated {
		logus.Infof(c, "reminder.Status:%v != models.ReminderStatusCreated", reminder.Data.Status)
		return nil
	}

	transfer, err := facade2debtus.Transfers.GetTransferByID(c, nil, reminder.Data.TransferID)
	if err != nil {
		if dal.IsNotFound(err) {
			logus.Errorf(c, err.Error())
			if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
				if reminder, err = dtdal.Reminder.GetReminderByID(c, tx, reminderID); err != nil {
					return
				}
				reminder.Data.Status = "invalid:no-transfer"
				reminder.Data.DtUpdated = time.Now()
				reminder.Data.DtNext = time.Time{}
				if err = dtdal.Reminder.SaveReminder(c, tx, reminder); err != nil {
					return
				}
				return
			}); err != nil {
				return fmt.Errorf("failed to update reminder: %w", err)
			}
			return nil
		} else {
			return fmt.Errorf("failed to load transfer: %w", err)
		}
	}

	if !transfer.Data.IsOutstanding {
		logus.Infof(c, "TransferEntry(id=%v) is not outstanding, transfer.Amount=%v, transfer.AmountInCentsReturned=%v", reminder.Data.TransferID, transfer.Data.AmountInCents, transfer.Data.AmountReturned())

		if err := gaedal.DiscardReminder(c, reminderID, reminder.Data.TransferID, ""); err != nil {
			return fmt.Errorf("failed to discard a reminder for non outstanding transfer id=%v: %w", reminder.Data.TransferID, err)
		}
		return nil
	}

	if err = sendReminderToUser(c, reminderID, transfer); err != nil {
		logus.Errorf(c, "Failed to send reminder (id=%v) for transfer %v: %v", reminderID, reminder.Data.TransferID, err.Error())
	}

	return nil
}

var errReminderAlreadySentOrIsBeingSent = errors.New("Reminder already sent or is being sent")

func sendReminderToUser(c context.Context, reminderID string, transfer models.TransferEntry) (err error) {

	var reminder models.Reminder

	// If sending notification failed do not try to resend - to prevent spamming.
	if err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
		if reminder, err = dtdal.Reminder.GetReminderByID(c, tx, reminderID); err != nil {
			return fmt.Errorf("failed to get reminder by id=%v: %w", reminderID, err)
		}
		if reminder.Data.Status != models.ReminderStatusCreated {
			return errReminderAlreadySentOrIsBeingSent
		}
		reminder.Data.Status = models.ReminderStatusSending
		if err = dtdal.Reminder.SaveReminder(tc, tx, reminder); err != nil { // TODO: User dtdal.Reminder.SaveReminder()
			return fmt.Errorf("failed to save reminder with new status to db: %w", err)
		}
		return
	}, nil); err != nil {
		if errors.Is(err, errReminderAlreadySentOrIsBeingSent) {
			logus.Infof(c, err.Error())
		} else {
			err = fmt.Errorf("failed to update reminder status to '%v': %w", models.ReminderStatusSending, err)
			logus.Errorf(c, err.Error())
		}
		return
	} else {
		logus.Infof(c, "Updated Reminder(id=%v) status to '%v'.", reminderID, models.ReminderStatusSending)
	}

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}
	user := models.NewAppUser(reminder.Data.UserID, nil)
	if err = db.Get(c, user.Record); err != nil {
		return err
	}

	var reminderIsSent, channelDisabledByUser bool
	if user.Data.HasAccount(telegram.PlatformID, "") {
		var (
			tgChatID int64
			tgBotID  string
		)
		if transferUserInfo := transfer.Data.UserInfoByUserID(reminder.Data.UserID); transferUserInfo.TgChatID != 0 {
			tgChatID = transferUserInfo.TgChatID
			tgBotID = transferUserInfo.TgBotID
		} else {
			var tgChat botsfwtgmodels.TgChatData
			_, tgChat, err = gaedal.GetTelegramChatByUserID(c, reminder.Data.UserID) // TODO: replace with DAL method
			if err != nil {
				if dal.IsNotFound(err) { // TODO: Get rid of datastore reference
					err = fmt.Errorf("failed to call gaedal.GetTelegramChatByUserID(userID=%v): %w", reminder.Data.UserID, err)
					return
				}
			} else {
				tgChatID = tgChat.BaseTgChatData().TelegramUserID
				tgBotID = tgChat.BaseTgChatData().BotID
			}
		}
		if tgChatID != 0 {
			if reminderIsSent, channelDisabledByUser, err = sendReminderByTelegram(c, transfer, reminder, tgChatID, tgBotID); err != nil {
				return
			} else if !reminderIsSent && !channelDisabledByUser {
				logus.Warningf(c, "Reminder is not sent to Telegram, err=%v", err)
			}
		}
	}
	if !reminderIsSent { // TODO: This is wrong to send same reminder by email if Telegram failed, complex and will screw up stats <= Are you sure?
		if user.Data.EmailAddress != "" {
			if err = sendReminderByEmail(c, reminder, user.Data.EmailAddress, transfer, *user.Data); err != nil {
				logus.Errorf(c, "Failure in sendReminderByEmail()")
			}
		} else {
			if !channelDisabledByUser {
				logus.Errorf(c, "Can't send reminder")
			}
			err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
				if reminder, err = dtdal.Reminder.GetReminderByID(c, tx, reminderID); err != nil {
					return err
				}
				reminder.Data.Status = models.ReminderStatusFailed
				return dtdal.Reminder.SaveReminder(c, tx, reminder)
			}, nil)
			if err != nil {
				logus.Errorf(c, fmt.Errorf("failed to set reminder status to '%v': %w", models.ReminderStatusFailed, err).Error())
			} else {
				logus.Infof(c, "Reminder status set to '%v'", reminder.Data.Status)
			}
		}
	}
	return nil // TODO: Handle errors!
}
