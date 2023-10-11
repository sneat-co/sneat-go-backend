// Copyright 2020 https://sneat.app/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either logistus or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sneatgaeapp

import (
	"github.com/sneat-co/sneat-core-modules/contactus"
	"github.com/sneat-co/sneat-core-modules/invitus"
	"github.com/sneat-co/sneat-core-modules/teamus"
	"github.com/sneat-co/sneat-core-modules/userus"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp/bots"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp/pages"
	"github.com/sneat-co/sneat-go-core/modules"
	"github.com/sneat-co/sneat-go-modules/assetus"
	"github.com/sneat-co/sneat-go-modules/generic"
	"github.com/sneat-co/sneat-go-modules/healthcheck"
	"github.com/sneat-co/sneat-go-modules/listus"
	"github.com/sneat-co/sneat-go-modules/retrospectus"
	"github.com/sneat-co/sneat-go-modules/schedulus"
	"github.com/sneat-co/sneat-go-modules/scrumus"
	"github.com/sneat-co/sneat-go-modules/sportus"
	"github.com/strongo/log"
	golog "log"
	"net/http"
)

func Start() {
	defaultLogger := golog.Default()
	log.AddLogger(log.NewPrinter("log.Default()", func(format string, a ...any) (n int, err error) {
		defaultLogger.Printf(format, a...)
		return 0, nil
	}))

	initInfrastructure()

	httpRouter := initHTTPRouter(globalOptionsHandler)

	bots.InitializeBots(httpRouter) // TODO: should be part of module registration?

	// A shortcut to map handlers to httpRouter
	var handle = func(method, path string, handler http.HandlerFunc) {
		httpRouter.HandlerFunc(method, path, wrapHTTPHandlerFunc(handler))
	}

	initHtmlPageHandlers(handle)

	healthcheck.InitHealthCheck(handle)

	registerModules(handle)

	// Ready to serve
	serve(httpRouter)
	//appengine.Main()
}

func initHtmlPageHandlers(handle modules.HTTPHandleFunc) {
	handle(http.MethodGet, "/", pages.IndexHandler)
}

func initInfrastructure() {
	initSentry()   // Errors logging
	initFirebase() // Connection to Firebase
	initEmail()    // Settings for sending out emails
}

func registerModules(handle modules.HTTPHandleFunc) {
	args := modules.NewModuleRegistrationArgs(handle)
	mods := []modules.Module{
		userus.Module(),
		assetus.Module(),
		teamus.Module(),
		listus.Module(),
		schedulus.Module(),
		contactus.Module(),
		invitus.Module(),
		scrumus.Module(),
		retrospectus.Module(),
		sportus.Module(),
		generic.Module(),
	}
	for _, m := range mods {
		m.Register(args)
	}
}
