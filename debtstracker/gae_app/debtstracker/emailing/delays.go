package emailing

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/emails"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"time"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

const SEND_EMAIL_TASK = "send-email"

func DelaySendEmail(c context.Context, id int64) error {
	return delayEmail.EnqueueWork(c, delaying.With(common.QUEUE_EMAILS, SEND_EMAIL_TASK, 0), id)
}

var ErrEmailIsInWrongStatus = errors.New("email is already sending or sent")

func delayedSendEmail(c context.Context, id int64) (err error) {
	logus.Debugf(c, "delayedSendEmail(%v)", id)

	var email models.Email

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return err
	}

	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if email, err = dtdal.Email.GetEmailByID(c, tx, id); err != nil {
			return err
		}
		if email.Data.Status != "queued" {
			return fmt.Errorf("%w: expected 'queued' got email.Status=%s", ErrEmailIsInWrongStatus, email.Data.Status)
		}
		email.Data.Status = "sending"
		return dtdal.Email.UpdateEmail(c, tx, email)
	}, nil); err != nil {
		err = fmt.Errorf("failed to update email status to 'queued': %w", err)
		if dal.IsNotFound(err) {
			logus.Warningf(c, err.Error())
			return nil // Do not retry
		} else if errors.Is(err, ErrEmailIsInWrongStatus) {
			logus.Warningf(c, err.Error())
			return nil // Do not retry
		}
		logus.Errorf(c, err.Error())
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
	if sentMessageID, err = SendEmail(c, emailMessage); err != nil {
		logus.Errorf(c, "Failed to send email: %v", err)

		if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
			if email, err = dtdal.Email.GetEmailByID(c, tx, id); err != nil {
				return err
			}
			if email.Data.Status != "sending" {
				return fmt.Errorf("%w: expected 'sending' got email.Status=%s", ErrEmailIsInWrongStatus, email.Data.Status)
			}
			email.Data.Status = "error"
			email.Data.Error = err.Error()
			return dtdal.Email.UpdateEmail(c, tx, email)
		}); err != nil {
			logus.Errorf(c, err.Error())
		}
		return nil // Do not retry
	}

	logus.Infof(c, "Sent email, message ID: %v", sentMessageID)

	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if email, err = dtdal.Email.GetEmailByID(c, tx, id); err != nil {
			return err
		}
		if email.Data.Status != "sending" {
			return fmt.Errorf("%w: expected 'sending' got email.Status=%s", ErrEmailIsInWrongStatus, email.Data.Status)
		}
		email.Data.Status = "sent"
		email.Data.DtSent = time.Now()
		email.Data.AwsSesMessageID = sentMessageID
		return dtdal.Email.UpdateEmail(c, tx, email)
	}); err != nil {
		logus.Errorf(c, err.Error())
		err = nil // Do not retry!
	}
	return nil // Do not retry!
}
