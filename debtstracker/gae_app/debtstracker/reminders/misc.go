package reminders

import (
	"context"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/strongo/log"
	"google.golang.org/appengine/v2"
	"net/http"
)

type TransferReminderTo int

const (
	TransferReminderToCreator TransferReminderTo = iota
	TransferReminderToCounterparty
)

func TestEmail(c context.Context, w http.ResponseWriter, r *http.Request) {
	//sg_DEV_key := "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";
	//	sg_DEV_key := ;
	if err := SendEmail(c, "", "Testing SendGrid 2", "Simple Text"); err != nil {
		fmt.Fprint(w, err)
	}
}

func allowOrigin(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func InviteFriend(w http.ResponseWriter, r *http.Request) {
	allowOrigin(w)
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}
	fromName := r.Form["from_name"][0]

	if err := SendEmail(
		appengine.NewContext(r),
		fromName,
		"Check this app",
		"<p>See this phone app - <a href=https://debtstracker.io/#utm_medium=email&utm_campaign=invite-from-app>https://debtstracker.io/</a> - runs on iOS, Android & Windows Phone.</p>"+
			"<p>--<br>Sent by "+fromName+" from DebtsTracker.IO app</p>"); err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, "Email sent")
	}
}

func SendEmail(c context.Context, fromName, subject, html string) (err error) {
	if fromName == "" {
		fromName = "DebtsTracker.IO"
	}
	sgClient := sendgrid.NewSendClient("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")
	// set http.Client to use the appengine client
	sendgrid.DefaultClient.HTTPClient = dtdal.HttpClient(c) //Just perform this swap, and you are good to go.
	from := mail.NewEmail(fromName, "hello@debtstracker.io")
	to := mail.NewEmail("Example User", "test@example.com")
	message := mail.NewSingleEmail(from, subject, to, "", html)
	log.Infof(c, "Sending from %v email message: %v", fromName, html)
	_, err = sgClient.Send(message)
	return
}

func SendReceipt(c context.Context, w http.ResponseWriter, r *http.Request) {
	log.Infof(c, "sendReceipt() started")
	err := r.ParseForm()
	if err != nil {
		m := "Failed to parse form: %v"
		log.Infof(c, m, err)
		w.WriteHeader(500)
		fmt.Fprintf(w, m, err)
		return
	}
	log.Infof(c, "Form parsed: %v", r.FormValue("from_name"))
	fromName := r.Form.Get("from_name")
	//	fromEmail := r.Form["from_email"][0]
	//	toName := r.Form["to_name"][0]
	//	toEmail := r.Form["to_email"][0]
	amount := r.Form.Get("value")
	currency := r.Form.Get("currency")
	subject := "Receipt for friend's loan money transfer"
	message := "<p>You've got " + amount + currency + " from " + fromName + "</p><p>--<br>Sent via <a href='https://debtstracker.io/#utm_source=app&utm_medium=email&utm_campaign=receipt&utm_content=footer'><b>DebtsTracker.IO</b></a> - available at <a href=https://itunes.apple.com/en/app/debttracker-pro/id303497125>Apple AppStore</a> & <a href=https://play.google.com/store/apps/details?id=com.stellar.debtsfree&hl=en>Google Play</a></p>"
	allowOrigin(w)
	if err := SendEmail(c, fromName, subject, message); err != nil {
		log.Infof(c, "Failed to send email: %v", err)
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, "Email sent")
	}
}
