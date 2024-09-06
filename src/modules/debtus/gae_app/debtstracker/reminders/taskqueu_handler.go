package reminders

import (
	"context"
	"errors"
	"fmt"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	gaedal2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal/gaedal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"net/http"
	"strconv"
	"time"
)

func SendReminderHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	//func sendNotificationForDueTransfer(ctx context.Context, key *datastore.Key) {
	err := r.ParseForm()
	if err != nil {
		logus.Errorf(ctx, "Failed to parse form")
		return
	}
	reminderID := r.FormValue("id")
	if reminderID == "" {
		logus.Errorf(ctx, "Failed to convert reminder ContactID to int")
		return
	}
	if err = sendReminder(ctx, reminderID); err != nil {
		logus.Errorf(ctx, err.Error())
		if !dal.IsNotFound(err) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func sendReminder(ctx context.Context, reminderID string) (err error) {
	logus.Debugf(ctx, "sendReminder(reminderID=%v)", reminderID)
	if reminderID == "" {
		return errors.New("reminderID == 0")
	}

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	reminder, err := dtdal.Reminder.GetReminderByID(ctx, db, reminderID)
	if err != nil {
		return err
	}
	if reminder.Data.Status != models4debtus.ReminderStatusCreated {
		logus.Infof(ctx, "reminder.Status:%v != models.ReminderStatusCreated", reminder.Data.Status)
		return nil
	}

	transfer, err := facade4debtus.Transfers.GetTransferByID(ctx, nil, reminder.Data.TransferID)
	if err != nil {
		if dal.IsNotFound(err) {
			logus.Errorf(ctx, err.Error())
			if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
				if reminder, err = dtdal.Reminder.GetReminderByID(ctx, tx, reminderID); err != nil {
					return
				}
				reminder.Data.Status = "invalid:no-transfer"
				reminder.Data.DtUpdated = time.Now()
				reminder.Data.DtNext = time.Time{}
				if err = dtdal.Reminder.SaveReminder(ctx, tx, reminder); err != nil {
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
		logus.Infof(ctx, "TransferEntry(id=%v) is not outstanding, transfer.Amount=%v, transfer.AmountInCentsReturned=%v", reminder.Data.TransferID, transfer.Data.AmountInCents, transfer.Data.AmountReturned())

		if err := gaedal2.DiscardReminder(ctx, reminderID, reminder.Data.TransferID, ""); err != nil {
			return fmt.Errorf("failed to discard a reminder for non outstanding transfer id=%v: %w", reminder.Data.TransferID, err)
		}
		return nil
	}

	if err = sendReminderToUser(ctx, reminderID, transfer); err != nil {
		logus.Errorf(ctx, "Failed to send reminder (id=%v) for transfer %v: %v", reminderID, reminder.Data.TransferID, err.Error())
	}

	return nil
}

var errReminderAlreadySentOrIsBeingSent = errors.New("reminder already sent or is being sent")

func sendReminderToUser(ctx context.Context, reminderID string, transfer models4debtus.TransferEntry) (err error) {

	var reminder models4debtus.Reminder

	// If sending notification failed do not try to resend - to prevent spamming.
	if err = facade.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if reminder, err = dtdal.Reminder.GetReminderByID(ctx, tx, reminderID); err != nil {
			return fmt.Errorf("failed to get reminder by id=%v: %w", reminderID, err)
		}
		if reminder.Data.Status != models4debtus.ReminderStatusCreated {
			return errReminderAlreadySentOrIsBeingSent
		}
		reminder.Data.Status = models4debtus.ReminderStatusSending
		if err = dtdal.Reminder.SaveReminder(tctx, tx, reminder); err != nil { // TODO: UserEntry dtdal.Reminder.SaveReminder()
			return fmt.Errorf("failed to save reminder with new status to db: %w", err)
		}
		return
	}, nil); err != nil {
		if errors.Is(err, errReminderAlreadySentOrIsBeingSent) {
			logus.Infof(ctx, err.Error())
		} else {
			err = fmt.Errorf("failed to update reminder status to '%v': %w", models4debtus.ReminderStatusSending, err)
			logus.Errorf(ctx, err.Error())
		}
		return
	} else {
		logus.Infof(ctx, "Updated Reminder(id=%v) status to '%v'.", reminderID, models4debtus.ReminderStatusSending)
	}

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}
	user := dbo4userus.NewUserEntry(reminder.Data.UserID)
	if err = dal4userus.GetUser(ctx, db, user); err != nil {
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
			_, tgChat, err = gaedal2.GetTelegramChatByUserID(ctx, reminder.Data.UserID) // TODO: replace with DAL method
			if err != nil {
				if dal.IsNotFound(err) { // TODO: Get rid of datastore reference
					err = fmt.Errorf("failed to call gaedal.GetTelegramChatByUserID(userID=%v): %w", reminder.Data.UserID, err)
					return
				}
			} else {
				if tgChatID, err = strconv.ParseInt(tgChat.BaseTgChatData().BotUserIDs[0], 10, 64); err != nil {
					err = fmt.Errorf("failed to parse tgChat.BaseTgChatData().BotUserID=%v: %w", tgChat.BaseTgChatData().BotUserIDs[0], err)
					return
				}
				tgBotID = "TODO:setup_bot_id_for_reminder"
			}
		}
		if tgChatID != 0 {
			if reminderIsSent, channelDisabledByUser, err = sendReminderByTelegram(ctx, transfer, reminder, tgChatID, tgBotID); err != nil {
				return
			} else if !reminderIsSent && !channelDisabledByUser {
				logus.Warningf(ctx, "Reminder is not sent to Telegram, err=%v", err)
			}
		}
	}
	if !reminderIsSent { // TODO: This is wrong to send same reminder by email if Telegram failed, complex and will screw up stats <= Are you sure?
		if user.Data.Email != "" {
			if err = sendReminderByEmail(ctx, reminder, user.Data.Email, transfer, user); err != nil {
				logus.Errorf(ctx, "Failure in sendReminderByEmail()")
			}
		} else {
			if !channelDisabledByUser {
				logus.Errorf(ctx, "Can't send reminder")
			}
			err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
				if reminder, err = dtdal.Reminder.GetReminderByID(ctx, tx, reminderID); err != nil {
					return err
				}
				reminder.Data.Status = models4debtus.ReminderStatusFailed
				return dtdal.Reminder.SaveReminder(ctx, tx, reminder)
			}, nil)
			if err != nil {
				logus.Errorf(ctx, fmt.Errorf("failed to set reminder status to '%v': %w", models4debtus.ReminderStatusFailed, err).Error())
			} else {
				logus.Infof(ctx, "Reminder status set to '%v'", reminder.Data.Status)
			}
		}
	}
	return nil // TODO: Handle errors!
}
