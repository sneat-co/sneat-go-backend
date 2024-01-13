package healthcheck

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/emails"
	"net/http"
	"time"
)

// httpGetTestEmail sends a test email
func httpGetTestEmail(w http.ResponseWriter, _ *http.Request) {
	message := emails.Email{
		From:    "DailyScrum.app@sneat.team",
		To:      []string{"alexander.trakhimenok@gmail.com"},
		Subject: fmt.Sprintf("Hi, it's %v", time.Now()),
		Text:    "Howdy, is it time to sleep?",
		HTML:    "Howdy, is it <b>time to sleep</b>?",
		ReplyTo: nil,
	}
	output, err := emails.Send(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Failed to send email: %v", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Email sent, message ContactID: %v", output.MessageID())
}
