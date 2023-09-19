package sneatgaeapp

import (
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/email2awsses"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/emails"
	"github.com/sneat-co/sneat-go-core/emails/email2console"
	"os"
)

func initEmail() {
	var client emails.Client
	if core.IsInProd() || os.Getenv("SNEAT_SEND_EMAIL") == "true" {
		client = email2awsses.NewEmailClient("", nil)
	} else {
		client = email2console.NewClient()
	}
	emails.Init(client)
}
