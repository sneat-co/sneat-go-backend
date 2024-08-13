package reminders

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/emailing"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/emails"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"time"
)

func sendReminderByEmail(c context.Context, reminder models4debtus.Reminder, emailTo string, transfer models4debtus.TransferEntry, user dbo4userus.UserEntry) (err error) {
	logus.Debugf(c, "sendReminderByEmail(reminder.ContactID=%v, emailTo=%v)", reminder.ID, emailTo)

	emailMessage := emails.Email{
		From: common4debtus.FROM_REMINDER,
		To: []string{
			emailTo, // Required
		},
		Subject: "Due payment notification",
		Text:    fmt.Sprintf("Hi %v, you have a due payment to %v: %v%v.", transfer.Data.Counterparty().ContactName, user.Data.Names.UserName, transfer.Data.AmountInCents, transfer.Data.Currency),
	}

	var emailClient emails.Client

	if emailClient, err = emailing.GetEmailClient(c); err != nil {
		return
	}

	var sent emails.Sent
	sent, err = emailClient.Send(emailMessage)

	sentAt := time.Now()

	var errDetails string
	if err != nil {
		errDetails = err.Error()
	}
	var emailMessageID string
	if sent != nil {
		emailMessageID = sent.MessageID()
	}

	if err = dtdal.Reminder.SetReminderIsSent(c, reminder.ID, sentAt, 0, emailMessageID, i18n.LocaleCodeEnUS, errDetails); err != nil {
		if err = dtdal.Reminder.DelaySetReminderIsSent(c, reminder.ID, sentAt, 0, emailMessageID, i18n.LocaleCodeEnUS, errDetails); err != nil {
			return fmt.Errorf("failed to delay setting reminder as sent: %w", err)
		}
	}

	if err != nil {
		// Print the error, cast err to awserr.Error to get the ByCode and
		// Message from an error.
		return fmt.Errorf("failed to send email using AWS SES: %w", err)
	}

	// Pretty-print the response data.
	logus.Debugf(c, "AWS SES output (for Reminder=%v): %v", reminder.ID, sent)
	return nil
}
