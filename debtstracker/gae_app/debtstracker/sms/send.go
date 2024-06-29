package sms

import (
	"context"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/strongo/gotwilio"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"strings"
)

func SendSms(c context.Context, isLive bool, toPhoneNumber, smsText string) (isTestSender bool, smsResponse *gotwilio.SmsResponse, twilioException *gotwilio.Exception, err error) {
	var (
		accountSid   string
		accountToken string
		fromNumber   string
	)

	if isLive {
		accountSid = common.TWILIO_LIVE_ACCOUNT_SID
		accountToken = common.TWILIO_LIVE_ACCOUNT_TOKEN
		fromNumber = common.TWILIO_LIVE_FROM_US
	} else {
		accountSid = common.TWILIO_TEST_ACCOUNT_SID
		accountToken = common.TWILIO_TEST_ACCOUNT_TOKEN
		fromNumber = common.TWILIO_TEST_FROM
	}

	twilio := gotwilio.NewTwilioClientCustomHTTP(accountSid, accountToken, dtdal.HttpClient(c))

	if smsResponse, twilioException, err = twilio.SendSMS(
		fromNumber,
		toPhoneNumber,
		smsText,
		"https://debtstracker-io.appspot.com/webooks/twilio/sms/status?sender=callback-url",
		common.TWILIO_APPLICATION_SID,
	); err != nil {
		return
	}

	if twilioException != nil && twilioException.Code == 21211 && len(toPhoneNumber) == 12 && strings.HasPrefix(toPhoneNumber, "+8") { // is not a valid phone number
		correctedPhoneNumber := strings.Replace(toPhoneNumber, "+8", "+7", 1)
		logus.Warningf(c, "%v. Will try to send after changing phone number from %v to %v", twilioException.Message, toPhoneNumber, correctedPhoneNumber)
		smsResponse, twilioException, err = twilio.SendSMS(
			fromNumber,
			correctedPhoneNumber,
			smsText,
			"https://debtstracker-io.appspot.com/webooks/twilio/sms/status?sender=callback-url",
			common.TWILIO_APPLICATION_SID,
		)
	}
	return
}

func TwilioExceptionToMessage(_ strongoapp.ExecutionContext, t i18n.SingleLocaleTranslator, ex *gotwilio.Exception) (messageText string, tryAnotherNumber bool) {
	switch ex.Code {
	case 21211: // Is not a valid phone number. https://www.twilio.com/docs/errors/21211
		tryAnotherNumber = true
		messageText = t.Translate(trans.MESSAGE_TEXT_INVALID_PHONE_NUMBER)
	case 21614: // Is is not a mobile number https://www.twilio.com/docs/errors/21614}
		tryAnotherNumber = true
		messageText = t.Translate(trans.MESSAGE_TEXT_PHONE_NUMBER_IS_NOT_SMS_CAPABLE)
	case 21612: // is not currently reachable using the 'From' phone number via SMS. https://www.twilio.com/docs/errors/21612
		tryAnotherNumber = true
		messageText = t.Translate("is not currently reachable using the 'From' phone number via SMS")
	case 21408: // Permission to send an SMS has not been enabled for the region indicated by the 'To' number: https://www.twilio.com/docs/errors/21408
		tryAnotherNumber = true
		messageText = t.Translate("Permission to send an SMS has not been enabled for the region indicated by the 'To' number")
	case 21610: // The message From/To pair violates a blacklist rule. https://www.twilio.com/docs/errors/21610
		tryAnotherNumber = true
		messageText = t.Translate("The message From/To pair violates a blacklist rule.")
	}
	return
}
