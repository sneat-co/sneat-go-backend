package emailing

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/core/queues"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/common4all"
	"github.com/sneat-co/sneat-go-core/emails"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"time"
)

const SendEmailTaskCode = "send-email"

func DelaySendEmail(ctx context.Context, id int64) error {
	return delayEmail.EnqueueWork(ctx, delaying.With(queues.QueueEmails, SendEmailTaskCode, 0), id)
}

var ErrEmailIsInWrongStatus = errors.New("email is already sending or sent")

func delayedSendEmail(ctx context.Context, id int64) (err error) {
	logus.Debugf(ctx, "delayedSendEmail(%v)", id)

	var email models4auth.Email

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return err
	}

	if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if email, err = common4all.Email.GetEmailByID(ctx, tx, id); err != nil {
			return err
		}
		if email.Data.Status != "queued" {
			return fmt.Errorf("%w: expected 'queued' got email.Status=%s", ErrEmailIsInWrongStatus, email.Data.Status)
		}
		email.Data.Status = "sending"
		return common4all.Email.UpdateEmail(ctx, tx, email)
	}, nil); err != nil {
		err = fmt.Errorf("failed to update email status to 'queued': %w", err)
		if dal.IsNotFound(err) {
			logus.Warningf(ctx, err.Error())
			return nil // Do not retry
		} else if errors.Is(err, ErrEmailIsInWrongStatus) {
			logus.Warningf(ctx, err.Error())
			return nil // Do not retry
		}
		logus.Errorf(ctx, err.Error())
		return err // Retry
	}

	var sentMessageID string
	emailMessage := emails.Email{
		From:    email.Data.From,
		To:      []string{email.Data.To},
		Subject: email.Data.Subject,
		Text:    email.Data.BodyText,
		HTML:    email.Data.BodyHtml,
	}
	if sentMessageID, err = SendEmail(ctx, emailMessage); err != nil {
		logus.Errorf(ctx, "Failed to send email: %v", err)

		if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
			if email, err = common4all.Email.GetEmailByID(ctx, tx, id); err != nil {
				return err
			}
			if email.Data.Status != "sending" {
				return fmt.Errorf("%w: expected 'sending' got email.Status=%s", ErrEmailIsInWrongStatus, email.Data.Status)
			}
			email.Data.Status = "error"
			email.Data.Error = err.Error()
			return common4all.Email.UpdateEmail(ctx, tx, email)
		}); err != nil {
			logus.Errorf(ctx, err.Error())
		}
		return nil // Do not retry
	}

	logus.Infof(ctx, "Sent email, message ContactID: %v", sentMessageID)

	if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if email, err = common4all.Email.GetEmailByID(ctx, tx, id); err != nil {
			return err
		}
		if email.Data.Status != "sending" {
			return fmt.Errorf("%w: expected 'sending' got email.Status=%s", ErrEmailIsInWrongStatus, email.Data.Status)
		}
		email.Data.Status = "sent"
		email.Data.DtSent = time.Now()
		email.Data.AwsSesMessageID = sentMessageID
		return common4all.Email.UpdateEmail(ctx, tx, email)
	}); err != nil {
		logus.Errorf(ctx, err.Error())
		err = nil // Do not retry!
	}
	return nil // Do not retry!
}
