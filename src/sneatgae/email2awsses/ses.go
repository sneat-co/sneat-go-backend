package email2awsses

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sneat-co/sneat-go-core/emails"
)

// NewEmailClient creates new client
func NewEmailClient(awsRegion string, credentials *credentials.Credentials) emails.Client {
	return &sesClient{region: awsRegion, credentials: credentials}
}

type sesClient struct {
	region      string
	credentials *credentials.Credentials
}

// Send create and send text or html email to single recipient.
// @returns (resp emails.Sent, err error)
func (client *sesClient) Send(email emails.Email) (emails.Sent, error) {
	// start a new aws session
	config := &aws.Config{
		Credentials: client.credentials,
	}
	if client.region != "" {
		config.Region = aws.String(client.region)
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, emails.NewSendEmailError("failed to create AWS session", err)
	}

	// start a new ses session
	svc := ses.New(sess)

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{ // Required
			ToAddresses: make([]*string, len(email.To)),
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
			},
			Subject: &ses.Content{ // Required
				Data:    aws.String(email.Subject), // Required
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(email.From), // Required
	}
	if email.Text != "" {
		params.Message.Body.Text = &ses.Content{
			Data:    aws.String(email.Text), // Required
			Charset: aws.String("UTF-8"),
		}
	}

	if email.HTML != "" {
		params.Message.Body.Html = &ses.Content{
			Data:    aws.String(email.HTML), // Required
			Charset: aws.String("UTF-8"),
		}
	}

	for i, v := range email.To {
		params.Destination.ToAddresses[i] = &v
	}

	if len(email.ReplyTo) > 0 {
		params.ReplyToAddresses = make([]*string, 0, len(email.ReplyTo))
		for i, v := range email.ReplyTo {
			if v == "" {
				return nil, fmt.Errorf("empty ReplyTo[%v]", i)
			}
			params.ReplyToAddresses = append(params.ReplyToAddresses, &v)
		}
	}

	// send email

	output, err := svc.SendEmail(params)
	if err != nil {
		return nil, emails.NewSendEmailError("failed to send email throw AWS SES", err)
	}
	return sent{messageID: *output.MessageId}, nil
}
