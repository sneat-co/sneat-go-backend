package emailing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-core/emails"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/i18n"
)

func CreateEmailRecordAndQueueForSending(ctx context.Context, emailEntity *models4auth.EmailData) (id int64, err error) {
	var email models4auth.Email

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		emailEntity.Status = "queued"
		if email, err = facade4debtus.Email.InsertEmail(ctx, tx, emailEntity); err != nil {
			err = fmt.Errorf("%w: Failed to insert Email record", err)
			return err
		}
		if err = DelaySendEmail(ctx, email.ID); err != nil {
			err = fmt.Errorf("%w: Failed to delay sending", err)
		}
		return err
	}); err != nil {
		return
	}

	return email.ID, err
}

func GetEmailText(ctx context.Context, translator i18n.SingleLocaleTranslator, templateName string, templateParams interface{}) (string, error) {
	return common4debtus.TextTemplates.RenderTemplate(ctx, translator, templateName, templateParams)
}

func GetEmailHtml(ctx context.Context, translator i18n.SingleLocaleTranslator, templateName string, templateParams interface{}) (s string, err error) {
	var buffer bytes.Buffer
	err = common4debtus.HtmlTemplates.RenderTemplate(ctx, &buffer, translator, templateName, templateParams)
	return buffer.String(), err
}

var GetEmailClient = func(_ context.Context) (emails.Client, error) {
	panic("GetEmailClient is not initialized")
}

func SendEmail(ctx context.Context, email emails.Email) (messageID string, err error) {
	if email.Text == "" && email.HTML == "" {
		panic(`email.Text == "" && email.HTML == ""`)
	}
	var emailClient emails.Client
	emailClient, err = GetEmailClient(ctx)
	if err != nil {
		return "", err
	}
	var sent emails.Sent
	if sent, err = emailClient.Send(email); err != nil {
		return "", err
	}
	return sent.MessageID(), err
	//
	//var awsSession *session.Session
	//if awsSession, err = anybot.NewAwsSession(); err != nil {
	//	return
	//}
	//svc := ses.New(awsSession)
	//params := &ses.SendEmailInput{
	//	Destination: &ses.Destination{ // Required
	//		ToAddresses: []*string{
	//			aws.String(to), // Required
	//		},
	//	},
	//	Message: &ses.Message{ // Required
	//		Body: &ses.Body{ // Required
	//		},
	//		Subject: &ses.Content{ // Required
	//			Data:    aws.String(subject), // Required
	//			Charset: aws.String("utf-8"),
	//		},
	//	},
	//	Source: aws.String(from), // Required
	//	ReplyToAddresses: []*string{
	//		aws.String(from), // Required
	//		// More values...
	//	},
	//	//ReturnPath:    aws.String("Address"),
	//	//ReturnPathArn: aws.String("AmazonResourceName"),
	//	//SourceArn:     aws.String("AmazonResourceName"),
	//}
	//if bodyText != "" {
	//	params.Message.Body.Text = &ses.Content{
	//		Data:    aws.String(bodyText), // Required
	//		Charset: aws.String("utf-8"),
	//	}
	//}
	//if bodyHtml != "" {
	//	params.Message.Body.Html = &ses.Content{
	//		Data:    aws.String(bodyHtml), // Required
	//		Charset: aws.String("utf-8"),
	//	}
	//}
	//
	////http.DefaultClient = urlfetch.Client(ctx)
	////http.DefaultTransport = &urlfetch.Transport{Context: ctx, AllowInvalidServerCertificate: false}
	//logus.Debugf(ctx, "Sending email through AWS SES: %v", params)
	//
	//resp, err := svc.SendEmail(params)
	//
	//if err != nil {
	//	// Print the error, cast err to awserr.Error to get the ByCode and
	//	// Message from an error.
	//	originalErr := err
	//	errMessage := err.Error()
	//	if to != strings.ToLower(to) && strings.Index(errMessage, "Email address is not verified") > 0 && strings.Index(errMessage, to) > 0 {
	//		params.Destination.ToAddresses[0] = aws.String(strings.ToLower(to))
	//		resp, err = svc.SendEmail(params)
	//		if err != nil {
	//			logus.Errorf(ctx, "Failed to send ToLower(email): %v", err)
	//			return "", originalErr
	//		}
	//	} else {
	//		logus.Errorf(ctx, "Failed to send email using AWS SES: %v", err)
	//		return "", fmt.Errorf("failed to send email: %w", err)
	//	}
	//}
	//
	//// Pretty-print the response data.
	//logus.Debugf(ctx, "AWS SES output: %v", resp)
	//return *resp.MessageId, err
}
