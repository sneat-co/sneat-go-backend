package webhooks

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/logus"
	"net/http"
	"time"
)

func TwilioWebhook(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	err := r.ParseForm()
	if err != nil {
		logus.Errorf(c, "Failed to parse POST form: %v", err)
		return
	}
	logus.Infof(c, "BODY: %v", r.Form)
	smsSid := r.PostFormValue("SmsSid")
	messageStatus := r.PostFormValue("MessageStatus")

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logus.Errorf(c, "Failed to get database: %v", err)
		return
	}
	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		var twilioSms = models.NewTwilioSms(smsSid, nil)
		err := tx.Get(tc, twilioSms.Record)
		if err != nil {
			return err
		}
		if twilioSms.Data.Status != messageStatus {
			twilioSms.Data.Status = messageStatus
			switch messageStatus {
			case "sent":
				twilioSms.Data.DtSent = time.Now()
			case "delivered":
				twilioSms.Data.DtDelivered = time.Now()
			}
			return tx.Set(tc, twilioSms.Record)
		}
		return nil
	}, nil)

	if err != nil {
		if dal.IsNotFound(err) {
			logus.Infof(c, "Unknown SMS: %v", smsSid)
		} else {
			logus.Errorf(c, "Failed to process SMS update: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		logus.Infof(c, "Success")
		w.WriteHeader(http.StatusOK)
	}
}
