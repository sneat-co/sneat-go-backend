package analytics

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/strongo/gamp"
	"github.com/strongo/logus"
	"net/http"
)

const (
	BASE_HOST = ".debtstracker.io"
)

const (
	EventCategoryReminders  = "reminders"
	EventActionReminderSent = "reminder-sent"
)

const (
	EventCategoryTransfers    = "transfers"
	EventActionDebtDueDateSet = "debt-due-date-set"
)

func SendSingleMessage(c context.Context, m gamp.Message) (err error) {
	if c == nil {
		return errors.New("Parameter 'c context.Context' is nil")
	}
	gaMeasurement := gamp.NewBufferedClient("", dtdal.HttpClient(c), nil)
	if err = gaMeasurement.Queue(m); err != nil {
		return err
	}
	if err = gaMeasurement.Flush(); err != nil {
		return err
	}
	var buffer bytes.Buffer
	_, _ = m.Write(&buffer)
	logus.Debugf(c, "Sent single message to GA: "+buffer.String())
	return nil
}

func getGaCommon(r *http.Request, userID string, userLanguage, platform string) gamp.Common {
	var userAgent string
	if r != nil {
		userAgent = r.UserAgent()
	} else {
		userAgent = "appengine"
	}

	return gamp.Common{
		TrackingID:    common.GA_TRACKING_ID,
		UserID:        userID,
		UserLanguage:  userLanguage,
		UserAgent:     userAgent,
		DataSource:    "backend",
		ApplicationID: "io.debtstracker.gae",
	}
}

func ReminderSent(c context.Context, userID string, userLanguage, platform string) {
	gaCommon := getGaCommon(nil, userID, userLanguage, platform)
	if err := SendSingleMessage(c, gamp.NewEvent(EventCategoryReminders, EventActionReminderSent, gaCommon)); err != nil {
		logus.Errorf(c, fmt.Errorf("failed to send even to GA: %w", err).Error())
	}
}

func ReceiptSentFromBot(whc botsfw.WebhookContext, channel string) error {
	ga := whc.GA()
	return ga.Queue(ga.GaEventWithLabel("receipts", "receipt-sent", channel))
}

func ReceiptSentFromApi(c context.Context, r *http.Request, userID string, userLanguage, platform, channel string) {
	gaCommon := getGaCommon(r, userID, userLanguage, platform)
	_ = SendSingleMessage(c, gamp.NewEventWithLabel(
		"receipts",
		"receipt-sent",
		channel,
		gaCommon,
	))
}
