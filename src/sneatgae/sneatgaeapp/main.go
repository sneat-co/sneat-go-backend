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
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/modules"
	"github.com/sneat-co/sneat-go-backend/src/modules/healthcheck"
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp/pages"
	"github.com/sneat-co/sneat-go-core/emails"
	"github.com/sneat-co/sneat-go-core/extension"
	"github.com/sneat-co/sneat-go-core/monitoring"
	"github.com/strongo/delaying"
	"net/http"
)

func CreateHttpRouter() *httprouter.Router {
	return initHTTPRouter(globalOptionsHandler)
}

// Start to be called from main.go and will start an HTTP server using http.ListenAndServe
func Start(
	reportPanic monitoring.PanicCapturer,
	wrapHandler HandlerWrapper,
	httpRouter *httprouter.Router,
	emailClient emails.Client,
	extraModule ...extension.Config,
) {
	if reportPanic != nil {
		ReportPanic = reportPanic
	}
	if wrapHandler == nil {
		wrapHandler = noWrapper
	}

	initInfrastructure(emailClient)

	//botscore.InitializeBots(httpRouter) // TODO: should be part of module registration?

	// A shortcut to map handlers to httpRouter
	var handle = func(method, path string, handler http.HandlerFunc) { // TODO: change from HandlerFunc to Handler?
		httpRouter.HandlerFunc(method, path, wrapHTTPHandler(handler, wrapHandler))
	}

	initHtmlPageHandlers(handle)

	healthcheck.InitHealthCheck(handle)

	RegisterModules(handle, extraModule)

	// Ready to serve
	serve(httpRouter)
	//appengine.Main()
}

func initHtmlPageHandlers(handle extension.HTTPHandleFunc) {
	handle(http.MethodGet, "/", pages.IndexHandler)
}

func initInfrastructure(emailClient emails.Client) {
	logFirebaseEmulatorVars() // Connection to Firebase
	emails.Init(emailClient)
}

func RegisterModules(handle extension.HTTPHandleFunc, extraModule []extension.Config) {
	args := extension.NewModuleRegistrationArgs(handle, delaying.MustRegisterFunc)
	standardModules := modules.Extensions()
	for _, m := range standardModules {
		m.Register(args)
	}
	for _, m := range extraModule {
		m.Register(args)
	}
}
