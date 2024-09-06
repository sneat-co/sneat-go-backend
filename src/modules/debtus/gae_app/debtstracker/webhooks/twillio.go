package webhooks

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
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
	if db, err = facade.GetSneatDB(c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logus.Errorf(c, "Failed to get database: %v", err)
		return
	}
	err = db.RunReadwriteTransaction(c, func(tctx context.Context, tx dal.ReadwriteTransaction) error {
		var twilioSms = models4debtus.NewTwilioSms(smsSid, nil)
		err := tx.Get(tctx, twilioSms.Record)
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
			return tx.Set(tctx, twilioSms.Record)
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
