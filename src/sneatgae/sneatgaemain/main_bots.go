package main

import (
	"flag"
	gae "github.com/bots-go-framework/bots-host-gae"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/botscore"
	"github.com/sneat-co/sneat-go-core/facade"
	"log"
)

var withBotsFlag = flag.Bool("with-bots", false, "Run bots")

func initBots(httpRouter *httprouter.Router) {
	botHost := gae.BotHost()
	flag.Parse()
	args := flag.Args()
	log.Println("Args:", args)
	//isDevAppServer := strings.HasPrefix(os.Getenv("GCLOUD_PROJECT"), "demo-")
	//log.Println("ENV:\n\t", strings.Join(os.Environ(), "\n\t"))
	if /*!isDevAppServer ||*/ withBotsFlag != nil && *withBotsFlag {
		log.Println("Initializing bots...")

		botscore.GetDb = facade.GetDatabase
		botscore.InitializeBots(botHost, httpRouter) // TODO: should be part of module registration?
	}
}
