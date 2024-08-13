package main

import (
	"github.com/bots-go-framework/bots-host-gae"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2"
)

func main() {
	logus.AddLogEntryHandler(logus.NewStandardGoLogger())
	gaeapp.Init(gae.BotHost())
	appengine.Main()
}
